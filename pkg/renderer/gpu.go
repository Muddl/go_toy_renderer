//go:build !headless

package renderer

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/muddl/go_toy_renderer/pkg/camera"
	"github.com/muddl/go_toy_renderer/pkg/gpu"
	"github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/overlay"
	"github.com/muddl/go_toy_renderer/pkg/render"
	"github.com/muddl/go_toy_renderer/pkg/scene"
)

const gpuTargetFPS = 60

// GPUBackend opens a GLFW window (with NoAPI client) and renders via wgpu-native.
// The GPU init chain and geometry pipeline are delegated to pkg/gpu.Device.
// A live performance overlay (F3 toggle) is composited by pkg/overlay.OverlayPass.
type GPUBackend struct {
	width, height int
	window        *glfw.Window
	device        *gpu.Device
	ovr           overlay.Overlay
	ovLayer       *overlay.TextLayer
	ovPass        *overlay.OverlayPass
	sampler       overlay.Sampler
	gpuScene      *scene.Scene // optional; if set, used for multi-mesh rendering
	// Previous-frame timing used for overlay display (current frame not yet complete
	// when the overlay texture must be uploaded before Device.RenderFrame).
	prevFrameDur time.Duration
	prevPassDur  time.Duration
}

// SetScene configures the GPU backend to use a scene.Scene for multi-mesh
// rendering with per-mesh transforms. If set, RenderFrame ignores the
// render.Scene meshes and uses the scene nodes instead.
func (g *GPUBackend) SetScene(s *scene.Scene) {
	g.gpuScene = s
}

// Init opens a GLFW window with ClientAPI=NoAPI, extracts the platform-specific
// native window handle, and initialises pkg/gpu.Device (which owns the full
// wgpu chain: instance → surface → adapter → device → queue → pipeline).
// The overlay pass is created after the device is ready; failure is non-fatal.
func (g *GPUBackend) Init(width, height int) error {
	g.width = width
	g.height = height

	// Initialise GLFW. Note: GLFW requires a display server even with NoAPI.
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("glfw init: %w", err)
	}

	// Disable the OpenGL context — wgpu manages its own context.
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	win, err := glfw.CreateWindow(width, height, "go_toy_renderer (GPU)", nil, nil)
	if err != nil {
		glfw.Terminate()
		return fmt.Errorf("create window: %w", err)
	}
	g.window = win

	win.SetKeyCallback(func(w *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		if action != glfw.Press {
			return
		}
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
		case glfw.KeyF3:
			g.ovr.CycleLevel()
		}
	})

	// Extract the platform-specific native window handle (HWND on Windows,
	// X11 Display+Window on Linux, CAMetalLayer on macOS).
	handle := getNativeWindowHandle(win)

	// Delegate the full wgpu chain (instance/surface/adapter/device/queue/pipeline)
	// to pkg/gpu.Device. On failure, Shutdown cleans up partial resources.
	g.device = gpu.New()
	if err := g.device.Init(uint32(width), uint32(height), handle); err != nil {
		g.device.Shutdown()
		glfw.Terminate()
		return fmt.Errorf("gpu device init: %w", err)
	}

	// Allocate the overlay text layer and attempt to create the GPU overlay
	// pipeline. Overlay failure is non-fatal — the renderer continues without it.
	g.ovLayer = overlay.NewTextLayer(width, height)
	g.ovPass = overlay.NewOverlayPass(g.device)
	if err := g.ovPass.CreateOverlayPipeline(); err != nil {
		// Overlay unavailable (e.g., shader compile error); disable gracefully.
		g.ovPass = nil
	} else {
		g.device.SetOverlayRenderer(g.ovPass)
	}

	return nil
}

// defaultCamera returns the default camera matching the Phase 12 hardcoded
// camera parameters: pos=(3,2,5) looking at origin, fov=60 deg.
func defaultCamera(width, height int) camera.Camera {
	return camera.New(
		math.Vec3{X: 3, Y: 2, Z: 5}, // position
		math.Vec3{},                   // target = origin
		math.Vec3{Y: 1},              // up
		60.0,                          // fov degrees
		float64(width)/float64(height), // aspect
		0.1,                           // near
		100.0,                         // far
	)
}

// RenderFrame uploads the first mesh in scene to the GPU (cached after the first
// frame), uploads camera and mesh uniforms, renders the geometry pass (with overlay
// composited into the same pass by pkg/gpu.Device if SetOverlayRenderer was called),
// and presents the frame.
// Returns ErrWindowClosed when the window has been dismissed.
func (g *GPUBackend) RenderFrame(scene *render.Scene) error {
	if g.window == nil {
		return fmt.Errorf("gpu backend: not initialized — call Init first")
	}
	if g.window.ShouldClose() {
		return ErrWindowClosed
	}

	frameStart := time.Now()

	// --- Camera uniform (once per frame) ---
	cam := defaultCamera(g.width, g.height)
	g.device.UpdateCameraUniforms(cam.ViewProjectionMatrixWebGPU())

	// --- Build draw nodes ---
	geoStart := time.Now()
	var nodes []gpu.DrawNode
	if g.gpuScene != nil {
		nodes = make([]gpu.DrawNode, 0, len(g.gpuScene.Nodes))
		for i := range g.gpuScene.Nodes {
			n := &g.gpuScene.Nodes[i]
			nodes = append(nodes, gpu.DrawNode{
				Mesh:  n.Mesh,
				Model: n.Transform.ModelMatrix(),
			})
		}
	} else {
		nodes = make([]gpu.DrawNode, 0, len(scene.Meshes))
		for _, mesh := range scene.Meshes {
			nodes = append(nodes, gpu.DrawNode{
				Mesh:  mesh,
				Model: math.NewIdentity(),
			})
		}
	}
	geoElapsed := time.Since(geoStart)

	verts, tris := 0, 0
	for i := range nodes {
		if nodes[i].Mesh != nil {
			verts += len(nodes[i].Mesh.Vertices)
			tris += len(nodes[i].Mesh.Indices) / 3
		}
	}

	// Build overlay metrics from the PREVIOUS frame's timing so that meaningful
	// values are available before Device.RenderFrame (which includes GPU work and
	// present). On frame 0 both durations are zero — acceptable for one frame.
	msPerDur := func(d time.Duration) float64 { return float64(d) / float64(time.Millisecond) }
	m := overlay.Metrics{
		FPS:              g.sampler.FPS(),
		FrameTimeMS:      msPerDur(g.prevFrameDur),
		CPUFrameTimeMS:   msPerDur(g.prevFrameDur - g.prevPassDur),
		GeometryUploadMS: msPerDur(geoElapsed),
		RenderPassMS:     msPerDur(g.prevPassDur),
		Backend:          "GPU",
		VertexCount:      verts,
		TriangleCount:    tris,
	}

	// Render the overlay text layer into the CPU pixel buffer, then upload to
	// the GPU texture so it is ready before Device.RenderFrame samples it.
	if g.ovPass != nil {
		g.ovLayer.Render(g.ovr.Level(), m)
		_ = g.ovPass.UpdateOverlayTexture(g.ovLayer)
	}

	// --- GPU render pass (geometry + overlay composited in same pass) ---
	passStart := time.Now()
	if err := g.device.RenderFrameMulti(nodes); err != nil {
		return err
	}
	g.prevPassDur = time.Since(passStart)

	glfw.PollEvents()

	if elapsed := time.Since(frameStart); elapsed < time.Second/gpuTargetFPS {
		time.Sleep(time.Second/gpuTargetFPS - elapsed)
	}

	// Record full frame duration (CPU + GPU + sleep) for the sampler and for
	// display on the next frame.
	g.prevFrameDur = time.Since(frameStart)
	g.sampler.AddSample(g.prevFrameDur)

	return nil
}

// Shutdown releases GPU resources, the overlay pass, and terminates GLFW.
func (g *GPUBackend) Shutdown() {
	if g.ovPass != nil {
		g.ovPass.Release()
		g.ovPass = nil
	}
	if g.device != nil {
		g.device.Shutdown()
		g.device = nil
	}
	glfw.Terminate()
}

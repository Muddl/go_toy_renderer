//go:build !headless

package renderer

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/muddl/go_toy_renderer/pkg/gpu"
	"github.com/muddl/go_toy_renderer/pkg/overlay"
	"github.com/muddl/go_toy_renderer/pkg/render"
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

// RenderFrame uploads the first mesh in scene to the GPU (cached after the first
// frame), renders the geometry pass (with overlay composited into the same pass
// by pkg/gpu.Device if SetOverlayRenderer was called), and presents the frame.
// Returns ErrWindowClosed when the window has been dismissed.
func (g *GPUBackend) RenderFrame(scene *render.Scene) error {
	if g.window == nil {
		return fmt.Errorf("gpu backend: not initialized — call Init first")
	}
	if g.window.ShouldClose() {
		return ErrWindowClosed
	}

	frameStart := time.Now()

	// --- Geometry upload ---
	geoStart := time.Now()
	if len(scene.Meshes) > 0 {
		if err := g.device.LoadGeometry(scene.Meshes[0]); err != nil {
			return fmt.Errorf("gpu backend: load geometry: %w", err)
		}
	}
	geoElapsed := time.Since(geoStart)

	g.sampler.AddSample(time.Since(frameStart))

	verts, tris := 0, 0
	for _, mesh := range scene.Meshes {
		verts += len(mesh.Vertices)
		tris += len(mesh.Indices) / 3
	}

	// Build overlay metrics with data available before the render pass.
	// RenderPassMS is approximated from the previous frame (zero on first frame).
	m := overlay.Metrics{
		FPS:              g.sampler.FPS(),
		FrameTimeMS:      float64(time.Since(frameStart)) / float64(time.Millisecond),
		CPUFrameTimeMS:   float64(time.Since(frameStart)) / float64(time.Millisecond),
		GeometryUploadMS: float64(geoElapsed) / float64(time.Millisecond),
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
	if err := g.device.RenderFrame(); err != nil {
		return err
	}
	_ = time.Since(passStart) // RenderPassMS available for next frame's metrics

	glfw.PollEvents()

	if elapsed := time.Since(frameStart); elapsed < time.Second/gpuTargetFPS {
		time.Sleep(time.Second/gpuTargetFPS - elapsed)
	}
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

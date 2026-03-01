//go:build !headless

package renderer

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/muddl/go_toy_renderer/pkg/gpu"
	"github.com/muddl/go_toy_renderer/pkg/render"
)

const gpuTargetFPS = 60

// GPUBackend opens a GLFW window (with NoAPI client) and renders via wgpu-native.
// The GPU init chain and Hello Triangle are delegated to pkg/gpu.Device.
type GPUBackend struct {
	width, height int
	window        *glfw.Window
	device        *gpu.Device
}

// Init opens a GLFW window with ClientAPI=NoAPI, extracts the platform-specific
// native window handle, and initialises pkg/gpu.Device (which owns the full
// wgpu chain: instance → surface → adapter → device → queue → pipeline).
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
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
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

	return nil
}

// RenderFrame uploads the first mesh in scene to the GPU (cached after the first
// frame) and renders one indexed draw call. Returns ErrWindowClosed when the
// window has been dismissed.
func (g *GPUBackend) RenderFrame(scene *render.Scene) error {
	if g.window.ShouldClose() {
		return ErrWindowClosed
	}

	frameStart := time.Now()

	// Upload geometry from the first scene mesh (multi-mesh support in Phase 14).
	if len(scene.Meshes) > 0 {
		if err := g.device.LoadGeometry(scene.Meshes[0]); err != nil {
			return fmt.Errorf("gpu backend: load geometry: %w", err)
		}
	}

	if err := g.device.RenderFrame(); err != nil {
		return err
	}

	glfw.PollEvents()

	if elapsed := time.Since(frameStart); elapsed < time.Second/gpuTargetFPS {
		time.Sleep(time.Second/gpuTargetFPS - elapsed)
	}
	return nil
}

// Shutdown releases GPU resources and terminates GLFW.
func (g *GPUBackend) Shutdown() {
	if g.device != nil {
		g.device.Shutdown()
		g.device = nil
	}
	glfw.Terminate()
}

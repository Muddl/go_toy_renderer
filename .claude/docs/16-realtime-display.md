# Real-time Display

**Status:** 📋 Planned — Phase 9 (Window & Real-time Display)

---

## Overview

Phase 9 replaces the single-shot PNG output with a live windowed renderer at 60 fps. The original `cmd/renderer` (PNG output) is preserved unchanged. A new entry point `cmd/renderer-rt` opens a window and runs a frame loop.

The windowing layer is implemented using **GLFW** (`github.com/go-gl/glfw/v3.3/glfw`), which is cross-platform (Windows, macOS, Linux), integrates directly with wgpu surface creation, and supports keyboard/mouse input.

---

## Package Layout

```
pkg/window/
├── window.go     # Window struct, Open, Close, Run, Resize callback
├── input.go      # Keyboard and mouse state polling
└── camera.go     # First-person camera controller (WASD + mouse look)

cmd/renderer-rt/
└── main.go       # Entry point: create window, renderer, run loop
```

---

## Window Lifecycle

```go
// pkg/window/window.go
type Window struct {
    handle *glfw.Window
    width  int
    height int
}

type Config struct {
    Title   string
    Width   int
    Height  int
    Vsync   bool
}

func Open(cfg Config) (*Window, error)
func (w *Window) Close()
func (w *Window) Run(onUpdate func(dt float64), onRender func()) // blocks until closed
func (w *Window) Size() (width, height int)
func (w *Window) ShouldClose() bool
```

**Callbacks registered at startup:**
- `SetFramebufferSizeCallback` → calls `renderer.Resize(w, h)` and recreates depth texture
- `SetKeyCallback` → populates `input.KeyState`
- `SetCursorPosCallback` → accumulates mouse delta for camera look
- `SetScrollCallback` → camera zoom

---

## Frame Loop

```go
// cmd/renderer-rt/main.go

win, _ := window.Open(window.Config{Title: "go_toy_renderer", Width: 1280, Height: 720})
defer win.Close()

renderer, _ := renderer.New(renderer.Auto, 1280, 720)
defer renderer.Destroy()

cam := camera.New( /* starting position */ )
ctrl := window.NewCameraController(cam, moveSpeed: 5.0, lookSensitivity: 0.002)

prevTime := glfw.GetTime()

win.Run(
    func(dt float64) {                     // onUpdate
        ctrl.Update(win.Input(), dt)       // move camera
    },
    func() {                               // onRender
        if err := renderer.Render(scene, cam); err != nil {
            log.Fatal(err)
        }
    },
)
```

---

## CPU Blit Path (Phase 9 Bridge)

Before the GPU pipeline is ready (Phase 11+), the CPU renderer runs inside the window using a blit path:

1. `CPURenderer.Render()` writes to `framebuffer.Framebuffer` (CPU memory)
2. Pixel data is uploaded to a `wgpu.Texture` via `Queue.WriteTexture()`
3. A full-screen triangle pass samples the texture and blits it to the swap chain back buffer
4. No rasterization happens on the GPU — it is only a presentation blit

This provides immediate feedback (real-time camera navigation with CPU rendering) before any GPU geometry pipeline work begins.

---

## Camera Controller (pkg/window/camera.go)

```go
type CameraController struct {
    cam             *camera.Camera
    moveSpeed       float64 // m/s
    lookSensitivity float64 // rad/pixel
    yaw, pitch      float64 // accumulated angles
}

func (cc *CameraController) Update(input *Input, dt float64) {
    // WASD: move forward/back/left/right
    // Q/E or Space/Ctrl: move up/down
    // Mouse delta: rotate yaw (left/right) and pitch (up/down)
    // Clamp pitch to [-89°, 89°] to avoid gimbal lock
}
```

**Input bindings:**

| Key / Input | Action |
|-------------|--------|
| W / S       | Move forward / back |
| A / D       | Strafe left / right |
| Space / Ctrl | Move up / down |
| Mouse drag (RMB held) | Look around |
| Scroll wheel | Zoom (adjust FOV) |
| ESC         | Release mouse, or close window |
| F11         | Toggle fullscreen |

---

## Frame Timing

```go
// Inside Run()
now := glfw.GetTime()
dt  := now - prevTime
prevTime = now

// FPS counter: update title every second
frameCount++
if now-lastTitleUpdate > 1.0 {
    win.handle.SetTitle(fmt.Sprintf("go_toy_renderer — %.1f fps / %.2f ms", fps, ms))
    lastTitleUpdate = now
    frameCount = 0
}
```

**Vsync:** `glfw.SwapInterval(1)` enables vsync (default). Pass `--no-vsync` flag to disable for benchmarking.

---

## Window Resize Handling

When the framebuffer size changes:

1. GLFW fires `FramebufferSizeCallback` with new `(w, h)`
2. `renderer.Resize(w, h)` is called:
   - CPU path: reallocates `framebuffer.Framebuffer`
   - GPU path: reconfigures swap chain, recreates depth texture
3. Camera aspect ratio is updated: `cam.Aspect = float64(w) / float64(h)`
4. Next `Render()` call uses the new size

---

## Backend Selection Flag

```
go run ./cmd/renderer-rt [flags]

Flags:
  --backend string   Rendering backend: cpu, gpu, or auto (default "auto")
  --width  int       Window width  (default 1280)
  --height int       Window height (default 720)
  --no-vsync         Disable vsync (for benchmarking)
  --msaa   int       MSAA sample count: 1, 4 (default 1, requires GPU)
```

---

## Testing

Window tests are challenging to automate in CI (no display server). Strategy:

1. **Headless mode** — `--headless` flag renders N frames off-screen (to a texture), then saves the last frame as PNG for comparison against the CPU golden image
2. **Window unit tests** — test input processing and camera update logic without opening a window (inject synthetic input events)
3. **CI gate** — CI runs headless mode on Linux with a virtual framebuffer (`Xvfb`) or uses the off-screen wgpu adapter

---

## Common Gotchas

- **GLFW must run on the main thread.** All GLFW calls (including window creation, event polling, and `SwapBuffers`) must happen on the goroutine that called `glfw.Init()`. Use `runtime.LockOSThread()` in `main()`.
- **macOS Cocoa threading:** On macOS, GLFW requires `runtime.LockOSThread()` and the `-tags cocoa` build constraint. Ensure this is documented in the build prerequisites.
- **Mouse capture:** When right-clicking to look around, capture the cursor with `glfw.CursorDisabled` to prevent it from reaching window borders. Release on ESC.
- **Window vs framebuffer size:** On high-DPI (Retina) displays, `GetWindowSize()` returns logical pixels but `GetFramebufferSize()` returns physical pixels. Always use framebuffer size for viewport/swap chain dimensions.
- **Swap chain recreation:** On Windows, any swap chain resize or window move can invalidate the swap chain. Handle `wgpu.SurfaceGetCurrentTexture` returning `SuboptimalSurface` by immediately reconfiguring.

---

## References

- [go-gl/glfw](https://github.com/go-gl/glfw) — Go bindings for GLFW 3
- [GLFW Documentation](https://www.glfw.org/docs/latest/)
- [wgpu Surface](https://docs.rs/wgpu/latest/wgpu/struct.Surface.html) — swap chain reference

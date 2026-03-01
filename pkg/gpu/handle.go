// Package gpu manages the wgpu initialization chain and GPU render pipeline.
// It wraps go-webgpu/webgpu (Zero-CGo FFI) to provide Instance → Adapter →
// Device → Queue → Surface → RenderPipeline setup and a Hello Triangle render
// path.
//
// GPU integration tests are gated by the GPU_TESTS=1 environment variable.
// The shared library must be loadable via WGPU_NATIVE_PATH at runtime.
// In headless (CI) builds, the Device is a no-op stub and goffi is never
// imported, allowing the package to compile with CGO_ENABLED=1.
package gpu

// NativeWindowHandle carries the platform-specific window handle values
// needed to create a wgpu surface. Populate the fields for your platform and
// pass to Device.Init. All uintptr fields are safe to pass across CGo
// boundaries as they are just integer-sized values.
type NativeWindowHandle struct {
	// Windows (Win32): HWND from GLFW GetWin32Window.
	HWND uintptr
	// Linux (X11): Display pointer and Window XID from GLFW GetX11Display/GetX11Window.
	X11Display uintptr
	X11Window  uint64
	// macOS: CAMetalLayer pointer obtained from the NSView backing layer.
	MetalLayer uintptr
}

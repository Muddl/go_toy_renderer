//go:build darwin && !headless

package renderer

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/muddl/go_toy_renderer/pkg/gpu"
)

// getNativeWindowHandle extracts the CAMetalLayer pointer from the GLFW window.
// GetCocoaWindow returns an NSWindow pointer; on macOS, wgpu requires a
// CAMetalLayer which must be obtained from the NSWindow's content view backing
// layer. A full Objective-C bridge is deferred to Phase 12 (macOS support);
// for now the NSWindow pointer is passed directly as a placeholder.
func getNativeWindowHandle(win *glfw.Window) gpu.NativeWindowHandle {
	return gpu.NativeWindowHandle{
		MetalLayer: uintptr(win.GetCocoaWindow()),
	}
}

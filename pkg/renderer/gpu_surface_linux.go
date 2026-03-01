//go:build linux && !headless

package renderer

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/muddl/go_toy_renderer/pkg/gpu"
)

// getNativeWindowHandle extracts the X11 Display pointer and Window XID from
// the GLFW window for wgpu surface creation on Linux.
func getNativeWindowHandle(win *glfw.Window) gpu.NativeWindowHandle {
	return gpu.NativeWindowHandle{
		X11Display: uintptr(unsafe.Pointer(glfw.GetX11Display())),
		X11Window:  uint64(win.GetX11Window()),
	}
}

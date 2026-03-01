//go:build windows && !headless

package renderer

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/muddl/go_toy_renderer/pkg/gpu"
)

// getNativeWindowHandle extracts the Win32 HWND from the GLFW window.
// The HWND is obtained via GLFW's GetWin32Window and cast to uintptr for
// safe cross-CGo-boundary transfer to pkg/gpu.
func getNativeWindowHandle(win *glfw.Window) gpu.NativeWindowHandle {
	return gpu.NativeWindowHandle{
		HWND: uintptr(unsafe.Pointer(win.GetWin32Window())),
	}
}

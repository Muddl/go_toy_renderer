//go:build windows

package gpu

import "github.com/go-webgpu/webgpu/wgpu"

// createPlatformSurface creates a wgpu surface from a Win32 HWND.
// handle.HWND must be non-zero (obtained from GLFW GetWin32Window).
func createPlatformSurface(inst *wgpu.Instance, handle NativeWindowHandle) (*wgpu.Surface, error) {
	return inst.CreateSurfaceFromWindowsHWND(0, handle.HWND)
}

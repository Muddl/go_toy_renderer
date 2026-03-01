//go:build windows && !headless

package gpu

import (
	"fmt"

	"github.com/go-webgpu/webgpu/wgpu"
)

// createPlatformSurface creates a wgpu surface from a Win32 HWND.
// handle.HWND must be non-zero (obtained from GLFW GetWin32Window).
func createPlatformSurface(inst *wgpu.Instance, handle NativeWindowHandle) (*wgpu.Surface, error) {
	s, err := inst.CreateSurfaceFromWindowsHWND(0, handle.HWND)
	if err != nil {
		return nil, fmt.Errorf("create surface from Windows HWND: %w", err)
	}
	return s, nil
}

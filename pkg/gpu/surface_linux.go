//go:build linux

package gpu

import (
	"fmt"

	"github.com/go-webgpu/webgpu/wgpu"
)

// createPlatformSurface creates a wgpu surface from an X11 Display and Window.
// handle.X11Display must be a non-nil *C.Display pointer (obtained from GLFW
// GetX11Display) and handle.X11Window must be the XID (from GetX11Window).
func createPlatformSurface(inst *wgpu.Instance, handle NativeWindowHandle) (*wgpu.Surface, error) {
	s, err := inst.CreateSurfaceFromXlibWindow(handle.X11Display, handle.X11Window)
	if err != nil {
		return nil, fmt.Errorf("create surface from X11 window: %w", err)
	}
	return s, nil
}

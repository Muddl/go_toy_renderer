//go:build darwin && !headless

package gpu

import (
	"fmt"

	"github.com/go-webgpu/webgpu/wgpu"
)

// createPlatformSurface creates a wgpu surface from a CAMetalLayer pointer.
// handle.MetalLayer must be the result of calling createMetalLayer in
// pkg/renderer (requires Objective-C bridge to create CAMetalLayer from NSView).
func createPlatformSurface(inst *wgpu.Instance, handle NativeWindowHandle) (*wgpu.Surface, error) {
	s, err := inst.CreateSurfaceFromMetalLayer(handle.MetalLayer)
	if err != nil {
		return nil, fmt.Errorf("create surface from Metal layer: %w", err)
	}
	return s, nil
}

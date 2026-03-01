//go:build darwin

package gpu

import "github.com/go-webgpu/webgpu/wgpu"

// createPlatformSurface creates a wgpu surface from a CAMetalLayer pointer.
// handle.MetalLayer must be the result of calling createMetalLayer in
// pkg/renderer (requires Objective-C bridge to create CAMetalLayer from NSView).
func createPlatformSurface(inst *wgpu.Instance, handle NativeWindowHandle) (*wgpu.Surface, error) {
	return inst.CreateSurfaceFromMetalLayer(handle.MetalLayer)
}

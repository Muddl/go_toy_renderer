//go:build headless

package gpu

// OverlayRenderer is a no-op interface in headless builds.
// The real interface (using wgpu.RenderPassEncoder) lives in overlay.go.
type OverlayRenderer interface{}

// SetOverlayRenderer is a no-op in headless builds — no GPU device is available.
func (d *Device) SetOverlayRenderer(_ OverlayRenderer) {}

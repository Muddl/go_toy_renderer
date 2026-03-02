//go:build !headless

package gpu

import (
	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
)

// OverlayRenderer is an optional interface for compositing a debug overlay
// into the active render pass after the main geometry draw call.
// Implemented by pkg/overlay.OverlayPass (!headless builds).
type OverlayRenderer interface {
	RenderIntoPass(pass *wgpu.RenderPassEncoder) error
}

// SetOverlayRenderer installs an optional overlay renderer called after the
// geometry DrawIndexed in each frame's render pass. Pass nil to remove.
func (d *Device) SetOverlayRenderer(r OverlayRenderer) { d.overlayRenderer = r }

// DeviceHandle returns the underlying wgpu.Device for use by extension
// packages (e.g. pkg/overlay). Returns nil if Init has not completed.
func (d *Device) DeviceHandle() *wgpu.Device { return d.device }

// QueueHandle returns the underlying wgpu.Queue for use by extension
// packages. Returns nil if Init has not completed.
func (d *Device) QueueHandle() *wgpu.Queue { return d.queue }

// SurfaceFormat returns the swap-chain texture format chosen during Init.
// Returns TextureFormatUndefined if Init has not completed.
func (d *Device) SurfaceFormat() gputypes.TextureFormat { return d.format }

// Dimensions returns the framebuffer width and height in pixels as set by Init.
func (d *Device) Dimensions() (width, height uint32) { return d.width, d.height }

// Package gpu manages the wgpu initialization chain and GPU render pipeline.
// It wraps go-webgpu/webgpu (Zero-CGo FFI) to provide Instance → Adapter →
// Device → Queue → Surface → RenderPipeline setup and a Hello Triangle render
// path.
//
// GPU integration tests are gated by the GPU_TESTS=1 environment variable.
// The shared library must be loadable via WGPU_NATIVE_PATH at runtime.
package gpu

import (
	"errors"

	"github.com/go-webgpu/webgpu/wgpu"
)

// Device holds the wgpu initialization chain and render resources.
// Create with [New] and initialise with [Device.Init] before calling
// [Device.RenderFrame].
type Device struct {
	instance *wgpu.Instance
	adapter  *wgpu.Adapter
	device   *wgpu.Device
	queue    *wgpu.Queue
	surface  *wgpu.Surface
	pipeline *wgpu.RenderPipeline
	width    uint32
	height   uint32
	ready    bool
}

// New returns a new, uninitialised Device.
func New() *Device { return &Device{} }

// IsReady returns true if Init completed successfully and the device is ready to render.
func (d *Device) IsReady() bool { return d.ready }

// HasQueue returns true if the wgpu queue was acquired during Init.
// Used in integration tests to verify the init chain ran up to queue acquisition.
func (d *Device) HasQueue() bool { return d.queue != nil }

// Init initialises wgpu and prepares the Device for rendering.
// surface must be a valid platform surface obtained from the GLFW window.
// width and height are the initial swap-chain dimensions in pixels.
//
// Phase 2 implements the full init chain; this stub returns nil.
func (d *Device) Init(width, height uint32, surface *wgpu.Surface) error {
	d.width = width
	d.height = height
	d.surface = surface
	d.ready = true
	return nil
}

// RenderFrame submits one frame to the GPU.
// Returns an error if Init has not been called successfully.
//
// Phase 3 implements the Hello Triangle render pass; this stub returns an
// error until the pipeline is wired up.
func (d *Device) RenderFrame() error {
	if !d.ready {
		return errors.New("gpu: Device not initialised — call Init first")
	}
	return errors.New("gpu: render pipeline not yet implemented (Phase 3)")
}

// Shutdown releases all wgpu resources held by the Device.
// Safe to call on an uninitialised Device and to call multiple times.
func (d *Device) Shutdown() {
	if d.pipeline != nil {
		d.pipeline.Release()
		d.pipeline = nil
	}
	if d.surface != nil {
		d.surface.Release()
		d.surface = nil
	}
	if d.queue != nil {
		d.queue.Release()
		d.queue = nil
	}
	if d.device != nil {
		d.device.Release()
		d.device = nil
	}
	if d.adapter != nil {
		d.adapter.Release()
		d.adapter = nil
	}
	if d.instance != nil {
		d.instance.Release()
		d.instance = nil
	}
	d.ready = false
}

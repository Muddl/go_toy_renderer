//go:build headless

package gpu

import (
	"errors"

	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// Device is a compile-time stub for headless (CI) builds where no GPU or
// wgpu library is available. The real Device (requiring CGO_ENABLED=0 via
// goffi) is in gpu.go and compiled only without the headless build tag.
type Device struct{}

// New returns a new Device stub. In headless builds this is always a no-op.
func New() *Device { return &Device{} }

// IsReady always returns false in headless builds.
func (d *Device) IsReady() bool { return false }

// HasQueue always returns false in headless builds.
func (d *Device) HasQueue() bool { return false }

// Init returns an error in headless builds — no GPU or wgpu library present.
func (d *Device) Init(_, _ uint32, _ NativeWindowHandle) error {
	return errors.New("gpu: not available in headless builds")
}

// RenderFrame returns an error in headless builds.
func (d *Device) RenderFrame() error {
	return errors.New("gpu: not available in headless builds")
}

// LoadGeometry returns an error in headless builds — no GPU device available.
func (d *Device) LoadGeometry(_ *geometry.Mesh) error {
	return errors.New("gpu: not available in headless builds")
}

// UpdateCameraUniforms is a no-op in headless builds.
func (d *Device) UpdateCameraUniforms(_ math.Mat4x4) {}

// UpdateMeshUniforms is a no-op in headless builds.
func (d *Device) UpdateMeshUniforms(_ math.Mat4x4) {}

// Shutdown is a no-op in headless builds.
func (d *Device) Shutdown() {}

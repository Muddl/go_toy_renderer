//go:build !headless

package overlay_test

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/gpu"
	"github.com/muddl/go_toy_renderer/pkg/overlay"
)

// TestOverlayPass_CreateOverlayPipeline_ErrorBeforeDeviceInit verifies that
// CreateOverlayPipeline returns an error when the gpu.Device has not been
// initialized (IsReady == false). No GPU hardware is required.
func TestOverlayPass_CreateOverlayPipeline_ErrorBeforeDeviceInit(t *testing.T) {
	d := gpu.New()
	op := overlay.NewOverlayPass(d)
	if err := op.CreateOverlayPipeline(); err == nil {
		t.Fatal("CreateOverlayPipeline() before Device.Init should return an error")
	}
}

// TestOverlayPass_UpdateOverlayTexture_ErrorWithNilLayer verifies that
// UpdateOverlayTexture returns an error when passed a nil TextLayer.
func TestOverlayPass_UpdateOverlayTexture_ErrorWithNilLayer(t *testing.T) {
	d := gpu.New()
	op := overlay.NewOverlayPass(d)
	if err := op.UpdateOverlayTexture(nil); err == nil {
		t.Fatal("UpdateOverlayTexture(nil) should return an error")
	}
}

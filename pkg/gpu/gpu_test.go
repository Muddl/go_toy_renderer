package gpu_test

import (
	"os"
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/gpu"
)

// TestDevice_New_ReturnsNonNil verifies New() returns a usable Device pointer.
func TestDevice_New_ReturnsNonNil(t *testing.T) {
	d := gpu.New()
	if d == nil {
		t.Fatal("gpu.New() returned nil")
	}
}

// TestDevice_Init_SkipsWithoutGPUTests verifies Init is skipped when GPU_TESTS=1
// is not set. This gate prevents GPU init in CI environments without hardware.
func TestDevice_Init_SkipsWithoutGPUTests(t *testing.T) {
	if os.Getenv("GPU_TESTS") != "1" {
		t.Skip("set GPU_TESTS=1 to run GPU integration tests")
	}
	// Full GPU init test is in the Phase 2 integration test.
	// This placeholder ensures the skip guard is exercised.
	t.Fatal("GPU_TESTS=1 set but no GPU fixture configured in Phase 1")
}

// TestDevice_Shutdown_IsIdempotent verifies Shutdown can be called on an
// uninitialised Device without panicking.
func TestDevice_Shutdown_IsIdempotent(t *testing.T) {
	d := gpu.New()
	d.Shutdown() // must not panic
	d.Shutdown() // second call also safe
}

// TestDevice_RenderFrame_ReturnsErrorBeforeInit verifies that RenderFrame
// returns an error when called before Init.
func TestDevice_RenderFrame_ReturnsErrorBeforeInit(t *testing.T) {
	d := gpu.New()
	if err := d.RenderFrame(); err == nil {
		t.Fatal("RenderFrame() before Init should return an error")
	}
}

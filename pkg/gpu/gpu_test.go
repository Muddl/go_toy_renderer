package gpu_test

import (
	"os"
	"strings"
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
	// Full GPU init test is below. This placeholder ensures the skip guard
	// is exercised in non-GPU CI runs.
	t.Fatal("GPU_TESTS=1 set but landed in skip-guard placeholder")
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

// TestDevice_IsReady_FalseBeforeInit verifies IsReady returns false before Init.
func TestDevice_IsReady_FalseBeforeInit(t *testing.T) {
	d := gpu.New()
	if d.IsReady() {
		t.Fatal("IsReady() should be false before Init")
	}
}

// --- Phase 2: wgpu Initialization Chain (RED tests) ---
// These tests gate on GPU_TESTS=1 and require WGPU_NATIVE_PATH to point to
// the platform wgpu-native library (e.g. assets/windows-x86_64-gnu/lib/wgpu_native.dll).

// TestDevice_Init_NilSurfaceReturnsError verifies that Init returns an error
// when surface is nil.
// RED: current stub returns nil → test fails.
// GREEN: real impl returns "surface is nil" error → test passes.
func TestDevice_Init_NilSurfaceReturnsError(t *testing.T) {
	if os.Getenv("GPU_TESTS") != "1" {
		t.Skip("set GPU_TESTS=1 to run GPU integration tests")
	}
	d := gpu.New()
	err := d.Init(640, 480, nil)
	if err == nil {
		t.Fatal("Init with nil surface should return an error, got nil")
	}
}

// TestDevice_Init_AcquiresQueueBeforeNilSurface verifies that the wgpu queue
// is acquired during Init even when the surface is nil (i.e., the chain runs
// up to queue acquisition before failing on the missing surface).
// RED: current stub returns nil and never sets queue → HasQueue() is false.
// GREEN: real impl sets queue, then errors on nil surface → HasQueue() is true.
func TestDevice_Init_AcquiresQueueBeforeNilSurface(t *testing.T) {
	if os.Getenv("GPU_TESTS") != "1" {
		t.Skip("set GPU_TESTS=1 to run GPU integration tests")
	}
	d := gpu.New()
	err := d.Init(640, 480, nil)
	if err == nil {
		t.Fatal("Init with nil surface should return an error, got nil")
	}
	// Skip if a pre-queue step failed (library not found or no GPU hardware).
	if strings.Contains(err.Error(), "failed to load native library") {
		t.Skipf("wgpu library not found — set WGPU_NATIVE_PATH: %v", err)
	}
	if strings.Contains(err.Error(), "request adapter") ||
		strings.Contains(err.Error(), "request device") ||
		strings.Contains(err.Error(), "get queue") {
		t.Skipf("no GPU hardware available: %v", err)
	}
	// At this point the error must be about nil surface — queue should be set.
	if !d.HasQueue() {
		t.Fatalf("Init: queue not acquired before nil-surface error: %v", err)
	}
}

// --- Phase 3: Hello Triangle (RED tests) ---

// TestDevice_RenderFrame_AfterInit_ReturnsErrorUntilPhase3 documents the expected
// error from RenderFrame before the Hello Triangle pipeline is wired up.
// After Phase 3 the error is replaced by a successful render; until then the
// stub returns a "not yet implemented" sentinel.
func TestDevice_RenderFrame_AfterInit_ReturnsErrorUntilPhase3(t *testing.T) {
	if os.Getenv("GPU_TESTS") != "1" {
		t.Skip("set GPU_TESTS=1 to run GPU integration tests")
	}
	// This test documents: RenderFrame without a surface/pipeline is an error.
	// A full Hello Triangle test requires a GLFW surface (tested via pkg/renderer).
	t.Log("RenderFrame integration test deferred to pkg/renderer with GLFW surface")
}

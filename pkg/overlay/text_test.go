package overlay_test

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/overlay"
)

// TestNewTextLayer_Dimensions_ArePositive verifies the layer has non-zero dimensions.
func TestNewTextLayer_Dimensions_ArePositive(t *testing.T) {
	tl := overlay.NewTextLayer(1280, 720)
	if tl.Width() <= 0 {
		t.Errorf("Width() = %d, want > 0", tl.Width())
	}
	if tl.Height() <= 0 {
		t.Errorf("Height() = %d, want > 0", tl.Height())
	}
}

// TestNewTextLayer_Pixels_SizeMatchesDimensions verifies the pixel buffer has
// exactly Width × Height × 4 bytes.
func TestNewTextLayer_Pixels_SizeMatchesDimensions(t *testing.T) {
	tl := overlay.NewTextLayer(1280, 720)
	want := tl.Width() * tl.Height() * 4
	if got := len(tl.Pixels()); got != want {
		t.Errorf("len(Pixels()) = %d, want %d (W*H*4)", got, want)
	}
}

// TestTextLayer_Render_LevelOff_ReturnsAllZero verifies that LevelOff clears
// the pixel buffer to transparent black.
func TestTextLayer_Render_LevelOff_ReturnsAllZero(t *testing.T) {
	tl := overlay.NewTextLayer(1280, 720)
	pixels := tl.Render(overlay.LevelOff, overlay.Metrics{})
	for i, b := range pixels {
		if b != 0 {
			t.Errorf("pixel[%d] = %d, want 0 for LevelOff", i, b)
			break
		}
	}
}

// TestTextLayer_Render_LevelFPS_HasNonZeroPixels verifies that Level 1 draws
// at least some opaque pixels.
func TestTextLayer_Render_LevelFPS_HasNonZeroPixels(t *testing.T) {
	tl := overlay.NewTextLayer(1280, 720)
	m := overlay.Metrics{FPS: 60.0}
	pixels := tl.Render(overlay.LevelFPS, m)
	hasOpaque := false
	for i := 3; i < len(pixels); i += 4 { // alpha channel
		if pixels[i] != 0 {
			hasOpaque = true
			break
		}
	}
	if !hasOpaque {
		t.Error("Render(LevelFPS) produced no opaque pixels — expected text to be drawn")
	}
}

// TestTextLayer_Render_LevelTimings_HasMorePixelsThanFPS verifies that Level 2
// draws more content than Level 1.
func TestTextLayer_Render_LevelTimings_HasMorePixelsThanFPS(t *testing.T) {
	tl := overlay.NewTextLayer(1280, 720)
	m := overlay.Metrics{FPS: 60.0, FrameTimeMS: 16.67, CPUFrameTimeMS: 14.0}

	pFPS := countOpaquePixels(tl.Render(overlay.LevelFPS, m))
	pTimings := countOpaquePixels(tl.Render(overlay.LevelTimings, m))

	if pTimings <= pFPS {
		t.Errorf("LevelTimings opaque pixels (%d) not greater than LevelFPS (%d)", pTimings, pFPS)
	}
}

// TestTextLayer_Render_LevelDetail_HasMostPixels verifies that Level 3 draws
// more content than Level 2.
func TestTextLayer_Render_LevelDetail_HasMostPixels(t *testing.T) {
	tl := overlay.NewTextLayer(1280, 720)
	m := overlay.Metrics{
		FPS: 60.0, FrameTimeMS: 16.67, CPUFrameTimeMS: 14.0,
		GeometryUploadMS: 1.0, RenderPassMS: 10.0, PresentMS: 0.5,
		Backend: "GPU", VertexCount: 8, TriangleCount: 12,
	}

	pTimings := countOpaquePixels(tl.Render(overlay.LevelTimings, m))
	pDetail := countOpaquePixels(tl.Render(overlay.LevelDetail, m))

	if pDetail <= pTimings {
		t.Errorf("LevelDetail opaque pixels (%d) not greater than LevelTimings (%d)", pDetail, pTimings)
	}
}

// TestTextLayer_Render_SmallFramebuffer_DoesNotPanic verifies the layer clips
// safely when the framebuffer is very small.
func TestTextLayer_Render_SmallFramebuffer_DoesNotPanic(t *testing.T) {
	tl := overlay.NewTextLayer(10, 10)
	// Should not panic regardless of level or metrics.
	tl.Render(overlay.LevelDetail, overlay.Metrics{FPS: 99, Backend: "cpu"})
}

// BenchmarkTextLayer_Render_LevelOff measures the overhead of Render when the
// overlay is hidden (most common steady-state). Must be ≤ 0.1 ms per frame.
func BenchmarkTextLayer_Render_LevelOff(b *testing.B) {
	tl := overlay.NewTextLayer(640, 480)
	m := overlay.Metrics{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tl.Render(overlay.LevelOff, m)
	}
}

// BenchmarkTextLayer_Render_LevelDetail measures worst-case overhead (all
// fields populated). Must be ≤ 0.1 ms per frame.
func BenchmarkTextLayer_Render_LevelDetail(b *testing.B) {
	tl := overlay.NewTextLayer(640, 480)
	m := overlay.Metrics{
		FPS: 60, FrameTimeMS: 16.67, CPUFrameTimeMS: 14.0,
		GeometryUploadMS: 1.0, RenderPassMS: 10.0, PresentMS: 0.5,
		Backend: "GPU", VertexCount: 8, TriangleCount: 12,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tl.Render(overlay.LevelDetail, m)
	}
}

func countOpaquePixels(pixels []byte) int {
	count := 0
	for i := 3; i < len(pixels); i += 4 {
		if pixels[i] != 0 {
			count++
		}
	}
	return count
}

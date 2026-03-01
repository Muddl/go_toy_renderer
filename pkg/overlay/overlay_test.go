package overlay_test

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/overlay"
)

// TestLevel_CycleLevel_WrapsAround verifies CycleLevel advances through all
// levels and wraps from LevelDetail back to LevelOff.
func TestLevel_CycleLevel_WrapsAround(t *testing.T) {
	o := &overlay.Overlay{}
	want := []overlay.Level{
		overlay.LevelOff,
		overlay.LevelFPS,
		overlay.LevelTimings,
		overlay.LevelDetail,
		overlay.LevelOff, // wraparound
	}
	for i, w := range want {
		if got := o.Level(); got != w {
			t.Fatalf("step %d: Level() = %d, want %d", i, got, w)
		}
		if i < len(want)-1 {
			o.CycleLevel()
		}
	}
}

// TestLevel_Constants_HaveExpectedValues ensures the level constants are ordered
// 0–3 so the modulo wraparound logic works correctly.
func TestLevel_Constants_HaveExpectedValues(t *testing.T) {
	cases := []struct {
		name string
		got  overlay.Level
		want overlay.Level
	}{
		{"LevelOff", overlay.LevelOff, 0},
		{"LevelFPS", overlay.LevelFPS, 1},
		{"LevelTimings", overlay.LevelTimings, 2},
		{"LevelDetail", overlay.LevelDetail, 3},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %d, want %d", c.name, c.got, c.want)
		}
	}
}

// TestMetrics_ZeroValue_IsValid verifies the zero-value Metrics is usable
// without initialisation.
func TestMetrics_ZeroValue_IsValid(t *testing.T) {
	var m overlay.Metrics
	if m.FPS != 0 {
		t.Errorf("FPS = %f, want 0", m.FPS)
	}
	if m.FrameTimeMS != 0 {
		t.Errorf("FrameTimeMS = %f, want 0", m.FrameTimeMS)
	}
	if m.Backend != "" {
		t.Errorf("Backend = %q, want empty", m.Backend)
	}
	if m.VertexCount != 0 {
		t.Errorf("VertexCount = %d, want 0", m.VertexCount)
	}
	if m.TriangleCount != 0 {
		t.Errorf("TriangleCount = %d, want 0", m.TriangleCount)
	}
}

// TestMetrics_AllFields_Assignable verifies every field in Metrics can be set
// (documents the complete public API without hard-coding internal layout).
func TestMetrics_AllFields_Assignable(t *testing.T) {
	m := overlay.Metrics{
		FPS:              60.0,
		FrameTimeMS:      16.67,
		CPUFrameTimeMS:   14.0,
		GeometryUploadMS: 1.2,
		RenderPassMS:     10.5,
		PresentMS:        0.8,
		Backend:          "gpu",
		VertexCount:      8,
		TriangleCount:    12,
	}
	if m.FPS != 60.0 {
		t.Errorf("FPS not assigned correctly")
	}
	if m.Backend != "gpu" {
		t.Errorf("Backend not assigned correctly")
	}
}

package overlay_test

import (
	"testing"
	"time"

	"github.com/muddl/go_toy_renderer/pkg/overlay"
)

// TestSampler_ZeroValue_FPSReturnsZero verifies a freshly created Sampler
// returns 0 FPS before any samples are added.
func TestSampler_ZeroValue_FPSReturnsZero(t *testing.T) {
	var s overlay.Sampler
	if fps := s.FPS(); fps != 0 {
		t.Errorf("FPS() on empty Sampler = %f, want 0", fps)
	}
}

// TestSampler_AddSample_SingleFrame_FPSMatchesDuration verifies that one
// sample of 16 ms gives approximately 62.5 FPS.
func TestSampler_AddSample_SingleFrame_FPSMatchesDuration(t *testing.T) {
	var s overlay.Sampler
	s.AddSample(16 * time.Millisecond)
	fps := s.FPS()
	if fps < 62 || fps > 63 {
		t.Errorf("FPS() after 16 ms sample = %.2f, want ~62.5", fps)
	}
}

// TestSampler_AddSample_SixtyFrames_FPSStabilises verifies that after 60
// identical 16.67 ms frames the rolling average converges to ~60 FPS.
func TestSampler_AddSample_SixtyFrames_FPSStabilises(t *testing.T) {
	var s overlay.Sampler
	frame := time.Second / 60
	for i := 0; i < 60; i++ {
		s.AddSample(frame)
	}
	fps := s.FPS()
	if fps < 59 || fps > 61 {
		t.Errorf("FPS() after 60 × 16.67 ms = %.2f, want ~60", fps)
	}
}

// TestSampler_AddSample_OldestSampleEvicted verifies that after more than 60
// samples, early samples no longer influence the average.
func TestSampler_AddSample_OldestSampleEvicted(t *testing.T) {
	var s overlay.Sampler
	// Fill with slow frames (10 FPS each = 100 ms per frame).
	for i := 0; i < 60; i++ {
		s.AddSample(100 * time.Millisecond)
	}
	// Replace all with fast frames (120 FPS each ≈ 8.33 ms).
	for i := 0; i < 60; i++ {
		s.AddSample(time.Second / 120)
	}
	fps := s.FPS()
	if fps < 115 || fps > 125 {
		t.Errorf("FPS() after eviction = %.2f, want ~120", fps)
	}
}

// Package overlay provides a live performance debug overlay for renderer-rt.
// The overlay is toggled via F3 and cycles through three granularity levels.
// All types in this file are pure Go with no build-tag restrictions so they
// compile under the headless CI build.
package overlay

// Level controls how much information the overlay displays.
type Level uint8

const (
	LevelOff     Level = 0 // overlay hidden
	LevelFPS     Level = 1 // FPS only
	LevelTimings Level = 2 // FPS + frame time + total CPU frame time
	LevelDetail  Level = 3 // + per-phase timings, backend, mesh stats
)

// Metrics holds performance counters captured once per frame by each backend.
// Fields are populated progressively: Level 1 uses FPS; Level 2 adds timing
// fields; Level 3 adds the per-phase breakdown and scene statistics.
type Metrics struct {
	// Level 1+
	FPS float64 // rolling-average frames per second

	// Level 2+
	FrameTimeMS    float64 // wall-clock frame duration in milliseconds
	CPUFrameTimeMS float64 // CPU work time (frame duration minus sleep) in ms

	// Level 3+
	GeometryUploadMS float64 // time spent uploading geometry to GPU (ms)
	RenderPassMS     float64 // time spent recording/submitting the render pass (ms)
	PresentMS        float64 // time spent in surface Present/SwapBuffers (ms)
	Backend          string  // "cpu" or "gpu"
	VertexCount      int     // total vertices in the active scene
	TriangleCount    int     // total triangles in the active scene
}

// Overlay manages the current debug level and the last recorded Metrics.
type Overlay struct {
	level   Level
	metrics Metrics
}

// Level returns the current display level.
func (o *Overlay) Level() Level { return o.level }

// CycleLevel advances the level by one, wrapping from LevelDetail back to
// LevelOff. Intended to be called from the F3 key callback.
func (o *Overlay) CycleLevel() {
	o.level = (o.level + 1) % 4
}

// Update stores the latest Metrics for use during the next Render call.
func (o *Overlay) Update(m Metrics) {
	o.metrics = m
}

// Metrics returns the most recently stored Metrics snapshot.
func (o *Overlay) Metrics() Metrics {
	return o.metrics
}

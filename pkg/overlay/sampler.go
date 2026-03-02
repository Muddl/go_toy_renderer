package overlay

import "time"

const samplerSize = 60 // rolling window: last 60 frame durations

// Sampler maintains a circular buffer of recent frame durations and computes
// a rolling-average FPS. The zero value is ready to use.
type Sampler struct {
	buf   [samplerSize]time.Duration
	head  int // index of the next write slot
	count int // number of valid entries (0..samplerSize)
}

// AddSample records one frame duration into the circular buffer.
func (s *Sampler) AddSample(d time.Duration) {
	s.buf[s.head] = d
	s.head = (s.head + 1) % samplerSize
	if s.count < samplerSize {
		s.count++
	}
}

// FPS returns the rolling-average frames-per-second across the buffered
// samples. Returns 0 if no samples have been added yet.
func (s *Sampler) FPS() float64 {
	if s.count == 0 {
		return 0
	}
	var total time.Duration
	for i := 0; i < s.count; i++ {
		total += s.buf[i]
	}
	avg := total / time.Duration(s.count)
	if avg <= 0 {
		return 0
	}
	return float64(time.Second) / float64(avg)
}

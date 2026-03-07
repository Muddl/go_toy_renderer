package gpu

import (
	"fmt"
	"time"
)

const printInterval = time.Second

// FrameTimer accumulates GPU frame-time samples and returns a formatted
// summary string once per second via MaybePrint.
type FrameTimer struct {
	samples   []uint64
	lastPrint time.Time
}

// NewFrameTimer returns a FrameTimer ready to use.
func NewFrameTimer() *FrameTimer {
	return &FrameTimer{lastPrint: time.Now()}
}

// Record adds one frame-time sample in nanoseconds.
func (ft *FrameTimer) Record(ns uint64) {
	ft.samples = append(ft.samples, ns)
}

// MaybePrint returns a formatted summary string if at least one second has
// elapsed since the last print, then resets the accumulator. Returns "" if
// it is too early to print.
func (ft *FrameTimer) MaybePrint() string {
	if time.Since(ft.lastPrint) < printInterval {
		return ""
	}
	if len(ft.samples) == 0 {
		ft.lastPrint = time.Now()
		return ""
	}
	var sum uint64
	for _, s := range ft.samples {
		sum += s
	}
	avgMs := float64(sum) / float64(len(ft.samples)) / 1e6
	msg := fmt.Sprintf("GPU frame time: %.2f ms avg (%d samples)", avgMs, len(ft.samples))
	ft.samples = ft.samples[:0]
	ft.lastPrint = time.Now()
	return msg
}

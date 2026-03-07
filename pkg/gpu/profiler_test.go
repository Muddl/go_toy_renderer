package gpu

import (
	"strings"
	"testing"
	"time"
)

func TestFrameTimer_NoOutputBeforeInterval(t *testing.T) {
	ft := NewFrameTimer()
	ft.Record(1_000_000) // 1 ms
	if msg := ft.MaybePrint(); msg != "" {
		t.Errorf("expected no output before 1 second, got %q", msg)
	}
}

func TestFrameTimer_OutputAfterInterval(t *testing.T) {
	ft := NewFrameTimer()
	// Back-date lastPrint so the interval has elapsed.
	ft.lastPrint = time.Now().Add(-2 * time.Second)
	ft.Record(500_000) // 0.5 ms
	msg := ft.MaybePrint()
	if msg == "" {
		t.Fatal("expected non-empty output after interval elapsed")
	}
	if !strings.Contains(msg, "GPU frame") {
		t.Errorf("expected message to contain 'GPU frame', got %q", msg)
	}
}

func TestFrameTimer_AveragesMultipleSamples(t *testing.T) {
	ft := NewFrameTimer()
	ft.lastPrint = time.Now().Add(-2 * time.Second)
	ft.Record(1_000_000) // 1 ms
	ft.Record(3_000_000) // 3 ms — average should be 2 ms
	msg := ft.MaybePrint()
	if !strings.Contains(msg, "2.00") {
		t.Errorf("expected average 2.00 ms in output, got %q", msg)
	}
}

func TestFrameTimer_ResetsAfterPrint(t *testing.T) {
	ft := NewFrameTimer()
	ft.lastPrint = time.Now().Add(-2 * time.Second)
	ft.Record(1_000_000)
	ft.MaybePrint() // consume the output
	// Immediately calling MaybePrint again should return "" (interval just reset).
	if msg := ft.MaybePrint(); msg != "" {
		t.Errorf("expected empty output right after print, got %q", msg)
	}
}

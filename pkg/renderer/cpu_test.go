//go:build headless

package renderer

import (
	"testing"
)

// In headless builds there is no display; CPUBackend.Init must return an error.
func TestCPUBackend_Init_Headless_ReturnsError(t *testing.T) {
	b := &CPUBackend{}
	err := b.Init(640, 480)
	if err == nil {
		t.Fatal("CPUBackend.Init should return an error in headless mode")
	}
}

func TestCPUBackend_RenderFrame_Headless_ReturnsError(t *testing.T) {
	b := &CPUBackend{}
	err := b.RenderFrame(nil)
	if err == nil {
		t.Fatal("CPUBackend.RenderFrame should return an error in headless mode")
	}
}

func TestCPUBackend_Shutdown_Headless_NoPanic(t *testing.T) {
	b := &CPUBackend{}
	b.Shutdown() // must not panic
}

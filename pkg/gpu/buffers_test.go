//go:build !headless

package gpu_test

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/gpu"
)

// TestVertexBuffer_ReturnsErrorWithoutInit verifies that VertexBuffer returns an
// error when called with an uninitialized Device (no wgpu device handle).
// This test runs without GPU hardware — it only verifies the guard condition.
func TestVertexBuffer_ReturnsErrorWithoutInit(t *testing.T) {
	d := gpu.New()
	_, err := gpu.VertexBuffer(d, nil)
	if err == nil {
		t.Fatal("VertexBuffer on uninitialized Device should return an error")
	}
}

// TestIndexBuffer_ReturnsErrorWithoutInit verifies that IndexBuffer returns an
// error when called with an uninitialized Device (no wgpu device handle).
// This test runs without GPU hardware — it only verifies the guard condition.
func TestIndexBuffer_ReturnsErrorWithoutInit(t *testing.T) {
	d := gpu.New()
	_, err := gpu.IndexBuffer(d, nil)
	if err == nil {
		t.Fatal("IndexBuffer on uninitialized Device should return an error")
	}
}

package renderer

import (
	"strings"
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/render"
)

// --- GPUBackend tests ---

func TestGPUBackend_Init_ReturnsNil(t *testing.T) {
	gpu := &GPUBackend{}
	if err := gpu.Init(640, 480); err != nil {
		t.Errorf("GPUBackend.Init() = %v, want nil", err)
	}
}

func TestGPUBackend_RenderFrame_ReturnsNotImplemented(t *testing.T) {
	gpu := &GPUBackend{}
	err := gpu.RenderFrame(render.NewScene())
	if err == nil {
		t.Fatal("GPUBackend.RenderFrame() should return an error")
	}
	if !strings.Contains(err.Error(), "GPU not yet implemented") {
		t.Errorf("GPUBackend.RenderFrame() error = %q, want to contain 'GPU not yet implemented'", err.Error())
	}
}

func TestGPUBackend_Shutdown_NoPanic(t *testing.T) {
	gpu := &GPUBackend{}
	gpu.Shutdown() // must not panic
}

// --- Factory tests ---

func TestNew_CPU_ReturnsCPUBackend(t *testing.T) {
	r, err := New("cpu")
	if err != nil {
		t.Fatalf("New(\"cpu\") error = %v", err)
	}
	if _, ok := r.(*CPUBackend); !ok {
		t.Errorf("New(\"cpu\") returned %T, want *CPUBackend", r)
	}
}

func TestNew_GPU_ReturnsGPUBackend(t *testing.T) {
	r, err := New("gpu")
	if err != nil {
		t.Fatalf("New(\"gpu\") error = %v", err)
	}
	if _, ok := r.(*GPUBackend); !ok {
		t.Errorf("New(\"gpu\") returned %T, want *GPUBackend", r)
	}
}

func TestNew_Auto_ReturnsCPUBackend(t *testing.T) {
	r, err := New("auto")
	if err != nil {
		t.Fatalf("New(\"auto\") error = %v", err)
	}
	if _, ok := r.(*CPUBackend); !ok {
		t.Errorf("New(\"auto\") returned %T, want *CPUBackend (GPU not yet available)", r)
	}
}

func TestNew_Unknown_ReturnsError(t *testing.T) {
	_, err := New("vulkan")
	if err == nil {
		t.Error("New(\"vulkan\") should return an error for unknown backend")
	}
}

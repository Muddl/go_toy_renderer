package renderer

import "fmt"

// New creates and returns a Renderer for the requested backend name.
// Valid values are "cpu", "gpu", and "auto".
//
// "auto" returns a CPUBackend for Phase 10; GPU-first with CPU fallback
// will be implemented in Phase 11 once wgpu-native bindings are available.
func New(backend string) (Renderer, error) {
	switch backend {
	case "cpu":
		return &CPUBackend{}, nil
	case "gpu":
		return &GPUBackend{}, nil
	case "auto":
		// Phase 11: try GPU first; fall back to CPU on any init error.
		// newAutoBackend is split by build tag: non-headless returns an autoBackend
		// that defers GPU→CPU fallback to Init; headless returns *CPUBackend directly.
		return newAutoBackend(), nil
	default:
		return nil, fmt.Errorf("unknown backend %q: must be cpu, gpu, or auto", backend)
	}
}

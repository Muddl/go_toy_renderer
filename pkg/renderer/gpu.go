package renderer

import (
	"errors"

	"github.com/muddl/go_toy_renderer/pkg/render"
)

// GPUBackend is a stub that satisfies Renderer until wgpu-native bindings
// are introduced in Phase 11. Init succeeds (seam for Phase 11), but
// RenderFrame always returns an error because the GPU pipeline is not yet
// implemented.
type GPUBackend struct{}

// Init is a no-op placeholder for the future wgpu-native initialisation.
func (g *GPUBackend) Init(_, _ int) error { return nil }

// RenderFrame always returns an error because GPU rendering is not yet
// implemented. Phase 11 will replace this body with wgpu draw calls.
func (g *GPUBackend) RenderFrame(_ *render.Scene) error {
	return errors.New("GPU not yet implemented")
}

// Shutdown is a no-op until Phase 11 allocates GPU resources.
func (g *GPUBackend) Shutdown() {}

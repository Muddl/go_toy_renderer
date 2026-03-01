//go:build headless

package renderer

import (
	"errors"

	"github.com/muddl/go_toy_renderer/pkg/render"
)

// GPUBackend is a stub for headless (CI) builds where no window or GPU is available.
type GPUBackend struct{}

// Init is a no-op in headless mode — no GPU context can be created.
func (g *GPUBackend) Init(_, _ int) error { return nil }

// RenderFrame always returns an error in headless mode because no GPU pipeline exists.
func (g *GPUBackend) RenderFrame(_ *render.Scene) error {
	return errors.New("GPU not yet implemented")
}

// Shutdown is a no-op in headless mode.
func (g *GPUBackend) Shutdown() {}

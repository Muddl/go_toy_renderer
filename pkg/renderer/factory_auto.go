//go:build !headless

package renderer

import (
	"fmt"

	"github.com/muddl/go_toy_renderer/pkg/render"
	"github.com/muddl/go_toy_renderer/pkg/scene"
)

// autoBackend tries GPUBackend first; if Init fails it falls back to CPUBackend.
// This is used by the "auto" selector in New().
type autoBackend struct {
	active Renderer
}

// newAutoBackend creates an autoBackend for non-headless builds.
func newAutoBackend() Renderer {
	return &autoBackend{}
}

// Init attempts GPU init first. On any error it falls back to CPUBackend.
func (a *autoBackend) Init(width, height int) error {
	g := &GPUBackend{}
	if err := g.Init(width, height); err == nil {
		a.active = g
		return nil
	}
	c := &CPUBackend{}
	if err := c.Init(width, height); err != nil {
		return fmt.Errorf("auto: both GPU and CPU init failed: %w", err)
	}
	a.active = c
	return nil
}

// RenderFrame delegates to whichever backend was selected during Init.
func (a *autoBackend) RenderFrame(scene *render.Scene) error {
	return a.active.RenderFrame(scene)
}

// SetScene forwards to the active backend if it supports scene-aware rendering.
func (a *autoBackend) SetScene(s *scene.Scene) {
	if setter, ok := a.active.(interface{ SetScene(*scene.Scene) }); ok {
		setter.SetScene(s)
	}
}

// Shutdown releases the active backend's resources.
func (a *autoBackend) Shutdown() {
	if a.active != nil {
		a.active.Shutdown()
	}
}

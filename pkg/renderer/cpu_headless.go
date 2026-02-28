//go:build headless

package renderer

import (
	"errors"

	"github.com/muddl/go_toy_renderer/pkg/render"
)

// CPUBackend is a stub for headless (CI) builds where no window is available.
type CPUBackend struct{}

// Init always returns an error in headless mode because no window can be created.
func (c *CPUBackend) Init(_, _ int) error {
	return errors.New("windowed rendering not available in headless build")
}

// RenderFrame always returns an error in headless mode.
func (c *CPUBackend) RenderFrame(_ *render.Scene) error {
	return errors.New("windowed rendering not available in headless build")
}

// Shutdown is a no-op in headless mode.
func (c *CPUBackend) Shutdown() {}

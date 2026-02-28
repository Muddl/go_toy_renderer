// Package renderer defines the Renderer interface and provides concrete
// backend implementations (CPU and GPU) for the real-time rendering loop.
package renderer

import (
	"errors"

	"github.com/muddl/go_toy_renderer/pkg/render"
)

// Renderer is the interface satisfied by all rendering backends.
// The main loop calls Init once, RenderFrame each tick, and Shutdown on exit.
type Renderer interface {
	// Init sets up the backend for a window of the given size.
	Init(width, height int) error
	// RenderFrame renders one frame of the scene.
	// Returns ErrWindowClosed when the display window has been dismissed.
	RenderFrame(scene *render.Scene) error
	// Shutdown releases all backend resources.
	Shutdown()
}

// ErrWindowClosed is returned by RenderFrame when the window has been closed
// or the user has pressed ESC. The main loop should exit cleanly on this error.
var ErrWindowClosed = errors.New("window closed")

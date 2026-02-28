package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/render"
	"github.com/muddl/go_toy_renderer/pkg/renderer"
)

func init() {
	// GLFW requires all calls on the OS thread.
	runtime.LockOSThread()
}

func main() {
	cfg, err := parseConfig(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if runErr := runRenderer(cfg); runErr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", runErr)
		os.Exit(1)
	}
}

// runRenderer initializes the chosen backend, builds the scene, and runs the
// frame loop until the window is closed or an unrecoverable error occurs.
// Shutdown is guaranteed to run via defer.
func runRenderer(cfg Config) error {
	r, err := renderer.New(cfg.Backend)
	if err != nil {
		return fmt.Errorf("renderer: %w", err)
	}

	if err = r.Init(cfg.Width, cfg.Height); err != nil {
		return fmt.Errorf("renderer init: %w", err)
	}
	defer r.Shutdown()

	scene := render.NewScene()
	scene.AddMesh(geometry.NewCube())

	for {
		if frameErr := r.RenderFrame(scene); frameErr != nil {
			if errors.Is(frameErr, renderer.ErrWindowClosed) {
				return nil
			}
			return fmt.Errorf("render frame: %w", frameErr)
		}
	}
}

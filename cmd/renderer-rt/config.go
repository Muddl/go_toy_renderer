package main

import (
	"errors"
	"flag"
	"fmt"
)

// Config holds parsed command-line configuration.
type Config struct {
	Width   int
	Height  int
	Backend string
}

// parseConfig parses command-line arguments into a Config.
// Supported flags: --width, --height, --backend.
func parseConfig(args []string) (Config, error) {
	fs := flag.NewFlagSet("renderer-rt", flag.ContinueOnError)

	width := fs.Int("width", 1280, "window width in pixels")
	height := fs.Int("height", 720, "window height in pixels")
	backend := fs.String("backend", "cpu", "rendering backend: cpu, gpu, or auto")

	if err := fs.Parse(args); err != nil {
		return Config{}, fmt.Errorf("flag parse: %w", err)
	}

	return Config{
		Width:   *width,
		Height:  *height,
		Backend: *backend,
	}, nil
}

// validateBackend returns an error if the backend is unsupported in the current phase.
// "cpu" and "auto" proceed; "gpu" is not yet implemented until Phase 10.
func validateBackend(backend string) error {
	switch backend {
	case "cpu", "auto":
		return nil
	case "gpu":
		return errors.New("--backend gpu is not yet implemented (available in Phase 10+)")
	default:
		return fmt.Errorf("unknown backend %q: must be one of cpu, auto, gpu", backend)
	}
}

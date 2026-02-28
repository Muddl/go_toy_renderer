package main

import (
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

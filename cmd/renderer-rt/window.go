//go:build !headless

package main

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// run initialises a GLFW window and enters the event loop.
// It returns when the window is closed or ESC is pressed.
func run(cfg Config) error {
	if err := glfw.Init(); err != nil {
		return fmt.Errorf("glfw init: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(cfg.Width, cfg.Height, "go_toy_renderer", nil, nil)
	if err != nil {
		return fmt.Errorf("create window: %w", err)
	}
	window.MakeContextCurrent()

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			w.SetShouldClose(true)
		}
	})

	for !window.ShouldClose() {
		glfw.PollEvents()
	}

	return nil
}

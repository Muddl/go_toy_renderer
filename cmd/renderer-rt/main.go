package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/render"
	"github.com/muddl/go_toy_renderer/pkg/renderer"
	"github.com/muddl/go_toy_renderer/pkg/scene"
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

	if initErr := r.Init(cfg.Width, cfg.Height); initErr != nil {
		return fmt.Errorf("renderer init: %w", initErr)
	}
	defer r.Shutdown()

	// Build GPU scene with cube + tetrahedron at different positions.
	gpuScene := buildDemoScene()

	// If the backend supports scene-aware rendering, use it.
	if gpuBackend, ok := r.(interface{ SetScene(*scene.Scene) }); ok {
		gpuBackend.SetScene(gpuScene)
	}

	// CPU fallback: render.Scene with just a cube.
	cpuScene := render.NewScene()
	cpuScene.AddMesh(geometry.NewCube())

	for {
		if frameErr := r.RenderFrame(cpuScene); frameErr != nil {
			if errors.Is(frameErr, renderer.ErrWindowClosed) {
				return nil
			}
			return fmt.Errorf("render frame: %w", frameErr)
		}
	}
}

// buildDemoScene creates a scene with a cube at (-1.5, 0, 0) and a
// tetrahedron at (1.5, 0, 0) to demonstrate per-mesh transforms.
func buildDemoScene() *scene.Scene {
	s := scene.NewScene()

	// Cube on the left.
	cubeTransform := scene.NewTransform()
	cubeTransform.Position = math.Vec3{X: -1.5, Y: 0, Z: 0}
	s.AddNode(scene.Node{
		Mesh:      geometry.NewCube(),
		Transform: cubeTransform,
	})

	// Tetrahedron on the right.
	tetraTransform := scene.NewTransform()
	tetraTransform.Position = math.Vec3{X: 1.5, Y: 0, Z: 0}
	s.AddNode(scene.Node{
		Mesh:      geometry.NewTetrahedron(),
		Transform: tetraTransform,
	})

	return s
}

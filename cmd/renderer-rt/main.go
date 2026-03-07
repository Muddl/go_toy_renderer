package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/loader"
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

	// Build GPU scene with meshes + camera cylinder.
	gpuScene, camCylIdx := buildDemoScene()

	// If the backend supports scene-aware rendering, use it.
	if gpuBackend, ok := r.(interface{ SetScene(*scene.Scene) }); ok {
		gpuBackend.SetScene(gpuScene)
	}
	if gpuBackend, ok := r.(interface{ SetCameraNodeIndex(int) }); ok {
		gpuBackend.SetCameraNodeIndex(camCylIdx)
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

// buildDemoScene creates a scene with a Utah teapot, cube, ground plane,
// and a camera cylinder. The teapot is loaded from assets/models/teapot.obj;
// if the file is not found the slot is filled with a tetrahedron fallback.
// Returns the scene and the index of the camera cylinder node.
func buildDemoScene() (*scene.Scene, int) {
	s := scene.NewScene()

	// Utah Teapot at centre — the star of the show.
	teapotMesh := loadTeapotOrFallback()
	teapotTransform := scene.NewTransform()
	teapotTransform.Scale = math.Vec3{X: 0.5, Y: 0.5, Z: 0.5}
	s.AddNode(scene.Node{
		Mesh:      teapotMesh,
		Transform: teapotTransform,
	})

	// Cube on the left.
	cubeTransform := scene.NewTransform()
	cubeTransform.Position = math.Vec3{X: -3.5, Y: 0, Z: 0}
	s.AddNode(scene.Node{
		Mesh:      geometry.NewCube(),
		Transform: cubeTransform,
	})

	// Ground plane beneath the meshes (static — no spin animation).
	planeTransform := scene.NewTransform()
	planeTransform.Position = math.Vec3{X: 0, Y: -1, Z: 0}
	s.AddNode(scene.Node{
		Mesh:      geometry.NewPlane(10, 10),
		Transform: planeTransform,
		Static:    true,
	})

	// Small cylinder at the camera position (updated each frame by the renderer).
	camCylTransform := scene.NewTransform()
	camCylTransform.Scale = math.Vec3{X: 0.3, Y: 0.3, Z: 0.3}
	camCylIdx := len(s.Nodes)
	s.AddNode(scene.Node{
		Mesh:      geometry.NewCylinder(12),
		Transform: camCylTransform,
		Static:    true,
	})

	return s, camCylIdx
}

// loadTeapotOrFallback loads the Utah Teapot OBJ from assets/models/teapot.obj.
// If the file is not found or cannot be parsed, a tetrahedron is returned as a
// fallback so the demo still runs in environments without the asset file.
func loadTeapotOrFallback() *geometry.Mesh {
	meshes, err := loader.LoadOBJ("assets/models/teapot.obj")
	if err == nil && len(meshes) > 0 {
		return meshes[0]
	}
	fmt.Fprintf(os.Stderr, "teapot: %v — using fallback tetrahedron\n", err)
	return geometry.NewTetrahedron()
}

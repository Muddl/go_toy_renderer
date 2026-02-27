// Package main provides the entry point for the toy 3D software renderer.
// It renders a colored cube with perspective projection and saves the result as output.png.
package main

import (
	"log"

	"github.com/muddl/go_toy_renderer/pkg/camera"
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
	math "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/render"
	"github.com/muddl/go_toy_renderer/pkg/shader"
)

func main() {
	const width, height = 640, 480

	// Build scene: a single multi-colored cube at world origin.
	scene := render.NewScene()
	scene.AddMesh(geometry.NewCube())

	// Camera: positioned upper-right-front so multiple faces are visible.
	cam := camera.New(
		math.Vec3{X: 3, Y: 2, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0,
		float64(width)/float64(height),
		0.1,
		100.0,
	)

	// Allocate framebuffer and render.
	fb := framebuffer.New(width, height)
	render.Render(scene, cam, fb, shader.VertexColor)

	// Write the result to disk.
	if err := fb.SavePNG("output.png"); err != nil {
		log.Fatalf("save PNG: %v", err)
	}

	log.Println("Rendered output.png")
}

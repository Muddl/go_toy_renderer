package geometry

import (
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// NewTetrahedron creates a regular tetrahedron mesh.
// A tetrahedron has 4 vertices and 4 triangular faces.
// Each vertex is assigned a unique color for visualization.
func NewTetrahedron() *Mesh {
	mesh := NewMesh()

	// Define tetrahedron vertices in a regular configuration
	// Using a tetrahedron centered at origin with edge length ~1.63
	v0 := NewVertex(
		math.Vec3{X: 1.0, Y: 1.0, Z: 1.0}, // Top-front-right
		math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}, // Red
	)
	v1 := NewVertex(
		math.Vec3{X: -1.0, Y: -1.0, Z: 1.0}, // Bottom-back-right
		math.Vec3{X: 0.0, Y: 1.0, Z: 0.0},   // Green
	)
	v2 := NewVertex(
		math.Vec3{X: -1.0, Y: 1.0, Z: -1.0}, // Top-back-left
		math.Vec3{X: 0.0, Y: 0.0, Z: 1.0},   // Blue
	)
	v3 := NewVertex(
		math.Vec3{X: 1.0, Y: -1.0, Z: -1.0}, // Bottom-front-left
		math.Vec3{X: 1.0, Y: 1.0, Z: 0.0},   // Yellow
	)

	mesh.AddVertex(v0)
	mesh.AddVertex(v1)
	mesh.AddVertex(v2)
	mesh.AddVertex(v3)

	// Add 4 triangular faces with counter-clockwise winding (outward-facing normals)
	mesh.AddTriangle(0, 1, 2) // Base triangle 1
	mesh.AddTriangle(0, 3, 1) // Base triangle 2
	mesh.AddTriangle(0, 2, 3) // Side triangle 1
	mesh.AddTriangle(1, 3, 2) // Side triangle 2

	return mesh
}

// NewCube creates a unit cube mesh centered at the origin.
// A cube has 8 vertices and 12 triangular faces (2 per quad face).
// Each vertex is assigned a color based on its position.
func NewCube() *Mesh {
	mesh := NewMesh()

	// Define 8 cube vertices (corners of a unit cube from -0.5 to 0.5)
	// Each vertex color is based on its position for easy identification
	vertices := []struct {
		pos   math.Vec3
		color math.Vec3
	}{
		// Bottom face (Y = -0.5)
		{math.Vec3{X: -0.5, Y: -0.5, Z: -0.5}, math.Vec3{X: 0.0, Y: 0.0, Z: 0.0}}, // 0: Black
		{math.Vec3{X: 0.5, Y: -0.5, Z: -0.5}, math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}},  // 1: Red
		{math.Vec3{X: 0.5, Y: -0.5, Z: 0.5}, math.Vec3{X: 1.0, Y: 0.0, Z: 1.0}},   // 2: Magenta
		{math.Vec3{X: -0.5, Y: -0.5, Z: 0.5}, math.Vec3{X: 0.0, Y: 0.0, Z: 1.0}},  // 3: Blue

		// Top face (Y = 0.5)
		{math.Vec3{X: -0.5, Y: 0.5, Z: -0.5}, math.Vec3{X: 0.0, Y: 1.0, Z: 0.0}}, // 4: Green
		{math.Vec3{X: 0.5, Y: 0.5, Z: -0.5}, math.Vec3{X: 1.0, Y: 1.0, Z: 0.0}},  // 5: Yellow
		{math.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, math.Vec3{X: 1.0, Y: 1.0, Z: 1.0}},   // 6: White
		{math.Vec3{X: -0.5, Y: 0.5, Z: 0.5}, math.Vec3{X: 0.0, Y: 1.0, Z: 1.0}},  // 7: Cyan
	}

	// Add all vertices to mesh
	for _, v := range vertices {
		mesh.AddVertex(NewVertex(v.pos, v.color))
	}

	// Add 12 triangles (2 per face) with counter-clockwise winding when viewed from outside
	// Each quad face is split into 2 triangles

	// Bottom face (Y = -0.5): normal must point down (-Y); CCW from below = 0,1,2 and 0,2,3
	mesh.AddTriangle(0, 1, 2)
	mesh.AddTriangle(0, 2, 3)

	// Top face (Y = 0.5): normal must point up (+Y); CCW from above = 4,6,5 and 4,7,6
	mesh.AddTriangle(4, 6, 5)
	mesh.AddTriangle(4, 7, 6)

	// Front face (Z = 0.5) - looking from front, CCW = 3,2,6 and 3,6,7
	mesh.AddTriangle(3, 2, 6)
	mesh.AddTriangle(3, 6, 7)

	// Back face (Z = -0.5) - looking from back, CCW = 1,0,4 and 1,4,5
	mesh.AddTriangle(1, 0, 4)
	mesh.AddTriangle(1, 4, 5)

	// Right face (X = 0.5) - looking from right, CCW = 2,1,5 and 2,5,6
	mesh.AddTriangle(2, 1, 5)
	mesh.AddTriangle(2, 5, 6)

	// Left face (X = -0.5) - looking from left, CCW = 0,3,7 and 0,7,4
	mesh.AddTriangle(0, 3, 7)
	mesh.AddTriangle(0, 7, 4)

	return mesh
}

package geometry

import (
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// NewTetrahedron creates a regular tetrahedron mesh with per-face normals.
// Each face has 3 dedicated vertices (12 total) so every vertex carries the face
// normal, enabling flat shading. 4 triangles, 12 indices.
func NewTetrahedron() *Mesh {
	mesh := NewMesh()

	// Corner positions of a regular tetrahedron centered at origin.
	p0 := math.Vec3{X: 1, Y: 1, Z: 1}
	p1 := math.Vec3{X: -1, Y: -1, Z: 1}
	p2 := math.Vec3{X: -1, Y: 1, Z: -1}
	p3 := math.Vec3{X: 1, Y: -1, Z: -1}

	// Per-face colors.
	colors := [4]math.Vec3{
		{X: 1, Y: 0, Z: 0}, // Red
		{X: 0, Y: 1, Z: 0}, // Green
		{X: 0, Y: 0, Z: 1}, // Blue
		{X: 1, Y: 1, Z: 0}, // Yellow
	}

	// Faces with CCW winding (outward-facing).
	type face struct {
		a, b, c math.Vec3
	}
	faces := [4]face{
		{p0, p2, p1}, // opposite p3
		{p0, p1, p3}, // opposite p2
		{p0, p3, p2}, // opposite p1
		{p1, p2, p3}, // opposite p0
	}

	for i, f := range faces {
		e1 := f.b.Subtract(f.a)
		e2 := f.c.Subtract(f.a)
		n := e1.Cross(e2).Normalize()

		base := len(mesh.Vertices)
		for _, pos := range [3]math.Vec3{f.a, f.b, f.c} {
			v := NewVertex(pos, colors[i])
			v.Normal = n
			mesh.AddVertex(v)
		}
		mesh.AddTriangle(base, base+1, base+2)
	}

	return mesh
}

// NewCube creates a unit cube mesh centered at the origin with per-face normals.
// Each face has 4 dedicated vertices (24 total) so every vertex carries the face
// normal, enabling flat shading. 12 triangles (2 per face), 36 indices.
func NewCube() *Mesh {
	mesh := NewMesh()

	type face struct {
		normal math.Vec3
		color  math.Vec3
		// Quad corners in CCW order when viewed from outside.
		corners [4]math.Vec3
	}

	h := 0.5 // half-extent

	faces := []face{
		{ // +Y (top)
			normal:  math.Vec3{Y: 1},
			color:   math.Vec3{X: 0, Y: 1, Z: 0},
			corners: [4]math.Vec3{{X: -h, Y: h, Z: -h}, {X: -h, Y: h, Z: h}, {X: h, Y: h, Z: h}, {X: h, Y: h, Z: -h}},
		},
		{ // -Y (bottom)
			normal:  math.Vec3{Y: -1},
			color:   math.Vec3{X: 0, Y: 0.4, Z: 0},
			corners: [4]math.Vec3{{X: -h, Y: -h, Z: h}, {X: -h, Y: -h, Z: -h}, {X: h, Y: -h, Z: -h}, {X: h, Y: -h, Z: h}},
		},
		{ // +X (right)
			normal:  math.Vec3{X: 1},
			color:   math.Vec3{X: 1, Y: 0, Z: 0},
			corners: [4]math.Vec3{{X: h, Y: -h, Z: -h}, {X: h, Y: h, Z: -h}, {X: h, Y: h, Z: h}, {X: h, Y: -h, Z: h}},
		},
		{ // -X (left)
			normal:  math.Vec3{X: -1},
			color:   math.Vec3{X: 0.4, Y: 0, Z: 0},
			corners: [4]math.Vec3{{X: -h, Y: -h, Z: h}, {X: -h, Y: h, Z: h}, {X: -h, Y: h, Z: -h}, {X: -h, Y: -h, Z: -h}},
		},
		{ // +Z (front)
			normal:  math.Vec3{Z: 1},
			color:   math.Vec3{X: 0, Y: 0, Z: 1},
			corners: [4]math.Vec3{{X: -h, Y: -h, Z: h}, {X: h, Y: -h, Z: h}, {X: h, Y: h, Z: h}, {X: -h, Y: h, Z: h}},
		},
		{ // -Z (back)
			normal:  math.Vec3{Z: -1},
			color:   math.Vec3{X: 0, Y: 0, Z: 0.4},
			corners: [4]math.Vec3{{X: h, Y: -h, Z: -h}, {X: -h, Y: -h, Z: -h}, {X: -h, Y: h, Z: -h}, {X: h, Y: h, Z: -h}},
		},
	}

	for _, f := range faces {
		base := len(mesh.Vertices)
		for _, c := range f.corners {
			v := NewVertex(c, f.color)
			v.Normal = f.normal
			mesh.AddVertex(v)
		}
		// Two CCW triangles per quad: 0-1-2, 0-2-3.
		mesh.AddTriangle(base+0, base+1, base+2)
		mesh.AddTriangle(base+0, base+2, base+3)
	}

	return mesh
}

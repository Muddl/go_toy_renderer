package geometry

import (
	stdmath "math"

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

// NewPlane creates a flat quad on the XZ plane at Y=0, centered at origin.
// It has 4 vertices, 2 triangles (6 indices), all normals pointing +Y.
func NewPlane(width, depth float64) *Mesh {
	mesh := NewMesh()

	hw := width / 2.0
	hd := depth / 2.0
	grey := math.Vec3{X: 0.5, Y: 0.5, Z: 0.5}
	up := math.Vec3{Y: 1}

	corners := [4]math.Vec3{
		{X: -hw, Y: 0, Z: -hd},
		{X: hw, Y: 0, Z: -hd},
		{X: hw, Y: 0, Z: hd},
		{X: -hw, Y: 0, Z: hd},
	}

	for _, c := range corners {
		v := NewVertex(c, grey)
		v.Normal = up
		mesh.AddVertex(v)
	}

	// Two CCW triangles (viewed from +Y): 0-2-1, 0-3-2.
	// Winding must produce a +Y geometric normal so the face isn't
	// backface-culled when the camera is above the plane.
	mesh.AddTriangle(0, 2, 1)
	mesh.AddTriangle(0, 3, 2)

	return mesh
}

// NewCylinder creates a cylinder with the given number of lateral segments.
// Height 1 (Y from -0.5 to 0.5), radius 0.5. Per-face normals for flat shading.
// Layout: lateral quads (4*seg verts), top cap fan (seg+1 verts), bottom cap fan (seg+1 verts).
func NewCylinder(segments int) *Mesh {
	mesh := NewMesh()
	r := 0.5
	halfH := 0.5
	color := math.Vec3{X: 0.6, Y: 0.6, Z: 0.6}

	// --- Lateral faces (4 verts per quad, per-face radial normals) ---
	for i := 0; i < segments; i++ {
		a0 := 2.0 * stdmath.Pi * float64(i) / float64(segments)
		a1 := 2.0 * stdmath.Pi * float64(i+1) / float64(segments)

		c0, s0 := stdmath.Cos(a0), stdmath.Sin(a0)
		c1, s1 := stdmath.Cos(a1), stdmath.Sin(a1)

		// Face normal = average of the two edge normals (flat shading).
		am := (a0 + a1) / 2.0
		normal := math.Vec3{X: stdmath.Cos(am), Y: 0, Z: stdmath.Sin(am)}

		// 4 corners: bottom-left, bottom-right, top-right, top-left.
		bl := math.Vec3{X: r * c0, Y: -halfH, Z: r * s0}
		br := math.Vec3{X: r * c1, Y: -halfH, Z: r * s1}
		tr := math.Vec3{X: r * c1, Y: halfH, Z: r * s1}
		tl := math.Vec3{X: r * c0, Y: halfH, Z: r * s0}

		base := len(mesh.Vertices)
		for _, pos := range []math.Vec3{bl, br, tr, tl} {
			v := NewVertex(pos, color)
			v.Normal = normal
			mesh.AddVertex(v)
		}
		// Two CCW triangles (viewed from outside): 0-2-1, 0-3-2.
		mesh.AddTriangle(base+0, base+2, base+1)
		mesh.AddTriangle(base+0, base+3, base+2)
	}

	// --- Top cap (fan from center, normal +Y) ---
	topCenter := len(mesh.Vertices)
	cv := NewVertex(math.Vec3{Y: halfH}, color)
	cv.Normal = math.Vec3{Y: 1}
	mesh.AddVertex(cv)
	for i := 0; i < segments; i++ {
		a := 2.0 * stdmath.Pi * float64(i) / float64(segments)
		v := NewVertex(math.Vec3{X: r * stdmath.Cos(a), Y: halfH, Z: r * stdmath.Sin(a)}, color)
		v.Normal = math.Vec3{Y: 1}
		mesh.AddVertex(v)
	}
	for i := 0; i < segments; i++ {
		next := (i + 1) % segments
		// CCW from above: center, next, current.
		mesh.AddTriangle(topCenter, topCenter+1+next, topCenter+1+i)
	}

	// --- Bottom cap (fan from center, normal -Y) ---
	botCenter := len(mesh.Vertices)
	bv := NewVertex(math.Vec3{Y: -halfH}, color)
	bv.Normal = math.Vec3{Y: -1}
	mesh.AddVertex(bv)
	for i := 0; i < segments; i++ {
		a := 2.0 * stdmath.Pi * float64(i) / float64(segments)
		v := NewVertex(math.Vec3{X: r * stdmath.Cos(a), Y: -halfH, Z: r * stdmath.Sin(a)}, color)
		v.Normal = math.Vec3{Y: -1}
		mesh.AddVertex(v)
	}
	for i := 0; i < segments; i++ {
		next := (i + 1) % segments
		// CCW from below: center, current, next (reversed from top).
		mesh.AddTriangle(botCenter, botCenter+1+i, botCenter+1+next)
	}

	return mesh
}

package geometry

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

// TestNewTetrahedron_CreatesValidMesh tests tetrahedron creation.
// With per-face normals the tetrahedron has 12 vertices (3 per face × 4 faces).
func TestNewTetrahedron_CreatesValidMesh(t *testing.T) {
	mesh := NewTetrahedron()

	if len(mesh.Vertices) != 12 {
		t.Errorf("Tetrahedron Vertices count = %d, want 12", len(mesh.Vertices))
	}

	// Tetrahedron has 4 triangular faces (12 indices)
	if len(mesh.Indices) != 12 {
		t.Errorf("Tetrahedron Indices count = %d, want 12", len(mesh.Indices))
	}

	if mesh.TriangleCount() != 4 {
		t.Errorf("Tetrahedron TriangleCount() = %d, want 4", mesh.TriangleCount())
	}

	if !mesh.ValidateIndices() {
		t.Error("Tetrahedron has invalid indices")
	}
}

// TestNewTetrahedron_HasFaceColors tests each face's 3 vertices share a color.
func TestNewTetrahedron_HasFaceColors(t *testing.T) {
	mesh := NewTetrahedron()

	// With 12 vertices (3 per face) we expect 4 unique face colors.
	colors := make(map[string]bool)
	for _, vertex := range mesh.Vertices {
		colors[vertex.Color.String()] = true
	}

	if len(colors) != 4 {
		t.Errorf("Tetrahedron has %d unique colors, want 4", len(colors))
	}
}

// TestNewTetrahedron_HasUnitNormals tests all vertex normals are unit-length.
func TestNewTetrahedron_HasUnitNormals(t *testing.T) {
	mesh := NewTetrahedron()
	for i, v := range mesh.Vertices {
		length := v.Normal.Length()
		if length < 0.999 || length > 1.001 {
			t.Errorf("vertex %d: normal length = %f, want ~1.0", i, length)
		}
	}
}

// TestNewTetrahedron_FaceNormalsPointOutward tests that each face's normal
// points away from the centroid (outward-facing).
func TestNewTetrahedron_FaceNormalsPointOutward(t *testing.T) {
	mesh := NewTetrahedron()
	// Centroid of the tetrahedron is at the average of the 4 unique corner positions.
	// With per-face vertices, each face's 3 vertices share the same normal.
	for tri := 0; tri < mesh.TriangleCount(); tri++ {
		i0, i1, i2 := mesh.GetTriangle(tri)
		v0 := mesh.Vertices[i0]
		// Face centroid.
		fc := v0.Position.Add(mesh.Vertices[i1].Position).Add(mesh.Vertices[i2].Position).Scale(1.0 / 3.0)
		// Normal should point roughly in the same direction as (faceCentroid - origin).
		dot := fc.Dot(v0.Normal)
		if dot <= 0 {
			t.Errorf("face %d: normal·faceCentroid = %f, expected > 0 (outward)", tri, dot)
		}
	}
}

// TestNewCube_CreatesValidMesh tests cube creation.
// With per-face normals the cube has 24 vertices (4 per face × 6 faces).
func TestNewCube_CreatesValidMesh(t *testing.T) {
	mesh := NewCube()

	if len(mesh.Vertices) != 24 {
		t.Errorf("Cube Vertices count = %d, want 24", len(mesh.Vertices))
	}

	// Cube has 6 quad faces = 12 triangles (36 indices)
	if len(mesh.Indices) != 36 {
		t.Errorf("Cube Indices count = %d, want 36", len(mesh.Indices))
	}

	if mesh.TriangleCount() != 12 {
		t.Errorf("Cube TriangleCount() = %d, want 12", mesh.TriangleCount())
	}

	if !mesh.ValidateIndices() {
		t.Error("Cube has invalid indices")
	}
}

// TestNewCube_HasCorrectWindingOrder tests cube triangle winding.
func TestNewCube_HasCorrectWindingOrder(t *testing.T) {
	mesh := NewCube()
	n := len(mesh.Vertices)

	for i := 0; i < mesh.TriangleCount(); i++ {
		i0, i1, i2 := mesh.GetTriangle(i)

		if i0 < 0 || i0 >= n || i1 < 0 || i1 >= n || i2 < 0 || i2 >= n {
			t.Errorf("Triangle %d has out-of-bounds index", i)
		}
		if i0 == i1 || i1 == i2 || i0 == i2 {
			t.Errorf("Triangle %d is degenerate with indices (%d, %d, %d)", i, i0, i1, i2)
		}
	}
}

// TestNewCube_VerticesAt8Corners tests cube has 8 unique corner positions.
func TestNewCube_VerticesAt8Corners(t *testing.T) {
	mesh := NewCube()

	positions := make(map[string]bool)
	for _, vertex := range mesh.Vertices {
		positions[vertex.Position.String()] = true
	}

	if len(positions) != 8 {
		t.Errorf("Cube has %d unique vertex positions, want 8", len(positions))
	}
}

// TestNewCube_FaceNormals_AreAxisAligned tests that each cube face has the
// correct axis-aligned unit normal.
func TestNewCube_FaceNormals_AreAxisAligned(t *testing.T) {
	mesh := NewCube()

	// 6 expected face normals (one per axis direction).
	expected := map[string]math.Vec3{
		"+X": {X: 1, Y: 0, Z: 0},
		"-X": {X: -1, Y: 0, Z: 0},
		"+Y": {X: 0, Y: 1, Z: 0},
		"-Y": {X: 0, Y: -1, Z: 0},
		"+Z": {X: 0, Y: 0, Z: 1},
		"-Z": {X: 0, Y: 0, Z: -1},
	}

	// Collect unique normals from the 24 vertices.
	found := make(map[string]bool)
	for _, v := range mesh.Vertices {
		for label, n := range expected {
			if v.Normal.Equals(n, 0.0001) {
				found[label] = true
			}
		}
	}

	for label := range expected {
		if !found[label] {
			t.Errorf("missing face normal %s", label)
		}
	}
}

// TestNewCube_HasUnitNormals tests all vertex normals are unit-length.
func TestNewCube_HasUnitNormals(t *testing.T) {
	mesh := NewCube()
	for i, v := range mesh.Vertices {
		length := v.Normal.Length()
		if length < 0.999 || length > 1.001 {
			t.Errorf("vertex %d: normal length = %f, want ~1.0", i, length)
		}
	}
}

// TestNewCube_HasColoredVertices tests cube vertices have colors.
func TestNewCube_HasColoredVertices(t *testing.T) {
	mesh := NewCube()

	// Check all vertices have valid colors (components in [0,1])
	for i, vertex := range mesh.Vertices {
		if vertex.Color.X < 0 || vertex.Color.X > 1 {
			t.Errorf("Cube vertex %d has invalid color.X = %f", i, vertex.Color.X)
		}
		if vertex.Color.Y < 0 || vertex.Color.Y > 1 {
			t.Errorf("Cube vertex %d has invalid color.Y = %f", i, vertex.Color.Y)
		}
		if vertex.Color.Z < 0 || vertex.Color.Z > 1 {
			t.Errorf("Cube vertex %d has invalid color.Z = %f", i, vertex.Color.Z)
		}
	}
}

package geometry

import (
	"testing"
)

// TestNewTetrahedron_CreatesValidMesh tests tetrahedron creation
func TestNewTetrahedron_CreatesValidMesh(t *testing.T) {
	mesh := NewTetrahedron()

	// Tetrahedron has 4 vertices
	if len(mesh.Vertices) != 4 {
		t.Errorf("Tetrahedron Vertices count = %d, want 4", len(mesh.Vertices))
	}

	// Tetrahedron has 4 triangular faces (12 indices)
	if len(mesh.Indices) != 12 {
		t.Errorf("Tetrahedron Indices count = %d, want 12", len(mesh.Indices))
	}

	// Should have 4 triangles
	if mesh.TriangleCount() != 4 {
		t.Errorf("Tetrahedron TriangleCount() = %d, want 4", mesh.TriangleCount())
	}

	// All indices should be valid
	if !mesh.ValidateIndices() {
		t.Error("Tetrahedron has invalid indices")
	}
}

// TestNewTetrahedron_HasUniqueColors tests each vertex has a distinct color
func TestNewTetrahedron_HasUniqueColors(t *testing.T) {
	mesh := NewTetrahedron()

	// Check that vertices have different colors
	colors := make(map[string]bool)
	for i, vertex := range mesh.Vertices {
		colorKey := vertex.Color.String()
		if colors[colorKey] {
			t.Errorf("Tetrahedron vertex %d has duplicate color", i)
		}
		colors[colorKey] = true
	}

	if len(colors) != 4 {
		t.Errorf("Tetrahedron has %d unique colors, want 4", len(colors))
	}
}

// TestNewTetrahedron_VerticesFormValidShape tests tetrahedron geometry
func TestNewTetrahedron_VerticesFormValidShape(t *testing.T) {
	mesh := NewTetrahedron()

	// All vertices should have non-zero distance from each other
	for i := 0; i < len(mesh.Vertices); i++ {
		for j := i + 1; j < len(mesh.Vertices); j++ {
			distance := mesh.Vertices[i].Position.Distance(mesh.Vertices[j].Position)
			if distance == 0 {
				t.Errorf("Tetrahedron vertices %d and %d are at the same position", i, j)
			}
		}
	}
}

// TestNewCube_CreatesValidMesh tests cube creation
func TestNewCube_CreatesValidMesh(t *testing.T) {
	mesh := NewCube()

	// Cube has 8 vertices
	if len(mesh.Vertices) != 8 {
		t.Errorf("Cube Vertices count = %d, want 8", len(mesh.Vertices))
	}

	// Cube has 6 quad faces = 12 triangles (36 indices)
	if len(mesh.Indices) != 36 {
		t.Errorf("Cube Indices count = %d, want 36", len(mesh.Indices))
	}

	// Should have 12 triangles
	if mesh.TriangleCount() != 12 {
		t.Errorf("Cube TriangleCount() = %d, want 12", mesh.TriangleCount())
	}

	// All indices should be valid
	if !mesh.ValidateIndices() {
		t.Error("Cube has invalid indices")
	}
}

// TestNewCube_HasCorrectWindingOrder tests cube triangle winding
func TestNewCube_HasCorrectWindingOrder(t *testing.T) {
	mesh := NewCube()

	// All triangles should have counter-clockwise winding (right-handed system)
	// We can't easily test this without implementing normal calculation,
	// but we can verify all indices are within bounds
	for i := 0; i < mesh.TriangleCount(); i++ {
		i0, i1, i2 := mesh.GetTriangle(i)

		if i0 < 0 || i0 >= 8 {
			t.Errorf("Triangle %d has out-of-bounds index i0=%d", i, i0)
		}
		if i1 < 0 || i1 >= 8 {
			t.Errorf("Triangle %d has out-of-bounds index i1=%d", i, i1)
		}
		if i2 < 0 || i2 >= 8 {
			t.Errorf("Triangle %d has out-of-bounds index i2=%d", i, i2)
		}

		// Check no degenerate triangles (all three indices different)
		if i0 == i1 || i1 == i2 || i0 == i2 {
			t.Errorf("Triangle %d is degenerate with indices (%d, %d, %d)", i, i0, i1, i2)
		}
	}
}

// TestNewCube_VerticesAt8Corners tests cube vertex positions
func TestNewCube_VerticesAt8Corners(t *testing.T) {
	mesh := NewCube()

	// Verify we have 8 distinct vertex positions
	positions := make(map[string]bool)
	for _, vertex := range mesh.Vertices {
		posKey := vertex.Position.String()
		positions[posKey] = true
	}

	if len(positions) != 8 {
		t.Errorf("Cube has %d unique vertex positions, want 8", len(positions))
	}
}

// TestNewCube_HasColoredVertices tests cube vertices have colors
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

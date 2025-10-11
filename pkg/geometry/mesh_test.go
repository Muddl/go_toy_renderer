package geometry

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

// TestMesh_Creation_CreatesEmptyMesh tests creating an empty mesh
func TestMesh_Creation_CreatesEmptyMesh(t *testing.T) {
	mesh := NewMesh()

	if len(mesh.Vertices) != 0 {
		t.Errorf("NewMesh().Vertices length = %d, want 0", len(mesh.Vertices))
	}

	if len(mesh.Indices) != 0 {
		t.Errorf("NewMesh().Indices length = %d, want 0", len(mesh.Indices))
	}
}

// TestMesh_AddVertex_AddsVertexToMesh tests adding vertices to a mesh
func TestMesh_AddVertex_AddsVertexToMesh(t *testing.T) {
	mesh := NewMesh()

	v1 := NewVertex(math.Vec3{X: 1.0, Y: 2.0, Z: 3.0}, math.Vec3{X: 1.0, Y: 0.0, Z: 0.0})
	v2 := NewVertex(math.Vec3{X: 4.0, Y: 5.0, Z: 6.0}, math.Vec3{X: 0.0, Y: 1.0, Z: 0.0})

	mesh.AddVertex(v1)
	mesh.AddVertex(v2)

	if len(mesh.Vertices) != 2 {
		t.Errorf("After adding 2 vertices, Vertices length = %d, want 2", len(mesh.Vertices))
	}

	if !mesh.Vertices[0].Equals(v1, 0.0001) {
		t.Errorf("Vertices[0] = %v, want %v", mesh.Vertices[0], v1)
	}

	if !mesh.Vertices[1].Equals(v2, 0.0001) {
		t.Errorf("Vertices[1] = %v, want %v", mesh.Vertices[1], v2)
	}
}

// TestMesh_AddTriangle_AddsTriangleIndices tests adding triangle indices
func TestMesh_AddTriangle_AddsTriangleIndices(t *testing.T) {
	mesh := NewMesh()

	mesh.AddTriangle(0, 1, 2)
	mesh.AddTriangle(0, 2, 3)

	if len(mesh.Indices) != 6 {
		t.Errorf("After adding 2 triangles, Indices length = %d, want 6", len(mesh.Indices))
	}

	expectedIndices := []int{0, 1, 2, 0, 2, 3}
	for i, expected := range expectedIndices {
		if mesh.Indices[i] != expected {
			t.Errorf("Indices[%d] = %d, want %d", i, mesh.Indices[i], expected)
		}
	}
}

// TestMesh_TriangleCount_ReturnsNumberOfTriangles tests triangle counting
func TestMesh_TriangleCount_ReturnsNumberOfTriangles(t *testing.T) {
	tests := []struct {
		name     string
		indices  []int
		expected int
	}{
		{
			name:     "empty mesh",
			indices:  []int{},
			expected: 0,
		},
		{
			name:     "one triangle",
			indices:  []int{0, 1, 2},
			expected: 1,
		},
		{
			name:     "two triangles",
			indices:  []int{0, 1, 2, 0, 2, 3},
			expected: 2,
		},
		{
			name:     "three triangles",
			indices:  []int{0, 1, 2, 2, 3, 4, 4, 5, 6},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mesh := Mesh{Indices: tt.indices}
			count := mesh.TriangleCount()

			if count != tt.expected {
				t.Errorf("TriangleCount() = %d, want %d", count, tt.expected)
			}
		})
	}
}

// TestMesh_GetTriangle_ReturnsTriangleIndices tests retrieving triangle indices
func TestMesh_GetTriangle_ReturnsTriangleIndices(t *testing.T) {
	mesh := NewMesh()
	mesh.AddTriangle(0, 1, 2)
	mesh.AddTriangle(3, 4, 5)
	mesh.AddTriangle(6, 7, 8)

	tests := []struct {
		name          string
		triangleIndex int
		expectedI0    int
		expectedI1    int
		expectedI2    int
	}{
		{
			name:          "first triangle",
			triangleIndex: 0,
			expectedI0:    0,
			expectedI1:    1,
			expectedI2:    2,
		},
		{
			name:          "second triangle",
			triangleIndex: 1,
			expectedI0:    3,
			expectedI1:    4,
			expectedI2:    5,
		},
		{
			name:          "third triangle",
			triangleIndex: 2,
			expectedI0:    6,
			expectedI1:    7,
			expectedI2:    8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i0, i1, i2 := mesh.GetTriangle(tt.triangleIndex)

			if i0 != tt.expectedI0 {
				t.Errorf("GetTriangle(%d) i0 = %d, want %d", tt.triangleIndex, i0, tt.expectedI0)
			}
			if i1 != tt.expectedI1 {
				t.Errorf("GetTriangle(%d) i1 = %d, want %d", tt.triangleIndex, i1, tt.expectedI1)
			}
			if i2 != tt.expectedI2 {
				t.Errorf("GetTriangle(%d) i2 = %d, want %d", tt.triangleIndex, i2, tt.expectedI2)
			}
		})
	}
}

// TestMesh_GetTriangleVertices_ReturnsVertexData tests retrieving triangle vertex data
func TestMesh_GetTriangleVertices_ReturnsVertexData(t *testing.T) {
	mesh := NewMesh()

	// Add vertices
	v0 := NewVertex(math.Vec3{X: 0.0, Y: 0.0, Z: 0.0}, math.Vec3{X: 1.0, Y: 0.0, Z: 0.0})
	v1 := NewVertex(math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}, math.Vec3{X: 0.0, Y: 1.0, Z: 0.0})
	v2 := NewVertex(math.Vec3{X: 0.0, Y: 1.0, Z: 0.0}, math.Vec3{X: 0.0, Y: 0.0, Z: 1.0})

	mesh.AddVertex(v0)
	mesh.AddVertex(v1)
	mesh.AddVertex(v2)

	// Add triangle
	mesh.AddTriangle(0, 1, 2)

	// Get triangle vertices
	vert0, vert1, vert2 := mesh.GetTriangleVertices(0)

	if !vert0.Equals(v0, 0.0001) {
		t.Errorf("GetTriangleVertices(0) vert0 = %v, want %v", vert0, v0)
	}
	if !vert1.Equals(v1, 0.0001) {
		t.Errorf("GetTriangleVertices(0) vert1 = %v, want %v", vert1, v1)
	}
	if !vert2.Equals(v2, 0.0001) {
		t.Errorf("GetTriangleVertices(0) vert2 = %v, want %v", vert2, v2)
	}
}

// TestMesh_ValidateIndices_ChecksIndexBounds tests index validation
func TestMesh_ValidateIndices_ChecksIndexBounds(t *testing.T) {
	tests := []struct {
		name        string
		vertices    []Vertex
		indices     []int
		expectValid bool
	}{
		{
			name:        "empty mesh is valid",
			vertices:    []Vertex{},
			indices:     []int{},
			expectValid: true,
		},
		{
			name: "valid indices",
			vertices: []Vertex{
				NewVertex(math.Vec3{}, math.Vec3{}),
				NewVertex(math.Vec3{}, math.Vec3{}),
				NewVertex(math.Vec3{}, math.Vec3{}),
			},
			indices:     []int{0, 1, 2},
			expectValid: true,
		},
		{
			name: "index out of bounds",
			vertices: []Vertex{
				NewVertex(math.Vec3{}, math.Vec3{}),
				NewVertex(math.Vec3{}, math.Vec3{}),
			},
			indices:     []int{0, 1, 5},
			expectValid: false,
		},
		{
			name: "negative index",
			vertices: []Vertex{
				NewVertex(math.Vec3{}, math.Vec3{}),
			},
			indices:     []int{-1, 0, 0},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mesh := Mesh{
				Vertices: tt.vertices,
				Indices:  tt.indices,
			}

			valid := mesh.ValidateIndices()

			if valid != tt.expectValid {
				t.Errorf("ValidateIndices() = %v, want %v", valid, tt.expectValid)
			}
		})
	}
}

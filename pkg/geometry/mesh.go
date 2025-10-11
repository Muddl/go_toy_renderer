package geometry

// Mesh represents a 3D mesh composed of vertices and triangle indices.
// The indices form triangles by referencing vertices in groups of three.
type Mesh struct {
	Vertices []Vertex // Vertex data (positions, colors, etc.)
	Indices  []int    // Triangle indices (every 3 indices forms a triangle)
}

// NewMesh creates an empty mesh with no vertices or indices.
func NewMesh() *Mesh {
	return &Mesh{
		Vertices: make([]Vertex, 0),
		Indices:  make([]int, 0),
	}
}

// AddVertex adds a vertex to the mesh and returns its index.
func (m *Mesh) AddVertex(v Vertex) int {
	m.Vertices = append(m.Vertices, v)
	return len(m.Vertices) - 1
}

// AddTriangle adds a triangle to the mesh using three vertex indices.
// The indices should reference valid vertices in the Vertices slice.
func (m *Mesh) AddTriangle(i0, i1, i2 int) {
	m.Indices = append(m.Indices, i0, i1, i2)
}

// TriangleCount returns the number of triangles in the mesh.
// Each triangle is composed of 3 indices.
func (m *Mesh) TriangleCount() int {
	return len(m.Indices) / 3
}

// GetTriangle returns the three vertex indices for the specified triangle.
// triangleIndex should be in the range [0, TriangleCount()).
func (m *Mesh) GetTriangle(triangleIndex int) (i0, i1, i2 int) {
	baseIndex := triangleIndex * 3
	return m.Indices[baseIndex], m.Indices[baseIndex+1], m.Indices[baseIndex+2]
}

// GetTriangleVertices returns the three vertices for the specified triangle.
// triangleIndex should be in the range [0, TriangleCount()).
func (m *Mesh) GetTriangleVertices(triangleIndex int) (v0, v1, v2 Vertex) {
	i0, i1, i2 := m.GetTriangle(triangleIndex)
	return m.Vertices[i0], m.Vertices[i1], m.Vertices[i2]
}

// ValidateIndices checks if all indices in the mesh reference valid vertices.
// Returns true if all indices are within bounds [0, len(Vertices)), false otherwise.
func (m *Mesh) ValidateIndices() bool {
	vertexCount := len(m.Vertices)

	for _, index := range m.Indices {
		if index < 0 || index >= vertexCount {
			return false
		}
	}

	return true
}

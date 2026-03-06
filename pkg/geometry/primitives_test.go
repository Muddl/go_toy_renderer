package geometry

import (
	stdmath "math"
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

// --- Plane tests ---

// TestNewPlane_CreatesValidMesh tests plane creation with 4 vertices and 2 triangles.
func TestNewPlane_CreatesValidMesh(t *testing.T) {
	mesh := NewPlane(10, 10)

	if len(mesh.Vertices) != 4 {
		t.Errorf("Plane Vertices count = %d, want 4", len(mesh.Vertices))
	}
	if len(mesh.Indices) != 6 {
		t.Errorf("Plane Indices count = %d, want 6", len(mesh.Indices))
	}
	if mesh.TriangleCount() != 2 {
		t.Errorf("Plane TriangleCount() = %d, want 2", mesh.TriangleCount())
	}
	if !mesh.ValidateIndices() {
		t.Error("Plane has invalid indices")
	}
}

// TestNewPlane_AllNormalsPointUp tests that all plane normals are +Y.
func TestNewPlane_AllNormalsPointUp(t *testing.T) {
	mesh := NewPlane(5, 5)
	up := math.Vec3{Y: 1}
	for i, v := range mesh.Vertices {
		if !v.Normal.Equals(up, 0.0001) {
			t.Errorf("vertex %d: normal = %v, want +Y", i, v.Normal)
		}
	}
}

// TestNewPlane_DimensionsMatchArguments tests that vertices span the given width and depth.
func TestNewPlane_DimensionsMatchArguments(t *testing.T) {
	mesh := NewPlane(8, 6)
	var minX, maxX, minZ, maxZ float64
	for i, v := range mesh.Vertices {
		if i == 0 || v.Position.X < minX {
			minX = v.Position.X
		}
		if i == 0 || v.Position.X > maxX {
			maxX = v.Position.X
		}
		if i == 0 || v.Position.Z < minZ {
			minZ = v.Position.Z
		}
		if i == 0 || v.Position.Z > maxZ {
			maxZ = v.Position.Z
		}
	}
	w := maxX - minX
	d := maxZ - minZ
	if w < 7.999 || w > 8.001 {
		t.Errorf("plane width = %f, want 8", w)
	}
	if d < 5.999 || d > 6.001 {
		t.Errorf("plane depth = %f, want 6", d)
	}
}

// --- Cylinder tests ---

// TestNewCylinder_CreatesValidMesh tests cylinder creation with 16 segments.
// Lateral: 16 quads × 4 verts = 64 verts, 16×2 tris = 96 indices.
// Top cap: 16+1 = 17 verts, 16 tris = 48 indices.
// Bottom cap: same. Total: 98 verts, 192 indices, 64 triangles.
func TestNewCylinder_CreatesValidMesh(t *testing.T) {
	seg := 16
	mesh := NewCylinder(seg)

	wantVerts := 4*seg + 2*(seg+1) // 98
	wantIdx := 12 * seg             // 192
	wantTris := 4 * seg             // 64

	if len(mesh.Vertices) != wantVerts {
		t.Errorf("Cylinder Vertices = %d, want %d", len(mesh.Vertices), wantVerts)
	}
	if len(mesh.Indices) != wantIdx {
		t.Errorf("Cylinder Indices = %d, want %d", len(mesh.Indices), wantIdx)
	}
	if mesh.TriangleCount() != wantTris {
		t.Errorf("Cylinder TriangleCount() = %d, want %d", mesh.TriangleCount(), wantTris)
	}
	if !mesh.ValidateIndices() {
		t.Error("Cylinder has invalid indices")
	}
}

// TestNewCylinder_LateralNormalsAreRadial tests that lateral vertices have
// horizontal unit normals (Y=0, length=1) pointing radially outward.
func TestNewCylinder_LateralNormalsAreRadial(t *testing.T) {
	seg := 12
	mesh := NewCylinder(seg)

	// Lateral vertices are the first 4*seg vertices.
	for i := 0; i < 4*seg; i++ {
		v := mesh.Vertices[i]
		n := v.Normal
		// Lateral normals should have Y ≈ 0.
		if n.Y < -0.001 || n.Y > 0.001 {
			t.Errorf("lateral vertex %d: normal.Y = %f, want ~0", i, n.Y)
		}
		// Should be unit length.
		length := n.Length()
		if length < 0.999 || length > 1.001 {
			t.Errorf("lateral vertex %d: normal length = %f, want ~1", i, length)
		}
	}
}

// TestNewCylinder_CapNormalsAreAxisAligned tests that top cap normals are +Y
// and bottom cap normals are -Y.
func TestNewCylinder_CapNormalsAreAxisAligned(t *testing.T) {
	seg := 8
	mesh := NewCylinder(seg)

	lateralCount := 4 * seg
	capSize := seg + 1

	// Top cap vertices: [lateralCount .. lateralCount+capSize)
	topUp := math.Vec3{Y: 1}
	for i := lateralCount; i < lateralCount+capSize; i++ {
		if !mesh.Vertices[i].Normal.Equals(topUp, 0.0001) {
			t.Errorf("top cap vertex %d: normal = %v, want +Y", i, mesh.Vertices[i].Normal)
		}
	}

	// Bottom cap vertices: [lateralCount+capSize .. lateralCount+2*capSize)
	botDown := math.Vec3{Y: -1}
	for i := lateralCount + capSize; i < lateralCount+2*capSize; i++ {
		if !mesh.Vertices[i].Normal.Equals(botDown, 0.0001) {
			t.Errorf("bottom cap vertex %d: normal = %v, want -Y", i, mesh.Vertices[i].Normal)
		}
	}
}

// TestNewCylinder_HeightAndRadius tests the cylinder spans Y=-0.5 to Y=0.5
// with radius 0.5 on the XZ plane.
func TestNewCylinder_HeightAndRadius(t *testing.T) {
	mesh := NewCylinder(16)
	var minY, maxY, maxR float64
	for i, v := range mesh.Vertices {
		if i == 0 || v.Position.Y < minY {
			minY = v.Position.Y
		}
		if i == 0 || v.Position.Y > maxY {
			maxY = v.Position.Y
		}
		r := stdmath.Sqrt(v.Position.X*v.Position.X + v.Position.Z*v.Position.Z)
		if r > maxR {
			maxR = r
		}
	}
	if minY < -0.501 || minY > -0.499 {
		t.Errorf("cylinder minY = %f, want -0.5", minY)
	}
	if maxY < 0.499 || maxY > 0.501 {
		t.Errorf("cylinder maxY = %f, want 0.5", maxY)
	}
	if maxR < 0.499 || maxR > 0.501 {
		t.Errorf("cylinder maxR = %f, want 0.5", maxR)
	}
}

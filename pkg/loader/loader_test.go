package loader_test

import (
	"path/filepath"
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/loader"
)

const eps = 0.0001

func TestLoadOBJ_Triangle_ParsesPositions(t *testing.T) {
	meshes, err := loader.LoadOBJ(filepath.Join("testdata", "triangle.obj"))
	if err != nil {
		t.Fatalf("LoadOBJ returned error: %v", err)
	}
	if len(meshes) != 1 {
		t.Fatalf("expected 1 mesh, got %d", len(meshes))
	}
	m := meshes[0]
	if len(m.Vertices) != 3 {
		t.Fatalf("expected 3 vertices, got %d", len(m.Vertices))
	}
	// First vertex position should be (0, 0, 0)
	v0 := m.Vertices[0]
	if !approxEq(v0.Position.X, 0.0) || !approxEq(v0.Position.Y, 0.0) || !approxEq(v0.Position.Z, 0.0) {
		t.Errorf("vertex 0 position: got (%.4f, %.4f, %.4f), want (0, 0, 0)",
			v0.Position.X, v0.Position.Y, v0.Position.Z)
	}
	// Second vertex position should be (1, 0, 0)
	v1 := m.Vertices[1]
	if !approxEq(v1.Position.X, 1.0) || !approxEq(v1.Position.Y, 0.0) || !approxEq(v1.Position.Z, 0.0) {
		t.Errorf("vertex 1 position: got (%.4f, %.4f, %.4f), want (1, 0, 0)",
			v1.Position.X, v1.Position.Y, v1.Position.Z)
	}
}

func TestLoadOBJ_Triangle_ParsesNormals(t *testing.T) {
	meshes, err := loader.LoadOBJ(filepath.Join("testdata", "triangle.obj"))
	if err != nil {
		t.Fatalf("LoadOBJ returned error: %v", err)
	}
	m := meshes[0]
	for i, v := range m.Vertices {
		if !approxEq(v.Normal.X, 0.0) || !approxEq(v.Normal.Y, 0.0) || !approxEq(v.Normal.Z, 1.0) {
			t.Errorf("vertex %d normal: got (%.4f, %.4f, %.4f), want (0, 0, 1)",
				i, v.Normal.X, v.Normal.Y, v.Normal.Z)
		}
	}
}

func TestLoadOBJ_Triangle_ParsesUVs(t *testing.T) {
	meshes, err := loader.LoadOBJ(filepath.Join("testdata", "triangle.obj"))
	if err != nil {
		t.Fatalf("LoadOBJ returned error: %v", err)
	}
	m := meshes[0]
	wantUVs := [][2]float64{{0.0, 0.0}, {1.0, 0.0}, {0.0, 1.0}}
	for i, v := range m.Vertices {
		if !approxEq(v.UV[0], wantUVs[i][0]) || !approxEq(v.UV[1], wantUVs[i][1]) {
			t.Errorf("vertex %d UV: got (%.4f, %.4f), want (%.4f, %.4f)",
				i, v.UV[0], v.UV[1], wantUVs[i][0], wantUVs[i][1])
		}
	}
}

func TestLoadOBJ_Triangle_ParsesFaceIndices(t *testing.T) {
	meshes, err := loader.LoadOBJ(filepath.Join("testdata", "triangle.obj"))
	if err != nil {
		t.Fatalf("LoadOBJ returned error: %v", err)
	}
	m := meshes[0]
	if m.TriangleCount() != 1 {
		t.Fatalf("expected 1 triangle, got %d", m.TriangleCount())
	}
	i0, i1, i2 := m.GetTriangle(0)
	if i0 != 0 || i1 != 1 || i2 != 2 {
		t.Errorf("triangle indices: got (%d, %d, %d), want (0, 1, 2)", i0, i1, i2)
	}
}

func TestLoadOBJ_Quad_TwoTriangles(t *testing.T) {
	meshes, err := loader.LoadOBJ(filepath.Join("testdata", "quad.obj"))
	if err != nil {
		t.Fatalf("LoadOBJ returned error: %v", err)
	}
	m := meshes[0]
	if m.TriangleCount() != 2 {
		t.Fatalf("expected 2 triangles, got %d", m.TriangleCount())
	}
}

func TestLoadOBJ_MissingFile_ReturnsError(t *testing.T) {
	_, err := loader.LoadOBJ("testdata/nonexistent.obj")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadOBJ_ValidatesIndices(t *testing.T) {
	meshes, err := loader.LoadOBJ(filepath.Join("testdata", "triangle.obj"))
	if err != nil {
		t.Fatalf("LoadOBJ returned error: %v", err)
	}
	if !meshes[0].ValidateIndices() {
		t.Error("mesh indices failed validation")
	}
}

func approxEq(a, b float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d < eps
}

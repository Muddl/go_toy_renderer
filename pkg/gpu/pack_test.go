package gpu_test

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/gpu"
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// TestPackVertices_EmptySlice_ReturnsEmpty verifies that PackVertices returns an
// empty slice when given nil input, with no allocation.
func TestPackVertices_EmptySlice_ReturnsEmpty(t *testing.T) {
	result := gpu.PackVertices(nil)
	if len(result) != 0 {
		t.Errorf("expected empty slice, got len=%d", len(result))
	}
}

// TestPackVertices_SingleVertex_PacksPositionColorNormalUV verifies the float32
// layout for a single vertex: px,py,pz, cr,cg,cb, nx,ny,nz, u,v (11 floats, stride 44 bytes).
func TestPackVertices_SingleVertex_PacksPositionColorNormalUV(t *testing.T) {
	v := geometry.NewVertex(
		math.Vec3{X: 1.0, Y: 2.0, Z: 3.0},
		math.Vec3{X: 0.5, Y: 0.25, Z: 0.75},
	)
	v.Normal = math.Vec3{X: 0, Y: 1, Z: 0}
	v.UV = [2]float64{0.25, 0.75}
	result := gpu.PackVertices([]geometry.Vertex{v})

	if len(result) != 11 {
		t.Fatalf("expected 11 floats for one vertex, got %d", len(result))
	}
	expected := [11]float32{1.0, 2.0, 3.0, 0.5, 0.25, 0.75, 0, 1, 0, 0.25, 0.75}
	for i, want := range expected {
		if result[i] != want {
			t.Errorf("result[%d] = %v, want %v", i, result[i], want)
		}
	}
}

// TestPackVertices_Stride_Is11FloatsPerVertex verifies that consecutive vertices
// are laid out without gaps: stride = 11 floats = 44 bytes.
func TestPackVertices_Stride_Is11FloatsPerVertex(t *testing.T) {
	v0 := geometry.NewVertex(math.Vec3{X: 1, Y: 0, Z: 0}, math.Vec3{X: 1, Y: 0, Z: 0})
	v1 := geometry.NewVertex(math.Vec3{X: 2, Y: 0, Z: 0}, math.Vec3{X: 0, Y: 1, Z: 0})
	result := gpu.PackVertices([]geometry.Vertex{v0, v1})

	if len(result) != 22 {
		t.Fatalf("two vertices should produce 22 floats, got %d", len(result))
	}
	// Second vertex starts at index 11 (stride = 11 floats).
	if result[11] != float32(2.0) {
		t.Errorf("second vertex position.X = %v, want 2.0", result[11])
	}
}

// TestPackIndices_EmptySlice_ReturnsEmpty verifies PackIndices returns an empty
// slice for nil input.
func TestPackIndices_EmptySlice_ReturnsEmpty(t *testing.T) {
	result := gpu.PackIndices(nil)
	if len(result) != 0 {
		t.Errorf("expected empty slice, got len=%d", len(result))
	}
}

// TestPackIndices_ConvertsIntToUint32 verifies that each int index is converted
// to the corresponding uint32 value without truncation.
func TestPackIndices_ConvertsIntToUint32(t *testing.T) {
	indices := []int{0, 1, 2, 3}
	result := gpu.PackIndices(indices)

	if len(result) != 4 {
		t.Fatalf("expected 4 uint32 values, got %d", len(result))
	}
	expected := [4]uint32{0, 1, 2, 3}
	for i, want := range expected {
		if result[i] != want {
			t.Errorf("result[%d] = %d, want %d", i, result[i], want)
		}
	}
}

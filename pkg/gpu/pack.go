package gpu

import "github.com/muddl/go_toy_renderer/pkg/geometry"

// PackVertices converts a slice of geometry.Vertex into []float32 for GPU upload.
// Each vertex is packed as: px, py, pz, cr, cg, cb (6 float32 values, stride 24 bytes).
func PackVertices(vertices []geometry.Vertex) []float32 {
	out := make([]float32, len(vertices)*6)
	for i, v := range vertices {
		base := i * 6
		out[base+0] = float32(v.Position.X)
		out[base+1] = float32(v.Position.Y)
		out[base+2] = float32(v.Position.Z)
		out[base+3] = float32(v.Color.X)
		out[base+4] = float32(v.Color.Y)
		out[base+5] = float32(v.Color.Z)
	}
	return out
}

// PackIndices converts a slice of int indices into []uint32 for GPU index buffer upload.
// Mesh indices are validated to be in [0, len(Vertices)) by Mesh.ValidateIndices,
// so the int→uint32 conversion cannot overflow in practice.
func PackIndices(indices []int) []uint32 {
	out := make([]uint32, len(indices))
	for i, idx := range indices {
		out[i] = uint32(idx) //nolint:gosec // bounded by mesh vertex count
	}
	return out
}

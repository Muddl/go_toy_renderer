package gpu

import "github.com/muddl/go_toy_renderer/pkg/geometry"

// PackVertices converts a slice of geometry.Vertex into []float32 for GPU upload.
// Each vertex is packed as: px, py, pz, cr, cg, cb, nx, ny, nz, u, v
// (11 float32 values, stride 44 bytes).
func PackVertices(vertices []geometry.Vertex) []float32 {
	out := make([]float32, len(vertices)*11)
	for i, v := range vertices {
		base := i * 11
		out[base+0] = float32(v.Position.X)
		out[base+1] = float32(v.Position.Y)
		out[base+2] = float32(v.Position.Z)
		out[base+3] = float32(v.Color.X)
		out[base+4] = float32(v.Color.Y)
		out[base+5] = float32(v.Color.Z)
		out[base+6] = float32(v.Normal.X)
		out[base+7] = float32(v.Normal.Y)
		out[base+8] = float32(v.Normal.Z)
		out[base+9] = float32(v.UV[0])
		out[base+10] = float32(v.UV[1])
	}
	return out
}

// PackIndices converts a slice of int indices into []uint32 for GPU index buffer upload.
// Mesh indices are validated to be in [0, len(Vertices)) by Mesh.ValidateIndices,
// so the int→uint32 conversion cannot overflow in practice.
func PackIndices(indices []int) []uint32 {
	out := make([]uint32, len(indices))
	for i, idx := range indices {
		out[i] = uint32(idx)
	}
	return out
}

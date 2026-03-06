package gpu

import (
	"encoding/binary"
	gomath "math"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

// CameraUniforms holds the camera view-projection matrix for GPU upload.
// Layout: mat4x4<f32> (64 bytes, column-major, 16-byte aligned).
type CameraUniforms struct {
	ViewProj math.Mat4x4
}

// Bytes serializes the view-projection matrix to 64 bytes of little-endian float32
// in column-major order, matching the WGSL mat4x4<f32> layout.
func (cu CameraUniforms) Bytes() []byte {
	return mat4ToBytes(cu.ViewProj)
}

// MeshUniforms holds the per-mesh model matrix for GPU upload.
// Layout: mat4x4<f32> (64 bytes, column-major, 16-byte aligned).
type MeshUniforms struct {
	Model math.Mat4x4
}

// Bytes serializes the model matrix to 64 bytes of little-endian float32
// in column-major order, matching the WGSL mat4x4<f32> layout.
func (mu MeshUniforms) Bytes() []byte {
	return mat4ToBytes(mu.Model)
}

// mat4ToBytes converts a Mat4x4 (column-major float64) to 64 bytes of
// little-endian float32.
func mat4ToBytes(m math.Mat4x4) []byte {
	out := make([]byte, 64)
	for i := 0; i < 16; i++ {
		binary.LittleEndian.PutUint32(out[i*4:], gomath.Float32bits(float32(m[i])))
	}
	return out
}

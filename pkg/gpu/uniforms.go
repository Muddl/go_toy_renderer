package gpu

import (
	"encoding/binary"
	gomath "math"

	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// DrawNode describes a mesh to draw with its model and normal matrices.
type DrawNode struct {
	Mesh         *geometry.Mesh
	Model        math.Mat4x4
	NormalMatrix math.Mat4x4
}

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

// MeshUniforms holds the per-mesh model and normal matrices for GPU upload.
// Layout: model mat4x4<f32> (64 bytes) + normalMatrix mat4x4<f32> (64 bytes) = 128 bytes.
type MeshUniforms struct {
	Model        math.Mat4x4
	NormalMatrix math.Mat4x4
}

// Bytes serializes model + normalMatrix to 128 bytes of little-endian float32.
func (mu MeshUniforms) Bytes() []byte {
	out := make([]byte, 128)
	copy(out[:64], mat4ToBytes(mu.Model))
	copy(out[64:], mat4ToBytes(mu.NormalMatrix))
	return out
}

// LightUniforms holds directional light parameters for GPU upload.
// Layout: 4 × vec4<f32> = 64 bytes (each vec3 padded to 16 bytes).
//   - direction: normalized light direction (world-space)
//   - color: light RGB intensity
//   - ambient: ambient light RGB
//   - cameraPos: camera world position (for specular)
type LightUniforms struct {
	Direction math.Vec3
	Color     math.Vec3
	Ambient   math.Vec3
	CameraPos math.Vec3
}

// Bytes serializes the light uniforms to 64 bytes (4 × vec4, each vec3 + 4 bytes padding).
func (lu LightUniforms) Bytes() []byte {
	out := make([]byte, 64)
	vecs := [4]math.Vec3{lu.Direction, lu.Color, lu.Ambient, lu.CameraPos}
	for i, v := range vecs {
		off := i * 16
		binary.LittleEndian.PutUint32(out[off:], gomath.Float32bits(float32(v.X)))
		binary.LittleEndian.PutUint32(out[off+4:], gomath.Float32bits(float32(v.Y)))
		binary.LittleEndian.PutUint32(out[off+8:], gomath.Float32bits(float32(v.Z)))
		// out[off+12..off+15] remains zero (padding)
	}
	return out
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

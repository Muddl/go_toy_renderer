package gpu

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

func TestCameraUniforms_ByteSize(t *testing.T) {
	cu := CameraUniforms{ViewProj: math.NewIdentity()}
	data := cu.Bytes()
	if len(data) != 64 {
		t.Fatalf("CameraUniforms.Bytes() = %d bytes, want 64", len(data))
	}
}

func TestMeshUniforms_ByteSize(t *testing.T) {
	mu := MeshUniforms{Model: math.NewIdentity()}
	data := mu.Bytes()
	if len(data) != 64 {
		t.Fatalf("MeshUniforms.Bytes() = %d bytes, want 64", len(data))
	}
}

func TestCameraUniforms_Alignment(t *testing.T) {
	// 64 bytes = 4 columns * 4 rows * 4 bytes/float32 = 64; 64 % 16 == 0
	cu := CameraUniforms{ViewProj: math.NewIdentity()}
	data := cu.Bytes()
	if len(data)%16 != 0 {
		t.Fatalf("CameraUniforms not 16-byte aligned: %d bytes", len(data))
	}
}

func TestMeshUniforms_Alignment(t *testing.T) {
	mu := MeshUniforms{Model: math.NewIdentity()}
	data := mu.Bytes()
	if len(data)%16 != 0 {
		t.Fatalf("MeshUniforms not 16-byte aligned: %d bytes", len(data))
	}
}

func TestCameraUniforms_IdentityBytes(t *testing.T) {
	cu := CameraUniforms{ViewProj: math.NewIdentity()}
	data := cu.Bytes()
	// Column-major identity: diag elements at byte offsets 0, 20, 40, 60
	// Each float32 = 4 bytes; identity has 1.0 on diagonal.
	// Column 0: [1,0,0,0] -> bytes 0..15; float32(1.0) at byte 0
	// Column 1: [0,1,0,0] -> bytes 16..31; float32(1.0) at byte 20
	// Column 2: [0,0,1,0] -> bytes 32..47; float32(1.0) at byte 40
	// Column 3: [0,0,0,1] -> bytes 48..63; float32(1.0) at byte 60
	diagOffsets := []int{0, 20, 40, 60}
	for _, off := range diagOffsets {
		// IEEE 754 float32(1.0) = 0x3F800000 -> LE: 0x00, 0x00, 0x80, 0x3F
		if data[off] != 0x00 || data[off+1] != 0x00 || data[off+2] != 0x80 || data[off+3] != 0x3F {
			t.Fatalf("identity diagonal at byte %d: got [%02x %02x %02x %02x], want [00 00 80 3f]",
				off, data[off], data[off+1], data[off+2], data[off+3])
		}
	}
}

package main

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	pkgmath "github.com/muddl/go_toy_renderer/pkg/math"
)

func TestFramebufferToRGBA_BlackFrame(t *testing.T) {
	fb := framebuffer.New(2, 2)
	got := framebufferToRGBA(fb)

	if len(got) != 2*2*4 {
		t.Fatalf("len = %d, want %d", len(got), 2*2*4)
	}
	for i, b := range got {
		if i%4 == 3 {
			if b != 255 {
				t.Errorf("byte[%d] (alpha) = %d, want 255", i, b)
			}
		} else {
			if b != 0 {
				t.Errorf("byte[%d] = %d, want 0", i, b)
			}
		}
	}
}

func TestFramebufferToRGBA_ColourConversion(t *testing.T) {
	fb := framebuffer.New(1, 1)
	fb.ColorBuffer[0] = pkgmath.Vec3{X: 1.0, Y: 0.5, Z: 0.0}

	got := framebufferToRGBA(fb)

	if len(got) != 4 {
		t.Fatalf("len = %d, want 4", len(got))
	}
	if got[0] != 255 {
		t.Errorf("R = %d, want 255", got[0])
	}
	if got[1] != 127 {
		t.Errorf("G = %d, want 127", got[1])
	}
	if got[2] != 0 {
		t.Errorf("B = %d, want 0", got[2])
	}
	if got[3] != 255 {
		t.Errorf("A = %d, want 255", got[3])
	}
}

func TestFramebufferToRGBA_Clamping(t *testing.T) {
	fb := framebuffer.New(1, 1)
	fb.ColorBuffer[0] = pkgmath.Vec3{X: 1.5, Y: -0.5, Z: 0.5}

	got := framebufferToRGBA(fb)

	if got[0] != 255 {
		t.Errorf("R (>1.0) = %d, want 255", got[0])
	}
	if got[1] != 0 {
		t.Errorf("G (<0.0) = %d, want 0", got[1])
	}
	if got[2] != 127 {
		t.Errorf("B (0.5) = %d, want 127", got[2])
	}
}

func TestFramebufferToRGBA_RowOrder(t *testing.T) {
	// Verify row-major order: pixel (x,y) is at offset (y*width+x)*4.
	fb := framebuffer.New(2, 2)
	fb.ColorBuffer[0] = pkgmath.Vec3{X: 1, Y: 0, Z: 0} // (0,0) red
	fb.ColorBuffer[3] = pkgmath.Vec3{X: 0, Y: 0, Z: 1} // (1,1) blue

	got := framebufferToRGBA(fb)

	if got[0] != 255 || got[1] != 0 || got[2] != 0 {
		t.Errorf("pixel (0,0): got R=%d G=%d B=%d, want 255 0 0", got[0], got[1], got[2])
	}
	if got[12] != 0 || got[13] != 0 || got[14] != 255 {
		t.Errorf("pixel (1,1): got R=%d G=%d B=%d, want 0 0 255", got[12], got[13], got[14])
	}
}

//go:build !headless

package overlay_test

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	mathpkg "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/overlay"
)

const (
	testFBWidth  = 640
	testFBHeight = 480
)

// grayBackground is a mid-gray used to initialise the test framebuffer so that
// composited pixels are clearly distinguishable from the original color.
var grayBackground = mathpkg.Vec3{X: 0.5, Y: 0.5, Z: 0.5}

// makeGrayFB creates a framebuffer of size w×h and fills it with grayBackground.
func makeGrayFB(w, h int) *framebuffer.Framebuffer {
	fb := framebuffer.New(w, h)
	fb.Clear(grayBackground, 1.0)
	return fb
}

// TestComposeIntoFramebuffer_OverlayPixels_WriteAtTopLeft verifies that
// compositing a rendered TextLayer writes changed pixels in the top-left
// region of the framebuffer where the overlay text sits.
func TestComposeIntoFramebuffer_OverlayPixels_WriteAtTopLeft(t *testing.T) {
	fb := makeGrayFB(testFBWidth, testFBHeight)
	tl := overlay.NewTextLayer(testFBWidth, testFBHeight)
	tl.Render(overlay.LevelFPS, overlay.Metrics{FPS: 60.0})

	overlay.ComposeIntoFramebuffer(tl, fb)

	// At least one pixel in the overlay region should no longer be mid-gray.
	changed := false
	for y := 0; y < tl.Height(); y++ {
		for x := 0; x < tl.Width(); x++ {
			p := fb.GetPixel(x, y)
			if p.X != grayBackground.X || p.Y != grayBackground.Y || p.Z != grayBackground.Z {
				changed = true
				break
			}
		}
		if changed {
			break
		}
	}
	if !changed {
		t.Error("ComposeIntoFramebuffer: no framebuffer pixels were changed — expected overlay text to be composited")
	}
}

// TestComposeIntoFramebuffer_OpaqueTextPixel_IsGreen verifies that a fully
// opaque (A=255) green text pixel replaces the framebuffer color with pure
// green (R=0, G=1, B=0) regardless of the background.
func TestComposeIntoFramebuffer_OpaqueTextPixel_IsGreen(t *testing.T) {
	fb := makeGrayFB(testFBWidth, testFBHeight)
	tl := overlay.NewTextLayer(testFBWidth, testFBHeight)
	pixels := tl.Render(overlay.LevelFPS, overlay.Metrics{FPS: 60.0})

	overlay.ComposeIntoFramebuffer(tl, fb)

	// Find the first fully-opaque (A==255) pixel in the TextLayer and confirm
	// the corresponding framebuffer pixel became pure green.
	for i := 3; i < len(pixels); i += 4 {
		if pixels[i] == 0xFF {
			pixIdx := i / 4
			x := pixIdx % tl.Width()
			y := pixIdx / tl.Width()
			got := fb.GetPixel(x, y)
			const eps = 1.0 / 255.0
			if got.X > eps || got.Z > eps {
				t.Errorf("opaque text pixel (%d,%d): R=%.4f B=%.4f want ~0 (green text)", x, y, got.X, got.Z)
			}
			if got.Y < 1.0-eps {
				t.Errorf("opaque text pixel (%d,%d): G=%.4f want ~1.0", x, y, got.Y)
			}
			return
		}
	}
	t.Error("no fully-opaque pixel found in TextLayer — check Render(LevelFPS)")
}

// TestComposeIntoFramebuffer_OutsideOverlay_Unchanged verifies that pixels
// beyond the TextLayer dimensions are not modified by compositing.
func TestComposeIntoFramebuffer_OutsideOverlay_Unchanged(t *testing.T) {
	fb := makeGrayFB(testFBWidth, testFBHeight)
	tl := overlay.NewTextLayer(testFBWidth, testFBHeight)
	tl.Render(overlay.LevelFPS, overlay.Metrics{FPS: 60.0})

	overlay.ComposeIntoFramebuffer(tl, fb)

	// Pixels outside the TextLayer bounds must remain untouched.
	checkX := tl.Width()  // first column beyond the overlay
	checkY := tl.Height() // first row beyond the overlay
	if checkX >= testFBWidth || checkY >= testFBHeight {
		t.Skip("TextLayer fills entire framebuffer — outside-overlay check not applicable")
	}
	p := fb.GetPixel(checkX, checkY)
	if p.X != grayBackground.X || p.Y != grayBackground.Y || p.Z != grayBackground.Z {
		t.Errorf("pixel (%d,%d) outside overlay changed: got (%.3f,%.3f,%.3f), want gray",
			checkX, checkY, p.X, p.Y, p.Z)
	}
}

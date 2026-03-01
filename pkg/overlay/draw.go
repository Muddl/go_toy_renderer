//go:build !headless

package overlay

import (
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	mathpkg "github.com/muddl/go_toy_renderer/pkg/math"
)

// ComposeIntoFramebuffer alpha-blends the TextLayer pixel buffer over the
// framebuffer's ColorBuffer. The TextLayer is placed at the top-left corner
// (0, 0) of the framebuffer; only pixels within the layer's bounds are
// touched. Pixels with alpha == 0 are skipped entirely.
func ComposeIntoFramebuffer(layer *TextLayer, fb *framebuffer.Framebuffer) {
	w := layer.Width()
	h := layer.Height()
	pixels := layer.Pixels()

	for y := 0; y < h && y < fb.Height; y++ {
		for x := 0; x < w && x < fb.Width; x++ {
			idx := (y*w + x) * 4
			srcA := float64(pixels[idx+3]) / 255.0
			if srcA == 0 {
				continue
			}
			srcR := float64(pixels[idx+0]) / 255.0
			srcG := float64(pixels[idx+1]) / 255.0
			srcB := float64(pixels[idx+2]) / 255.0

			dst := fb.GetPixel(x, y)
			invA := 1.0 - srcA
			fb.ColorBuffer[y*fb.Width+x] = mathpkg.Vec3{
				X: srcA*srcR + invA*dst.X,
				Y: srcA*srcG + invA*dst.Y,
				Z: srcA*srcB + invA*dst.Z,
			}
		}
	}
}

package renderer

import (
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
)

// framebufferToRGBA converts a Framebuffer's float64 color buffer to a packed
// RGBA byte slice suitable for upload to an OpenGL texture.
// Each color component is clamped to [0,1] then mapped to [0,255].
// Alpha is always 255. Layout: R,G,B,A per pixel, row-major (y*width+x).
func framebufferToRGBA(fb *framebuffer.Framebuffer) []byte {
	out := make([]byte, fb.Width*fb.Height*4)
	for i, c := range fb.ColorBuffer {
		out[i*4+0] = clampF64ToByte(c.X)
		out[i*4+1] = clampF64ToByte(c.Y)
		out[i*4+2] = clampF64ToByte(c.Z)
		out[i*4+3] = 255
	}
	return out
}

// clampF64ToByte maps a float64 in [0,1] to uint8 in [0,255] with clamping.
func clampF64ToByte(v float64) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 255
	}
	return uint8(v * 255)
}

// Package framebuffer provides a 2D pixel store with color and depth buffers
// for use as the output stage of a software rendering pipeline.
package framebuffer

import (
	"image"
	"image/color"
	"image/png"
	"os"

	math "github.com/muddl/go_toy_renderer/pkg/math"
)

// Framebuffer holds per-pixel color (RGB float64 in [0,1]) and depth values.
// Buffer layout: index = y*Width + x (top-left origin, +X right, +Y down).
type Framebuffer struct {
	Width       int
	Height      int
	ColorBuffer []math.Vec3
	DepthBuffer []float64
}

// New creates a Framebuffer of the given dimensions.
// All color pixels are initialized to black (0,0,0) and all depths to 1.0 (far plane).
func New(width, height int) *Framebuffer {
	size := width * height
	colorBuffer := make([]math.Vec3, size)
	depthBuffer := make([]float64, size)
	for i := range depthBuffer {
		depthBuffer[i] = 1.0
	}
	return &Framebuffer{
		Width:       width,
		Height:      height,
		ColorBuffer: colorBuffer,
		DepthBuffer: depthBuffer,
	}
}

// Clear resets all pixels to the given color and all depth values to depth.
func (fb *Framebuffer) Clear(c math.Vec3, depth float64) {
	for i := range fb.ColorBuffer {
		fb.ColorBuffer[i] = c
		fb.DepthBuffer[i] = depth
	}
}

// SetPixel writes color and depth at (x, y) only if the new depth is strictly
// less than the current depth (depth test: closer overwrites farther).
// Out-of-bounds coordinates are silently ignored.
func (fb *Framebuffer) SetPixel(x, y int, c math.Vec3, depth float64) {
	if !fb.inBounds(x, y) {
		return
	}
	idx := y*fb.Width + x
	if depth < fb.DepthBuffer[idx] {
		fb.ColorBuffer[idx] = c
		fb.DepthBuffer[idx] = depth
	}
}

// GetPixel returns the color at (x, y). Returns a zero Vec3 for out-of-bounds coordinates.
func (fb *Framebuffer) GetPixel(x, y int) math.Vec3 {
	if !fb.inBounds(x, y) {
		return math.Vec3{}
	}
	return fb.ColorBuffer[y*fb.Width+x]
}

// GetDepth returns the depth at (x, y). Returns 1.0 for out-of-bounds coordinates.
func (fb *Framebuffer) GetDepth(x, y int) float64 {
	if !fb.inBounds(x, y) {
		return 1.0
	}
	return fb.DepthBuffer[y*fb.Width+x]
}

// SavePNG writes the framebuffer as a PNG image to the given file path.
// Float RGB values in [0,1] are converted to 8-bit (0-255) with clamping.
func (fb *Framebuffer) SavePNG(filename string) error {
	img := image.NewNRGBA(image.Rect(0, 0, fb.Width, fb.Height))
	for y := 0; y < fb.Height; y++ {
		for x := 0; x < fb.Width; x++ {
			c := fb.ColorBuffer[y*fb.Width+x]
			img.SetNRGBA(x, y, color.NRGBA{
				R: clampToByte(c.X),
				G: clampToByte(c.Y),
				B: clampToByte(c.Z),
				A: 255,
			})
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

// inBounds reports whether (x, y) is within the framebuffer dimensions.
func (fb *Framebuffer) inBounds(x, y int) bool {
	return x >= 0 && x < fb.Width && y >= 0 && y < fb.Height
}

// clampToByte converts a float64 color component in [0,1] to a uint8 in [0,255].
// Values outside [0,1] are clamped.
func clampToByte(v float64) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 255
	}
	return uint8(v * 255)
}

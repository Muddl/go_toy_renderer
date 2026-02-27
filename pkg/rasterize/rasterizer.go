// Package rasterize converts screen-space triangles into pixels stored in a framebuffer.
// It uses barycentric coordinates for coverage testing and attribute interpolation.
package rasterize

import (
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	math "github.com/muddl/go_toy_renderer/pkg/math"
)

// ScreenVertex represents a vertex in 2D screen space with depth and color attributes.
// X and Y are pixel coordinates; pixel (ix, iy) has its center at (float64(ix)+0.5, float64(iy)+0.5).
// Z is the depth value for depth testing (0=closest, 1=farthest).
// Color holds RGB components, each in [0, 1].
type ScreenVertex struct {
	X, Y  float64
	Z     float64
	Color math.Vec3
}

// Triangle fills the pixels of fb covered by the triangle defined by v0, v1, v2.
// Screen coordinates are in pixels with origin at top-left, +X right, +Y down.
// Depth and color are linearly interpolated using barycentric coordinates.
// Degenerate triangles (zero area) are silently skipped.
// Out-of-bounds pixels are silently ignored via fb.SetPixel.
func Triangle(v0, v1, v2 ScreenVertex, fb *framebuffer.Framebuffer) {
	// Compute signed area (× 2) of the triangle.
	// Positive = CCW winding, negative = CW winding (screen space, Y-down).
	area := edgeFunction(v0.X, v0.Y, v1.X, v1.Y, v2.X, v2.Y)

	// Skip degenerate (zero-area) triangles; use area² to avoid math.Abs import.
	if area*area < 1e-16 {
		return
	}

	invArea := 1.0 / area

	// Compute tight bounding box around the triangle.
	minX := int(minF(v0.X, minF(v1.X, v2.X)))
	minY := int(minF(v0.Y, minF(v1.Y, v2.Y)))
	maxX := int(maxF(v0.X, maxF(v1.X, v2.X))) + 1
	maxY := int(maxF(v0.Y, maxF(v1.Y, v2.Y))) + 1

	// Clamp bounding box to framebuffer dimensions.
	if minX < 0 {
		minX = 0
	}
	if minY < 0 {
		minY = 0
	}
	if maxX > fb.Width {
		maxX = fb.Width
	}
	if maxY > fb.Height {
		maxY = fb.Height
	}

	for iy := minY; iy < maxY; iy++ {
		for ix := minX; ix < maxX; ix++ {
			// Sample at pixel center.
			px := float64(ix) + 0.5
			py := float64(iy) + 0.5

			// Compute barycentric weights for each vertex.
			w0 := edgeFunction(v1.X, v1.Y, v2.X, v2.Y, px, py)
			w1 := edgeFunction(v2.X, v2.Y, v0.X, v0.Y, px, py)
			w2 := edgeFunction(v0.X, v0.Y, v1.X, v1.Y, px, py)

			if !insideTriangle(area, w0, w1, w2) {
				continue
			}

			// Normalize barycentric coordinates.
			b0 := w0 * invArea
			b1 := w1 * invArea
			b2 := w2 * invArea

			// Interpolate depth and color across the triangle.
			depth := b0*v0.Z + b1*v1.Z + b2*v2.Z
			color := math.Vec3{
				X: b0*v0.Color.X + b1*v1.Color.X + b2*v2.Color.X,
				Y: b0*v0.Color.Y + b1*v1.Color.Y + b2*v2.Color.Y,
				Z: b0*v0.Color.Z + b1*v1.Color.Z + b2*v2.Color.Z,
			}

			fb.SetPixel(ix, iy, color, depth)
		}
	}
}

// insideTriangle reports whether the barycentric weights indicate a point lies inside
// the triangle, supporting both CCW (area > 0) and CW (area < 0) winding.
func insideTriangle(area, w0, w1, w2 float64) bool {
	if area > 0 {
		return w0 >= 0 && w1 >= 0 && w2 >= 0
	}
	return w0 <= 0 && w1 <= 0 && w2 <= 0
}

// edgeFunction computes the signed 2D cross product of edge a→b with vector a→p.
// The result is positive when p is to the left of (or on) the directed edge a→b.
func edgeFunction(ax, ay, bx, by, px, py float64) float64 {
	return (bx-ax)*(py-ay) - (by-ay)*(px-ax)
}

// minF returns the smaller of two float64 values.
func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// maxF returns the larger of two float64 values.
func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

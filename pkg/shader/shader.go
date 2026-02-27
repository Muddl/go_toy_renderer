// Package shader provides per-pixel color calculation functions for the rendering pipeline.
//
// A shader is a function that receives interpolated per-fragment attributes and returns
// the final RGB color for that pixel. The simplest shaders pass color through unchanged;
// more complex ones can apply lighting, texture lookups, or other effects.
//
// For MVP, three shaders are provided:
//   - VertexColor: pass-through interpolated vertex color (demonstrates interpolation)
//   - NewFlatColor: constant color regardless of attributes (useful for debugging)
//   - Depth: visualize interpolated depth as grayscale (useful for debugging)
package shader

import "github.com/muddl/go_toy_renderer/pkg/math"

// Attributes holds interpolated per-fragment data passed to a Func.
// Values are set by the rasterizer via barycentric interpolation.
type Attributes struct {
	// Color is the interpolated vertex color, components in [0, 1].
	Color math.Vec3
	// Depth is the interpolated fragment depth, in [0, 1].
	Depth float64
	// Future fields: Position math.Vec3, Normal math.Vec3, UV math.Vec3
}

// Func computes the final RGB color for a fragment given its interpolated attributes.
// It is called once per visible pixel inside a rasterized triangle.
type Func func(Attributes) math.Vec3

// VertexColor returns the interpolated vertex color unchanged.
// Use this to verify that color interpolation across triangles is working correctly.
func VertexColor(attr Attributes) math.Vec3 {
	return attr.Color
}

// NewFlatColor returns a Func that always outputs the given constant color,
// regardless of the fragment's interpolated attributes.
// Useful for solid-color fills and debugging the rendering pipeline.
func NewFlatColor(color math.Vec3) Func {
	return func(_ Attributes) math.Vec3 {
		return color
	}
}

// Depth converts the fragment's interpolated depth to a grayscale color.
// A depth of 0 (near plane) maps to black; a depth of 1 (far plane) maps to white.
// Use this to visualize the depth buffer and verify depth interpolation.
func Depth(attr Attributes) math.Vec3 {
	g := attr.Depth
	return math.Vec3{X: g, Y: g, Z: g}
}

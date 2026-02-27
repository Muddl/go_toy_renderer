package rasterize

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	math "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/shader"
)

// BenchmarkTriangleShaded_Small measures rasterization of a small triangle
// (roughly 10×10 pixel bounding box).
func BenchmarkTriangleShaded_Small(b *testing.B) {
	fb := framebuffer.New(100, 100)
	v0 := ScreenVertex{X: 40, Y: 40, Z: 0.5, Color: math.Vec3{X: 1, Y: 0, Z: 0}}
	v1 := ScreenVertex{X: 60, Y: 40, Z: 0.5, Color: math.Vec3{X: 0, Y: 1, Z: 0}}
	v2 := ScreenVertex{X: 50, Y: 55, Z: 0.5, Color: math.Vec3{X: 0, Y: 0, Z: 1}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TriangleShaded(v0, v1, v2, shader.VertexColor, fb)
	}
}

// BenchmarkTriangleShaded_Large measures rasterization of a triangle that covers
// most of the framebuffer, stressing the inner loop.
func BenchmarkTriangleShaded_Large(b *testing.B) {
	fb := framebuffer.New(640, 480)
	v0 := ScreenVertex{X: 0, Y: 480, Z: 0.5, Color: math.Vec3{X: 1, Y: 0, Z: 0}}
	v1 := ScreenVertex{X: 640, Y: 480, Z: 0.5, Color: math.Vec3{X: 0, Y: 1, Z: 0}}
	v2 := ScreenVertex{X: 320, Y: 0, Z: 0.5, Color: math.Vec3{X: 0, Y: 0, Z: 1}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TriangleShaded(v0, v1, v2, shader.VertexColor, fb)
	}
}

// BenchmarkEdgeFunction measures the raw cost of the edge function calculation.
func BenchmarkEdgeFunction(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = edgeFunction(10.5, 20.5, 30.5, 40.5, 15.5, 25.5)
	}
}

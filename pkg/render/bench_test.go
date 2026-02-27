package render

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/camera"
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
	math "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/shader"
)

// BenchmarkRender_Cube_640x480 measures the full render pipeline cost for
// a cube at 640×480 resolution — the standard MVP output resolution.
func BenchmarkRender_Cube_640x480(b *testing.B) {
	scene := NewScene()
	scene.AddMesh(geometry.NewCube())
	cam := camera.New(
		math.Vec3{X: 3, Y: 2, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0, float64(640)/float64(480), 0.1, 100.0,
	)
	fb := framebuffer.New(640, 480)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Render(scene, cam, fb, shader.VertexColor)
	}
}

// BenchmarkRender_Triangle_100x100 measures the render pipeline cost for
// a single triangle in a small framebuffer.
func BenchmarkRender_Triangle_100x100(b *testing.B) {
	scene := goldenTriangleScene()
	cam := goldenCamera()
	fb := framebuffer.New(100, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Render(scene, cam, fb, shader.VertexColor)
	}
}

// BenchmarkTransformVertex measures the per-vertex MVP transform cost.
func BenchmarkTransformVertex(b *testing.B) {
	cam := goldenCamera()
	mvp := cam.ViewProjectionMatrix()
	v := geometry.NewVertex(math.Vec3{X: 0.5, Y: 0.5, Z: 0}, math.Vec3{X: 1, Y: 0, Z: 0})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = transformVertex(v, mvp, 640, 480)
	}
}

// Package render orchestrates the 3D rendering pipeline.
// It transforms mesh geometry from world space through camera space to screen space,
// rasterizes each triangle with barycentric interpolation, applies a per-pixel shader,
// and writes results to a framebuffer.
//
// Usage:
//
//	scene := render.NewScene()
//	scene.AddMesh(geometry.NewCube())
//	cam := camera.New(...)
//	fb := framebuffer.New(640, 480)
//	render.Render(scene, cam, fb, shader.VertexColor)
//	fb.SavePNG("output.png")
package render

import (
	"github.com/muddl/go_toy_renderer/pkg/camera"
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
	math "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/rasterize"
	"github.com/muddl/go_toy_renderer/pkg/shader"
)

// Scene holds a collection of meshes to be rendered together.
// For MVP all meshes use an identity model matrix (world space equals local space).
type Scene struct {
	Meshes []*geometry.Mesh
}

// NewScene creates an empty scene with no meshes.
func NewScene() *Scene {
	return &Scene{Meshes: make([]*geometry.Mesh, 0)}
}

// AddMesh appends m to the scene's mesh list.
func (s *Scene) AddMesh(m *geometry.Mesh) {
	s.Meshes = append(s.Meshes, m)
}

// Render clears fb, then transforms and rasterizes every mesh in scene using
// the given camera and shaderFn.
//
// The framebuffer is cleared to black (depth 1.0) before rendering begins.
// Triangles where any vertex has clip-space w ≤ 0 (behind the camera plane) are
// silently skipped. Degenerate triangles are silently skipped by the rasterizer.
//
// For MVP, an identity model matrix is assumed for all meshes.
func Render(scene *Scene, cam camera.Camera, fb *framebuffer.Framebuffer, shaderFn shader.Func) {
	fb.Clear(math.Vec3{}, 1.0)

	mvp := cam.ViewProjectionMatrix()
	w := fb.Width
	h := fb.Height

	for _, mesh := range scene.Meshes {
		// Transform every vertex to screen space once per mesh.
		sverts := make([]rasterize.ScreenVertex, len(mesh.Vertices))
		valid := make([]bool, len(mesh.Vertices))

		for i, v := range mesh.Vertices {
			sv, ok := transformVertex(v, mvp, w, h)
			sverts[i] = sv
			valid[i] = ok
		}

		// Rasterize each triangle whose vertices are all in front of the camera.
		for t := 0; t < mesh.TriangleCount(); t++ {
			i0, i1, i2 := mesh.GetTriangle(t)
			if !valid[i0] || !valid[i1] || !valid[i2] {
				continue
			}

			rasterize.TriangleShaded(sverts[i0], sverts[i1], sverts[i2], shaderFn, fb)
		}
	}
}

// transformVertex converts v from world space to a rasterize.ScreenVertex in screen space
// using the given MVP matrix and framebuffer dimensions.
// Returns (sv, true) on success, or (zero, false) if the vertex is on or behind the
// camera plane (clip-space w ≤ 0), which would produce invalid perspective divide results.
func transformVertex(v geometry.Vertex, mvp math.Mat4x4, width, height int) (rasterize.ScreenVertex, bool) {
	cx, cy, cz, cw := mvp.MultiplyVec4(v.Position.X, v.Position.Y, v.Position.Z, 1.0)
	if cw <= 0 {
		return rasterize.ScreenVertex{}, false
	}

	// Perspective divide: clip space → NDC (Normalized Device Coordinates).
	// NDC x,y ∈ [-1, 1]; NDC z ∈ [-1, 1] with -1 = near plane, +1 = far plane.
	ndcX := cx / cw
	ndcY := cy / cw
	ndcZ := cz / cw

	// Map NDC z [-1, 1] → depth buffer [0, 1].
	// Near plane (ndcZ = -1) → depth 0 (closest); far plane (ndcZ = +1) → depth 1.
	depth := (ndcZ + 1.0) / 2.0

	// Viewport transform: NDC → screen pixels.
	// NDC +Y is up; screen +Y is down, so Y is flipped.
	screenX := (ndcX + 1.0) / 2.0 * float64(width)
	screenY := (1.0 - ndcY) / 2.0 * float64(height)

	return rasterize.ScreenVertex{
		X:     screenX,
		Y:     screenY,
		Z:     depth,
		Color: v.Color,
	}, true
}

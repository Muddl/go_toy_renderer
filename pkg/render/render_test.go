package render

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/camera"
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
	math "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/shader"
)

// testCamera returns a camera at (0,0,5) looking at the origin.
// The frustum covers the region around the origin at z=0.
func testCamera() camera.Camera {
	return camera.New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		60.0, 1.0, 0.1, 100.0,
	)
}

// newFacingTriangle returns a mesh with one CCW triangle centered at origin,
// facing the camera that is on the +Z axis.
func newFacingTriangle() *geometry.Mesh {
	m := geometry.NewMesh()
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	m.AddVertex(geometry.NewVertex(math.Vec3{X: -1, Y: -1, Z: 0}, red))
	m.AddVertex(geometry.NewVertex(math.Vec3{X: 1, Y: -1, Z: 0}, red))
	m.AddVertex(geometry.NewVertex(math.Vec3{X: 0, Y: 1, Z: 0}, red))
	m.AddTriangle(0, 1, 2)

	return m
}

// countNonBackground counts pixels that differ from the background (black).
func countNonBackground(fb *framebuffer.Framebuffer) int {
	count := 0
	black := math.Vec3{}
	for _, c := range fb.ColorBuffer {
		if c != black {
			count++
		}
	}

	return count
}

// --- Scene tests ---

func TestScene_NewScene_IsEmpty(t *testing.T) {
	s := NewScene()
	if len(s.Meshes) != 0 {
		t.Errorf("NewScene Meshes length = %d, want 0", len(s.Meshes))
	}
}

func TestScene_AddMesh_IncreasesCount(t *testing.T) {
	s := NewScene()
	s.AddMesh(geometry.NewMesh())
	s.AddMesh(geometry.NewMesh())
	if len(s.Meshes) != 2 {
		t.Errorf("after AddMesh×2: len(Meshes) = %d, want 2", len(s.Meshes))
	}
}

func TestScene_AddMesh_StoresMesh(t *testing.T) {
	s := NewScene()
	m := geometry.NewMesh()
	s.AddMesh(m)
	if s.Meshes[0] != m {
		t.Error("AddMesh did not store the given mesh pointer")
	}
}

// --- Render tests ---

func TestRender_EmptyScene_DoesNotPanic(t *testing.T) {
	fb := framebuffer.New(10, 10)
	cam := testCamera()
	scene := NewScene()
	// Must not panic.
	Render(scene, cam, fb, shader.VertexColor)
}

func TestRender_EmptyScene_FramebufferCleared(t *testing.T) {
	fb := framebuffer.New(10, 10)
	// Pre-fill with non-zero color to ensure Render clears it.
	for i := range fb.ColorBuffer {
		fb.ColorBuffer[i] = math.Vec3{X: 1, Y: 1, Z: 1}
	}

	cam := testCamera()
	scene := NewScene()
	Render(scene, cam, fb, shader.VertexColor)

	// Render should have cleared the framebuffer first; all pixels should be black.
	if got := countNonBackground(fb); got != 0 {
		t.Errorf("empty scene: %d non-black pixels, want 0", got)
	}
}

func TestRender_VisibleTriangle_RendersNonBlankPixels(t *testing.T) {
	fb := framebuffer.New(100, 100)
	cam := testCamera()
	scene := NewScene()
	scene.AddMesh(newFacingTriangle())

	Render(scene, cam, fb, shader.VertexColor)

	if got := countNonBackground(fb); got == 0 {
		t.Error("visible triangle rendered 0 pixels, expected > 0")
	}
}

func TestRender_FlatColorShader_AllRenderedPixelsSameColor(t *testing.T) {
	fb := framebuffer.New(100, 100)
	cam := testCamera()
	scene := NewScene()
	scene.AddMesh(newFacingTriangle())

	blue := math.Vec3{X: 0, Y: 0, Z: 1}
	Render(scene, cam, fb, shader.NewFlatColor(blue))

	for y := 0; y < fb.Height; y++ {
		for x := 0; x < fb.Width; x++ {
			c := fb.GetPixel(x, y)
			if c == (math.Vec3{}) {
				continue // background pixel
			}
			if c.X != blue.X || c.Y != blue.Y || c.Z != blue.Z {
				t.Errorf("pixel (%d,%d) = %v, want flat blue %v", x, y, c, blue)
				return
			}
		}
	}
}

func TestRender_DepthOrdering_CloserTriangleWins(t *testing.T) {
	// Two triangles at the same screen position, different depths.
	// The closer one (z=1, closer to camera at z=5) should win.
	fb := framebuffer.New(100, 100)
	cam := testCamera()
	scene := NewScene()

	farMesh := geometry.NewMesh()
	farColor := math.Vec3{X: 1, Y: 0, Z: 0} // red = far
	farMesh.AddVertex(geometry.NewVertex(math.Vec3{X: -2, Y: -2, Z: 0}, farColor))
	farMesh.AddVertex(geometry.NewVertex(math.Vec3{X: 2, Y: -2, Z: 0}, farColor))
	farMesh.AddVertex(geometry.NewVertex(math.Vec3{X: 0, Y: 2, Z: 0}, farColor))
	farMesh.AddTriangle(0, 1, 2)

	nearMesh := geometry.NewMesh()
	nearColor := math.Vec3{X: 0, Y: 0, Z: 1} // blue = near
	nearMesh.AddVertex(geometry.NewVertex(math.Vec3{X: -2, Y: -2, Z: 1}, nearColor))
	nearMesh.AddVertex(geometry.NewVertex(math.Vec3{X: 2, Y: -2, Z: 1}, nearColor))
	nearMesh.AddVertex(geometry.NewVertex(math.Vec3{X: 0, Y: 2, Z: 1}, nearColor))
	nearMesh.AddTriangle(0, 1, 2)

	// Add far mesh first; near mesh should overwrite it.
	scene.AddMesh(farMesh)
	scene.AddMesh(nearMesh)

	Render(scene, cam, fb, shader.VertexColor)

	// The center pixel should be blue (near), not red (far).
	centerX, centerY := fb.Width/2, fb.Height/2
	c := fb.GetPixel(centerX, centerY)
	if c == (math.Vec3{}) {
		t.Fatalf("center pixel (%d,%d) is black; expected colored pixel", centerX, centerY)
	}
	if c.Z < 0.5 {
		t.Errorf("center pixel = %v, expected blue (near mesh); got more red than blue", c)
	}
}

func TestRender_BehindCameraTriangle_IsSkipped(t *testing.T) {
	// Triangle entirely behind the camera (z > 5 while camera is at z=5 looking at origin).
	fb := framebuffer.New(100, 100)
	cam := testCamera()
	scene := NewScene()

	behindMesh := geometry.NewMesh()
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	// Vertices behind the camera (z=10 when camera is at z=5 looking toward -Z)
	behindMesh.AddVertex(geometry.NewVertex(math.Vec3{X: -1, Y: -1, Z: 10}, red))
	behindMesh.AddVertex(geometry.NewVertex(math.Vec3{X: 1, Y: -1, Z: 10}, red))
	behindMesh.AddVertex(geometry.NewVertex(math.Vec3{X: 0, Y: 1, Z: 10}, red))
	behindMesh.AddTriangle(0, 1, 2)
	scene.AddMesh(behindMesh)

	Render(scene, cam, fb, shader.VertexColor)

	if got := countNonBackground(fb); got != 0 {
		t.Errorf("behind-camera triangle set %d pixels, want 0", got)
	}
}

// --- transformVertex tests ---

func TestTransformVertex_OriginMapsToFramebufferCenter(t *testing.T) {
	// With camera at (0,0,5) looking at origin, a vertex at world origin
	// should project to approximately the center of the framebuffer.
	cam := testCamera()
	mvp := cam.ViewProjectionMatrix()
	w, h := 100, 100

	v := geometry.NewVertex(math.Vec3{X: 0, Y: 0, Z: 0}, math.Vec3{X: 1, Y: 0, Z: 0})
	sv, ok := transformVertex(v, mvp, w, h)

	if !ok {
		t.Fatal("transformVertex returned ok=false for vertex at origin (should be in front of camera)")
	}

	centerX := float64(w) / 2.0
	centerY := float64(h) / 2.0
	eps := 2.0 // allow 2-pixel tolerance

	if sv.X < centerX-eps || sv.X > centerX+eps {
		t.Errorf("screen X = %.2f, want ~%.2f", sv.X, centerX)
	}
	if sv.Y < centerY-eps || sv.Y > centerY+eps {
		t.Errorf("screen Y = %.2f, want ~%.2f", sv.Y, centerY)
	}
}

func TestTransformVertex_BehindCamera_ReturnsFalse(t *testing.T) {
	// A vertex behind the camera should return ok=false.
	cam := testCamera() // camera at z=5 looking toward -Z
	mvp := cam.ViewProjectionMatrix()

	// Place vertex at z=10, which is behind the camera.
	v := geometry.NewVertex(math.Vec3{X: 0, Y: 0, Z: 10}, math.Vec3{})
	_, ok := transformVertex(v, mvp, 100, 100)

	if ok {
		t.Error("transformVertex returned ok=true for vertex behind camera, want false")
	}
}

func TestTransformVertex_DepthInRange(t *testing.T) {
	// Vertex at world origin with default camera should have depth in [0, 1].
	cam := testCamera()
	mvp := cam.ViewProjectionMatrix()

	v := geometry.NewVertex(math.Vec3{X: 0, Y: 0, Z: 0}, math.Vec3{})
	sv, ok := transformVertex(v, mvp, 100, 100)

	if !ok {
		t.Fatal("transformVertex returned ok=false unexpectedly")
	}
	if sv.Z < 0 || sv.Z > 1 {
		t.Errorf("depth = %.4f, want in [0, 1]", sv.Z)
	}
}

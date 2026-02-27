package render

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/camera"
	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	"github.com/muddl/go_toy_renderer/pkg/geometry"
	math "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/shader"
)

// update is set via -update flag to regenerate golden image files.
var update = flag.Bool("update", false, "update golden image files")

// goldenTriangleScene returns a simple deterministic RGB triangle centered at origin.
func goldenTriangleScene() *Scene {
	scene := NewScene()
	m := geometry.NewMesh()
	// Red at bottom-left, green at bottom-right, blue at top.
	m.AddVertex(geometry.NewVertex(math.Vec3{X: -1, Y: -1, Z: 0}, math.Vec3{X: 1, Y: 0, Z: 0}))
	m.AddVertex(geometry.NewVertex(math.Vec3{X: 1, Y: -1, Z: 0}, math.Vec3{X: 0, Y: 1, Z: 0}))
	m.AddVertex(geometry.NewVertex(math.Vec3{X: 0, Y: 1, Z: 0}, math.Vec3{X: 0, Y: 0, Z: 1}))
	m.AddTriangle(0, 1, 2)
	scene.AddMesh(m)
	return scene
}

// goldenCamera returns a fixed camera looking at the origin from the +Z axis.
func goldenCamera() camera.Camera {
	return camera.New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		60.0, 1.0, 0.1, 100.0,
	)
}

// fbToImage converts a framebuffer to an *image.NRGBA for PNG comparison.
func fbToImage(fb *framebuffer.Framebuffer) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, fb.Width, fb.Height))
	for y := 0; y < fb.Height; y++ {
		for x := 0; x < fb.Width; x++ {
			c := fb.GetPixel(x, y)
			clamp := func(v float64) uint8 {
				if v <= 0 {
					return 0
				}
				if v >= 1 {
					return 255
				}
				return uint8(v * 255)
			}
			img.SetNRGBA(x, y, color.NRGBA{R: clamp(c.X), G: clamp(c.Y), B: clamp(c.Z), A: 255})
		}
	}
	return img
}

// TestRender_GoldenImage_Triangle renders a deterministic RGB triangle and compares
// pixel data against a stored reference image in testdata/.
// Run with -update to regenerate the reference.
func TestRender_GoldenImage_Triangle(t *testing.T) {
	const (
		width      = 100
		height     = 100
		goldenPath = "testdata/golden_triangle.png"
	)

	fb := framebuffer.New(width, height)
	Render(goldenTriangleScene(), goldenCamera(), fb, shader.VertexColor)

	got := fbToImage(fb)
	var gotBuf bytes.Buffer
	if err := png.Encode(&gotBuf, got); err != nil {
		t.Fatalf("encode rendered image: %v", err)
	}

	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatalf("create testdata dir: %v", err)
		}
		if err := os.WriteFile(goldenPath, gotBuf.Bytes(), 0o644); err != nil { //nolint:gosec // test helper writing test data
			t.Fatalf("write golden file: %v", err)
		}
		t.Logf("updated golden image: %s", goldenPath)
		return
	}

	wantBytes, err := os.ReadFile(goldenPath) //nolint:gosec // well-known test data path
	if err != nil {
		t.Fatalf("read golden image %s: %v (run with -update to generate)", goldenPath, err)
	}
	if !bytes.Equal(gotBuf.Bytes(), wantBytes) {
		t.Errorf("rendered image differs from golden image %s; run with -update to refresh", goldenPath)
	}
}

// --- Integration tests ---

// TestRender_Integration_TriangleCoversCenter verifies that a triangle whose
// vertices straddle the screen center produces colored pixels at the center.
func TestRender_Integration_TriangleCoversCenter(t *testing.T) {
	fb := framebuffer.New(100, 100)
	Render(goldenTriangleScene(), goldenCamera(), fb, shader.VertexColor)

	center := fb.GetPixel(50, 50)
	if center == (math.Vec3{}) {
		t.Error("center pixel is black; expected the triangle to cover the screen center")
	}
}

// TestRender_Integration_BackgroundIsBlack checks that pixels outside the
// triangle remain at the clear color (black).
func TestRender_Integration_BackgroundIsBlack(t *testing.T) {
	fb := framebuffer.New(100, 100)
	Render(goldenTriangleScene(), goldenCamera(), fb, shader.VertexColor)

	// Top-left corner should not be covered by the triangle.
	corner := fb.GetPixel(0, 0)
	if corner != (math.Vec3{}) {
		t.Errorf("corner pixel = %v, want black background", corner)
	}
}

// TestRender_Integration_CubeProducesPixels verifies that rendering a cube
// produces a non-trivial number of colored pixels across the framebuffer.
func TestRender_Integration_CubeProducesPixels(t *testing.T) {
	fb := framebuffer.New(200, 200)
	scene := NewScene()
	scene.AddMesh(geometry.NewCube())

	cam := camera.New(
		math.Vec3{X: 3, Y: 2, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0,
		float64(200)/float64(200),
		0.1,
		100.0,
	)
	Render(scene, cam, fb, shader.VertexColor)

	colored := countNonBackground(fb)
	if colored < 500 {
		t.Errorf("cube render produced only %d colored pixels (< 500), expected many more", colored)
	}
}

// TestRender_Integration_DepthShaderProducesGrayscale checks that the Depth shader
// produces only grayscale pixels (R==G==B for every colored pixel).
func TestRender_Integration_DepthShaderProducesGrayscale(t *testing.T) {
	fb := framebuffer.New(100, 100)
	Render(goldenTriangleScene(), goldenCamera(), fb, shader.Depth)

	const eps = 1.0 / 255.0
	for y := 0; y < fb.Height; y++ {
		for x := 0; x < fb.Width; x++ {
			c := fb.GetPixel(x, y)
			if c == (math.Vec3{}) {
				continue // background pixel
			}
			if abs64(c.X-c.Y) > eps || abs64(c.Y-c.Z) > eps {
				t.Errorf("pixel (%d,%d) = %v is not grayscale (R≠G or G≠B)", x, y, c)
				return
			}
		}
	}
}

// TestRender_Integration_FlatColorShaderFillsTriangle verifies that every
// non-background pixel in a flat-colored render matches the chosen color.
func TestRender_Integration_FlatColorShaderFillsTriangle(t *testing.T) {
	fb := framebuffer.New(100, 100)
	want := math.Vec3{X: 0.5, Y: 0.2, Z: 0.8}
	Render(goldenTriangleScene(), goldenCamera(), fb, shader.NewFlatColor(want))

	const eps = 1.0 / 255.0 // tolerance for float→byte→float round-trip
	for y := 0; y < fb.Height; y++ {
		for x := 0; x < fb.Width; x++ {
			c := fb.GetPixel(x, y)
			if c == (math.Vec3{}) {
				continue // background
			}
			if abs64(c.X-want.X) > eps || abs64(c.Y-want.Y) > eps || abs64(c.Z-want.Z) > eps {
				t.Errorf("pixel (%d,%d) = %v, want flat %v", x, y, c, want)
				return
			}
		}
	}
}

// TestRender_Integration_MultipleMeshes ensures that adding two meshes to a scene
// results in more colored pixels than either mesh alone.
func TestRender_Integration_MultipleMeshes(t *testing.T) {
	cam := goldenCamera()

	fbSingle := framebuffer.New(100, 100)
	sceneSingle := NewScene()
	sceneSingle.AddMesh(goldenTriangleScene().Meshes[0])
	Render(sceneSingle, cam, fbSingle, shader.VertexColor)
	pixelsSingle := countNonBackground(fbSingle)

	// Second triangle: shifted left so it adds new pixels.
	m2 := geometry.NewMesh()
	green := math.Vec3{X: 0, Y: 1, Z: 0}
	m2.AddVertex(geometry.NewVertex(math.Vec3{X: -2, Y: -1, Z: 0}, green))
	m2.AddVertex(geometry.NewVertex(math.Vec3{X: 0, Y: -1, Z: 0}, green))
	m2.AddVertex(geometry.NewVertex(math.Vec3{X: -1, Y: 1, Z: 0}, green))
	m2.AddTriangle(0, 1, 2)

	fbDouble := framebuffer.New(100, 100)
	sceneDouble := NewScene()
	sceneDouble.AddMesh(goldenTriangleScene().Meshes[0])
	sceneDouble.AddMesh(m2)
	Render(sceneDouble, cam, fbDouble, shader.VertexColor)
	pixelsDouble := countNonBackground(fbDouble)

	if pixelsDouble <= pixelsSingle {
		t.Errorf("two meshes produced %d pixels, single mesh produced %d; expected more from two meshes",
			pixelsDouble, pixelsSingle)
	}
}

// TestRender_Integration_DifferentCameraAngles verifies that rotating the camera
// to different positions changes which pixels are rendered.
func TestRender_Integration_DifferentCameraAngles(t *testing.T) {
	scene := NewScene()
	scene.AddMesh(geometry.NewCube())

	makeCamera := func(pos math.Vec3) camera.Camera {
		return camera.New(pos, math.Vec3{}, math.Vec3{X: 0, Y: 1, Z: 0}, 45.0, 1.0, 0.1, 100.0)
	}

	// Render from front and from side; results should differ.
	fb1 := framebuffer.New(50, 50)
	Render(scene, makeCamera(math.Vec3{X: 0, Y: 0, Z: 5}), fb1, shader.VertexColor)

	fb2 := framebuffer.New(50, 50)
	Render(scene, makeCamera(math.Vec3{X: 5, Y: 0, Z: 0}), fb2, shader.VertexColor)

	same := true
	for i := range fb1.ColorBuffer {
		if fb1.ColorBuffer[i] != fb2.ColorBuffer[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("renders from different camera angles produced identical framebuffers")
	}
}

// abs64 returns the absolute value of x.
func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

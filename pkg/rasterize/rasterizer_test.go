package rasterize

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/framebuffer"
	math "github.com/muddl/go_toy_renderer/pkg/math"
	"github.com/muddl/go_toy_renderer/pkg/shader"
)

const epsilon = 0.001

func approxEqual(a, b float64) bool {
	if a > b {
		return a-b <= epsilon
	}
	return b-a <= epsilon
}

func vec3Equal(a, b math.Vec3) bool {
	return approxEqual(a.X, b.X) && approxEqual(a.Y, b.Y) && approxEqual(a.Z, b.Z)
}

// countPixelsSet counts how many pixels in the framebuffer have non-black color.
func countPixelsSet(fb *framebuffer.Framebuffer) int {
	count := 0
	black := math.Vec3{}
	for _, c := range fb.ColorBuffer {
		if !vec3Equal(c, black) {
			count++
		}
	}
	return count
}

// --- Degenerate triangles ---

func TestRasterizer_Triangle_DegenerateTriangleRendersNothing(t *testing.T) {
	fb := framebuffer.New(10, 10)
	// Collinear vertices — zero area
	v0 := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.5, Color: math.Vec3{X: 1, Y: 0, Z: 0}}
	v1 := ScreenVertex{X: 3.5, Y: 3.5, Z: 0.5, Color: math.Vec3{X: 0, Y: 1, Z: 0}}
	v2 := ScreenVertex{X: 5.5, Y: 5.5, Z: 0.5, Color: math.Vec3{X: 0, Y: 0, Z: 1}}

	Triangle(v0, v1, v2, fb)

	if got := countPixelsSet(fb); got != 0 {
		t.Errorf("degenerate triangle set %d pixels, want 0", got)
	}
}

// --- Single pixel coverage ---

func TestRasterizer_Triangle_SinglePixelCoverage(t *testing.T) {
	fb := framebuffer.New(5, 5)
	white := math.Vec3{X: 1, Y: 1, Z: 1}
	// Small triangle that covers only the center of pixel (0,0) at (0.5, 0.5)
	v0 := ScreenVertex{X: 0.1, Y: 0.1, Z: 0.5, Color: white}
	v1 := ScreenVertex{X: 0.9, Y: 0.1, Z: 0.5, Color: white}
	v2 := ScreenVertex{X: 0.5, Y: 0.9, Z: 0.5, Color: white}

	Triangle(v0, v1, v2, fb)

	// Pixel (0,0) should be set
	if !vec3Equal(fb.GetPixel(0, 0), white) {
		t.Errorf("pixel (0,0) = %v, want white", fb.GetPixel(0, 0))
	}
	// No other pixels should be set
	total := countPixelsSet(fb)
	if total != 1 {
		t.Errorf("single-pixel triangle set %d pixels, want 1", total)
	}
}

// --- Axis-aligned triangle ---

func TestRasterizer_Triangle_AxisAlignedTriangleFillsCorrectPixels(t *testing.T) {
	fb := framebuffer.New(4, 4)
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	// Right triangle: vertices at pixel corners, covers bottom-left half of 4x4 grid
	// v0=(0,0), v1=(4,0), v2=(0,4) — large CCW triangle covering 10 pixel centers
	// Pixel centers: (0.5,0.5), (1.5,0.5), ... are inside if x+y <= 3.5 (roughly)
	v0 := ScreenVertex{X: 0, Y: 0, Z: 0.5, Color: red}
	v1 := ScreenVertex{X: 4, Y: 0, Z: 0.5, Color: red}
	v2 := ScreenVertex{X: 0, Y: 4, Z: 0.5, Color: red}

	Triangle(v0, v1, v2, fb)

	// Pixels on the left/bottom of the diagonal should be set
	// Pixel (0,0) at center (0.5,0.5): inside (0.5+0.5 < 4)
	if !vec3Equal(fb.GetPixel(0, 0), red) {
		t.Errorf("pixel (0,0) should be set (inside triangle)")
	}
	// Pixel (3,3) at center (3.5,3.5): on or outside (3.5+3.5 >= 4)
	// Pixel (0,3) at center (0.5,3.5): inside (0.5+3.5 = 4, borderline)
	// The key check: at least some pixels in each row are set
	pixelsInRow0 := 0
	for x := 0; x < 4; x++ {
		if !vec3Equal(fb.GetPixel(x, 0), math.Vec3{}) {
			pixelsInRow0++
		}
	}
	if pixelsInRow0 < 3 {
		t.Errorf("row 0 has %d set pixels, want at least 3", pixelsInRow0)
	}
}

// --- Color interpolation at vertex ---

func TestRasterizer_Triangle_ColorAtVertexMatchesVertexColor(t *testing.T) {
	fb := framebuffer.New(10, 10)
	// v0 placed exactly at pixel (1,1) center (1.5, 1.5) with unique red color
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	green := math.Vec3{X: 0, Y: 1, Z: 0}
	blue := math.Vec3{X: 0, Y: 0, Z: 1}
	v0 := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.5, Color: red}
	v1 := ScreenVertex{X: 5.5, Y: 1.5, Z: 0.5, Color: green}
	v2 := ScreenVertex{X: 1.5, Y: 5.5, Z: 0.5, Color: blue}

	Triangle(v0, v1, v2, fb)

	// At pixel (1,1) (center 1.5,1.5 = v0 exactly), barycentric weight for v0 = 1
	// so color should equal v0.Color = red
	got := fb.GetPixel(1, 1)
	if !vec3Equal(got, red) {
		t.Errorf("pixel at v0 = %v, want red %v (barycentric weight = 1 at vertex)", got, red)
	}
}

// --- Color interpolation at centroid ---

func TestRasterizer_Triangle_ColorAtCentroidIsInterpolated(t *testing.T) {
	fb := framebuffer.New(20, 20)
	// Triangle with vertices at pixel centers; centroid pixel should get average color
	// v0=(1.5,1.5) red, v1=(7.5,1.5) green, v2=(4.5,7.5) blue
	// Centroid = ((1.5+7.5+4.5)/3, (1.5+1.5+7.5)/3) = (4.5, 3.5)
	// Pixel at centroid = (4, 3), center (4.5, 3.5) = exactly the centroid
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	green := math.Vec3{X: 0, Y: 1, Z: 0}
	blue := math.Vec3{X: 0, Y: 0, Z: 1}
	v0 := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.5, Color: red}
	v1 := ScreenVertex{X: 7.5, Y: 1.5, Z: 0.5, Color: green}
	v2 := ScreenVertex{X: 4.5, Y: 7.5, Z: 0.5, Color: blue}

	Triangle(v0, v1, v2, fb)

	// At centroid, each barycentric weight = 1/3, so:
	// color = (1/3)*red + (1/3)*green + (1/3)*blue = (1/3, 1/3, 1/3)
	expected := math.Vec3{X: 1.0 / 3.0, Y: 1.0 / 3.0, Z: 1.0 / 3.0}
	got := fb.GetPixel(4, 3)
	if !vec3Equal(got, expected) {
		t.Errorf("centroid pixel (4,3) = %v, want ~%v (1/3 each channel)", got, expected)
	}
}

// --- Depth interpolation ---

func TestRasterizer_Triangle_DepthInterpolation(t *testing.T) {
	fb := framebuffer.New(10, 10)
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	blue := math.Vec3{X: 0, Y: 0, Z: 1}

	// Triangle 1: at depth 0.8 (far)
	v0far := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.8, Color: red}
	v1far := ScreenVertex{X: 5.5, Y: 1.5, Z: 0.8, Color: red}
	v2far := ScreenVertex{X: 3.5, Y: 5.5, Z: 0.8, Color: red}
	Triangle(v0far, v1far, v2far, fb)

	// Triangle 2: same coverage but closer (depth 0.3)
	v0near := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.3, Color: blue}
	v1near := ScreenVertex{X: 5.5, Y: 1.5, Z: 0.3, Color: blue}
	v2near := ScreenVertex{X: 3.5, Y: 5.5, Z: 0.3, Color: blue}
	Triangle(v0near, v1near, v2near, fb)

	// Center pixel (3,3) should be blue (closer triangle wins)
	got := fb.GetPixel(3, 3)
	if !vec3Equal(got, blue) {
		t.Errorf("pixel (3,3) = %v, want blue (closer triangle at z=0.3)", got)
	}

	// Depth at center should be ~0.3
	depth := fb.GetDepth(3, 3)
	if !approxEqual(depth, 0.3) {
		t.Errorf("depth at (3,3) = %f, want ~0.3", depth)
	}
}

// --- Off-screen triangle ---

func TestRasterizer_Triangle_OffScreenTriangleRendersNothing(t *testing.T) {
	fb := framebuffer.New(10, 10)
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	// Triangle entirely at negative coordinates
	v0 := ScreenVertex{X: -5, Y: -5, Z: 0.5, Color: red}
	v1 := ScreenVertex{X: -1, Y: -5, Z: 0.5, Color: red}
	v2 := ScreenVertex{X: -3, Y: -1, Z: 0.5, Color: red}

	Triangle(v0, v1, v2, fb)

	if got := countPixelsSet(fb); got != 0 {
		t.Errorf("off-screen triangle set %d pixels, want 0", got)
	}
}

// --- Partially off-screen ---

func TestRasterizer_Triangle_PartiallyOffScreenNoPanic(t *testing.T) {
	fb := framebuffer.New(5, 5)
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	// Triangle that straddles the left boundary — no panic expected
	v0 := ScreenVertex{X: -2, Y: 0, Z: 0.5, Color: red}
	v1 := ScreenVertex{X: 3, Y: 0, Z: 0.5, Color: red}
	v2 := ScreenVertex{X: 0.5, Y: 4, Z: 0.5, Color: red}

	// Must not panic
	Triangle(v0, v1, v2, fb)

	// At least some pixels should be set (the in-bounds part)
	if got := countPixelsSet(fb); got == 0 {
		t.Error("partially off-screen triangle set 0 pixels, expected some in-bounds pixels")
	}
	// All set pixels must be within bounds
	for y := 0; y < fb.Height; y++ {
		for x := 0; x < fb.Width; x++ {
			// No out-of-bounds panic means the test passes by reaching here
			_ = fb.GetPixel(x, y)
		}
	}
}

// --- Clockwise winding ---

func TestRasterizer_Triangle_ClockwiseWindingRendersCorrectly(t *testing.T) {
	fbCCW := framebuffer.New(10, 10)
	fbCW := framebuffer.New(10, 10)
	white := math.Vec3{X: 1, Y: 1, Z: 1}

	// CCW triangle
	v0 := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.5, Color: white}
	v1 := ScreenVertex{X: 6.5, Y: 1.5, Z: 0.5, Color: white}
	v2 := ScreenVertex{X: 4.0, Y: 6.5, Z: 0.5, Color: white}
	Triangle(v0, v1, v2, fbCCW)

	// CW triangle (v1 and v2 swapped)
	Triangle(v0, v2, v1, fbCW)

	// Both should set the same number of pixels
	ccwCount := countPixelsSet(fbCCW)
	cwCount := countPixelsSet(fbCW)
	if ccwCount == 0 {
		t.Fatal("CCW triangle set 0 pixels")
	}
	if ccwCount != cwCount {
		t.Errorf("CCW set %d pixels, CW set %d pixels — should be equal", ccwCount, cwCount)
	}
}

// --- TriangleShaded tests ---

// TestRasterizer_TriangleShaded_VertexColorMatchesTriangle verifies that
// TriangleShaded with shader.VertexColor produces identical output to Triangle.
func TestRasterizer_TriangleShaded_VertexColorMatchesTriangle(t *testing.T) {
	fbA := framebuffer.New(10, 10)
	fbB := framebuffer.New(10, 10)
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	green := math.Vec3{X: 0, Y: 1, Z: 0}
	blue := math.Vec3{X: 0, Y: 0, Z: 1}

	v0 := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.5, Color: red}
	v1 := ScreenVertex{X: 7.5, Y: 1.5, Z: 0.5, Color: green}
	v2 := ScreenVertex{X: 4.5, Y: 7.5, Z: 0.5, Color: blue}

	Triangle(v0, v1, v2, fbA)
	TriangleShaded(v0, v1, v2, shader.VertexColor, fbB)

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			a := fbA.GetPixel(x, y)
			b := fbB.GetPixel(x, y)
			if !vec3Equal(a, b) {
				t.Errorf("pixel (%d,%d): Triangle=%v, TriangleShaded(VertexColor)=%v", x, y, a, b)
			}
		}
	}
}

// TestRasterizer_TriangleShaded_FlatColorOverridesVertexColor verifies that
// NewFlatColor shader makes all pixels the specified constant color regardless of vertex colors.
func TestRasterizer_TriangleShaded_FlatColorOverridesVertexColor(t *testing.T) {
	fb := framebuffer.New(10, 10)
	red := math.Vec3{X: 1, Y: 0, Z: 0}
	green := math.Vec3{X: 0, Y: 1, Z: 0}
	blue := math.Vec3{X: 0, Y: 0, Z: 1}
	flatBlue := shader.NewFlatColor(blue)

	// Vertices with red and green colors, but flat-color shader should output blue.
	v0 := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.5, Color: red}
	v1 := ScreenVertex{X: 7.5, Y: 1.5, Z: 0.5, Color: green}
	v2 := ScreenVertex{X: 4.5, Y: 7.5, Z: 0.5, Color: red}

	TriangleShaded(v0, v1, v2, flatBlue, fb)

	// Every set pixel should be blue.
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			c := fb.GetPixel(x, y)
			if c == (math.Vec3{}) {
				continue // unset pixel, skip
			}
			if !vec3Equal(c, blue) {
				t.Errorf("pixel (%d,%d) = %v, want blue %v", x, y, c, blue)
			}
		}
	}
	if countPixelsSet(fb) == 0 {
		t.Error("expected some pixels to be set, got 0")
	}
}

// TestRasterizer_TriangleShaded_DepthShaderOutputsGrayscale verifies that
// shader.Depth produces equal RGB components (grayscale) for every pixel.
func TestRasterizer_TriangleShaded_DepthShaderOutputsGrayscale(t *testing.T) {
	fb := framebuffer.New(10, 10)
	white := math.Vec3{X: 1, Y: 1, Z: 1}

	v0 := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.4, Color: white}
	v1 := ScreenVertex{X: 7.5, Y: 1.5, Z: 0.6, Color: white}
	v2 := ScreenVertex{X: 4.5, Y: 7.5, Z: 0.5, Color: white}

	TriangleShaded(v0, v1, v2, shader.Depth, fb)

	if countPixelsSet(fb) == 0 {
		t.Fatal("depth shader rendered 0 pixels")
	}
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			c := fb.GetPixel(x, y)
			if c == (math.Vec3{}) {
				continue
			}
			// Grayscale: R == G == B
			if !approxEqual(c.X, c.Y) || !approxEqual(c.Y, c.Z) {
				t.Errorf("pixel (%d,%d) = %v is not grayscale (depth shader)", x, y, c)
			}
		}
	}
}

// TestRasterizer_TriangleShaded_DegenerateTriangle verifies that degenerate
// (zero-area) triangles set no pixels regardless of the shader.
func TestRasterizer_TriangleShaded_DegenerateTriangle(t *testing.T) {
	fb := framebuffer.New(10, 10)
	red := math.Vec3{X: 1, Y: 0, Z: 0}

	// Collinear vertices — zero area.
	v0 := ScreenVertex{X: 1.5, Y: 1.5, Z: 0.5, Color: red}
	v1 := ScreenVertex{X: 3.5, Y: 3.5, Z: 0.5, Color: red}
	v2 := ScreenVertex{X: 5.5, Y: 5.5, Z: 0.5, Color: red}

	TriangleShaded(v0, v1, v2, shader.VertexColor, fb)

	if got := countPixelsSet(fb); got != 0 {
		t.Errorf("degenerate triangle set %d pixels with shader, want 0", got)
	}
}

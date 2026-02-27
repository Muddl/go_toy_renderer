package framebuffer

import (
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"testing"

	math "github.com/muddl/go_toy_renderer/pkg/math"
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

// --- NewFramebuffer ---

func TestFramebuffer_New_HasCorrectDimensions(t *testing.T) {
	fb := New(640, 480)

	if fb.Width != 640 {
		t.Errorf("Width = %d, want 640", fb.Width)
	}
	if fb.Height != 480 {
		t.Errorf("Height = %d, want 480", fb.Height)
	}
}

func TestFramebuffer_New_ColorBufferHasCorrectSize(t *testing.T) {
	fb := New(10, 20)

	want := 10 * 20
	if len(fb.ColorBuffer) != want {
		t.Errorf("len(ColorBuffer) = %d, want %d", len(fb.ColorBuffer), want)
	}
}

func TestFramebuffer_New_DepthBufferHasCorrectSize(t *testing.T) {
	fb := New(10, 20)

	want := 10 * 20
	if len(fb.DepthBuffer) != want {
		t.Errorf("len(DepthBuffer) = %d, want %d", len(fb.DepthBuffer), want)
	}
}

func TestFramebuffer_New_ColorBufferInitializedToBlack(t *testing.T) {
	fb := New(4, 4)
	black := math.Vec3{}

	for i, c := range fb.ColorBuffer {
		if !vec3Equal(c, black) {
			t.Errorf("ColorBuffer[%d] = %v, want black %v", i, c, black)
		}
	}
}

func TestFramebuffer_New_DepthBufferInitializedToOne(t *testing.T) {
	fb := New(4, 4)

	for i, d := range fb.DepthBuffer {
		if !approxEqual(d, 1.0) {
			t.Errorf("DepthBuffer[%d] = %f, want 1.0", i, d)
		}
	}
}

// --- Clear ---

func TestFramebuffer_Clear_SetsAllColorPixels(t *testing.T) {
	fb := New(5, 5)
	skyBlue := math.Vec3{X: 0.39, Y: 0.58, Z: 0.93}

	fb.Clear(skyBlue, 1.0)

	for i, c := range fb.ColorBuffer {
		if !vec3Equal(c, skyBlue) {
			t.Errorf("ColorBuffer[%d] = %v, want %v", i, c, skyBlue)
		}
	}
}

func TestFramebuffer_Clear_SetsAllDepthValues(t *testing.T) {
	fb := New(5, 5)
	fb.Clear(math.Vec3{}, 0.5)

	for i, d := range fb.DepthBuffer {
		if !approxEqual(d, 0.5) {
			t.Errorf("DepthBuffer[%d] = %f, want 0.5", i, d)
		}
	}
}

// --- SetPixel / GetPixel ---

func TestFramebuffer_SetPixel_WritesColorAtValidCoordinate(t *testing.T) {
	fb := New(10, 10)
	red := math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}

	fb.SetPixel(3, 4, red, 0.5)
	got := fb.GetPixel(3, 4)

	if !vec3Equal(got, red) {
		t.Errorf("GetPixel(3,4) = %v, want %v", got, red)
	}
}

func TestFramebuffer_SetPixel_OutOfBoundsIsIgnored(t *testing.T) {
	fb := New(10, 10)
	red := math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}

	// Should not panic
	fb.SetPixel(-1, 0, red, 0.5)
	fb.SetPixel(10, 0, red, 0.5)
	fb.SetPixel(0, -1, red, 0.5)
	fb.SetPixel(0, 10, red, 0.5)

	// All pixels should still be black (default)
	black := math.Vec3{}
	for i, c := range fb.ColorBuffer {
		if !vec3Equal(c, black) {
			t.Errorf("ColorBuffer[%d] = %v after out-of-bounds write, want black", i, c)
		}
	}
}

func TestFramebuffer_GetPixel_ReturnsCorrectColor(t *testing.T) {
	fb := New(10, 10)
	green := math.Vec3{X: 0.0, Y: 1.0, Z: 0.0}

	fb.SetPixel(7, 2, green, 0.3)
	got := fb.GetPixel(7, 2)

	if !vec3Equal(got, green) {
		t.Errorf("GetPixel(7,2) = %v, want %v", got, green)
	}
}

func TestFramebuffer_GetPixel_OutOfBoundsReturnsZero(t *testing.T) {
	fb := New(10, 10)

	got := fb.GetPixel(-1, 0)
	if !vec3Equal(got, math.Vec3{}) {
		t.Errorf("GetPixel(-1,0) = %v, want zero Vec3", got)
	}

	got = fb.GetPixel(10, 10)
	if !vec3Equal(got, math.Vec3{}) {
		t.Errorf("GetPixel(10,10) = %v, want zero Vec3", got)
	}
}

// --- Depth testing ---

func TestFramebuffer_SetPixel_CloserPixelOverwritesFarther(t *testing.T) {
	fb := New(10, 10)
	red := math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}
	blue := math.Vec3{X: 0.0, Y: 0.0, Z: 1.0}

	// Write farther pixel first
	fb.SetPixel(5, 5, red, 0.8)
	// Write closer pixel — should overwrite
	fb.SetPixel(5, 5, blue, 0.3)

	got := fb.GetPixel(5, 5)
	if !vec3Equal(got, blue) {
		t.Errorf("GetPixel(5,5) = %v, want blue (closer pixel)", got)
	}
}

func TestFramebuffer_SetPixel_FartherPixelDoesNotOverwrite(t *testing.T) {
	fb := New(10, 10)
	red := math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}
	blue := math.Vec3{X: 0.0, Y: 0.0, Z: 1.0}

	// Write closer pixel first
	fb.SetPixel(5, 5, red, 0.3)
	// Write farther pixel — should NOT overwrite
	fb.SetPixel(5, 5, blue, 0.9)

	got := fb.GetPixel(5, 5)
	if !vec3Equal(got, red) {
		t.Errorf("GetPixel(5,5) = %v, want red (closer pixel not overwritten)", got)
	}
}

func TestFramebuffer_SetPixel_EqualDepthIsIgnored(t *testing.T) {
	fb := New(10, 10)
	red := math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}
	blue := math.Vec3{X: 0.0, Y: 0.0, Z: 1.0}

	fb.SetPixel(5, 5, red, 0.5)
	// Same depth — should NOT overwrite
	fb.SetPixel(5, 5, blue, 0.5)

	got := fb.GetPixel(5, 5)
	if !vec3Equal(got, red) {
		t.Errorf("GetPixel(5,5) = %v, want red (equal depth does not overwrite)", got)
	}
}

func TestFramebuffer_SetPixel_UpdatesDepthBuffer(t *testing.T) {
	fb := New(10, 10)
	red := math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}

	fb.SetPixel(2, 3, red, 0.4)
	got := fb.GetDepth(2, 3)

	if !approxEqual(got, 0.4) {
		t.Errorf("GetDepth(2,3) = %f, want 0.4", got)
	}
}

// --- Buffer indexing ---

func TestFramebuffer_BufferIndexing_CorrectPixelSelected(t *testing.T) {
	fb := New(8, 6)
	// Write a unique color at (2, 1)
	// Expected index = 1*8 + 2 = 10
	unique := math.Vec3{X: 0.25, Y: 0.5, Z: 0.75}

	fb.SetPixel(2, 1, unique, 0.5)

	// Direct index check: index 10 should have the unique color
	idx := 1*8 + 2
	got := fb.ColorBuffer[idx]
	if !vec3Equal(got, unique) {
		t.Errorf("ColorBuffer[%d] = %v, want %v (indexing check)", idx, got, unique)
	}

	// Neighbors should be untouched (black)
	if !vec3Equal(fb.ColorBuffer[9], math.Vec3{}) {
		t.Errorf("ColorBuffer[9] should be black (left neighbor), got %v", fb.ColorBuffer[9])
	}
	if !vec3Equal(fb.ColorBuffer[11], math.Vec3{}) {
		t.Errorf("ColorBuffer[11] should be black (right neighbor), got %v", fb.ColorBuffer[11])
	}
}

// --- GetDepth ---

func TestFramebuffer_GetDepth_OutOfBoundsReturnsOne(t *testing.T) {
	fb := New(10, 10)

	got := fb.GetDepth(-1, 0)
	if !approxEqual(got, 1.0) {
		t.Errorf("GetDepth(-1,0) = %f, want 1.0", got)
	}

	got = fb.GetDepth(10, 10)
	if !approxEqual(got, 1.0) {
		t.Errorf("GetDepth(10,10) = %f, want 1.0", got)
	}
}

// --- SavePNG ---

func TestFramebuffer_SavePNG_CreatesFile(t *testing.T) {
	fb := New(4, 4)
	path := filepath.Join(t.TempDir(), "test.png")

	if err := fb.SavePNG(path); err != nil {
		t.Fatalf("SavePNG() error = %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("SavePNG() did not create the file")
	}
}

func TestFramebuffer_SavePNG_HasCorrectDimensions(t *testing.T) {
	fb := New(16, 8)
	path := filepath.Join(t.TempDir(), "test.png")

	if err := fb.SavePNG(path); err != nil {
		t.Fatalf("SavePNG() error = %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("os.Open() error = %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("image.Decode() error = %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 16 {
		t.Errorf("image width = %d, want 16", bounds.Dx())
	}
	if bounds.Dy() != 8 {
		t.Errorf("image height = %d, want 8", bounds.Dy())
	}
}

func TestFramebuffer_SavePNG_PixelColorsMatchBuffer(t *testing.T) {
	fb := New(10, 10)
	// Write a bright red pixel at (3, 5)
	fb.SetPixel(3, 5, math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}, 0.5)
	path := filepath.Join(t.TempDir(), "test.png")

	if err := fb.SavePNG(path); err != nil {
		t.Fatalf("SavePNG() error = %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("os.Open() error = %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("image.Decode() error = %v", err)
	}

	// Check the red pixel at (3, 5)
	r, g, b, _ := img.At(3, 5).RGBA()
	// RGBA() returns 16-bit values (0-65535), normalize to 8-bit
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)

	if r8 < 250 {
		t.Errorf("pixel(3,5).R = %d, want ~255", r8)
	}
	if g8 > 5 {
		t.Errorf("pixel(3,5).G = %d, want ~0", g8)
	}
	if b8 > 5 {
		t.Errorf("pixel(3,5).B = %d, want ~0", b8)
	}

	// Check a black pixel at (0, 0)
	r, g, b, _ = img.At(0, 0).RGBA()
	r8 = uint8(r >> 8)
	g8 = uint8(g >> 8)
	b8 = uint8(b >> 8)

	if r8 > 5 || g8 > 5 || b8 > 5 {
		t.Errorf("pixel(0,0) = (%d,%d,%d), want black", r8, g8, b8)
	}
}

func TestFramebuffer_SavePNG_ClampsBrightColors(t *testing.T) {
	fb := New(4, 4)
	// Set a color beyond [0,1] range — should clamp to 255
	fb.ColorBuffer[0] = math.Vec3{X: 2.0, Y: -0.5, Z: 1.5}
	path := filepath.Join(t.TempDir(), "test.png")

	if err := fb.SavePNG(path); err != nil {
		t.Fatalf("SavePNG() error = %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("os.Open() error = %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatalf("image.Decode() error = %v", err)
	}

	r, g, b, _ := img.At(0, 0).RGBA()
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)

	if r8 != 255 {
		t.Errorf("clamped R = %d, want 255 (from 2.0)", r8)
	}
	if g8 != 0 {
		t.Errorf("clamped G = %d, want 0 (from -0.5)", g8)
	}
	if b8 != 255 {
		t.Errorf("clamped B = %d, want 255 (from 1.5)", b8)
	}
}

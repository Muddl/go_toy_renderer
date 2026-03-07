package gpu_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/gpu"
)

// TestDevice_CreateTexture2D_ReturnsErrorBeforeInit verifies that
// CreateTexture2D returns an error when the device has not been initialised.
func TestDevice_CreateTexture2D_ReturnsErrorBeforeInit(t *testing.T) {
	d := gpu.New()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	_, err := d.CreateTexture2D(img)
	if err == nil {
		t.Fatal("CreateTexture2D before Init should return an error")
	}
}

// TestDevice_CreateTexture2D_NilImage_ReturnsError verifies that a nil image
// is rejected even on an uninitialised device.
func TestDevice_CreateTexture2D_NilImage_ReturnsError(t *testing.T) {
	d := gpu.New()
	_, err := d.CreateTexture2D(nil)
	if err == nil {
		t.Fatal("CreateTexture2D with nil image should return an error")
	}
}

// TestDevice_CreateTexture2D_ZeroSizeImage_ReturnsError verifies that a zero-size
// image is rejected.
func TestDevice_CreateTexture2D_ZeroSizeImage_ReturnsError(t *testing.T) {
	d := gpu.New()
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))
	_, err := d.CreateTexture2D(img)
	if err == nil {
		t.Fatal("CreateTexture2D with zero-size image should return an error")
	}
}

// TestTexture2D_ImageData_PreservesPixels verifies that Texture2DFromImage
// correctly stores image RGBA data in the expected byte layout.
func TestTexture2D_ImageData_PreservesPixels(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 1))
	img.SetRGBA(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	img.SetRGBA(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 128})

	data := gpu.RGBAImageBytes(img)
	if len(data) != 8 {
		t.Fatalf("expected 8 bytes for 2x1 RGBA image, got %d", len(data))
	}
	// First pixel: R=255, G=0, B=0, A=255
	if data[0] != 255 || data[1] != 0 || data[2] != 0 || data[3] != 255 {
		t.Errorf("pixel 0: got (%d,%d,%d,%d), want (255,0,0,255)",
			data[0], data[1], data[2], data[3])
	}
	// Second pixel: R=0, G=255, B=0, A=128
	if data[4] != 0 || data[5] != 255 || data[6] != 0 || data[7] != 128 {
		t.Errorf("pixel 1: got (%d,%d,%d,%d), want (0,255,0,128)",
			data[4], data[5], data[6], data[7])
	}
}

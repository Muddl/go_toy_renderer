//go:build headless

package gpu

import (
	"errors"
	"image"
)

// Texture2D is a stub for headless (CI) builds — no GPU available.
type Texture2D struct {
	Width  uint32
	Height uint32
}

// Release is a no-op in headless builds.
func (t *Texture2D) Release() {}

// CreateTexture2D always returns an error in headless builds.
func (d *Device) CreateTexture2D(img *image.RGBA) (*Texture2D, error) {
	if img == nil {
		return nil, errors.New("gpu: CreateTexture2D: image is nil")
	}
	bounds := img.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		return nil, errors.New("gpu: CreateTexture2D: image has zero dimension")
	}
	return nil, errors.New("gpu: not available in headless builds")
}

//go:build !headless

package gpu

import (
	"errors"
	"fmt"
	"image"

	"github.com/go-webgpu/webgpu/wgpu"
	"github.com/gogpu/gputypes"
)

// Texture2D is a GPU-resident RGBA texture created from an image.RGBA.
type Texture2D struct {
	texture *wgpu.Texture
	view    *wgpu.TextureView
	Width   uint32
	Height  uint32
}

// Release frees the GPU resources held by the Texture2D.
func (t *Texture2D) Release() {
	if t.view != nil {
		t.view.Release()
		t.view = nil
	}
	if t.texture != nil {
		t.texture.Destroy()
		t.texture.Release()
		t.texture = nil
	}
}

// CreateTexture2D creates a GPU texture from the given RGBA image and uploads
// the pixel data via Queue.WriteTexture. The caller must call Release() when done.
func (d *Device) CreateTexture2D(img *image.RGBA) (*Texture2D, error) {
	if img == nil {
		return nil, errors.New("gpu: CreateTexture2D: image is nil")
	}
	bounds := img.Bounds()
	w := uint32(bounds.Dx())
	h := uint32(bounds.Dy())
	if w == 0 || h == 0 {
		return nil, errors.New("gpu: CreateTexture2D: image has zero dimension")
	}
	if !d.ready {
		return nil, errors.New("gpu: CreateTexture2D: device not initialized")
	}

	// Create the destination texture (sampled + copy-dst).
	tex := d.device.CreateTexture(&wgpu.TextureDescriptor{
		Size:          gputypes.Extent3D{Width: w, Height: h, DepthOrArrayLayers: 1},
		MipLevelCount: 1,
		SampleCount:   1,
		Dimension:     gputypes.TextureDimension2D,
		Format:        gputypes.TextureFormatRGBA8Unorm,
		Usage:         gputypes.TextureUsageTextureBinding | gputypes.TextureUsageCopyDst,
	})
	if tex == nil {
		return nil, errors.New("gpu: CreateTexture2D: CreateTexture returned nil")
	}

	// BytesPerRow must be aligned to 256 bytes per WebGPU spec.
	rawBytesPerRow := w * 4
	alignedBytesPerRow := ((rawBytesPerRow + 255) / 256) * 256

	var uploadData []byte
	if alignedBytesPerRow == rawBytesPerRow {
		// No padding needed — use image Pix slice directly.
		uploadData = RGBAImageBytes(img)
	} else {
		// Pad each row to the aligned stride.
		uploadData = make([]byte, alignedBytesPerRow*h)
		src := RGBAImageBytes(img)
		for row := uint32(0); row < h; row++ {
			srcStart := row * rawBytesPerRow
			dstStart := row * alignedBytesPerRow
			copy(uploadData[dstStart:dstStart+rawBytesPerRow], src[srcStart:srcStart+rawBytesPerRow])
		}
	}

	d.queue.WriteTexture(
		&wgpu.TexelCopyTextureInfo{
			Texture:  tex.Handle(),
			MipLevel: 0,
			Origin:   gputypes.Origin3D{},
			Aspect:   wgpu.TextureAspectAll,
		},
		uploadData,
		&wgpu.TexelCopyBufferLayout{
			Offset:       0,
			BytesPerRow:  alignedBytesPerRow,
			RowsPerImage: h,
		},
		&gputypes.Extent3D{Width: w, Height: h, DepthOrArrayLayers: 1},
	)

	view := tex.CreateView(nil)
	if view == nil {
		tex.Destroy()
		tex.Release()
		return nil, fmt.Errorf("gpu: CreateTexture2D: CreateView returned nil")
	}

	return &Texture2D{texture: tex, view: view, Width: w, Height: h}, nil
}

package gpu

import "image"

// RGBAImageBytes returns the raw RGBA pixel bytes from img in row-major order.
// Each pixel is 4 bytes: R, G, B, A. The returned slice references the image
// Pix field directly and should not be modified after upload.
func RGBAImageBytes(img *image.RGBA) []byte {
	return img.Pix
}

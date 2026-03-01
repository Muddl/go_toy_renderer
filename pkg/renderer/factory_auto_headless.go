//go:build headless

package renderer

// newAutoBackend returns a CPUBackend in headless builds.
// No GPU hardware or display is available in headless (CI) mode.
func newAutoBackend() Renderer {
	return &CPUBackend{}
}

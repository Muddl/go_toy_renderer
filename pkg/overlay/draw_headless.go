//go:build headless

package overlay

import "github.com/muddl/go_toy_renderer/pkg/framebuffer"

// ComposeIntoFramebuffer is a no-op in headless builds (no window, no blit).
func ComposeIntoFramebuffer(_ *TextLayer, _ *framebuffer.Framebuffer) {}

package overlay

import "fmt"

const (
	glyphScale  = 2                              // pixel scale multiplier for legibility
	cellW       = GlyphWidth * glyphScale        // rendered cell width in pixels
	cellH       = GlyphHeight * glyphScale       // rendered cell height in pixels
	lineSpacing = cellH + 2                      // vertical distance between lines
	marginX     = 8                              // left margin in pixels
	marginY     = 8                              // top margin in pixels
)

// TextLayer renders overlay text into an RGBA byte slice at a fixed top-left
// position. It is sized to fit the widest possible Level 3 output.
type TextLayer struct {
	width  int
	height int
	pixels []byte // RGBA, row-major, top-left origin
}

// NewTextLayer creates a TextLayer sized to hold up to maxLines lines of text
// at the given framebuffer dimensions. Width is capped to a comfortable
// overlay column on the left edge.
func NewTextLayer(fbWidth, fbHeight int) *TextLayer {
	const maxCols = 32
	w := maxCols*cellW + marginX*2
	if w > fbWidth {
		w = fbWidth
	}
	const maxLines = 7
	h := maxLines*lineSpacing + marginY*2
	if h > fbHeight {
		h = fbHeight
	}
	return &TextLayer{
		width:  w,
		height: h,
		pixels: make([]byte, w*h*4),
	}
}

// Width returns the pixel width of the text layer.
func (tl *TextLayer) Width() int { return tl.width }

// Height returns the pixel height of the text layer.
func (tl *TextLayer) Height() int { return tl.height }

// Pixels returns the raw RGBA pixel buffer.
func (tl *TextLayer) Pixels() []byte { return tl.pixels }

// Render clears the layer and draws the lines appropriate for the given level
// and metrics, returning the RGBA pixel buffer.
func (tl *TextLayer) Render(level Level, m Metrics) []byte {
	// Clear to transparent black.
	for i := range tl.pixels {
		tl.pixels[i] = 0
	}
	if level == LevelOff {
		return tl.pixels
	}

	lines := buildLines(level, m)
	for i, line := range lines {
		tl.drawString(marginX, marginY+i*lineSpacing, line)
	}
	return tl.pixels
}

// buildLines returns the text lines appropriate for the given level.
func buildLines(level Level, m Metrics) []string {
	switch level {
	case LevelFPS:
		return []string{
			fmt.Sprintf("FPS: %.1f", m.FPS),
		}
	case LevelTimings:
		return []string{
			fmt.Sprintf("FPS: %.1f", m.FPS),
			fmt.Sprintf("FRAME: %.2f MS", m.FrameTimeMS),
			fmt.Sprintf("CPU:   %.2f MS", m.CPUFrameTimeMS),
		}
	case LevelDetail:
		return []string{
			fmt.Sprintf("FPS: %.1f", m.FPS),
			fmt.Sprintf("FRAME: %.2f MS", m.FrameTimeMS),
			fmt.Sprintf("CPU:   %.2f MS", m.CPUFrameTimeMS),
			fmt.Sprintf("GEO:   %.2f MS", m.GeometryUploadMS),
			fmt.Sprintf("PASS:  %.2f MS", m.RenderPassMS),
			fmt.Sprintf("PRES:  %.2f MS", m.PresentMS),
			fmt.Sprintf("BACKEND: %s  V:%d T:%d", m.Backend, m.VertexCount, m.TriangleCount),
		}
	default:
		return nil
	}
}

// drawString renders a string of ASCII characters into the layer starting at
// pixel (x, y). Characters outside the mapped glyph set are drawn as blanks.
func (tl *TextLayer) drawString(x, y int, s string) {
	cx := x
	for i := 0; i < len(s); i++ {
		tl.drawGlyph(cx, y, GlyphFor(s[i]))
		cx += cellW
		if cx+cellW > tl.width {
			break
		}
	}
}

// drawGlyph renders a single glyph at pixel (px, py) using white text on a
// semi-transparent dark background for contrast.
func (tl *TextLayer) drawGlyph(px, py int, g Glyph) {
	for row := 0; row < GlyphHeight; row++ {
		for col := 0; col < GlyphWidth; col++ {
			lit := g[row]&(1<<uint(GlyphWidth-1-col)) != 0
			for sy := 0; sy < glyphScale; sy++ {
				for sx := 0; sx < glyphScale; sx++ {
					dx := px + col*glyphScale + sx
					dy := py + row*glyphScale + sy
					if dx < 0 || dx >= tl.width || dy < 0 || dy >= tl.height {
						continue
					}
					idx := (dy*tl.width + dx) * 4
					if lit {
						// Bright green text (visible on both light and dark scenes).
						tl.pixels[idx+0] = 0x00 // R
						tl.pixels[idx+1] = 0xFF // G
						tl.pixels[idx+2] = 0x00 // B
						tl.pixels[idx+3] = 0xFF // A (opaque)
					} else {
						// Semi-transparent dark background behind text.
						tl.pixels[idx+0] = 0x00
						tl.pixels[idx+1] = 0x00
						tl.pixels[idx+2] = 0x00
						tl.pixels[idx+3] = 0xA0 // ~63% opacity
					}
				}
			}
		}
	}
}

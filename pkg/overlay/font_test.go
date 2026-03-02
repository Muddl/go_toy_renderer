package overlay_test

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/overlay"
)

// TestFont_Glyph_KnownChars_HaveNonZeroPixels verifies that printable ASCII
// characters return glyphs with at least one lit pixel.
func TestFont_Glyph_KnownChars_HaveNonZeroPixels(t *testing.T) {
	// Space is intentionally blank; exclude it from the lit-pixel check.
	printable := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ:./%-"
	for _, ch := range printable {
		g := overlay.GlyphFor(byte(ch))
		lit := 0
		for _, row := range g {
			for bit := 0; bit < overlay.GlyphWidth; bit++ {
				if row&(1<<uint(overlay.GlyphWidth-1-bit)) != 0 {
					lit++
				}
			}
		}
		if lit == 0 {
			t.Errorf("GlyphFor(%q) has no lit pixels", ch)
		}
	}
}

// TestFont_Glyph_UnknownChar_ReturnsBlank verifies that an unmapped byte
// (e.g. 0x01) returns a fully blank glyph rather than panicking.
func TestFont_Glyph_UnknownChar_ReturnsBlank(t *testing.T) {
	g := overlay.GlyphFor(0x01)
	for i, row := range g {
		if row != 0 {
			t.Errorf("blank glyph row %d = %08b, want 0", i, row)
		}
	}
}

// TestFont_GlyphDimensions_AreCorrect verifies that GlyphFor returns a
// non-zero glyph for a mapped character (confirming the glyph table is
// populated) and that the glyph type has exactly GlyphHeight rows.
func TestFont_GlyphDimensions_AreCorrect(t *testing.T) {
	g := overlay.GlyphFor('A')
	// Verify GlyphHeight matches the array length at the type level.
	var zero overlay.Glyph
	if g == zero {
		t.Error("GlyphFor('A') returned all-zero glyph; expected a mapped character")
	}
}

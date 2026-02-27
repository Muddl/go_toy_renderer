# Framebuffer Component

**Status:** ✅ **COMPLETE** (Phase 4 — 2026-02-27)
**Package:** `pkg/framebuffer/`
**Coverage:** 94.6% (21 tests)

## Purpose

Store pixel colors and depth values, handle depth testing, and export to image files.

## Core Type: Framebuffer

**Represents:** 2D grid of pixels with color and depth

**Data:**
- `Width, Height int` - dimensions in pixels
- `ColorBuffer []Vector3` - RGB values per pixel (or []Color)
- `DepthBuffer []float64` - depth value per pixel

**Buffer layout:** Linear array, indexed as `y*Width + x`

## Key Operations

### Creation
```
NewFramebuffer(width, height int) *Framebuffer
```
- Allocate color and depth buffers
- Initialize depth to far plane (1.0) or infinity
- Initialize color to background color (e.g., black or sky blue)

### Pixel Write with Depth Test
```
SetPixel(x, y int, color Vector3, depth float64)
```
**Steps:**
1. Check bounds (is pixel inside framebuffer?)
2. Get current depth at (x, y)
3. If new depth < current depth (closer):
   - Update depth buffer
   - Update color buffer
4. Otherwise: discard (occluded)

**Depth test:** Closer objects overwrite farther objects.

### Pixel Read
```
GetPixel(x, y int) Vector3
```
Retrieve final color at pixel (for debugging or effects).

### Clear
```
Clear(color Vector3, depth float64)
```
Reset all pixels to background color and depth.

## Image Export

### Format Support

**Minimum (MVP):**
- PNG or BMP output

**Recommended library:** `image` package from Go standard library + `image/png`

### Color Conversion

**Internal:** Float RGB (0.0-1.0 range)

**Output:** 8-bit RGB (0-255 range)

**Conversion:**
```
r_byte := clamp(int(r_float * 255), 0, 255)
```

**Gamma correction (optional for MVP):** Apply gamma 2.2 for more realistic output.

### API Example
```go
func (fb *Framebuffer) SavePNG(filename string) error
func (fb *Framebuffer) SaveBMP(filename string) error
```

## Testing Requirements

### Buffer Operations
- Pixel read/write at valid coordinates works
- Out-of-bounds access handled gracefully (return error or ignore)
- Clear resets all pixels to expected values
- Buffer indexing correct (x, y maps to right pixel)

### Depth Testing
- Closer pixel overwrites farther pixel
- Farther pixel does not overwrite closer pixel
- Equal depth handled consistently
- Depth buffer initialized to far plane value

### Image Export
- Output image has correct dimensions
- Pixel colors match framebuffer contents
- Color range clamping works (no overflow)
- File writes successfully to disk

## Design Considerations

### Memory Layout

**Options:**
1. Single linear array: `[]Vector3` of size width×height
2. 2D array: `[][]Vector3`

**Recommendation:** Linear array (option 1)
- Better cache locality
- Simpler memory management
- Standard in graphics

### Depth Buffer Precision

**float64 vs float32:**
- float64: More precision, slower, more memory
- float32: Less precision, faster, standard

**Recommendation for MVP:** float64 (Go's default, simpler)

Post-MVP: Consider float32 for performance.

### Coordinate Convention

**Recommendation:**
- (0, 0) = top-left
- +X = right
- +Y = down

Matches image formats and simplifies export.

### Background Color

**Options:**
- Black (0, 0, 0) - classic
- Cornflower blue (0.39, 0.58, 0.93) - DirectX tradition
- Custom color

**For MVP:** Configurable background color.

## Implemented API

```go
// Constructor
func New(width, height int) *Framebuffer

// Operations
func (fb *Framebuffer) Clear(color math.Vec3, depth float64)
func (fb *Framebuffer) SetPixel(x, y int, color math.Vec3, depth float64)
func (fb *Framebuffer) GetPixel(x, y int) math.Vec3
func (fb *Framebuffer) GetDepth(x, y int) float64
func (fb *Framebuffer) SavePNG(filename string) error
```

## API Example (Conceptual)

```go
// Create framebuffer
fb := framebuffer.New(800, 600)

// Clear to background
fb.Clear(
    math.NewVector3(0.1, 0.1, 0.2), // dark blue
    1.0,                             // far depth
)

// Write pixel with depth test
fb.SetPixel(400, 300,
    math.NewVector3(1, 0, 0), // red
    0.5,                       // depth
)

// Save result
err := fb.SavePNG("output.png")
```

## Common Gotchas

1. **Depth buffer not initialized** - random values cause glitches
2. **Reversed depth test** - farther overwrites closer (wrong!)
3. **Color clamping** - values outside [0,1] need clamping before export
4. **Coordinate flip** - Y-axis direction in screen vs image coordinates
5. **Buffer index out of bounds** - crashes or wrong pixels

## Future Enhancements (Post-MVP)

- Multi-sample anti-aliasing (MSAA) - multiple samples per pixel
- HDR framebuffer - float RGB values preserved for tone mapping
- Multiple render targets - separate buffers for different data
- Stencil buffer - for masking and effects
- Tile-based rendering - process framebuffer in chunks

# Rasterizer Component

**Status:** ✅ **Complete** (Phase 5 — 2026-02-27) — `pkg/rasterize/`, 9 tests, 100% coverage

## Purpose

Convert screen-space triangles into pixels (fragments) with interpolated attributes.

## Core Responsibility

**Input:** Triangle with 3 vertices in screen space (pixel coordinates) + attributes (color, depth)

**Output:** Set of pixel coordinates covered by triangle + interpolated attributes per pixel

## Key Algorithm: Triangle Rasterization

### Scanline Method (Recommended for MVP)

**Concept:** For each horizontal line (scanline) intersecting triangle, find left and right edges, fill pixels between.

**Steps:**
1. Sort vertices by Y coordinate (top to bottom)
2. For each scanline Y from top to bottom:
   - Find X intersection with left edge
   - Find X intersection with right edge
   - Fill pixels from Xleft to Xright
   - Interpolate attributes across scanline

**Pros:**
- Intuitive and straightforward
- Easy to implement
- Good for teaching

**Cons:**
- Slower than modern methods
- Special cases for flat tops/bottoms

### Alternative: Barycentric Coordinates

**Concept:** For each pixel in triangle bounding box, check if inside using barycentric weights.

**Pros:**
- Simpler code structure
- Easier attribute interpolation
- No edge-case handling

**Cons:**
- Tests many pixels outside triangle
- Can be slower for large triangles

**Recommendation:** Either works for MVP. Barycentric is cleaner code.

## Attribute Interpolation

### What to Interpolate
- Depth (Z) - for depth testing
- Color (RGB) - for shading
- Future: normals, UVs, etc.

### Perspective-Correct Interpolation

**Important:** Direct linear interpolation is wrong for 3D!

**Correct approach:**
1. Interpolate 1/Z linearly
2. Interpolate attribute/Z linearly
3. Divide by interpolated 1/Z to get final attribute

**For MVP:** Can start with linear interpolation if perspective incorrectness isn't noticeable.

## Edge Cases to Handle

### Degenerate Triangles
- Zero area (all vertices collinear)
- **Solution:** Skip rendering

### Partial Coverage
- Triangle partially off-screen
- **Solution:** Clip or just skip pixels outside framebuffer

### Thin Triangles
- Very narrow triangles might miss pixels
- **Solution:** Conservative rasterization (round outward) or accept some gaps

### Flat Triangles
- Top or bottom edge is horizontal
- **Solution:** Handle as special case in scanline, or naturally works in barycentric

## Testing Requirements

### Correctness Tests
- Single pixel triangle renders exactly one pixel
- Axis-aligned triangles fill expected rectangle
- Known triangle coordinates produce expected pixel coverage
- No gaps between adjacent triangles (shared edge)
- No overlaps on shared edges (tie-breaking rule)

### Edge Case Tests
- Degenerate triangle (zero area) renders nothing
- Triangle completely off-screen renders nothing
- Triangle partially off-screen clips correctly

### Interpolation Tests
- Color interpolation is correct at vertices
- Color interpolation is correct at triangle center
- Depth interpolation maintains ordering

## Implemented API

```go
// Package: pkg/rasterize

// ScreenVertex represents a vertex in 2D screen space with depth and color attributes.
// X and Y are pixel coordinates; pixel (ix, iy) has its center at (float64(ix)+0.5, float64(iy)+0.5).
// Z is the depth value for depth testing (0=closest, 1=farthest).
// Color holds RGB components, each in [0, 1].
type ScreenVertex struct {
    X, Y  float64
    Z     float64
    Color math.Vec3
}

// Triangle fills the pixels of fb covered by the triangle defined by v0, v1, v2.
// Depth and color are linearly interpolated using barycentric coordinates.
// Degenerate triangles (zero area) are silently skipped.
// Supports both CCW and CW vertex winding order.
func Triangle(v0, v1, v2 ScreenVertex, fb *framebuffer.Framebuffer)
```

**Algorithm:** Barycentric coordinates (cleaner code, winding-agnostic, natural attribute interpolation)

**Test coverage:** 9 tests, 100% coverage
- Degenerate triangle, single-pixel coverage, axis-aligned triangle
- Color at vertex (barycentric weight = 1), color at centroid (= average)
- Depth interpolation with depth test, off-screen and partially off-screen
- CW and CCW winding produce identical coverage

## Performance Considerations

### Bounding Box
Calculate tight bounding box around triangle to minimize pixel tests.

### Early-Out Tests
- Check if triangle completely off-screen
- Check if triangle behind near plane

### Fill Convention
Use consistent tie-breaking rule for pixels exactly on edge (prevents double-drawing).

**Standard:** Top-left rule - pixel centers on top/left edges are inside, others outside.

## Design Decisions

### Pixel Centers
**Recommendation:** Pixel centers at 0.5 offsets (e.g., pixel (0,0) has center at (0.5, 0.5)).

**Why:** Matches OpenGL and most graphics APIs.

### Coordinate Origin
**Recommendation:** (0, 0) at top-left, +X right, +Y down.

**Why:** Matches image formats and framebuffer conventions.

## Common Gotchas

1. **Integer vs float coordinates** - use float for interpolation, convert to int for pixel access
2. **Off-by-one errors** - careful with bounding box edges
3. **Depth precision** - depth values must have sufficient precision (float32 minimum)
4. **Attribute clamping** - interpolated colors might exceed [0,1] range
5. **Shared edge double-draw** - can cause artifacts if not handled

## Future Optimizations (Post-MVP)

- Hierarchical rasterization (tile-based)
- SIMD vectorization (process multiple pixels together)
- Multi-threading (different triangles on different threads)
- Early Z-test (depth test before shading)
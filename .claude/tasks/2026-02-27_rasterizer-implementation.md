# Task Summary: Phase 5 Rasterization Implementation

## Metadata
```yaml
date: 2026-02-27
phase: 5
type: feature
status: complete
branch: claude/phase-5-docs-update-PZENd
```

## Objective

Implement Phase 5 (Rasterization) of the go_toy_renderer: create `pkg/rasterize/` with a triangle rasterizer that converts screen-space vertices into framebuffer pixels with interpolated color and depth attributes.

Also update all project documentation to reflect Phase 5 completion.

## Actions

### Implementation
1. Created `pkg/rasterize/rasterizer_test.go` (TDD Red phase) — 9 tests covering all behaviors
2. Created `pkg/rasterize/rasterizer.go` (TDD Green phase) — full implementation
3. Renamed `RasterizeTriangle` → `Triangle` to satisfy revive linter (package context rule)
4. All tests pass: 9 tests, 100% coverage, 0 lint issues

### Documentation
5. Updated `.claude/docs/06-rasterizer-component.md` — marked complete, added actual API
6. Updated `.claude/docs/12-development-roadmap.md` — Phase 5 COMPLETED, progress 7/14
7. Updated `.claude/docs/10-mvp-features.md` — rasterizer features complete, progress 10/12
8. Updated `.claude/docs/README.md` — Phase 5 complete, next steps Phase 6
9. Updated `CLAUDE.md` — progress summary, recent updates, current project structure

## Key Decisions

### Barycentric vs Scanline Algorithm
**Decision:** Barycentric coordinates.
**Rationale:** Cleaner code structure (no edge case handling for flat tops/bottoms), natural attribute interpolation via weights, inherently winding-order agnostic.

### Function Name: `Triangle` not `RasterizeTriangle`
**Decision:** Named `Triangle` per Go convention (package provides context: `rasterize.Triangle`).
**Rationale:** The `revive` linter flags `rasterize.RasterizeTriangle` as repetitive. Go API convention is to avoid repeating package name in identifier.

### Degenerate Check: `area*area < 1e-16`
**Decision:** Used squared comparison instead of `math.Abs(area) < 1e-8`.
**Rationale:** Avoids importing stdlib `math` package alongside the project's `math` package (naming conflict), and is mathematically equivalent for practical screen-space coordinates.

### Winding Order: Support Both CCW and CW
**Decision:** `insideTriangle(area, w0, w1, w2)` helper checks sign against area sign.
**Rationale:** Screen-space winding depends on Y-axis convention; supporting both makes the rasterizer robust for future pipeline integration.

### Pixel Center Sampling: `float64(ix) + 0.5`
**Decision:** Sample at pixel center (OpenGL convention).
**Rationale:** Standard graphics convention; matches the framebuffer specification.

## Results

### Files Created
- `pkg/rasterize/rasterizer.go` — ScreenVertex type, Triangle(), edgeFunction, insideTriangle, minF, maxF
- `pkg/rasterize/rasterizer_test.go` — 9 tests, 100% coverage

### Test Results
```
ok  github.com/muddl/go_toy_renderer/pkg/rasterize  coverage: 100.0% of statements
ok  github.com/muddl/go_toy_renderer/pkg/camera     coverage: 100.0% of statements
ok  github.com/muddl/go_toy_renderer/pkg/framebuffer  coverage: 90.2% of statements
ok  github.com/muddl/go_toy_renderer/pkg/geometry   coverage: 100.0% of statements
ok  github.com/muddl/go_toy_renderer/pkg/math       coverage: 96.1% of statements
```

### Lint: 0 issues (golangci-lint v2)

### Files Updated
- `.claude/docs/06-rasterizer-component.md`
- `.claude/docs/12-development-roadmap.md`
- `.claude/docs/10-mvp-features.md`
- `.claude/docs/README.md`
- `CLAUDE.md`

## Next Steps

- **Phase 6 (Shading):** Implement `pkg/shader/` with shader interface and vertex color pass-through
- The rasterizer's `Triangle()` function accepts `ScreenVertex` with `Color` already interpolated, so Phase 6 shading will be applied before passing to `Triangle` (or the color becomes the "shaded" output)
- Phase 7 (Render Pipeline) will connect camera → geometry → rasterizer → framebuffer end-to-end

## Lessons Learned

1. **Revive package-name rule:** Go convention avoids repeating the package name in exported identifiers. `rasterize.Triangle` is cleaner than `rasterize.RasterizeTriangle`.
2. **Squared comparison for zero check:** `area*area < 1e-16` elegantly avoids importing stdlib math while achieving the same result as `|area| < 1e-8`.
3. **Winding order matters for rasterizer:** Supporting both CCW and CW in the inside test makes the component robust and avoids bugs later when world-space CCW triangles become CW in screen space (due to Y-axis flip).
4. **Pixel center precision:** Placing vertices exactly at pixel centers in tests (e.g., v0 at 1.5, 1.5) creates deterministic test cases where barycentric weight equals exactly 1.0 at the vertex, enabling precise color interpolation tests.

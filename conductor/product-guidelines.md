# Product Guidelines — go_toy_renderer

## Voice and Tone

**Technical and precise.** Favour accuracy over friendliness; assume the reader knows Go. Explain graphics concepts where the "why" matters (e.g. why column-major matrices, why right-handed coords), but don't over-explain language basics.

## Design Principles

1. **Correctness first, performance second.** Get the math right before optimising. Benchmarks come after passing tests.
2. **Simplicity over abstraction.** Don't introduce interfaces, generics, or helper utilities that serve only one use. Three similar lines of code is better than a premature abstraction.
3. **Testability at every layer.** Each pipeline stage must be independently testable with known inputs and expected outputs.
4. **No over-engineering beyond current phase.** Only implement what the current phase requires. Future phases will add what they need.
5. **TDD strictly enforced.** Failing test committed before any implementation. Red-Green-Refactor is the only workflow.
6. **Separation of concerns across pipeline stages.** Each package owns one stage: math, geometry, camera, rasterize, shader, render. Cross-package dependencies flow in one direction.
7. **Right-handed coordinate system.** +X=Right, +Y=Up, +Z=Out (OpenGL style). Never silently deviate.
8. **Column-major matrices.** Multiply on right: `result = matrix × vector`. Document any exception explicitly.

## Process Standards

- All work goes through feature branches and pull requests — **never commit directly to `main`**.
- Each phase is complete only when: implementation done, PR merged, task summary written, and all docs updated.
- GPU phases target cross-platform support (D3D12/Metal/Vulkan via `go-webgpu/webgpu` Zero-CGo FFI bindings to wgpu-native).
- HLSL shaders are compiled to WGSL via naga (not hand-written WGSL).

## Common Gotchas

Quick-reference for known pitfalls. See `conductor/architecture.md` for per-package detail.

### Math

- **Matrix multiply order** — easy to reverse; `result = matrix × vector` (column-major, right-multiply).
- **Perspective divide by zero** — reject any vertex where `W ≤ 0` before dividing.
- **Float precision in Normalize** — zero-length vector will produce NaN; guard if needed.
- **Left vs right-handed** — this project is right-handed (+Z out of screen); don't silently mix conventions.

### Transforms

- **VP matrix order** — `Projection.Multiply(View)`, NOT `View.Multiply(Projection)`.
- **Forgetting perspective divide** — NDC conversion requires `x/w, y/w, z/w` after clip.
- **Depth mapping** — `(ndcZ + 1) / 2` maps NDC z ∈ [-1,1] to depth ∈ [0,1]; depth 0=near, 1=far.
- **Y-axis flip** — NDC +Y is up; screen +Y is down. Use `(1 - ndcY) / 2 * height`.

### Rasterization

- **Integer vs float coordinates** — use float for interpolation, convert to int for pixel access.
- **Pixel center offset** — pixel `(ix, iy)` has center at `(float64(ix)+0.5, float64(iy)+0.5)`.
- **Color clamping** — interpolated values can exceed [0,1]; framebuffer clamps on export only.
- **Degenerate triangles** — area² < 1e-16; silently skip rather than panic.

### Pipeline

- **Triangles behind camera** — any vertex with `w ≤ 0` means reject the whole triangle.
- **Depth buffer initialisation** — must be 1.0 (far plane) before each frame; `Clear(color, 1.0)`.
- **Backface winding** — winding order affects culling; current rasterizer is winding-agnostic (both CCW and CW rendered).
- **Matrix combination order** — `MVP = ProjectionMatrix × ViewMatrix × ModelMatrix`.

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
- GPU phases target cross-platform support (D3D12/Metal/Vulkan via wgpu-native).
- HLSL shaders are compiled to WGSL via naga (not hand-written WGSL).

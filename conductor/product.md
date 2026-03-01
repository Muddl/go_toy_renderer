# Product Definition — go_toy_renderer

## Project Name

go_toy_renderer

## Description

A toy 3D renderer implemented in Go, progressing from a CPU software renderer to GPU-accelerated real-time rendering via WebGPU.

## Problem Statement

Understanding 3D graphics from first principles — implementing the full pipeline (CPU → GPU) rather than using a black-box engine.

## Target Users

The author — a developer learning 3D graphics and GPU programming through hands-on implementation.

## Key Goals

1. Learn the full 3D pipeline from scratch: math primitives → rasterization → GPU shaders.
2. Produce working, well-tested code at each phase; no phase is "done" until tests pass and docs are updated.
3. Build toward real-time GPU rendering via WebGPU (go-webgpu/webgpu Zero-CGo FFI over wgpu-native) with HLSL shaders compiled to WGSL.

## Current Status

- **MVP Complete (Phases 0–8):** CPU software renderer renders a colored cube to a 640×480 PNG.
- **GPU Roadmap Active (Phases 9–16):** Real-time windowed rendering, WebGPU backend, HLSL shader pipeline.

## Success Criteria

- **MVP (done):** Render a colored 3D object with correct perspective and depth ordering to a 640×480 PNG.
- **Phase 9:** Real-time window at 60 fps via GLFW with CPU blit.
- **Phase 11:** Hello Triangle rendered via go-webgpu/webgpu (WebGPU, Zero-CGo FFI).
- **Phase 13:** Full HLSL shader pipeline (vertex + fragment) running on the GPU.
- **Phase 16:** Textured OBJ model rendered in real time with Phong/PBR shading.

## Pipeline Overview

```
┌──────────────┐
│ Scene Setup  │  Define geometry, camera, lights
└──────┬───────┘
       │
┌──────▼───────┐
│ Vertex Stage │  Transform vertices (Local → World → View → Clip → NDC → Screen)
└──────┬───────┘
       │
┌──────▼───────┐
│ Primitive    │  Assemble triangles, cull backfaces, reject W≤0
│ Assembly     │
└──────┬───────┘
       │
┌──────▼───────┐
│ Rasterizer   │  Convert triangles to pixels, interpolate attributes
└──────┬───────┘
       │
┌──────▼───────┐
│ Fragment     │  Calculate per-pixel color (shading)
│ Shader       │
└──────┬───────┘
       │
┌──────▼───────┐
│ Framebuffer  │  Write pixels with depth testing (depth < current)
└──────┬───────┘
       │
┌──────▼───────┐
│ Image Output │  Save framebuffer to PNG / blit to window
└──────────────┘
```

**Vertex coordinate transform pipeline:**
```
clip  = mvp.MultiplyVec4(px, py, pz, 1.0)   // clip space
ndcX  = cx/cw, ndcY = cy/cw, ndcZ = cz/cw  // NDC [-1, 1]
depth = (ndcZ + 1) / 2                       // [0,1]: 0=near, 1=far
sX    = (ndcX + 1) / 2 * width              // screen X (px)
sY    = (1 - ndcY) / 2 * height             // screen Y (Y flipped)
```

## Coordinate System & Matrix Conventions

| Convention | Value | Rationale |
|-----------|-------|-----------|
| Handedness | Right-handed | Standard in math and OpenGL; +X=Right, +Y=Up, +Z=Out |
| Matrix storage | Column-major | Matches mathematical notation and GPU conventions |
| Multiply order | `result = matrix × vector` | Multiply on right |
| VP matrix | `Projection × View` | Pre-multiplied; applied per-vertex |
| Depth buffer | 0.0 = near, 1.0 = far | OpenGL convention; `(ndcZ + 1) / 2` mapping |
| Depth test | `depth < current` (closer wins) | Standard Z-buffer |
| Triangle winding | Counter-clockwise (front face) | Enables backface culling |
| Screen origin | Top-left (0,0); +Y down | Matches image format conventions |
| Pixel center | `(ix + 0.5, iy + 0.5)` | Matches OpenGL pixel-center convention |

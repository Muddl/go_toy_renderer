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
3. Build toward real-time GPU rendering via WebGPU (wgpu-native) with HLSL shaders compiled to WGSL.

## Current Status

- **MVP Complete (Phases 0–8):** CPU software renderer renders a colored cube to a 640×480 PNG.
- **GPU Roadmap Active (Phases 9–16):** Real-time windowed rendering, WebGPU backend, HLSL shader pipeline.

## Success Criteria

- **MVP (done):** Render a colored 3D object with correct perspective and depth ordering to a 640×480 PNG.
- **Phase 9:** Real-time window at 60 fps via GLFW with CPU blit.
- **Phase 11:** Hello Triangle rendered via wgpu-native (WebGPU).
- **Phase 13:** Full HLSL shader pipeline (vertex + fragment) running on the GPU.
- **Phase 16:** Textured OBJ model rendered in real time with Phong/PBR shading.

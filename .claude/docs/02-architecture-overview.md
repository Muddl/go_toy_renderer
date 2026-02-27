# Architecture Overview

## High-Level Design Philosophy

**Separation of Concerns:** Each stage of the rendering pipeline should be isolated and testable.

**Data Flow:** Geometry → Transformations → Projection → Rasterization → Framebuffer → Image Output

**Simplicity First:** Optimize for readability and correctness, not performance (for MVP).

## Core Pipeline Stages

```
┌──────────────┐
│ Scene Setup  │  Define geometry, camera, lights
└──────┬───────┘
       │
┌──────▼───────┐
│ Vertex Stage │  Transform vertices (Model → World → View → Clip → NDC → Screen)
└──────┬───────┘
       │
┌──────▼───────┐
│ Primitive    │  Assemble triangles, cull backfaces, clip to view frustum
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
│ Framebuffer  │  Write pixels with depth testing
└──────┬───────┘
       │
┌──────▼───────┐
│ Image Output │  Save framebuffer to file
└──────────────┘
```

## Package Organization

### Recommended Structure

```
pkg/
├── math/           # Vector, matrix, transforms
├── geometry/       # Mesh, vertex, triangle types
├── camera/         # Camera and view transforms
├── render/         # Core rendering pipeline
├── rasterize/      # Triangle rasterization
├── shader/         # Shading calculations
└── framebuffer/    # Pixel buffer with depth test

cmd/
└── renderer/       # Main application
```

## Key Design Decisions

### 1. Coordinate Systems

**Right-handed coordinate system:**
- +X = Right
- +Y = Up
- +Z = Out of screen (toward viewer)

**Why:** Standard in mathematics and many 3D tools (OpenGL style)

### 2. Matrix Convention

**Column-major matrices** with column vectors (multiply on right):
```
result = matrix * vector
```

**Why:** Matches mathematical notation and GPU conventions

### 3. Triangle Winding

**Counter-clockwise** front faces (when viewed from front)

**Why:** Standard convention, makes backface culling straightforward

### 4. Color Representation

**Float-based** during calculation, convert to bytes at output

**Why:** Avoids precision loss during shading calculations

### 5. Depth Buffering

**Standard Z-buffer:** 0.0 = near plane, 1.0 = far plane

**Why:** Matches OpenGL convention and the implemented `(ndcZ + 1) / 2` depth mapping. Depth buffer initialised to 1.0 (far) and closer pixels overwrite with the `depth < current` test.

## Data Flow Between Components

### Vertex Data Flow
```
Mesh (local space)
  → Model Matrix
  → World Space
  → View Matrix
  → Camera Space
  → Projection Matrix
  → Clip Space
  → Perspective Divide
  → NDC Space
  → Viewport Transform
  → Screen Space
```

### Per-Triangle Flow
```
3 Vertices (screen space + attributes)
  → Backface Culling (optional)
  → Clipping (if needed)
  → Rasterization (find pixels covered)
  → For each pixel:
      → Interpolate attributes (depth, color, etc.)
      → Depth test
      → Shade fragment
      → Write to framebuffer
```

## Testing Strategy Philosophy

- **Unit test math operations** extensively (they're the foundation)
- **Integration test pipeline stages** with known inputs/outputs
- **Golden image tests** for end-to-end validation
- **Benchmark critical paths** (rasterization, shading)

## Extensibility Considerations

While keeping MVP simple, design should allow for:

1. **Different shaders** - Interface for fragment shading
2. **Multiple geometries** - Scene can hold multiple meshes
3. **Parallel rendering** - Framebuffer regions can be rendered independently
4. **Additional attributes** - Vertex structure can expand (normals, UVs, etc.)

## GPU Architecture (Phases 9–16)

The CPU software renderer is the learning foundation. The GPU roadmap layers a hardware-accelerated backend on top of it without replacing the existing code.

### Backend Abstraction

```
┌─────────────────────────────┐
│         Renderer Interface   │  pkg/renderer/renderer.go
│  Render(scene, cam) error   │
└──────────┬──────────────────┘
           │
    ┌──────┴──────┐
    │             │
┌───▼───┐   ┌────▼────┐
│  CPU  │   │   GPU   │
│Renderer│  │Renderer │
└───────┘   └─────────┘
pkg/render  pkg/gpu
(existing)  (new)
```

### GPU Pipeline (WebGPU / wgpu-native)

```
┌──────────────┐
│  Scene Setup  │  geometry.Mesh → GPUScene (vertex/index buffers in VRAM)
└──────┬───────┘
       │
┌──────▼───────┐
│ Vertex Shader│  HLSL → WGSL; runs on GPU; applies MVP matrix from uniform buffer
│  (HLSL/WGSL) │
└──────┬───────┘
       │
┌──────▼───────┐
│ HW Rasterizer│  GPU hardware; barycentric interpolation handled automatically
│  (WebGPU)    │
└──────┬───────┘
       │
┌──────▼───────┐
│ Fragment     │  HLSL → WGSL; per-pixel shading (color, lighting, texture)
│ Shader       │
└──────┬───────┘
       │
┌──────▼───────┐
│ HW Depth Test│  GPU depth buffer; 32-bit float precision
└──────┬───────┘
       │
┌──────▼───────┐
│ Swap Chain   │  GPU back buffer; presented to window via D3D12/Metal/Vulkan
└──────────────┘
```

### Shader Authoring Flow

```
HLSL source (.hlsl)
       │
  go generate
       │
    DXC/naga
       │
WGSL output (.wgsl)   ← committed to repo
       │
  wgpu runtime
       │
  GPU execution
```

### Cross-Platform Backend Selection

| Platform | wgpu-native Backend | Notes |
|----------|-------------------|-------|
| Windows  | D3D12             | Preferred, then Vulkan |
| macOS    | Metal             | Only option |
| Linux    | Vulkan            | Preferred, then OpenGL |
| Web      | WebGPU (browser)  | Future — WASM target |

### Package Layout (Post-GPU)

```
pkg/
├── math/           # Vec3, Mat4x4 (unchanged)
├── geometry/       # Mesh, Vertex (+ Normal, UV in Phase 15/16)
├── camera/         # Camera, view/projection matrices (unchanged)
├── render/         # CPU pipeline (preserved as reference)
├── rasterize/      # CPU rasterizer (preserved)
├── shader/         # CPU shader functions (preserved)
├── framebuffer/    # CPU framebuffer (preserved)
├── renderer/       # NEW: Renderer interface + factory
├── gpu/            # NEW: WebGPU device, swapchain, pipeline, buffers
├── scene/          # NEW: Transform component, GPUScene
└── loader/         # NEW (Phase 16): OBJ file loading

cmd/
├── renderer/       # PNG output (unchanged)
└── renderer-rt/    # NEW: Real-time windowed renderer
```

## What's NOT in Architecture (For Now)

- Scene graph hierarchy (deferred to Phase 14+ with Transform component)
- Material system (deferred to Phase 16)
- Texture sampling (deferred to Phase 16)
- Shadow mapping (deferred to Phase 16 stretch goal)
- Post-processing effects (deferred to Phase 16)
- Animation/skinning (not planned)
- Resource management beyond basic mesh loading (deferred to Phase 16)
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

**Reverse Z-buffer** (1.0 = near, 0.0 = far) optional but recommended

**Why:** Better floating-point precision distribution (can start with standard 0=near, 1=far)

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

## What's NOT in Architecture (For Now)

- Scene graph hierarchy
- Material system
- Texture sampling
- Shadow mapping
- Post-processing effects
- Animation/skinning
- Resource management beyond basic mesh loading
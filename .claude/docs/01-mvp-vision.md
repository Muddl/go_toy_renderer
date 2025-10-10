# MVP Vision & Scope

## What We're Building

A **toy 3D software renderer** that demonstrates fundamental 3D graphics concepts without relying on GPU acceleration. This is a learning project that renders 3D scenes to a 2D framebuffer using CPU-only calculations.

## Why This Project Exists

**Learning Goals:**
- Understand 3D graphics pipeline from first principles
- Master coordinate transformations and perspective projection
- Implement rasterization algorithms
- Learn how GPUs work by doing it on CPU

**Non-Goals (Out of Scope for MVP):**
- Production-quality performance
- GPU acceleration
- Advanced lighting models (PBR, GI, etc.)
- Animation system
- Asset pipeline/editor

## Success Criteria for MVP

The MVP is successful when it can:

1. **Load a simple 3D model** (hardcoded or from OBJ file)
2. **Transform and project** vertices from 3D world space to 2D screen space
3. **Rasterize triangles** with basic shading
4. **Output an image file** showing the rendered scene
5. **Handle basic camera** positioning and orientation

## Visual Target

By the end of MVP, the renderer should produce an image of:
- A simple 3D object (cube, tetrahedron, or simple mesh)
- With flat or per-vertex shading
- From an adjustable camera position
- Saved as PNG/BMP file

**Example output quality target:** Similar to early 3D games (think Wolf3D/Doom era), not modern AAA graphics.

## Constraints

- **CPU-only rendering** - No OpenGL/Vulkan/DirectX
- **Single-threaded for MVP** - Concurrency can come later
- **Minimal dependencies** - Standard library preferred, image encoding library acceptable
- **Simple data formats** - OBJ files or hardcoded geometry
- **Fixed resolution** - 640x480 or similar is fine

## What Makes This "Bare Bones"

We're building the **absolute minimum** to demonstrate:
- 3D math works correctly
- Projection makes sense
- Triangles appear on screen
- Basic shading is visible

Everything else (textures, advanced lighting, shadows, etc.) comes after MVP.
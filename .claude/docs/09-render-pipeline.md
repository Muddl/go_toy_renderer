# Render Pipeline Component

## Purpose

Orchestrate the entire rendering process from scene to image output.

## Pipeline Overview

The render pipeline ties all components together in the correct order.

## Main Render Function Flow

```go
func Render(scene *Scene, camera *Camera, framebuffer *Framebuffer, shader ShaderFunc)
```

### Step-by-Step Pipeline

#### 1. Setup Phase
```
- Clear framebuffer (color + depth)
- Get view-projection matrix from camera
- Prepare for rendering
```

#### 2. Per-Mesh Processing
```
For each mesh in scene:
    - Get model matrix (mesh transform)
    - Compute MVP matrix = Projection × View × Model
```

#### 3. Vertex Processing (Vertex Shader Stage)
```
For each vertex in mesh:
    - Transform position: MVP × vertex.Position
    - Perform perspective divide (x/w, y/w, z/w)
    - Convert to screen coordinates (NDC → pixel coords)
    - Store transformed vertex with attributes
```

#### 4. Primitive Assembly
```
For each triangle (group of 3 indices):
    - Fetch transformed vertices
    - Optional: Backface culling
    - Optional: Clipping (if needed)
```

#### 5. Rasterization
```
For each triangle:
    - Rasterize to pixels
    - For each covered pixel:
        - Interpolate attributes (depth, color, etc.)
        - Call fragment shader
        - Depth test and write to framebuffer
```

#### 6. Output Phase
```
- Save framebuffer to image file
```

## Scene Representation

**For MVP, simple structure:**

```go
type Scene struct {
    Meshes []*Mesh
    // Future: lights, skybox, etc.
}
```

**Each mesh has:**
- Geometry (vertices, indices)
- Transform (model matrix - position, rotation, scale)

## Transform Pipeline Details

### Coordinate Space Transformations

```
Local Space (mesh vertices)
    ↓ Model Matrix
World Space
    ↓ View Matrix
Camera/View Space
    ↓ Projection Matrix
Clip Space
    ↓ Perspective Divide (x/w, y/w, z/w, w=1)
NDC Space (-1 to 1)
    ↓ Viewport Transform
Screen Space (pixel coordinates)
```

### Matrix Combination

**Efficient approach:** Pre-multiply matrices per mesh, not per vertex.

```go
mvp := projection.Multiply(view).Multiply(model)

// Then for each vertex:
clipSpace := mvp.TransformVector(vertex.Position)
```

## Optional Optimizations (Not Required for MVP)

### Backface Culling

**Purpose:** Don't render triangles facing away from camera.

**Method:**
1. Compute triangle normal in screen space
2. If normal points away from viewer (Z > 0): skip triangle

**Benefit:** ~50% fewer triangles to rasterize

**MVP:** Can skip, render all triangles.

### Frustum Clipping

**Purpose:** Handle triangles partially outside view frustum.

**Method:**
- Clip triangle against frustum planes
- Generate new triangles from clipped result

**MVP:** Can skip if you handle off-screen coordinates in rasterizer.

### Near Plane Clipping

**Important:** Triangles behind camera (negative W) must be clipped!

**Method:** Clip against W = 0 plane before perspective divide.

**MVP:** Minimum - reject triangles with all vertices W ≤ 0.

## Testing Requirements

### End-to-End Tests
- Render simple scene (cube), verify image produced
- Camera movement changes output image
- Different meshes produce different output
- Transformation matrices applied correctly

### Transform Tests
- Vertex at origin stays at origin with identity transform
- Translation moves vertices correctly
- Rotation rotates vertices correctly
- Full MVP chain produces expected screen coordinates

### Pipeline Integration Tests
- All components called in correct order
- Data flows correctly between stages
- Framebuffer contains expected pixel data after render

## API Example (Conceptual)

```go
// Create scene
scene := render.NewScene()
cube := geometry.CreateCube()
cube.SetPosition(0, 0, 0)
scene.AddMesh(cube)

// Create camera
camera := camera.New(
    math.NewVector3(0, 0, 5),
    math.NewVector3(0, 0, 0),
    math.NewVector3(0, 1, 0),
)
camera.SetPerspective(45, 800.0/600.0, 0.1, 100)

// Create framebuffer
fb := framebuffer.New(800, 600)

// Create shader
shader := shader.VertexColorShader

// Render!
renderer := render.NewRenderer()
renderer.Render(scene, camera, fb, shader)

// Save output
fb.SavePNG("output.png")
```

## Common Gotchas

1. **Matrix multiplication order** - Easy to reverse, hard to debug
2. **Perspective divide by zero** - W = 0 or negative
3. **Depth range** - Make sure depth maps correctly to [0, 1]
4. **Attribute interpolation** - Must be perspective-correct for 3D
5. **Winding order** - Affects backface culling result
6. **Transform not applied** - Forgot to multiply by model matrix

## Design Decisions

### Renderer State

**Stateless (recommended for MVP):**
```go
func Render(scene, camera, framebuffer, shader)
```

**Stateful:**
```go
renderer := NewRenderer()
renderer.SetCamera(camera)
renderer.SetShader(shader)
renderer.Render(scene)
```

**Recommendation:** Stateless is simpler for MVP.

### Error Handling

**Critical errors:**
- Invalid framebuffer dimensions
- Null mesh/camera
- Invalid matrix (NaN, infinity)

**Non-critical:**
- Degenerate triangles (skip silently)
- Triangles off-screen (skip or clip)

## Performance Considerations (Post-MVP)

- Pre-transform vertices once per mesh, not per triangle
- Batch triangles for cache efficiency
- Parallel rasterization (different triangles on different threads)
- Early Z-culling (depth test before shading)
- Bounding volume culling (skip entire meshes outside frustum)

## MVP Success Criteria

Pipeline is complete when:
1. Can render hardcoded mesh (cube/tetrahedron)
2. Transformations work (rotate, translate, scale)
3. Camera view changes affect output
4. Perspective looks correct (far = small, near = large)
5. Depth test prevents Z-fighting
6. Output image saved successfully

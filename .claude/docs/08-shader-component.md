# Shader Component

## Purpose

Calculate final pixel color based on geometry attributes and lighting (if any).

## Shader Concept

A shader is a function that takes interpolated vertex attributes and produces a final color.

**For MVP:** Shaders can be very simple - no complex lighting required.

## Shader Types for MVP

### 1. Flat Color Shader (Simplest)

**Input:** None (or vertex position)

**Output:** Constant color

**Use case:** Testing, debug visualization

```go
func FlatColorShader(attributes Attributes) Vector3 {
    return NewVector3(1, 0, 0) // always red
}
```

### 2. Vertex Color Shader (Recommended for MVP)

**Input:** Interpolated vertex color

**Output:** That color

**Use case:** Demonstrates interpolation working correctly

```go
func VertexColorShader(attributes Attributes) Vector3 {
    return attributes.Color // pass through
}
```

### 3. Depth Shader (Debug)

**Input:** Fragment depth

**Output:** Depth visualized as grayscale

**Use case:** Verify depth is interpolating correctly

```go
func DepthShader(attributes Attributes) Vector3 {
    gray := attributes.Depth
    return NewVector3(gray, gray, gray)
}
```

### 4. Simple Diffuse Lighting (Optional MVP)

**Input:** Interpolated normal, light direction

**Output:** Color modulated by lighting

**Use case:** More realistic rendering

```go
func DiffuseLightingShader(attributes Attributes) Vector3 {
    lightDir := NewVector3(0, 0, 1).Normalize()
    normal := attributes.Normal.Normalize()

    ndotl := max(normal.Dot(lightDir), 0)

    return attributes.Color.Scale(ndotl)
}
```

**Note:** Requires normal vectors in vertex data (not in minimal MVP).

## Shader Interface Design

### Option 1: Function Type

```go
type ShaderFunc func(Attributes) Vector3

type Attributes struct {
    Position Vector3
    Color    Vector3
    Depth    float64
    // Future: Normal, UV, etc.
}
```

**Pros:** Simple, flexible, easy to test
**Cons:** Can't hold state

### Option 2: Interface

```go
type Shader interface {
    Shade(attributes Attributes) Vector3
}
```

**Pros:** Can hold shader parameters (uniforms)
**Cons:** More boilerplate

**Recommendation:** Start with function type (Option 1).

## Attributes to Interpolate

**Minimum for MVP:**
- Color (Vector3)
- Depth (float64) - needed for depth test

**Add later:**
- Normal (Vector3) - for lighting
- UV (Vector2) - for textures
- Position (Vector3) - for effects

## Testing Requirements

### Shader Output Tests
- Flat shader returns constant color
- Vertex color shader returns input color unchanged
- Depth shader converts depth to grayscale correctly
- Colors are in valid [0, 1] range (or clamped)

### Integration Tests
- Different shaders produce different framebuffer output
- Shader receives correctly interpolated attributes
- Shader is called once per visible pixel

## Design Considerations

### Shader Uniforms

**Uniforms:** Parameters constant across all pixels (e.g., light position, time).

**For MVP:** Not needed, can hardcode in shader function.

**Post-MVP:** Add uniforms struct:
```go
type Uniforms struct {
    LightPosition Vector3
    LightColor    Vector3
    Time          float64
}

type ShaderFunc func(Attributes, Uniforms) Vector3
```

### Shader Selection

How does renderer know which shader to use?

**Options:**
1. Global shader set on renderer
2. Per-mesh shader assignment
3. Material system (overkill for MVP)

**Recommendation for MVP:** Single global shader, easy to swap.

## API Example (Conceptual)

```go
// Define a shader
myShader := func(attr shader.Attributes) math.Vector3 {
    // Simple vertex color pass-through
    return attr.Color
}

// Set on renderer
renderer.SetShader(myShader)

// Render with that shader
renderer.Render(scene, camera)
```

## Common Gotchas

1. **Forgetting to normalize vectors** - normals, light directions
2. **Color overflow** - lighting can produce values > 1.0, need clamping
3. **Negative lighting** - dot product can be negative, need max(0, value)
4. **Uninitialized attributes** - make sure all attributes are set
5. **Expensive shaders** - called per-pixel, performance critical

## Future Enhancements (Post-MVP)

### Lighting Models
- Phong lighting (diffuse + specular)
- Blinn-Phong (faster, similar)
- Physically-based rendering (PBR)

### Textures
- Texture sampling in shader
- Normal mapping
- Environment mapping

### Advanced Features
- Multiple lights
- Shadows (shadow mapping)
- Transparency (alpha blending)
- Post-processing (bloom, tone mapping)

## Recommended MVP Progression

1. **Start:** Flat color shader (verify rendering works)
2. **Next:** Vertex color shader (verify interpolation works)
3. **Optional:** Depth shader (debug depth buffer)
4. **Nice to have:** Simple diffuse lighting (looks better)

**Don't over-complicate:** MVP should prove the pipeline works, not be production-quality.

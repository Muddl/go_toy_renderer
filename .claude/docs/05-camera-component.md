# Camera Component

## Purpose

Define viewpoint and create view and projection transformations.

## Core Type: Camera

**Represents:** A virtual camera in 3D space

**Data needed:**
- `Position Vector3` - where camera is located
- `Target Vector3` - where camera looks at
- `Up Vector3` - camera orientation (usually (0,1,0))
- `FOV float64` - field of view in degrees
- `Aspect float64` - width/height ratio
- `Near float64` - near clipping plane distance
- `Far float64` - far clipping plane distance

## Key Operations

### View Matrix
```
Create view matrix using LookAt algorithm:
- Converts world space to camera space
- Camera is at origin, looking down -Z
```

**Formula:** Typically implemented as:
1. Calculate basis vectors (forward, right, up)
2. Build rotation matrix
3. Build translation matrix
4. Combine (rotation × translation)

### Projection Matrix
```
Create perspective projection matrix:
- Converts camera space to clip space
- Applies perspective distortion (far things smaller)
- Maps to [-1, 1] cube (NDC after perspective divide)
```

**Parameters:**
- Vertical FOV (45-90 degrees typical)
- Aspect ratio (screen width / height)
- Near plane (e.g., 0.1)
- Far plane (e.g., 1000.0)

**Important:** Near > 0 (can't be zero or negative).

### View-Projection Matrix
```
Combine view and projection for efficiency:
VP = Projection × View
```

Can pre-multiply to save computation per vertex.

## Camera Controls (Optional for MVP)

**Nice to have but not required:**
- Orbit around target
- FPS-style movement
- Zoom in/out

**For MVP:** Hardcoded camera position is acceptable.

## Testing Requirements

### View Matrix Tests
- Camera at origin looking down -Z produces identity-like matrix
- Camera position affects translation component
- LookAt produces correct basis vectors
- Up vector influences rotation

### Projection Matrix Tests
- FOV affects perspective strength
- Aspect ratio affects X/Y scaling
- Near/far planes map to correct NDC depth range
- Points at near plane map to Z=-1 (or 1 for reverse-Z)
- Points at far plane map to Z=1 (or 0 for reverse-Z)

### Integration Tests
- Point in front of camera projects to visible screen space
- Point behind camera has negative W (clipped)
- Point outside FOV projects outside NDC cube

## API Example (Conceptual)

```go
// Create camera
camera := camera.New(
    math.NewVector3(0, 0, 5),    // position
    math.NewVector3(0, 0, 0),    // target
    math.NewVector3(0, 1, 0),    // up
)

// Set projection parameters
camera.SetPerspective(
    45.0,           // FOV in degrees
    800.0 / 600.0,  // aspect ratio
    0.1,            // near
    100.0,          // far
)

// Get transformation matrices
viewMatrix := camera.ViewMatrix()
projMatrix := camera.ProjectionMatrix()
viewProjMatrix := camera.ViewProjectionMatrix()
```

## Design Considerations

### Camera Space Convention
**Standard:** Camera looks down **-Z axis**, with +Y up and +X right.

**Why:** Matches OpenGL convention, right-handed system.

### FOV Vertical vs Horizontal
**Recommendation:** Use **vertical FOV**.

**Why:**
- More common in 3D graphics
- Aspect ratio adjusts horizontal FOV automatically

### Near/Far Plane Values
**Typical values:**
- Near: 0.1 to 1.0
- Far: 100.0 to 1000.0

**Trade-off:**
- Large far/near ratio = worse depth precision
- Start with 100:1 ratio for testing

## Common Gotchas

1. **Near plane = 0** - causes divide by zero in projection
2. **Far < Near** - invalid, causes incorrect depth
3. **Aspect ratio mismatch** - stretches image
4. **FOV too wide** - causes extreme distortion (>90° is unusual)
5. **Up vector parallel to view direction** - undefined LookAt
6. **Flipped up vector** - camera appears upside down

## Future Enhancements (Post-MVP)

- Orthographic projection (for UI, CAD-style views)
- Camera frustum extraction (for culling)
- Camera ray casting (for picking)
- Multiple camera support (picture-in-picture)
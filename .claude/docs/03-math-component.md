# Math Component

## Purpose

Provide foundational 3D math operations: vectors, matrices, and transformations.

## Core Types

### Vector3
**Represents:** 3D point or direction

**Data:**
- `X, Y, Z float64`

**Operations needed:**
- Add, Subtract, Scale
- Dot product
- Cross product
- Length, Normalize
- Distance between points

### Vector4
**Represents:** Homogeneous coordinates (for matrix transforms)

**Data:**
- `X, Y, Z, W float64`

**Operations needed:**
- Convert from/to Vector3
- Basic arithmetic

### Matrix4x4
**Represents:** 4x4 transformation matrix

**Data:**
- `[16]float64` or `[4][4]float64` (choose one for consistency)

**Operations needed:**
- Identity matrix
- Matrix multiplication
- Transform vector (multiply)
- Inverse (needed for camera)
- Transpose

## Key Transformation Matrices

### Translation
```
Move point by offset (dx, dy, dz)
```

### Rotation
**For MVP, implement at least one:**
- Rotation around arbitrary axis
- OR: Separate X, Y, Z axis rotations

### Scale
```
Uniform or non-uniform scaling
```

### Look-At Matrix
```
Create view matrix from:
- Eye position (camera location)
- Target position (where camera looks)
- Up vector (camera orientation)
```

### Perspective Projection
```
Parameters:
- Field of view (FOV)
- Aspect ratio
- Near plane
- Far plane

Output: Clip space coordinates
```

### Viewport Transform
```
Convert NDC (-1 to 1) to screen coordinates (0 to width/height)

Parameters:
- Screen width
- Screen height
```

## Design Principles

**Immutability preferred:** Operations return new values rather than modifying in place (for thread safety and clarity).

**Small and focused:** Each type should do one thing well.

**Well-tested:** Math bugs are the hardest to debug visually.

## Testing Requirements

### Must Have
- Vector addition/subtraction correctness
- Cross product right-hand rule validation
- Matrix multiplication associativity
- Identity matrix behavior
- Known transformation results (e.g., 90° rotation)
- Perspective projection edge cases (near/far plane)

### Nice to Have
- Matrix inverse correctness (multiply by inverse = identity)
- Numerical stability tests
- Benchmark performance of hot paths

## Common Gotchas to Test

1. **Gimbal lock** in rotation composition
2. **Perspective divide by zero** (W = 0)
3. **Floating point precision** in normalize operations
4. **Matrix multiplication order** (easy to reverse)
5. **Left vs right-handed** coordinate system consistency

## API Example (Conceptual)

```go
// Vector operations
v1 := math.NewVector3(1, 0, 0)
v2 := math.NewVector3(0, 1, 0)
cross := v1.Cross(v2) // Should be (0, 0, 1)

// Matrix creation
model := math.TranslationMatrix(2, 0, 0)
view := math.LookAt(eye, target, up)
proj := math.Perspective(fov, aspect, near, far)

// Transformation
mvp := proj.Multiply(view).Multiply(model)
transformed := mvp.TransformVector(vertex)
```

## Performance Considerations (Post-MVP)

- Matrix operations are hot paths
- Consider SIMD-friendly layouts later
- Cache transform chains (don't recompute each frame)
- Benchmark before optimizing
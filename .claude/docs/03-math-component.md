# Math Component

## Status: Phase 1 Complete ✅

**Implemented:** Vec3, Mat4x4, basic transformations
**Remaining:** Advanced camera matrices (Phase 3)

## Purpose

Provide foundational 3D math operations: vectors, matrices, and transformations.

## Core Types

### Vec3 ✅ **IMPLEMENTED**
**Represents:** 3D point or direction

**Data:**
- `X, Y, Z float64`

**Operations implemented:**
- ✅ Add, Subtract, Scale
- ✅ Dot product
- ✅ Cross product (right-handed)
- ✅ Length, Normalize
- ✅ Distance between points

**Test coverage:** 46 tests, 100% coverage

**Location:** `pkg/math/vec3.go`, `pkg/math/vec3_test.go`

### Vector4
**Status:** Not implemented - may not be needed for MVP

**Note:** Homogeneous coordinates (w component) are handled implicitly in Mat4x4.MultiplyVec3, which treats Vec3 as a point with w=1. A separate Vec4 type may be added in the future if needed for advanced transformations.

### Mat4x4 ✅ **IMPLEMENTED**
**Represents:** 4x4 transformation matrix (column-major order)

**Data:**
- `[16]float64` array stored in column-major order
- Internal layout: [col0, col1, col2, col3]

**Operations implemented:**
- ✅ Identity matrix (NewIdentity)
- ✅ Zero matrix (NewZero)
- ✅ Matrix multiplication (Multiply)
- ✅ Transform vector (MultiplyVec3, treats as point with w=1)
- ✅ Transpose
- ✅ Element access (Get, Set with bounds checking)
- ⏳ Inverse (deferred - may be needed for Phase 3 Camera)

**Test coverage:** Comprehensive tests with table-driven approach, 100% coverage

**Location:** `pkg/math/mat4x4.go`, `pkg/math/mat4x4_test.go`

## Key Transformation Matrices

### Translation ✅ **IMPLEMENTED**
```go
NewTranslation(tx, ty, tz float64) Mat4x4
// Move point by offset (tx, ty, tz)
```

### Rotation ✅ **IMPLEMENTED**
```go
NewRotationX(angle float64) Mat4x4  // Rotate around X-axis
NewRotationY(angle float64) Mat4x4  // Rotate around Y-axis
NewRotationZ(angle float64) Mat4x4  // Rotate around Z-axis
// Angles in radians, right-handed coordinate system
```

### Scale ✅ **IMPLEMENTED**
```go
NewScale(sx, sy, sz float64) Mat4x4
// Uniform or non-uniform scaling
```

### Look-At Matrix ⏳ **PHASE 3**
```
Create view matrix from:
- Eye position (camera location)
- Target position (where camera looks)
- Up vector (camera orientation)

Status: Deferred to Phase 3 (Camera System)
```

### Perspective Projection ⏳ **PHASE 3**
```
Parameters:
- Field of view (FOV)
- Aspect ratio
- Near plane
- Far plane

Output: Clip space coordinates

Status: Deferred to Phase 3 (Camera System)
```

### Viewport Transform ⏳ **PHASE 3**
```
Convert NDC (-1 to 1) to screen coordinates (0 to width/height)

Parameters:
- Screen width
- Screen height

Status: Deferred to Phase 3 (Camera System)
```

## Design Principles

**Immutability preferred:** Operations return new values rather than modifying in place (for thread safety and clarity).

**Small and focused:** Each type should do one thing well.

**Well-tested:** Math bugs are the hardest to debug visually.

## Testing Requirements

### Phase 1 - Completed ✅
- ✅ Vector addition/subtraction correctness
- ✅ Cross product right-hand rule validation
- ✅ Matrix multiplication associativity
- ✅ Identity matrix behavior
- ✅ Known transformation results (e.g., 90° rotation)
- ✅ Zero-length vector normalization
- ✅ Edge cases (zero vectors, unit vectors, negative values)
- ✅ Epsilon-based float comparison (±0.0001)

### Phase 3 - Pending
- ⏳ Perspective projection edge cases (near/far plane)
- ⏳ Matrix inverse correctness (multiply by inverse = identity)

### Post-MVP - Future Enhancements
- ⏳ Numerical stability tests
- ⏳ Benchmark performance of hot paths

## Common Gotchas to Test

1. **Gimbal lock** in rotation composition
2. **Perspective divide by zero** (W = 0)
3. **Floating point precision** in normalize operations
4. **Matrix multiplication order** (easy to reverse)
5. **Left vs right-handed** coordinate system consistency

## API Example (Actual Implementation)

```go
import "github.com/muddl/go_toy_renderer/pkg/math"

// Vector operations
v1 := math.Vec3{X: 1, Y: 0, Z: 0}
v2 := math.Vec3{X: 0, Y: 1, Z: 0}
cross := v1.Cross(v2) // Returns Vec3{0, 0, 1}

// Vector math
sum := v1.Add(v2)                    // Vec3{1, 1, 0}
scaled := v1.Scale(2.0)              // Vec3{2, 0, 0}
length := v1.Length()                // 1.0
normalized := v1.Normalize()         // Unit vector

// Matrix creation (implemented)
model := math.NewTranslation(2, 0, 0)     // Translation matrix
rotation := math.NewRotationY(math.Pi/2)  // 90° rotation around Y
scale := math.NewScale(2, 2, 2)           // Uniform scale

// Matrix creation (Phase 3 - not yet implemented)
// view := math.NewLookAt(eye, target, up)
// proj := math.NewPerspective(fov, aspect, near, far)

// Matrix operations
identity := math.NewIdentity()
combined := rotation.Multiply(scale)  // Combine transformations
transformed := model.MultiplyVec3(v1) // Transform a point
transposed := model.Transpose()       // Matrix transpose
```

## Performance Considerations (Post-MVP)

- Matrix operations are hot paths
- Consider SIMD-friendly layouts later
- Cache transform chains (don't recompute each frame)
- Benchmark before optimizing
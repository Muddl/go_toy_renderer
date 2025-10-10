# Test Strategy

## Testing Philosophy

**For a renderer, correctness is paramount.** Visual bugs are hard to debug, so comprehensive testing is essential.

**Test pyramid for this project:**
```
        ┌─────────────┐
        │  Golden     │  ← Few, high-value
        │  Image      │
        ├─────────────┤
        │ Integration │  ← Some, key workflows
        │   Tests     │
        ├─────────────┤
        │    Unit     │  ← Many, fast, focused
        │   Tests     │
        └─────────────┘
```

---

## Unit Tests (Foundation)

**Goal:** Test individual components in isolation with known inputs/outputs.

### Math Component - CRITICAL

**Why critical:** Math bugs cascade through entire renderer.

**Test coverage:**
- Vector operations (add, subtract, scale, dot, cross, normalize)
- Matrix operations (multiply, transpose, inverse)
- Known transformations (90° rotation, identity, etc.)
- Edge cases (zero vectors, degenerate matrices)
- Floating point precision (near-equality tests)

**Example tests:**
```go
func TestVector3Cross(t *testing.T) {
    v1 := NewVector3(1, 0, 0)
    v2 := NewVector3(0, 1, 0)
    result := v1.Cross(v2)

    expected := NewVector3(0, 0, 1)
    if !result.Equals(expected, 0.0001) {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}

func TestMatrixMultiplyIdentity(t *testing.T) {
    m := RandomMatrix()
    identity := IdentityMatrix()

    result := m.Multiply(identity)

    if !result.Equals(m, 0.0001) {
        t.Error("M × I should equal M")
    }
}
```

**Test tools:**
- Table-driven tests for multiple cases
- Epsilon comparison for floats (~0.0001)
- Test all edge cases (zero, negative, very large values)

---

### Geometry Component

**Test coverage:**
- Mesh creation (valid vertex/index counts)
- Hardcoded primitives (cube has 8 verts, 12 tris)
- Index bounds checking
- OBJ loading (if implemented)

**Example tests:**
```go
func TestCreateCube(t *testing.T) {
    cube := CreateCube()

    if len(cube.Vertices) != 8 {
        t.Errorf("Cube should have 8 vertices, got %d", len(cube.Vertices))
    }

    if len(cube.Indices) != 36 { // 12 triangles × 3
        t.Errorf("Cube should have 36 indices, got %d", len(cube.Indices))
    }

    // Validate all indices are in range
    for i, idx := range cube.Indices {
        if idx < 0 || idx >= len(cube.Vertices) {
            t.Errorf("Index %d out of bounds: %d", i, idx)
        }
    }
}
```

---

### Camera Component

**Test coverage:**
- View matrix generation
- Projection matrix generation
- Known camera positions produce expected matrices
- Edge cases (FOV extremes, near=far, etc.)

**Example tests:**
```go
func TestLookAtIdentity(t *testing.T) {
    // Camera at origin looking down -Z should produce near-identity view
    camera := NewCamera(
        NewVector3(0, 0, 0),
        NewVector3(0, 0, -1),
        NewVector3(0, 1, 0),
    )

    view := camera.ViewMatrix()

    // Check specific matrix properties
    // (translation should be zero, rotation minimal, etc.)
}
```

---

### Framebuffer Component

**Test coverage:**
- Pixel read/write correctness
- Depth test logic (closer wins)
- Clear functionality
- Bounds checking
- Color clamping on export

**Example tests:**
```go
func TestDepthTest(t *testing.T) {
    fb := NewFramebuffer(100, 100)
    fb.Clear(NewVector3(0, 0, 0), 1.0)

    // Write pixel at depth 0.5
    fb.SetPixel(50, 50, NewVector3(1, 0, 0), 0.5)

    // Try to write at depth 0.7 (farther - should fail)
    fb.SetPixel(50, 50, NewVector3(0, 1, 0), 0.7)

    color := fb.GetPixel(50, 50)
    if !color.Equals(NewVector3(1, 0, 0), 0.0001) {
        t.Error("Farther pixel should not overwrite closer pixel")
    }

    // Write at depth 0.3 (closer - should succeed)
    fb.SetPixel(50, 50, NewVector3(0, 0, 1), 0.3)

    color = fb.GetPixel(50, 50)
    if !color.Equals(NewVector3(0, 0, 1), 0.0001) {
        t.Error("Closer pixel should overwrite farther pixel")
    }
}
```

---

### Rasterizer Component

**Test coverage:**
- Bounding box calculation
- Barycentric coordinate computation
- Single-pixel triangle
- Axis-aligned triangles
- Edge cases (degenerate, off-screen)

**Example tests:**
```go
func TestSinglePixelTriangle(t *testing.T) {
    // Triangle with all vertices at same pixel should fill exactly 1 pixel
    v := ScreenVertex{X: 10.5, Y: 10.5, ...}

    pixels := RasterizeTriangle(v, v, v)

    if len(pixels) != 1 {
        t.Errorf("Expected 1 pixel, got %d", len(pixels))
    }

    if pixels[0].X != 10 || pixels[0].Y != 10 {
        t.Error("Pixel at wrong location")
    }
}
```

---

## Integration Tests

**Goal:** Test component interactions with realistic data flows.

### Transform Pipeline

**Test coverage:**
- Vertex transforms through full pipeline
- Known 3D point → expected screen coordinate
- Matrix concatenation order
- Perspective divide correctness

**Example test:**
```go
func TestTransformPipeline(t *testing.T) {
    // Setup matrices
    model := IdentityMatrix()
    view := LookAt(...)
    proj := Perspective(...)

    mvp := proj.Multiply(view).Multiply(model)

    // Transform a known point
    worldPoint := NewVector3(0, 0, -5) // In front of camera

    clipSpace := mvp.TransformVector(worldPoint)
    ndcSpace := clipSpace.PerspectiveDivide()
    screenSpace := ViewportTransform(ndcSpace, width, height)

    // Verify point is in screen bounds and roughly centered
    if screenSpace.X < 0 || screenSpace.X > width {
        t.Error("Point transformed outside screen")
    }
}
```

---

### Render Complete Scene

**Test coverage:**
- Full render produces output without crashes
- Framebuffer has non-background pixels
- Depth buffer has varied values
- Output file is created

**Example test:**
```go
func TestRenderCube(t *testing.T) {
    scene := NewScene()
    scene.AddMesh(CreateCube())

    camera := NewCamera(...)
    fb := NewFramebuffer(100, 100)
    shader := VertexColorShader

    // Should not crash
    Render(scene, camera, fb, shader)

    // Framebuffer should have some non-background pixels
    hasColor := false
    for y := 0; y < 100; y++ {
        for x := 0; x < 100; x++ {
            if !fb.GetPixel(x, y).Equals(backgroundColor, 0.01) {
                hasColor = true
                break
            }
        }
    }

    if !hasColor {
        t.Error("Rendered framebuffer is all background color")
    }
}
```

---

## Golden Image Tests (Validation)

**Goal:** Compare rendered output against known-good reference images.

**When to use:**
- After MVP is working
- To prevent visual regressions
- To validate end-to-end correctness

**How it works:**
1. Render scene to framebuffer
2. Save as image
3. Compare against reference image pixel-by-pixel
4. Allow small epsilon for floating-point variance

**Tools:**
- Use image diff libraries
- Generate reference images from working implementation
- Store reference images in `testdata/` directory

**Example test:**
```go
func TestGoldenImageCube(t *testing.T) {
    // Render cube
    fb := renderTestScene("cube")

    // Load reference image
    reference := loadImage("testdata/cube_reference.png")

    // Compare
    diff := compareImages(fb, reference, epsilon=0.01)

    if diff > 0.001 { // Allow 0.1% pixels different
        t.Errorf("Image differs from reference by %.2f%%", diff*100)
        // Save diff image for debugging
        saveDiffImage("testdata/cube_diff.png", fb, reference)
    }
}
```

**Challenges:**
- Reference images must be committed to repo
- Different platforms might produce slightly different results
- Need to regenerate references when algorithm changes

---

## Benchmark Tests (Performance)

**Goal:** Track performance of critical code paths.

**What to benchmark:**
- Matrix multiplication
- Triangle rasterization
- Full frame render

**Example benchmark:**
```go
func BenchmarkMatrixMultiply(b *testing.B) {
    m1 := RandomMatrix()
    m2 := RandomMatrix()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = m1.Multiply(m2)
    }
}

func BenchmarkRenderFrame(b *testing.B) {
    scene := CreateTestScene()
    camera := CreateTestCamera()
    fb := NewFramebuffer(800, 600)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fb.Clear(backgroundColor, 1.0)
        Render(scene, camera, fb, VertexColorShader)
    }
}
```

**Run benchmarks:**
```bash
go test -bench=. -benchmem
```

---

## Test Organization

### Directory Structure
```
pkg/math/
    vector_test.go
    matrix_test.go
pkg/geometry/
    mesh_test.go
    primitives_test.go
pkg/render/
    pipeline_test.go
    integration_test.go
    golden_test.go
testdata/
    cube_reference.png
    depth_reference.png
```

### Test Naming Convention
```go
// Unit tests
func TestFunctionName(t *testing.T)

// Table-driven tests
func TestFunctionName_Cases(t *testing.T)

// Integration tests
func TestIntegration_ScenarioName(t *testing.T)

// Benchmarks
func BenchmarkOperationName(b *testing.B)
```

---

## Testing Best Practices

### 1. Isolate Tests
Each test should be independent. Don't rely on test order.

### 2. Use Table-Driven Tests
For testing multiple similar cases:
```go
func TestVectorAdd(t *testing.T) {
    cases := []struct {
        name     string
        v1, v2   Vector3
        expected Vector3
    }{
        {"zero vectors", Vec3(0,0,0), Vec3(0,0,0), Vec3(0,0,0)},
        {"positive", Vec3(1,2,3), Vec3(4,5,6), Vec3(5,7,9)},
        {"negative", Vec3(1,1,1), Vec3(-1,-1,-1), Vec3(0,0,0)},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            result := tc.v1.Add(tc.v2)
            if !result.Equals(tc.expected, 0.0001) {
                t.Errorf("got %v, want %v", result, tc.expected)
            }
        })
    }
}
```

### 3. Test Edge Cases
- Zero values
- Negative values
- Very large/small values
- Boundary conditions
- Invalid inputs

### 4. Use Helper Functions
```go
func assertVector3Equal(t *testing.T, got, want Vector3) {
    t.Helper()
    if !got.Equals(want, 0.0001) {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

---

## Test Coverage Goals

**Minimum for MVP:**
- Math package: >90% coverage
- Geometry: >80% coverage
- Rasterizer: >70% coverage
- Integration: Basic happy path tests

**Check coverage:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## MVP Testing Checklist

Before MVP is complete, verify:

- [ ] All math operations have unit tests
- [ ] Matrix transformations tested with known results
- [ ] Geometry primitives validated
- [ ] Framebuffer depth test verified
- [ ] Rasterization produces expected pixels
- [ ] End-to-end render completes without errors
- [ ] Output image is visually correct (manual check)
- [ ] At least one golden image test passes
- [ ] No crashes with valid inputs
- [ ] Coverage >70% on core packages

**When all checked: Testing is sufficient for MVP.**

# Geometry Component

## Purpose

Define geometric primitives and mesh data structures for 3D objects.

## Core Types

### Vertex
**Represents:** A single point in 3D space with attributes

**Minimum data for MVP:**
- `Position Vector3` - 3D location
- `Color Vector3` - RGB color (0-1 range)

**Future attributes (post-MVP):**
- Normal vector (for lighting)
- Texture coordinates (UV)
- Tangent/bitangent (for normal mapping)

### Triangle
**Represents:** A single triangle face

**Data:**
- Three vertex indices (into mesh vertex array)
- OR: Three vertex values directly

**Design choice:** Index-based is more memory efficient and standard.

### Mesh
**Represents:** A collection of vertices and triangles forming a 3D object

**Data:**
- `Vertices []Vertex` - vertex buffer
- `Indices []int` - index buffer (groups of 3 form triangles)

**Operations needed:**
- Get triangle by index
- Iterate over all triangles
- Calculate bounding box (for culling - optional MVP)

## Mesh Creation Helpers

### Hardcoded Primitives

**For testing and MVP demonstration:**

```
CreateCube() *Mesh
CreateTetrahedron() *Mesh
CreatePyramid() *Mesh
```

**Why:** Don't need file loading to test renderer.

### File Loading (Nice to Have for MVP)

```
LoadOBJ(filename string) (*Mesh, error)
```

**OBJ format:** Simple, text-based, widely supported
- Only need to parse: `v` (vertex), `f` (face)
- Can ignore: materials, normals, UVs for MVP

## Data Organization

### Index Buffer Advantages
- Vertices shared between triangles (saves memory)
- Standard in 3D graphics
- Matches how GPUs work

### Winding Order
All triangles should have **counter-clockwise** winding when viewed from front.

**Why:** Enables backface culling.

## Coordinate Space

Mesh vertices are in **local/object space**.

**Transformation to world space:** Apply model matrix.

## Testing Requirements

### Geometry Validation
- Mesh has valid indices (no out of bounds)
- Triangles have consistent winding
- Primitive creation produces expected vertex count
- OBJ loader handles malformed files gracefully

### Geometric Properties
- Cube has 8 vertices, 12 triangles
- Tetrahedron has 4 vertices, 4 triangles
- Vertex colors are in valid range (0-1)

## API Example (Conceptual)

```go
// Create test geometry
cube := geometry.CreateCube()

// Access triangles
for i := 0; i < len(cube.Indices)/3; i++ {
    idx0 := cube.Indices[i*3]
    idx1 := cube.Indices[i*3+1]
    idx2 := cube.Indices[i*3+2]

    v0 := cube.Vertices[idx0]
    v1 := cube.Vertices[idx1]
    v2 := cube.Vertices[idx2]

    // Process triangle...
}

// Load from file (optional)
mesh, err := geometry.LoadOBJ("model.obj")
```

## Design Considerations

### Memory Layout
For MVP: Simple slice of structs is fine.

Post-MVP: Consider struct-of-arrays for cache efficiency:
```go
type Mesh struct {
    PositionsX []float64
    PositionsY []float64
    PositionsZ []float64
    // ... separate arrays per attribute
}
```

### Vertex Attributes
Start minimal, add attributes as needed:
1. MVP: Position + Color
2. Add Normal (for lighting)
3. Add UV (for textures)
4. Add Tangent (for normal maps)

**Don't over-engineer:** Add complexity when you need it.

## Common Gotchas

1. **Index out of bounds** - validate indices on load
2. **Inconsistent winding** - some OBJ files are clockwise
3. **Degenerate triangles** - zero area (all vertices colinear)
4. **Large coordinate values** - can cause precision issues
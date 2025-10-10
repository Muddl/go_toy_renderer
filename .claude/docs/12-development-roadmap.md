# Development Roadmap

## Overview

This roadmap breaks MVP development into phases, each building on the previous. Complete each phase before moving to the next.

**Estimated total time:** 2-4 weeks for hobbyist (few hours per day)

---

## Phase 1: Math Foundation (Days 1-3)

**Goal:** Rock-solid math library that everything else depends on.

### Tasks
1. Implement Vector3 type and operations
   - Add, subtract, scale, dot, cross, normalize, length
   - Test extensively with known results

2. Implement Matrix4x4 type and operations
   - Identity, multiply, transpose
   - Test matrix multiplication properties

3. Create transformation matrices
   - Translation, rotation (axis-angle or Euler), scale
   - Test known transformations (90° rotation, etc.)

4. Implement camera matrices
   - LookAt (view matrix)
   - Perspective projection
   - Test with known camera positions

5. Viewport transformation
   - NDC to screen coordinates
   - Test coordinate mapping

### Completion Criteria
- [ ] All vector operations pass tests
- [ ] Matrix operations pass tests
- [ ] Can create all transformation matrices
- [ ] >90% test coverage in math package
- [ ] No known bugs

**When complete:** You can transform points from 3D to 2D. Foundation is solid.

---

## Phase 2: Geometry & Scene (Days 4-5)

**Goal:** Define 3D objects and organize them in a scene.

### Tasks
1. Implement Vertex type
   - Position (Vector3)
   - Color (Vector3)

2. Implement Mesh type
   - Vertex buffer (slice of vertices)
   - Index buffer (slice of ints)
   - Methods to access triangles

3. Create hardcoded primitives
   - Start with tetrahedron (4 vertices, simplest)
   - Then cube (8 vertices, more complex)
   - Assign different colors to vertices

4. Implement basic Scene type
   - List of meshes
   - Method to add/remove meshes

5. Test geometry
   - Verify vertex counts
   - Validate indices in bounds
   - Check winding order consistency

### Completion Criteria
- [ ] Can create tetrahedron and cube
- [ ] Meshes have valid vertex/index data
- [ ] Scene can hold multiple meshes
- [ ] Geometry tests pass

**When complete:** You have test objects to render.

---

## Phase 3: Camera System (Days 6-7)

**Goal:** Flexible viewpoint control.

### Tasks
1. Implement Camera type
   - Position, target, up vector
   - FOV, aspect, near, far

2. Implement view matrix generation
   - Use LookAt from Phase 1
   - Test camera transformations

3. Implement projection matrix generation
   - Use perspective projection from Phase 1
   - Test projection properties

4. Combine into view-projection matrix
   - Pre-multiply for efficiency
   - Test complete transformation

5. Create test camera configurations
   - Camera in front of object
   - Camera offset to side
   - Different FOV values

### Completion Criteria
- [ ] Camera generates correct view matrix
- [ ] Camera generates correct projection matrix
- [ ] Combined VP matrix transforms points correctly
- [ ] Camera tests pass

**When complete:** You can control viewpoint.

---

## Phase 4: Framebuffer (Days 8-9)

**Goal:** Store rendered pixels and output images.

### Tasks
1. Implement Framebuffer type
   - Color buffer (RGB per pixel)
   - Depth buffer (Z per pixel)

2. Implement Clear method
   - Reset to background color and depth

3. Implement SetPixel with depth test
   - Check if new pixel is closer
   - Update color and depth if closer

4. Implement GetPixel (for testing)
   - Retrieve color at coordinates

5. Implement image export
   - Convert float RGB to byte RGB
   - Save as PNG using Go image libraries
   - Test file output

### Completion Criteria
- [ ] Can create framebuffer of any size
- [ ] Clear works correctly
- [ ] Depth test prevents farther pixels from overwriting closer ones
- [ ] Can save valid PNG file
- [ ] Framebuffer tests pass

**When complete:** You can store and save rendered images.

---

## Phase 5: Rasterization (Days 10-12)

**Goal:** Convert triangles to pixels with interpolation.

### Tasks
1. Choose rasterization algorithm
   - Recommended: Barycentric (simpler code)
   - Alternative: Scanline (more classic)

2. Implement barycentric coordinate calculation
   - Test with known points inside/outside triangle

3. Implement triangle bounding box
   - Find min/max X and Y coordinates

4. Implement rasterization loop
   - For each pixel in bounding box
   - Check if inside triangle
   - Calculate barycentric weights

5. Implement attribute interpolation
   - Interpolate depth (Z)
   - Interpolate color (RGB)
   - Test interpolation correctness

6. Connect to framebuffer
   - For each pixel covered
   - Call framebuffer.SetPixel with depth test

7. Test rasterization
   - Single pixel triangle
   - Axis-aligned triangles
   - Arbitrary triangles
   - Verify color interpolation

### Completion Criteria
- [ ] Can rasterize any triangle to pixels
- [ ] Attribute interpolation is correct at vertices
- [ ] Depth interpolation works
- [ ] No crashes on degenerate triangles
- [ ] Rasterizer tests pass

**When complete:** You can draw triangles to framebuffer.

---

## Phase 6: Shading (Days 13-14)

**Goal:** Calculate per-pixel colors.

### Tasks
1. Define Shader interface or function type
   - Input: interpolated attributes
   - Output: final color

2. Implement vertex color shader
   - Simply return interpolated color
   - Test with known inputs

3. Implement flat color shader (for testing)
   - Return constant color
   - Use for debugging

4. Optional: Implement depth visualization shader
   - Convert depth to grayscale
   - Useful for debugging

5. Test shaders
   - Verify correct colors returned
   - Ensure called once per pixel

### Completion Criteria
- [ ] Shader interface defined
- [ ] At least one shader implemented (vertex color)
- [ ] Shader receives correct interpolated data
- [ ] Shader tests pass

**When complete:** You can control per-pixel appearance.

---

## Phase 7: Render Pipeline Integration (Days 15-18)

**Goal:** Connect all components into working renderer.

### Tasks
1. Implement vertex transformation stage
   - For each vertex in mesh:
   - Apply model matrix (identity for now)
   - Apply view-projection matrix
   - Perform perspective divide
   - Convert to screen coordinates

2. Implement primitive assembly
   - Group indices into triangles (every 3 indices)
   - Fetch transformed vertices for each triangle

3. Connect rasterization and shading
   - For each triangle:
   - Rasterize to pixels
   - For each pixel, call shader
   - Write to framebuffer

4. Implement main Render function
   - Clear framebuffer
   - Transform all vertices
   - Process all triangles
   - Return/save framebuffer

5. Test end-to-end
   - Render tetrahedron
   - Render cube
   - Save output image
   - Visually verify correctness

6. Debug visual issues
   - Check matrix multiplication order
   - Verify depth test direction
   - Check winding order
   - Fix any artifacts

### Completion Criteria
- [ ] Can render complete mesh to image
- [ ] Perspective looks correct
- [ ] Depth ordering is correct
- [ ] Colors interpolate smoothly
- [ ] Output image matches expectations
- [ ] No crashes on valid input

**When complete:** MVP is functionally complete!

---

## Phase 8: Testing & Polish (Days 19-21)

**Goal:** Ensure quality and reliability.

### Tasks
1. Write integration tests
   - Test complete render pipeline
   - Test various camera positions
   - Test multiple meshes

2. Create golden image tests
   - Render reference images
   - Compare output against references
   - Set up automated comparison

3. Improve test coverage
   - Aim for >80% overall
   - Focus on uncovered critical paths

4. Create demo application
   - cmd/renderer/main.go
   - Render cube from nice angle
   - Save as output.png

5. Write basic documentation
   - README with usage instructions
   - Code comments for public APIs
   - Example usage

6. Performance check
   - Benchmark critical paths
   - Identify bottlenecks (don't optimize yet)
   - Document performance characteristics

### Completion Criteria
- [ ] All tests passing
- [ ] >70% test coverage
- [ ] At least one golden image test
- [ ] Demo app produces nice output image
- [ ] README explains how to use
- [ ] Code is reasonably documented

**When complete:** MVP is ready to show!

---

## Beyond MVP (Future Phases)

### Phase 9: Enhancements (Optional)
- OBJ file loading
- Backface culling
- Multiple meshes in scene
- Wireframe rendering mode

### Phase 10: Lighting (Optional)
- Add normals to vertices
- Implement simple diffuse lighting
- Directional light source

### Phase 11: Optimization (Optional)
- Multi-threading (parallel triangle processing)
- SIMD vectorization
- Profile-guided optimization

### Phase 12: Advanced Features (Optional)
- Texture mapping
- Specular highlights
- Simple shadows
- Post-processing effects

---

## Daily Development Checklist

**Start of day:**
- [ ] Review what you completed yesterday
- [ ] Identify today's goal (1-2 tasks from current phase)
- [ ] Set up test environment

**During development:**
- [ ] Write tests BEFORE implementation (TDD)
- [ ] Commit working code frequently
- [ ] Run tests after each change
- [ ] Take breaks when stuck

**End of day:**
- [ ] All tests passing
- [ ] Code committed to git
- [ ] Update roadmap with progress
- [ ] Note any blockers or questions

---

## MVP Success Checklist

MVP is complete when:

- [ ] Math library fully implemented and tested
- [ ] Can create test geometry (cube/tetrahedron)
- [ ] Camera system works
- [ ] Framebuffer stores and saves images
- [ ] Rasterizer fills triangles correctly
- [ ] Shader calculates colors
- [ ] Render pipeline connects all components
- [ ] Can render colored 3D object with perspective
- [ ] Output image saved as PNG
- [ ] Tests pass (>70% coverage)
- [ ] Demo app works
- [ ] No critical bugs

**When all checked: MVP is done!**

---

## Estimated Timeline

**Conservative estimate (beginner):** 20-25 days
**Moderate estimate (some 3D experience):** 12-15 days
**Optimistic estimate (experienced):** 8-10 days

**Factors affecting timeline:**
- Prior 3D graphics experience
- Time available per day
- Debugging time needed
- Test-writing discipline
- Scope creep (adding non-MVP features)

**Tips to stay on track:**
1. Stick to MVP scope (resist feature creep!)
2. Write tests first (saves debugging time)
3. Commit working code frequently
4. Ask for help when stuck >1 hour
5. Celebrate completing each phase

**Most important:** Complete is better than perfect. Finish MVP first, then improve.

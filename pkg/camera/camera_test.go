package camera

import (
	stdmath "math"
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

const epsilon = 0.0001

// approxEqual compares two float64 values within epsilon tolerance.
func approxEqual(a, b float64) bool {
	return stdmath.Abs(a-b) <= epsilon
}

// TestCamera_New_StoresFieldsCorrectly tests that constructor stores all fields.
func TestCamera_New_StoresFieldsCorrectly(t *testing.T) {
	pos := math.Vec3{X: 1, Y: 2, Z: 3}
	target := math.Vec3{X: 0, Y: 0, Z: 0}
	up := math.Vec3{X: 0, Y: 1, Z: 0}

	cam := New(pos, target, up, 45.0, 16.0/9.0, 0.1, 100.0)

	if !cam.Position.Equals(pos, epsilon) {
		t.Errorf("Camera.Position = %v, want %v", cam.Position, pos)
	}
	if !cam.Target.Equals(target, epsilon) {
		t.Errorf("Camera.Target = %v, want %v", cam.Target, target)
	}
	if !cam.Up.Equals(up, epsilon) {
		t.Errorf("Camera.Up = %v, want %v", cam.Up, up)
	}
	if !approxEqual(cam.FOV, 45.0) {
		t.Errorf("Camera.FOV = %f, want 45.0", cam.FOV)
	}
	if !approxEqual(cam.Aspect, 16.0/9.0) {
		t.Errorf("Camera.Aspect = %f, want %f", cam.Aspect, 16.0/9.0)
	}
	if !approxEqual(cam.Near, 0.1) {
		t.Errorf("Camera.Near = %f, want 0.1", cam.Near)
	}
	if !approxEqual(cam.Far, 100.0) {
		t.Errorf("Camera.Far = %f, want 100.0", cam.Far)
	}
}

// TestCamera_ViewMatrix_CameraAtOriginLookingDownNegZ tests canonical camera orientation.
// Camera at (0,0,5) looking at origin should transform origin to (0,0,-5) in camera space.
func TestCamera_ViewMatrix_CameraAtOriginLookingDownNegZ(t *testing.T) {
	cam := New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0, 1.0, 0.1, 100.0,
	)

	view := cam.ViewMatrix()

	// World origin should map to (0, 0, -5) in camera space
	ox, oy, oz, ow := view.MultiplyVec4(0, 0, 0, 1)
	if !approxEqual(ow, 1.0) {
		t.Fatalf("ViewMatrix * origin produced unexpected w = %f, want 1.0", ow)
	}

	if !approxEqual(ox, 0.0) || !approxEqual(oy, 0.0) || !approxEqual(oz, -5.0) {
		t.Errorf("ViewMatrix * origin = (%f,%f,%f), want (0,0,-5)", ox, oy, oz)
	}
}

// TestCamera_ViewMatrix_CameraPositionMapsToOrigin tests that camera position maps to (0,0,0).
func TestCamera_ViewMatrix_CameraPositionMapsToOrigin(t *testing.T) {
	cam := New(
		math.Vec3{X: 3, Y: 4, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0, 1.0, 0.1, 100.0,
	)

	view := cam.ViewMatrix()

	// Camera position in world space should map to origin in camera space
	px, py, pz, pw := view.MultiplyVec4(3, 4, 5, 1)
	if !approxEqual(pw, 1.0) {
		t.Fatalf("ViewMatrix produced unexpected w for camera position = %f", pw)
	}

	if !approxEqual(px, 0.0) || !approxEqual(py, 0.0) || !approxEqual(pz, 0.0) {
		t.Errorf("ViewMatrix * camera position = (%f,%f,%f), want (0,0,0)", px, py, pz)
	}
}

// TestCamera_ViewMatrix_BasisVectorsAreOrthonormal verifies rotation block is an orthonormal frame.
func TestCamera_ViewMatrix_BasisVectorsAreOrthonormal(t *testing.T) {
	cam := New(
		math.Vec3{X: 3, Y: 4, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0, 1.0, 0.1, 100.0,
	)

	view := cam.ViewMatrix()

	// Extract the three basis row vectors
	right := math.Vec3{X: view.Get(0, 0), Y: view.Get(0, 1), Z: view.Get(0, 2)}
	upVec := math.Vec3{X: view.Get(1, 0), Y: view.Get(1, 1), Z: view.Get(1, 2)}
	fwd := math.Vec3{X: view.Get(2, 0), Y: view.Get(2, 1), Z: view.Get(2, 2)}

	// Each should have unit length
	if !approxEqual(right.Length(), 1.0) {
		t.Errorf("right vector length = %f, want 1.0", right.Length())
	}
	if !approxEqual(upVec.Length(), 1.0) {
		t.Errorf("up vector length = %f, want 1.0", upVec.Length())
	}
	if !approxEqual(fwd.Length(), 1.0) {
		t.Errorf("fwd vector length = %f, want 1.0", fwd.Length())
	}

	// Each pair should be perpendicular
	if !approxEqual(right.Dot(upVec), 0.0) {
		t.Errorf("right dot up = %f, want 0.0", right.Dot(upVec))
	}
	if !approxEqual(right.Dot(fwd), 0.0) {
		t.Errorf("right dot fwd = %f, want 0.0", right.Dot(fwd))
	}
	if !approxEqual(upVec.Dot(fwd), 0.0) {
		t.Errorf("up dot fwd = %f, want 0.0", upVec.Dot(fwd))
	}
}

// TestCamera_ProjectionMatrix_NearPlane_MapsToNegativeOne verifies near plane NDC mapping.
func TestCamera_ProjectionMatrix_NearPlane_MapsToNegativeOne(t *testing.T) {
	cam := New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		90.0, 1.0, 1.0, 100.0,
	)

	proj := cam.ProjectionMatrix()

	// Point at z=-near in camera space should map to NDC z = -1 after perspective divide
	_, _, pz, pw := proj.MultiplyVec4(0, 0, -cam.Near, 1)
	if stdmath.Abs(pw) < epsilon {
		t.Fatal("ProjectionMatrix produced near-zero w for near plane point")
	}
	ndcZ := pz / pw

	if !approxEqual(ndcZ, -1.0) {
		t.Errorf("Near plane NDC Z = %f, want -1.0", ndcZ)
	}
}

// TestCamera_ProjectionMatrix_FarPlane_MapsToPositiveOne verifies far plane NDC mapping.
func TestCamera_ProjectionMatrix_FarPlane_MapsToPositiveOne(t *testing.T) {
	cam := New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		90.0, 1.0, 1.0, 100.0,
	)

	proj := cam.ProjectionMatrix()

	// Point at z=-far in camera space should map to NDC z = 1 after perspective divide
	_, _, pz, pw := proj.MultiplyVec4(0, 0, -cam.Far, 1)
	if stdmath.Abs(pw) < epsilon {
		t.Fatal("ProjectionMatrix produced near-zero w for far plane point")
	}
	ndcZ := pz / pw

	if !approxEqual(ndcZ, 1.0) {
		t.Errorf("Far plane NDC Z = %f, want 1.0", ndcZ)
	}
}

// TestCamera_ProjectionMatrix_FOV90Aspect1_CorrectXScaling verifies x scaling at 90° FOV.
// With 90° vertical FOV and aspect=1, f=1. A point at (1,0,-1) should project to NDC x=1.
func TestCamera_ProjectionMatrix_FOV90Aspect1_CorrectXScaling(t *testing.T) {
	cam := New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		90.0, 1.0, 0.1, 100.0,
	)

	proj := cam.ProjectionMatrix()

	// f = 1/tan(45°) = 1; x_ndc = (f/aspect) * (x / -z) = 1 * (1/1) = 1
	px, _, _, pw := proj.MultiplyVec4(1, 0, -1, 1)
	if stdmath.Abs(pw) < epsilon {
		t.Fatal("ProjectionMatrix produced near-zero w")
	}
	ndcX := px / pw

	if !approxEqual(ndcX, 1.0) {
		t.Errorf("NDC X for (1,0,-1) with 90° FOV = %f, want 1.0", ndcX)
	}
}

// TestCamera_ProjectionMatrix_PointBehindCamera_HasNegativeW verifies clip space rejection.
func TestCamera_ProjectionMatrix_PointBehindCamera_HasNegativeW(t *testing.T) {
	cam := New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0, 1.0, 0.1, 100.0,
	)

	proj := cam.ProjectionMatrix()

	// A point at positive z in camera space is behind the camera; W should be negative
	_, _, _, pw := proj.MultiplyVec4(0, 0, 1.0, 1)

	if pw >= 0 {
		t.Errorf("Point behind camera should have negative W, got %f", pw)
	}
}

// TestCamera_ProjectionMatrix_NarrowFOV_HasLargerScale verifies FOV effect on scaling.
func TestCamera_ProjectionMatrix_NarrowFOV_HasLargerScale(t *testing.T) {
	cam45 := New(math.Vec3{}, math.Vec3{}, math.Vec3{X: 0, Y: 1}, 45.0, 1.0, 0.1, 100.0)
	cam90 := New(math.Vec3{}, math.Vec3{}, math.Vec3{X: 0, Y: 1}, 90.0, 1.0, 0.1, 100.0)

	proj45 := cam45.ProjectionMatrix()
	proj90 := cam90.ProjectionMatrix()

	// Y scale = f = 1/tan(fov/2); narrower FOV => larger f => more zoom
	scale45 := proj45.Get(1, 1)
	scale90 := proj90.Get(1, 1)

	if scale45 <= scale90 {
		t.Errorf("45° FOV Y scale (%f) should be > 90° FOV Y scale (%f)", scale45, scale90)
	}
}

// TestCamera_ProjectionMatrix_WideAspect_ReducesXScale verifies aspect ratio effect.
func TestCamera_ProjectionMatrix_WideAspect_ReducesXScale(t *testing.T) {
	cam := New(math.Vec3{}, math.Vec3{}, math.Vec3{X: 0, Y: 1}, 45.0, 16.0/9.0, 0.1, 100.0)
	proj := cam.ProjectionMatrix()

	// X scale = f/aspect, Y scale = f; aspect > 1 => X scale < Y scale
	xScale := proj.Get(0, 0)
	yScale := proj.Get(1, 1)

	if xScale >= yScale {
		t.Errorf("X scale (%f) should be < Y scale (%f) for aspect > 1", xScale, yScale)
	}
}

// TestCamera_ViewProjectionMatrix_EqualsProjectionTimesView verifies VP = Projection * View.
func TestCamera_ViewProjectionMatrix_EqualsProjectionTimesView(t *testing.T) {
	cam := New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0, 800.0/600.0, 0.1, 100.0,
	)

	view := cam.ViewMatrix()
	proj := cam.ProjectionMatrix()
	vp := cam.ViewProjectionMatrix()
	expected := proj.Multiply(view)

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			got := vp.Get(row, col)
			want := expected.Get(row, col)
			if !approxEqual(got, want) {
				t.Errorf("ViewProjectionMatrix[%d,%d] = %f, want %f", row, col, got, want)
			}
		}
	}
}

// TestCamera_ViewProjectionMatrix_WorldOriginProjectsToCenter tests full transform pipeline.
// A camera at (0,0,5) looking at origin: the world origin should project to NDC (0,0).
func TestCamera_ViewProjectionMatrix_WorldOriginProjectsToCenter(t *testing.T) {
	cam := New(
		math.Vec3{X: 0, Y: 0, Z: 5},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		90.0, 1.0, 0.1, 100.0,
	)

	vp := cam.ViewProjectionMatrix()

	// World origin is directly in front of the camera; should project to NDC center
	px, py, _, pw := vp.MultiplyVec4(0, 0, 0, 1)
	if pw <= 0 {
		t.Fatalf("World origin should have positive W after VP transform, got %f", pw)
	}

	ndcX := px / pw
	ndcY := py / pw

	if !approxEqual(ndcX, 0.0) || !approxEqual(ndcY, 0.0) {
		t.Errorf("World origin NDC = (%f, %f), want (0, 0)", ndcX, ndcY)
	}
}

// TestCamera_ViewMatrix_UpVectorParallelToViewDirection_DoesNotPanic verifies no panic on degenerate input.
func TestCamera_ViewMatrix_UpVectorParallelToViewDirection_DoesNotPanic(t *testing.T) {
	// Camera looking straight up with up=(0,1,0) — degenerate case
	cam := New(
		math.Vec3{X: 0, Y: 5, Z: 0},
		math.Vec3{X: 0, Y: 0, Z: 0},
		math.Vec3{X: 0, Y: 1, Z: 0},
		45.0, 1.0, 0.1, 100.0,
	)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ViewMatrix panicked on degenerate up vector: %v", r)
		}
	}()

	_ = cam.ViewMatrix()
}

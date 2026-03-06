package scene

import (
	stdmath "math"
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

const epsilon = 0.0001

func matApproxEqual(a, b math.Mat4x4, eps float64) bool {
	for i := 0; i < 16; i++ {
		if stdmath.Abs(a[i]-b[i]) > eps {
			return false
		}
	}
	return true
}

func TestTransform_ModelMatrix_Identity(t *testing.T) {
	tr := Transform{
		Position: math.Vec3{X: 0, Y: 0, Z: 0},
		Rotation: math.Vec3{X: 0, Y: 0, Z: 0},
		Scale:    math.Vec3{X: 1, Y: 1, Z: 1},
	}
	got := tr.ModelMatrix()
	want := math.NewIdentity()
	if !matApproxEqual(got, want, epsilon) {
		t.Fatalf("identity transform: got %v, want %v", got, want)
	}
}

func TestTransform_ModelMatrix_TranslationOnly(t *testing.T) {
	tr := Transform{
		Position: math.Vec3{X: 3, Y: -2, Z: 5},
		Rotation: math.Vec3{X: 0, Y: 0, Z: 0},
		Scale:    math.Vec3{X: 1, Y: 1, Z: 1},
	}
	got := tr.ModelMatrix()
	want := math.NewTranslation(3, -2, 5)
	if !matApproxEqual(got, want, epsilon) {
		t.Fatalf("translation only: got %v, want %v", got, want)
	}
}

func TestTransform_ModelMatrix_ScaleOnly(t *testing.T) {
	tr := Transform{
		Position: math.Vec3{X: 0, Y: 0, Z: 0},
		Rotation: math.Vec3{X: 0, Y: 0, Z: 0},
		Scale:    math.Vec3{X: 2, Y: 3, Z: 0.5},
	}
	got := tr.ModelMatrix()
	want := math.NewScale(2, 3, 0.5)
	if !matApproxEqual(got, want, epsilon) {
		t.Fatalf("scale only: got %v, want %v", got, want)
	}
}

func TestTransform_ModelMatrix_RotationY(t *testing.T) {
	angle := stdmath.Pi / 4.0 // 45 degrees
	tr := Transform{
		Position: math.Vec3{X: 0, Y: 0, Z: 0},
		Rotation: math.Vec3{X: 0, Y: angle, Z: 0},
		Scale:    math.Vec3{X: 1, Y: 1, Z: 1},
	}
	got := tr.ModelMatrix()
	want := math.NewRotationY(angle)
	if !matApproxEqual(got, want, epsilon) {
		t.Fatalf("rotation Y only: got %v, want %v", got, want)
	}
}

func TestTransform_ModelMatrix_TranslationScaleRotation(t *testing.T) {
	// Model = Translation * RotationZ * RotationY * RotationX * Scale
	tr := Transform{
		Position: math.Vec3{X: 1, Y: 2, Z: 3},
		Rotation: math.Vec3{X: 0.1, Y: 0.2, Z: 0.3},
		Scale:    math.Vec3{X: 2, Y: 2, Z: 2},
	}
	got := tr.ModelMatrix()

	// Manually compute the expected result.
	s := math.NewScale(2, 2, 2)
	rx := math.NewRotationX(0.1)
	ry := math.NewRotationY(0.2)
	rz := math.NewRotationZ(0.3)
	tl := math.NewTranslation(1, 2, 3)
	want := tl.Multiply(rz.Multiply(ry.Multiply(rx.Multiply(s))))
	if !matApproxEqual(got, want, epsilon) {
		t.Fatalf("combined transform:\ngot  %v\nwant %v", got, want)
	}
}

func TestNewTransform_DefaultsToIdentity(t *testing.T) {
	tr := NewTransform()
	got := tr.ModelMatrix()
	want := math.NewIdentity()
	if !matApproxEqual(got, want, epsilon) {
		t.Fatalf("NewTransform default: got %v, want %v", got, want)
	}
}

func TestNode_HasMeshAndTransform(t *testing.T) {
	tr := NewTransform()
	n := Node{Transform: tr}
	if n.Mesh != nil {
		t.Fatal("expected nil mesh for empty node")
	}
	got := n.Transform.ModelMatrix()
	want := math.NewIdentity()
	if !matApproxEqual(got, want, epsilon) {
		t.Fatal("expected identity model matrix for default transform")
	}
}

func TestScene_AddAndListNodes(t *testing.T) {
	s := NewScene()
	if len(s.Nodes) != 0 {
		t.Fatal("expected empty scene")
	}
	n := Node{Transform: NewTransform()}
	s.AddNode(n)
	if len(s.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(s.Nodes))
	}
}

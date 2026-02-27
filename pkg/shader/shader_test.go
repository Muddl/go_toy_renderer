// Package shader provides per-pixel color calculation for the rendering pipeline.
package shader

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

const epsilon = 0.001

func approxEqual(a, b float64) bool {
	d := a - b
	return d > -epsilon && d < epsilon
}

func vec3Equal(a, b math.Vec3) bool {
	return approxEqual(a.X, b.X) && approxEqual(a.Y, b.Y) && approxEqual(a.Z, b.Z)
}

// --- VertexColor tests ---

func TestVertexColor_ReturnsInputColor(t *testing.T) {
	attr := Attributes{Color: math.Vec3{X: 1.0, Y: 0.5, Z: 0.25}}
	got := VertexColor(attr)
	want := math.Vec3{X: 1.0, Y: 0.5, Z: 0.25}
	if !vec3Equal(got, want) {
		t.Errorf("VertexColor() = %v, want %v", got, want)
	}
}

func TestVertexColor_WithZeroColor(t *testing.T) {
	attr := Attributes{Color: math.Vec3{}}
	got := VertexColor(attr)
	want := math.Vec3{}
	if !vec3Equal(got, want) {
		t.Errorf("VertexColor() = %v, want %v", got, want)
	}
}

func TestVertexColor_IgnoresDepth(t *testing.T) {
	attr1 := Attributes{Color: math.Vec3{X: 0.8, Y: 0.2, Z: 0.6}, Depth: 0.0}
	attr2 := Attributes{Color: math.Vec3{X: 0.8, Y: 0.2, Z: 0.6}, Depth: 1.0}
	got1 := VertexColor(attr1)
	got2 := VertexColor(attr2)
	if !vec3Equal(got1, got2) {
		t.Errorf("VertexColor() changed with depth: %v vs %v", got1, got2)
	}
}

// --- NewFlatColor tests ---

func TestNewFlatColor_ReturnsConstantColor(t *testing.T) {
	red := math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}
	shader := NewFlatColor(red)
	attr := Attributes{Color: math.Vec3{X: 0.5, Y: 0.5, Z: 0.5}, Depth: 0.3}
	got := shader(attr)
	if !vec3Equal(got, red) {
		t.Errorf("NewFlatColor()() = %v, want %v", got, red)
	}
}

func TestNewFlatColor_IgnoresAttributes(t *testing.T) {
	blue := math.Vec3{X: 0.0, Y: 0.0, Z: 1.0}
	shader := NewFlatColor(blue)

	attr1 := Attributes{Color: math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}, Depth: 0.0}
	attr2 := Attributes{Color: math.Vec3{X: 0.0, Y: 1.0, Z: 0.0}, Depth: 1.0}

	got1 := shader(attr1)
	got2 := shader(attr2)
	if !vec3Equal(got1, blue) || !vec3Equal(got2, blue) {
		t.Errorf("NewFlatColor() should ignore attributes: got %v and %v, want %v", got1, got2, blue)
	}
}

func TestNewFlatColor_ReturnsFuncType(t *testing.T) {
	color := math.Vec3{X: 0.5, Y: 0.5, Z: 0.5}
	f := NewFlatColor(color)
	attr := Attributes{}
	got := f(attr)
	if !vec3Equal(got, color) {
		t.Errorf("NewFlatColor() result as Func = %v, want %v", got, color)
	}
}

// --- Depth tests ---

func TestDepth_MidDepth(t *testing.T) {
	attr := Attributes{Depth: 0.5}
	got := Depth(attr)
	want := math.Vec3{X: 0.5, Y: 0.5, Z: 0.5}
	if !vec3Equal(got, want) {
		t.Errorf("Depth(0.5) = %v, want %v", got, want)
	}
}

func TestDepth_ZeroDepth(t *testing.T) {
	attr := Attributes{Depth: 0.0}
	got := Depth(attr)
	want := math.Vec3{X: 0.0, Y: 0.0, Z: 0.0}
	if !vec3Equal(got, want) {
		t.Errorf("Depth(0.0) = %v, want %v", got, want)
	}
}

func TestDepth_FullDepth(t *testing.T) {
	attr := Attributes{Depth: 1.0}
	got := Depth(attr)
	want := math.Vec3{X: 1.0, Y: 1.0, Z: 1.0}
	if !vec3Equal(got, want) {
		t.Errorf("Depth(1.0) = %v, want %v", got, want)
	}
}

func TestDepth_IgnoresColor(t *testing.T) {
	attr1 := Attributes{Depth: 0.7, Color: math.Vec3{X: 1.0, Y: 0.0, Z: 0.0}}
	attr2 := Attributes{Depth: 0.7, Color: math.Vec3{X: 0.0, Y: 1.0, Z: 0.0}}
	got1 := Depth(attr1)
	got2 := Depth(attr2)
	if !vec3Equal(got1, got2) {
		t.Errorf("Depth() changed with color: %v vs %v", got1, got2)
	}
}

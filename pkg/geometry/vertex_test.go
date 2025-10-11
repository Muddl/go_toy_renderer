package geometry

import (
	"testing"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

// TestVertex_Creation_CreatesVertexWithPositionAndColor tests vertex creation
func TestVertex_Creation_CreatesVertexWithPositionAndColor(t *testing.T) {
	position := math.Vec3{X: 1.0, Y: 2.0, Z: 3.0}
	color := math.Vec3{X: 1.0, Y: 0.0, Z: 0.0} // Red

	vertex := Vertex{
		Position: position,
		Color:    color,
	}

	if !vertex.Position.Equals(position, 0.0001) {
		t.Errorf("Vertex.Position = %v, want %v", vertex.Position, position)
	}

	if !vertex.Color.Equals(color, 0.0001) {
		t.Errorf("Vertex.Color = %v, want %v", vertex.Color, color)
	}
}

// TestVertex_NewVertex_CreatesVertexWithGivenValues tests NewVertex constructor
func TestVertex_NewVertex_CreatesVertexWithGivenValues(t *testing.T) {
	tests := []struct {
		name     string
		position math.Vec3
		color    math.Vec3
	}{
		{
			name:     "origin with white color",
			position: math.Vec3{X: 0.0, Y: 0.0, Z: 0.0},
			color:    math.Vec3{X: 1.0, Y: 1.0, Z: 1.0},
		},
		{
			name:     "arbitrary position with red color",
			position: math.Vec3{X: 5.0, Y: -3.0, Z: 2.5},
			color:    math.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
		},
		{
			name:     "negative position with blue color",
			position: math.Vec3{X: -1.0, Y: -2.0, Z: -3.0},
			color:    math.Vec3{X: 0.0, Y: 0.0, Z: 1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vertex := NewVertex(tt.position, tt.color)

			if !vertex.Position.Equals(tt.position, 0.0001) {
				t.Errorf("NewVertex().Position = %v, want %v", vertex.Position, tt.position)
			}

			if !vertex.Color.Equals(tt.color, 0.0001) {
				t.Errorf("NewVertex().Color = %v, want %v", vertex.Color, tt.color)
			}
		})
	}
}

// TestVertex_Equals_ReturnsTrueForEqualVertices tests vertex equality
func TestVertex_Equals_ReturnsTrueForEqualVertices(t *testing.T) {
	v1 := NewVertex(
		math.Vec3{X: 1.0, Y: 2.0, Z: 3.0},
		math.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
	)
	v2 := NewVertex(
		math.Vec3{X: 1.0, Y: 2.0, Z: 3.0},
		math.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
	)

	if !v1.Equals(v2, 0.0001) {
		t.Errorf("Equals() = false, want true for equal vertices")
	}
}

// TestVertex_Equals_ReturnsFalseForDifferentVertices tests vertex inequality
func TestVertex_Equals_ReturnsFalseForDifferentVertices(t *testing.T) {
	tests := []struct {
		name string
		v1   Vertex
		v2   Vertex
	}{
		{
			name: "different positions",
			v1: NewVertex(
				math.Vec3{X: 1.0, Y: 2.0, Z: 3.0},
				math.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
			),
			v2: NewVertex(
				math.Vec3{X: 2.0, Y: 2.0, Z: 3.0},
				math.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
			),
		},
		{
			name: "different colors",
			v1: NewVertex(
				math.Vec3{X: 1.0, Y: 2.0, Z: 3.0},
				math.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
			),
			v2: NewVertex(
				math.Vec3{X: 1.0, Y: 2.0, Z: 3.0},
				math.Vec3{X: 0.0, Y: 1.0, Z: 0.0},
			),
		},
		{
			name: "different position and color",
			v1: NewVertex(
				math.Vec3{X: 1.0, Y: 2.0, Z: 3.0},
				math.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
			),
			v2: NewVertex(
				math.Vec3{X: 5.0, Y: 6.0, Z: 7.0},
				math.Vec3{X: 0.0, Y: 0.0, Z: 1.0},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.v1.Equals(tt.v2, 0.0001) {
				t.Errorf("Equals() = true, want false for %s", tt.name)
			}
		})
	}
}

// TestVertex_Equals_HandlesEpsilonComparison tests epsilon-based float comparison
func TestVertex_Equals_HandlesEpsilonComparison(t *testing.T) {
	v1 := NewVertex(
		math.Vec3{X: 1.0, Y: 2.0, Z: 3.0},
		math.Vec3{X: 1.0, Y: 0.0, Z: 0.0},
	)
	// Nearly equal within epsilon
	v2 := NewVertex(
		math.Vec3{X: 1.00001, Y: 2.00001, Z: 3.00001},
		math.Vec3{X: 1.00001, Y: 0.00001, Z: 0.00001},
	)

	if !v1.Equals(v2, 0.001) {
		t.Errorf("Equals() = false, want true for vertices within epsilon")
	}
}

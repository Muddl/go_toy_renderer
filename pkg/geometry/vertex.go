// Package geometry provides types for 3D geometry primitives
package geometry

import (
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// Vertex represents a single vertex in 3D space with position, color, and normal attributes.
type Vertex struct {
	Position math.Vec3 // Position in 3D space
	Color    math.Vec3 // RGB color (each component in range [0, 1])
	Normal   math.Vec3 // Surface normal (unit vector)
}

// NewVertex creates a new vertex with the specified position and color.
func NewVertex(position, color math.Vec3) Vertex {
	return Vertex{
		Position: position,
		Color:    color,
	}
}

// Equals checks if two vertices are approximately equal within the given epsilon tolerance.
func (v Vertex) Equals(other Vertex, epsilon float64) bool {
	return v.Position.Equals(other.Position, epsilon) &&
		v.Color.Equals(other.Color, epsilon)
}

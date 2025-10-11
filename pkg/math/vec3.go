// Package math provides mathematical primitives for 3D graphics.
// This includes vectors, matrices, and transformation utilities.
package math

import (
	"fmt"
	stdmath "math"
)

// Vec3 represents a 3D vector with X, Y, and Z components.
// Used for positions, directions, colors, and normals in 3D space.
type Vec3 struct {
	X, Y, Z float64
}

// Add returns the sum of two vectors (component-wise addition).
// This operation is commutative: v1.Add(v2) == v2.Add(v1)
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

// Subtract returns the difference of two vectors (component-wise subtraction).
// This operation represents the vector from 'other' to 'v'.
func (v Vec3) Subtract(other Vec3) Vec3 {
	return Vec3{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Scale returns the vector scaled by a scalar value (component-wise multiplication).
// Useful for scaling directions, adjusting magnitudes, and vector operations.
func (v Vec3) Scale(scalar float64) Vec3 {
	return Vec3{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

// Dot returns the dot product (scalar product) of two vectors.
// The dot product measures how much two vectors point in the same direction.
// Returns 0 for perpendicular vectors, positive for similar directions, negative for opposite.
func (v Vec3) Dot(other Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross returns the cross product of two vectors.
// The result is a vector perpendicular to both input vectors.
// The magnitude equals the area of the parallelogram formed by the vectors.
// Uses right-handed coordinate system: i×j=k, j×k=i, k×i=j
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Length returns the magnitude (length) of the vector.
// Calculated as sqrt(x² + y² + z²).
func (v Vec3) Length() float64 {
	return sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns a unit vector (length 1) in the same direction as v.
// If v is a zero vector, returns zero vector to avoid division by zero.
func (v Vec3) Normalize() Vec3 {
	length := v.Length()
	if length == 0 {
		return Vec3{0, 0, 0}
	}
	return v.Scale(1.0 / length)
}

// Distance returns the Euclidean distance between two points represented as vectors.
// Calculated as the length of the vector between the two points.
func (v Vec3) Distance(other Vec3) float64 {
	return v.Subtract(other).Length()
}

// Equals checks if two vectors are approximately equal within the given epsilon tolerance.
// This is useful for comparing floating-point vectors where exact equality is unreliable.
func (v Vec3) Equals(other Vec3, epsilon float64) bool {
	return stdmath.Abs(v.X-other.X) <= epsilon &&
		stdmath.Abs(v.Y-other.Y) <= epsilon &&
		stdmath.Abs(v.Z-other.Z) <= epsilon
}

// String returns a string representation of the vector for debugging and logging.
func (v Vec3) String() string {
	return fmt.Sprintf("Vec3(%.4f, %.4f, %.4f)", v.X, v.Y, v.Z)
}

// sqrt is a helper function wrapper around math.Sqrt for cleaner code.
// Using Go's built-in math.Sqrt for square root calculation.
// Import aliased as stdmath to avoid package name conflict.
func sqrt(x float64) float64 {
	return stdmath.Sqrt(x)
}

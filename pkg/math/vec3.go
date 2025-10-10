// Package math provides mathematical primitives for 3D graphics.
// This includes vectors, matrices, and transformation utilities.
package math

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

package math

import (
	stdmath "math"
)

// Mat4x4 represents a 4x4 matrix in column-major order.
// Used for transformations in 3D graphics (translation, rotation, scale, projection).
// Elements are stored as [m00, m10, m20, m30, m01, m11, m21, m31, ...]
// where mRC means row R, column C.
type Mat4x4 [16]float64

// NewIdentity creates a 4x4 identity matrix.
// The identity matrix has 1s on the diagonal and 0s elsewhere.
// Multiplying any matrix by the identity returns the original matrix.
func NewIdentity() Mat4x4 {
	return Mat4x4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

// NewZero creates a 4x4 zero matrix.
// All elements are 0.
func NewZero() Mat4x4 {
	return Mat4x4{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
}

// Get returns the matrix element at the specified row and column.
// Row and column indices are 0-based (0-3).
// Matrix is stored in column-major order.
func (m Mat4x4) Get(row, col int) float64 {
	return m[col*4+row]
}

// Set modifies the matrix element at the specified row and column.
// Row and column indices are 0-based (0-3).
// Matrix is stored in column-major order.
func (m *Mat4x4) Set(row, col int, value float64) {
	m[col*4+row] = value
}

// Multiply returns the product of two matrices (m * other).
// Matrix multiplication is not commutative: A*B ≠ B*A in general.
// Used for combining transformations: result = m * other.
func (m Mat4x4) Multiply(other Mat4x4) Mat4x4 {
	var result Mat4x4

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			sum := 0.0
			for k := 0; k < 4; k++ {
				sum += m.Get(row, k) * other.Get(k, col)
			}
			result.Set(row, col, sum)
		}
	}

	return result
}

// MultiplyVec4 multiplies the matrix by a 4-component homogeneous vector (x, y, z, w).
// Returns all four components without performing perspective divide.
// Use w=1 for points, w=0 for direction vectors.
// Required for projection transformations where the W component encodes depth.
func (m Mat4x4) MultiplyVec4(x, y, z, w float64) (float64, float64, float64, float64) {
	rx := m.Get(0, 0)*x + m.Get(0, 1)*y + m.Get(0, 2)*z + m.Get(0, 3)*w
	ry := m.Get(1, 0)*x + m.Get(1, 1)*y + m.Get(1, 2)*z + m.Get(1, 3)*w
	rz := m.Get(2, 0)*x + m.Get(2, 1)*y + m.Get(2, 2)*z + m.Get(2, 3)*w
	rw := m.Get(3, 0)*x + m.Get(3, 1)*y + m.Get(3, 2)*z + m.Get(3, 3)*w

	return rx, ry, rz, rw
}

// MultiplyVec3 multiplies the matrix by a 3D vector, treating it as a homogeneous coordinate (w=1).
// Returns the transformed vector (x, y, z components only).
// Used for transforming positions in 3D space.
func (m Mat4x4) MultiplyVec3(v Vec3) Vec3 {
	// Treat Vec3 as homogeneous coordinate with w=1
	x := m.Get(0, 0)*v.X + m.Get(0, 1)*v.Y + m.Get(0, 2)*v.Z + m.Get(0, 3)
	y := m.Get(1, 0)*v.X + m.Get(1, 1)*v.Y + m.Get(1, 2)*v.Z + m.Get(1, 3)
	z := m.Get(2, 0)*v.X + m.Get(2, 1)*v.Y + m.Get(2, 2)*v.Z + m.Get(2, 3)

	return Vec3{x, y, z}
}

// Transpose returns the transpose of the matrix (rows become columns).
// For a matrix M, transpose(M)[i][j] = M[j][i].
// The transpose of the transpose equals the original matrix.
func (m Mat4x4) Transpose() Mat4x4 {
	var result Mat4x4

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			result.Set(row, col, m.Get(col, row))
		}
	}

	return result
}

// Inverse returns the inverse of the matrix and true, or a zero matrix and false
// if the matrix is singular (determinant ≈ 0).
func (m Mat4x4) Inverse() (Mat4x4, bool) {
	// Cofactor expansion of a 4x4 matrix.
	// Using direct element access for clarity (column-major: m[col*4+row]).
	a := func(r, c int) float64 { return m[c*4+r] }

	s0 := a(0, 0)*a(1, 1) - a(1, 0)*a(0, 1)
	s1 := a(0, 0)*a(1, 2) - a(1, 0)*a(0, 2)
	s2 := a(0, 0)*a(1, 3) - a(1, 0)*a(0, 3)
	s3 := a(0, 1)*a(1, 2) - a(1, 1)*a(0, 2)
	s4 := a(0, 1)*a(1, 3) - a(1, 1)*a(0, 3)
	s5 := a(0, 2)*a(1, 3) - a(1, 2)*a(0, 3)

	c5 := a(2, 2)*a(3, 3) - a(3, 2)*a(2, 3)
	c4 := a(2, 1)*a(3, 3) - a(3, 1)*a(2, 3)
	c3 := a(2, 1)*a(3, 2) - a(3, 1)*a(2, 2)
	c2 := a(2, 0)*a(3, 3) - a(3, 0)*a(2, 3)
	c1 := a(2, 0)*a(3, 2) - a(3, 0)*a(2, 2)
	c0 := a(2, 0)*a(3, 1) - a(3, 0)*a(2, 1)

	det := s0*c5 - s1*c4 + s2*c3 + s3*c2 - s4*c1 + s5*c0
	if stdmath.Abs(det) < 1e-12 {
		return Mat4x4{}, false
	}

	invDet := 1.0 / det
	var inv Mat4x4
	inv.Set(0, 0, (a(1, 1)*c5-a(1, 2)*c4+a(1, 3)*c3)*invDet)
	inv.Set(0, 1, (-a(0, 1)*c5+a(0, 2)*c4-a(0, 3)*c3)*invDet)
	inv.Set(0, 2, (a(3, 1)*s5-a(3, 2)*s4+a(3, 3)*s3)*invDet)
	inv.Set(0, 3, (-a(2, 1)*s5+a(2, 2)*s4-a(2, 3)*s3)*invDet)

	inv.Set(1, 0, (-a(1, 0)*c5+a(1, 2)*c2-a(1, 3)*c1)*invDet)
	inv.Set(1, 1, (a(0, 0)*c5-a(0, 2)*c2+a(0, 3)*c1)*invDet)
	inv.Set(1, 2, (-a(3, 0)*s5+a(3, 2)*s2-a(3, 3)*s1)*invDet)
	inv.Set(1, 3, (a(2, 0)*s5-a(2, 2)*s2+a(2, 3)*s1)*invDet)

	inv.Set(2, 0, (a(1, 0)*c4-a(1, 1)*c2+a(1, 3)*c0)*invDet)
	inv.Set(2, 1, (-a(0, 0)*c4+a(0, 1)*c2-a(0, 3)*c0)*invDet)
	inv.Set(2, 2, (a(3, 0)*s4-a(3, 1)*s2+a(3, 3)*s0)*invDet)
	inv.Set(2, 3, (-a(2, 0)*s4+a(2, 1)*s2-a(2, 3)*s0)*invDet)

	inv.Set(3, 0, (-a(1, 0)*c3+a(1, 1)*c1-a(1, 2)*c0)*invDet)
	inv.Set(3, 1, (a(0, 0)*c3-a(0, 1)*c1+a(0, 2)*c0)*invDet)
	inv.Set(3, 2, (-a(3, 0)*s3+a(3, 1)*s1-a(3, 2)*s0)*invDet)
	inv.Set(3, 3, (a(2, 0)*s3-a(2, 1)*s1+a(2, 2)*s0)*invDet)

	return inv, true
}

// NormalMatrix computes transpose(inverse(model)) for transforming normals.
// If the model matrix is singular, returns identity as a safe fallback.
func NormalMatrix(model Mat4x4) Mat4x4 {
	inv, ok := model.Inverse()
	if !ok {
		return NewIdentity()
	}
	return inv.Transpose()
}

// NewTranslation creates a translation matrix that moves points by (tx, ty, tz).
// When applied to a vector (x, y, z), produces (x+tx, y+ty, z+tz).
func NewTranslation(tx, ty, tz float64) Mat4x4 {
	return Mat4x4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		tx, ty, tz, 1,
	}
}

// NewScale creates a scale matrix that scales by (sx, sy, sz).
// When applied to a vector (x, y, z), produces (x*sx, y*sy, z*sz).
func NewScale(sx, sy, sz float64) Mat4x4 {
	return Mat4x4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}

// NewRotationX creates a rotation matrix around the X axis by angle (in radians).
// Positive angle rotates counter-clockwise when looking from positive X towards origin.
// Uses right-handed coordinate system.
func NewRotationX(angle float64) Mat4x4 {
	c := stdmath.Cos(angle)
	s := stdmath.Sin(angle)

	return Mat4x4{
		1, 0, 0, 0,
		0, c, s, 0,
		0, -s, c, 0,
		0, 0, 0, 1,
	}
}

// NewRotationY creates a rotation matrix around the Y axis by angle (in radians).
// Positive angle rotates counter-clockwise when looking from positive Y towards origin.
// Uses right-handed coordinate system.
func NewRotationY(angle float64) Mat4x4 {
	c := stdmath.Cos(angle)
	s := stdmath.Sin(angle)

	return Mat4x4{
		c, 0, -s, 0,
		0, 1, 0, 0,
		s, 0, c, 0,
		0, 0, 0, 1,
	}
}

// NewRotationZ creates a rotation matrix around the Z axis by angle (in radians).
// Positive angle rotates counter-clockwise when looking from positive Z towards origin.
// Uses right-handed coordinate system.
func NewRotationZ(angle float64) Mat4x4 {
	c := stdmath.Cos(angle)
	s := stdmath.Sin(angle)

	return Mat4x4{
		c, s, 0, 0,
		-s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}

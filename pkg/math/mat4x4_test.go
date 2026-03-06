package math

import (
	"math"
	"testing"
)

// Helper function to compare Mat4x4 with epsilon tolerance.
func mat4Equals(m1, m2 Mat4x4, tolerance float64) bool {
	for i := 0; i < 16; i++ {
		if math.Abs(m1[i]-m2[i]) > tolerance {
			return false
		}
	}
	return true
}

// TestMat4x4_NewIdentity tests identity matrix creation.
func TestMat4x4_NewIdentity_CreatesIdentityMatrix(t *testing.T) {
	m := NewIdentity()

	// Identity matrix should have 1s on the diagonal and 0s elsewhere
	expected := Mat4x4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	if !mat4Equals(m, expected, epsilon) {
		t.Errorf("NewIdentity() = %v, want %v", m, expected)
	}
}

// TestMat4x4_NewZero tests zero matrix creation.
func TestMat4x4_NewZero_CreatesZeroMatrix(t *testing.T) {
	m := NewZero()

	// Zero matrix should have all zeros
	expected := Mat4x4{
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}

	if !mat4Equals(m, expected, epsilon) {
		t.Errorf("NewZero() = %v, want %v", m, expected)
	}
}

// TestMat4x4_Get tests element access.
func TestMat4x4_Get_ReturnsCorrectElement(t *testing.T) {
	// Matrix stored in column-major order: [col0, col1, col2, col3]
	// This represents the visual matrix:
	//  1  5  9 13
	//  2  6 10 14
	//  3  7 11 15
	//  4  8 12 16
	m := Mat4x4{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	tests := []struct {
		name     string
		row      int
		col      int
		expected float64
	}{
		{"top-left", 0, 0, 1},
		{"top-right", 0, 3, 13},
		{"bottom-left", 3, 0, 4},
		{"bottom-right", 3, 3, 16},
		{"middle", 1, 2, 10},
		{"second row first col", 1, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := m.Get(tt.row, tt.col)
			if math.Abs(result-tt.expected) > epsilon {
				t.Errorf("Get(%d, %d) = %f, want %f", tt.row, tt.col, result, tt.expected)
			}
		})
	}
}

// TestMat4x4_Set tests element modification.
func TestMat4x4_Set_ModifiesElement(t *testing.T) {
	m := NewZero()
	m.Set(1, 2, 42.0)

	result := m.Get(1, 2)
	expected := 42.0

	if math.Abs(result-expected) > epsilon {
		t.Errorf("After Set(1, 2, 42.0), Get(1, 2) = %f, want %f", result, expected)
	}

	// Verify other elements are still zero
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i == 1 && j == 2 {
				continue // Skip the modified element
			}
			if math.Abs(m.Get(i, j)) > epsilon {
				t.Errorf("Get(%d, %d) should be 0, got %f", i, j, m.Get(i, j))
			}
		}
	}
}

// TestMat4x4_Multiply_WithIdentity tests matrix multiplication with identity.
func TestMat4x4_Multiply_WithIdentity_ReturnsOriginal(t *testing.T) {
	m := Mat4x4{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}
	identity := NewIdentity()

	result := m.Multiply(identity)

	if !mat4Equals(result, m, epsilon) {
		t.Errorf("Multiply with identity should return original matrix")
	}

	// Test commutativity with identity
	result2 := identity.Multiply(m)
	if !mat4Equals(result2, m, epsilon) {
		t.Errorf("Identity multiply should be commutative")
	}
}

// TestMat4x4_Multiply_WithZero tests matrix multiplication with zero.
func TestMat4x4_Multiply_WithZero_ReturnsZero(t *testing.T) {
	m := Mat4x4{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}
	zero := NewZero()

	result := m.Multiply(zero)
	expected := NewZero()

	if !mat4Equals(result, expected, epsilon) {
		t.Errorf("Multiply with zero should return zero matrix")
	}
}

// TestMat4x4_Multiply_General tests general matrix multiplication.
func TestMat4x4_Multiply_General_ReturnsCorrectProduct(t *testing.T) {
	// Simple 2x2-like matrices in 4x4 form for easy verification
	m1 := Mat4x4{
		1, 0, 0, 0,
		0, 2, 0, 0,
		0, 0, 3, 0,
		0, 0, 0, 1,
	}
	m2 := Mat4x4{
		2, 0, 0, 0,
		0, 3, 0, 0,
		0, 0, 4, 0,
		0, 0, 0, 1,
	}

	result := m1.Multiply(m2)
	expected := Mat4x4{
		2, 0, 0, 0,
		0, 6, 0, 0,
		0, 0, 12, 0,
		0, 0, 0, 1,
	}

	if !mat4Equals(result, expected, epsilon) {
		t.Errorf("Multiply() = %v, want %v", result, expected)
	}
}

// TestMat4x4_MultiplyVec3 tests matrix-vector multiplication.
func TestMat4x4_MultiplyVec3_WithIdentity_ReturnsOriginal(t *testing.T) {
	identity := NewIdentity()
	v := Vec3{1.0, 2.0, 3.0}

	result := identity.MultiplyVec3(v)

	if !vec3Equals(result, v, epsilon) {
		t.Errorf("Identity * Vec3 = %v, want %v", result, v)
	}
}

// TestMat4x4_MultiplyVec3_Scale tests scaling transformation.
func TestMat4x4_MultiplyVec3_Scale_ScalesVector(t *testing.T) {
	// Scale matrix: 2x in X, 3x in Y, 4x in Z
	scale := Mat4x4{
		2, 0, 0, 0,
		0, 3, 0, 0,
		0, 0, 4, 0,
		0, 0, 0, 1,
	}
	v := Vec3{1.0, 1.0, 1.0}
	expected := Vec3{2.0, 3.0, 4.0}

	result := scale.MultiplyVec3(v)

	if !vec3Equals(result, expected, epsilon) {
		t.Errorf("Scale * Vec3 = %v, want %v", result, expected)
	}
}

// TestMat4x4_Transpose tests matrix transposition.
func TestMat4x4_Transpose_SwapsRowsAndColumns(t *testing.T) {
	// Original matrix in column-major:
	//  1  5  9 13
	//  2  6 10 14
	//  3  7 11 15
	//  4  8 12 16
	m := Mat4x4{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	result := m.Transpose()
	// Transposed matrix (rows become columns):
	//  1  2  3  4
	//  5  6  7  8
	//  9 10 11 12
	// 13 14 15 16
	expected := Mat4x4{
		1, 5, 9, 13,
		2, 6, 10, 14,
		3, 7, 11, 15,
		4, 8, 12, 16,
	}

	if !mat4Equals(result, expected, epsilon) {
		t.Errorf("Transpose() = %v, want %v", result, expected)
	}
}

// TestMat4x4_Transpose_Identity tests transposing identity matrix.
func TestMat4x4_Transpose_Identity_ReturnsIdentity(t *testing.T) {
	identity := NewIdentity()
	result := identity.Transpose()

	if !mat4Equals(result, identity, epsilon) {
		t.Errorf("Transpose of identity should be identity")
	}
}

// TestMat4x4_NewTranslation tests translation matrix creation.
func TestMat4x4_NewTranslation_CreatesTranslationMatrix(t *testing.T) {
	translation := NewTranslation(10, 20, 30)

	// Apply to origin
	origin := Vec3{0, 0, 0}
	result := translation.MultiplyVec3(origin)
	expected := Vec3{10, 20, 30}

	if !vec3Equals(result, expected, epsilon) {
		t.Errorf("Translation of origin = %v, want %v", result, expected)
	}

	// Apply to non-origin point
	point := Vec3{1, 2, 3}
	result2 := translation.MultiplyVec3(point)
	expected2 := Vec3{11, 22, 33}

	if !vec3Equals(result2, expected2, epsilon) {
		t.Errorf("Translation of point = %v, want %v", result2, expected2)
	}
}

// TestMat4x4_NewScale tests scale matrix creation.
func TestMat4x4_NewScale_CreatesScaleMatrix(t *testing.T) {
	scale := NewScale(2, 3, 4)

	point := Vec3{1, 1, 1}
	result := scale.MultiplyVec3(point)
	expected := Vec3{2, 3, 4}

	if !vec3Equals(result, expected, epsilon) {
		t.Errorf("Scale transformation = %v, want %v", result, expected)
	}
}

// TestMat4x4_NewRotationX tests X-axis rotation matrix.
func TestMat4x4_NewRotationX_RotatesAroundXAxis(t *testing.T) {
	// 90 degree rotation around X axis
	rotation := NewRotationX(math.Pi / 2)

	// Rotate point on Y axis, should move to Z axis
	point := Vec3{0, 1, 0}
	result := rotation.MultiplyVec3(point)
	expected := Vec3{0, 0, 1}

	if !vec3Equals(result, expected, epsilon) {
		t.Errorf("RotationX(90°) of Y-axis point = %v, want %v", result, expected)
	}
}

// TestMat4x4_NewRotationY tests Y-axis rotation matrix.
func TestMat4x4_NewRotationY_RotatesAroundYAxis(t *testing.T) {
	// 90 degree rotation around Y axis
	rotation := NewRotationY(math.Pi / 2)

	// Rotate point on Z axis, should move to X axis
	point := Vec3{0, 0, 1}
	result := rotation.MultiplyVec3(point)
	expected := Vec3{1, 0, 0}

	if !vec3Equals(result, expected, epsilon) {
		t.Errorf("RotationY(90°) of Z-axis point = %v, want %v", result, expected)
	}
}

// TestMat4x4_NewRotationZ tests Z-axis rotation matrix.
func TestMat4x4_NewRotationZ_RotatesAroundZAxis(t *testing.T) {
	// 90 degree rotation around Z axis
	rotation := NewRotationZ(math.Pi / 2)

	// Rotate point on X axis, should move to Y axis
	point := Vec3{1, 0, 0}
	result := rotation.MultiplyVec3(point)
	expected := Vec3{0, 1, 0}

	if !vec3Equals(result, expected, epsilon) {
		t.Errorf("RotationZ(90°) of X-axis point = %v, want %v", result, expected)
	}
}

// TestMat4x4_MultiplyVec4_WithIdentity tests that identity passes all components through.
func TestMat4x4_MultiplyVec4_WithIdentity_ReturnsOriginal(t *testing.T) {
	identity := NewIdentity()

	rx, ry, rz, rw := identity.MultiplyVec4(1.0, 2.0, 3.0, 1.0)

	if math.Abs(rx-1.0) > epsilon || math.Abs(ry-2.0) > epsilon ||
		math.Abs(rz-3.0) > epsilon || math.Abs(rw-1.0) > epsilon {
		t.Errorf("MultiplyVec4 with identity = (%f,%f,%f,%f), want (1,2,3,1)", rx, ry, rz, rw)
	}
}

// TestMat4x4_MultiplyVec4_TranslationMatrix tests that translation affects position with w=1.
func TestMat4x4_MultiplyVec4_TranslationMatrix_TranslatesPosition(t *testing.T) {
	translation := NewTranslation(5.0, 3.0, -2.0)

	rx, ry, rz, rw := translation.MultiplyVec4(1.0, 1.0, 1.0, 1.0)

	if math.Abs(rx-6.0) > epsilon || math.Abs(ry-4.0) > epsilon ||
		math.Abs(rz-(-1.0)) > epsilon || math.Abs(rw-1.0) > epsilon {
		t.Errorf("MultiplyVec4 with translation = (%f,%f,%f,%f), want (6,4,-1,1)", rx, ry, rz, rw)
	}
}

// TestMat4x4_MultiplyVec4_TranslationMatrix_DoesNotTranslateDirection tests w=0 direction vector.
func TestMat4x4_MultiplyVec4_TranslationMatrix_DoesNotTranslateDirection(t *testing.T) {
	// With w=0 (direction vector), translation should have no effect
	translation := NewTranslation(5.0, 3.0, -2.0)

	rx, ry, rz, rw := translation.MultiplyVec4(1.0, 0.0, 0.0, 0.0)

	if math.Abs(rx-1.0) > epsilon || math.Abs(ry-0.0) > epsilon ||
		math.Abs(rz-0.0) > epsilon || math.Abs(rw-0.0) > epsilon {
		t.Errorf("MultiplyVec4 direction with translation = (%f,%f,%f,%f), want (1,0,0,0)", rx, ry, rz, rw)
	}
}

func TestMat4x4_Inverse_Identity_ReturnsIdentity(t *testing.T) {
	inv, ok := NewIdentity().Inverse()
	if !ok {
		t.Fatal("identity should be invertible")
	}
	if !mat4Equals(inv, NewIdentity(), epsilon) {
		t.Errorf("inverse of identity should be identity")
	}
}

func TestMat4x4_Inverse_Translation_ReturnsNegatedTranslation(t *testing.T) {
	m := NewTranslation(3, 4, 5)
	inv, ok := m.Inverse()
	if !ok {
		t.Fatal("translation should be invertible")
	}
	product := m.Multiply(inv)
	if !mat4Equals(product, NewIdentity(), epsilon) {
		t.Errorf("M * M^-1 should be identity, got %v", product)
	}
}

func TestMat4x4_Inverse_Scale_ReturnsReciprocalScale(t *testing.T) {
	m := NewScale(2, 3, 4)
	inv, ok := m.Inverse()
	if !ok {
		t.Fatal("scale should be invertible")
	}
	product := m.Multiply(inv)
	if !mat4Equals(product, NewIdentity(), epsilon) {
		t.Errorf("M * M^-1 should be identity, got %v", product)
	}
}

func TestMat4x4_Inverse_Singular_ReturnsFalse(t *testing.T) {
	// Zero matrix is singular.
	_, ok := NewZero().Inverse()
	if ok {
		t.Error("zero matrix should not be invertible")
	}
}

func TestMat4x4_NormalMatrix_UniformScale(t *testing.T) {
	// For uniform scale, normalMatrix should undo the scale for normals.
	model := NewScale(2, 2, 2).Multiply(NewRotationY(0.5))
	nm := NormalMatrix(model)
	// Apply normal matrix to Y-axis normal.
	ny := nm.MultiplyVec3(Vec3{Y: 1})
	// Should still point roughly in Y direction, unit length after normalize.
	nyn := ny.Normalize()
	if math.Abs(nyn.Y-1.0) > 0.001 {
		t.Errorf("normal matrix of uniform scale+rotation should preserve Y normal, got %v", nyn)
	}
}

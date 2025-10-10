package math

import (
	"math"
	"testing"
)

// epsilon for floating point comparisons
const epsilon = 0.0001

// Helper function to compare Vec3 with epsilon tolerance
func vec3Equals(v1, v2 Vec3, tolerance float64) bool {
	return math.Abs(v1.X-v2.X) < tolerance &&
		math.Abs(v1.Y-v2.Y) < tolerance &&
		math.Abs(v1.Z-v2.Z) < tolerance
}

// TestVec3_Add tests vector addition operation
func TestVec3_Add_ReturnsSumOfVectors(t *testing.T) {
	tests := []struct {
		name     string
		v1       Vec3
		v2       Vec3
		expected Vec3
	}{
		{
			name:     "positive vectors",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{4.0, 5.0, 6.0},
			expected: Vec3{5.0, 7.0, 9.0},
		},
		{
			name:     "negative vectors",
			v1:       Vec3{-1.0, -2.0, -3.0},
			v2:       Vec3{-4.0, -5.0, -6.0},
			expected: Vec3{-5.0, -7.0, -9.0},
		},
		{
			name:     "mixed positive and negative",
			v1:       Vec3{1.0, -2.0, 3.0},
			v2:       Vec3{-4.0, 5.0, -6.0},
			expected: Vec3{-3.0, 3.0, -3.0},
		},
		{
			name:     "zero vector",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{0.0, 0.0, 0.0},
			expected: Vec3{1.0, 2.0, 3.0},
		},
		{
			name:     "both zero vectors",
			v1:       Vec3{0.0, 0.0, 0.0},
			v2:       Vec3{0.0, 0.0, 0.0},
			expected: Vec3{0.0, 0.0, 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Add(tt.v2)
			if !vec3Equals(result, tt.expected, epsilon) {
				t.Errorf("Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestVec3_Subtract tests vector subtraction operation
func TestVec3_Subtract_ReturnsDifferenceOfVectors(t *testing.T) {
	tests := []struct {
		name     string
		v1       Vec3
		v2       Vec3
		expected Vec3
	}{
		{
			name:     "positive vectors",
			v1:       Vec3{5.0, 7.0, 9.0},
			v2:       Vec3{1.0, 2.0, 3.0},
			expected: Vec3{4.0, 5.0, 6.0},
		},
		{
			name:     "negative vectors",
			v1:       Vec3{-1.0, -2.0, -3.0},
			v2:       Vec3{-4.0, -5.0, -6.0},
			expected: Vec3{3.0, 3.0, 3.0},
		},
		{
			name:     "mixed positive and negative",
			v1:       Vec3{1.0, -2.0, 3.0},
			v2:       Vec3{-4.0, 5.0, -6.0},
			expected: Vec3{5.0, -7.0, 9.0},
		},
		{
			name:     "subtract zero vector",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{0.0, 0.0, 0.0},
			expected: Vec3{1.0, 2.0, 3.0},
		},
		{
			name:     "subtract from itself",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{1.0, 2.0, 3.0},
			expected: Vec3{0.0, 0.0, 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Subtract(tt.v2)
			if !vec3Equals(result, tt.expected, epsilon) {
				t.Errorf("Subtract() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestVec3_Scale tests scalar multiplication operation
func TestVec3_Scale_ReturnsScaledVector(t *testing.T) {
	tests := []struct {
		name     string
		v        Vec3
		scalar   float64
		expected Vec3
	}{
		{
			name:     "scale by positive integer",
			v:        Vec3{1.0, 2.0, 3.0},
			scalar:   2.0,
			expected: Vec3{2.0, 4.0, 6.0},
		},
		{
			name:     "scale by negative integer",
			v:        Vec3{1.0, 2.0, 3.0},
			scalar:   -2.0,
			expected: Vec3{-2.0, -4.0, -6.0},
		},
		{
			name:     "scale by fraction",
			v:        Vec3{2.0, 4.0, 6.0},
			scalar:   0.5,
			expected: Vec3{1.0, 2.0, 3.0},
		},
		{
			name:     "scale by zero",
			v:        Vec3{1.0, 2.0, 3.0},
			scalar:   0.0,
			expected: Vec3{0.0, 0.0, 0.0},
		},
		{
			name:     "scale by one",
			v:        Vec3{1.0, 2.0, 3.0},
			scalar:   1.0,
			expected: Vec3{1.0, 2.0, 3.0},
		},
		{
			name:     "scale zero vector",
			v:        Vec3{0.0, 0.0, 0.0},
			scalar:   5.0,
			expected: Vec3{0.0, 0.0, 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Scale(tt.scalar)
			if !vec3Equals(result, tt.expected, epsilon) {
				t.Errorf("Scale() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestVec3_Construction tests basic Vec3 construction
func TestVec3_Construction_CreatesVectorWithCorrectValues(t *testing.T) {
	v := Vec3{1.0, 2.0, 3.0}

	if v.X != 1.0 {
		t.Errorf("Expected X = 1.0, got %f", v.X)
	}
	if v.Y != 2.0 {
		t.Errorf("Expected Y = 2.0, got %f", v.Y)
	}
	if v.Z != 3.0 {
		t.Errorf("Expected Z = 3.0, got %f", v.Z)
	}
}

// TestVec3_Dot tests dot product operation
func TestVec3_Dot_ReturnsScalarProduct(t *testing.T) {
	tests := []struct {
		name     string
		v1       Vec3
		v2       Vec3
		expected float64
	}{
		{
			name:     "perpendicular vectors (i and j)",
			v1:       Vec3{1.0, 0.0, 0.0},
			v2:       Vec3{0.0, 1.0, 0.0},
			expected: 0.0,
		},
		{
			name:     "parallel vectors (same direction)",
			v1:       Vec3{1.0, 0.0, 0.0},
			v2:       Vec3{2.0, 0.0, 0.0},
			expected: 2.0,
		},
		{
			name:     "parallel vectors (opposite direction)",
			v1:       Vec3{1.0, 0.0, 0.0},
			v2:       Vec3{-1.0, 0.0, 0.0},
			expected: -1.0,
		},
		{
			name:     "general vectors",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{4.0, 5.0, 6.0},
			expected: 32.0, // 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
		},
		{
			name:     "with zero vector",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{0.0, 0.0, 0.0},
			expected: 0.0,
		},
		{
			name:     "with negative components",
			v1:       Vec3{-1.0, 2.0, -3.0},
			v2:       Vec3{4.0, -5.0, 6.0},
			expected: -32.0, // -1*4 + 2*-5 + -3*6 = -4 - 10 - 18 = -32
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Dot(tt.v2)
			if math.Abs(result-tt.expected) > epsilon {
				t.Errorf("Dot() = %f, want %f", result, tt.expected)
			}
		})
	}
}

// TestVec3_Cross tests cross product operation
func TestVec3_Cross_ReturnsPerpendicularVector(t *testing.T) {
	tests := []struct {
		name     string
		v1       Vec3
		v2       Vec3
		expected Vec3
	}{
		{
			name:     "i cross j equals k (right-handed)",
			v1:       Vec3{1.0, 0.0, 0.0},
			v2:       Vec3{0.0, 1.0, 0.0},
			expected: Vec3{0.0, 0.0, 1.0},
		},
		{
			name:     "j cross k equals i (right-handed)",
			v1:       Vec3{0.0, 1.0, 0.0},
			v2:       Vec3{0.0, 0.0, 1.0},
			expected: Vec3{1.0, 0.0, 0.0},
		},
		{
			name:     "k cross i equals j (right-handed)",
			v1:       Vec3{0.0, 0.0, 1.0},
			v2:       Vec3{1.0, 0.0, 0.0},
			expected: Vec3{0.0, 1.0, 0.0},
		},
		{
			name:     "j cross i equals -k (anti-commutative)",
			v1:       Vec3{0.0, 1.0, 0.0},
			v2:       Vec3{1.0, 0.0, 0.0},
			expected: Vec3{0.0, 0.0, -1.0},
		},
		{
			name:     "parallel vectors (zero result)",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{2.0, 4.0, 6.0},
			expected: Vec3{0.0, 0.0, 0.0},
		},
		{
			name:     "general vectors",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{4.0, 5.0, 6.0},
			expected: Vec3{-3.0, 6.0, -3.0}, // (2*6-3*5, 3*4-1*6, 1*5-2*4)
		},
		{
			name:     "with zero vector",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{0.0, 0.0, 0.0},
			expected: Vec3{0.0, 0.0, 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Cross(tt.v2)
			if !vec3Equals(result, tt.expected, epsilon) {
				t.Errorf("Cross() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestVec3_Length tests vector magnitude calculation
func TestVec3_Length_ReturnsVectorMagnitude(t *testing.T) {
	tests := []struct {
		name     string
		v        Vec3
		expected float64
	}{
		{
			name:     "unit vector i",
			v:        Vec3{1.0, 0.0, 0.0},
			expected: 1.0,
		},
		{
			name:     "unit vector j",
			v:        Vec3{0.0, 1.0, 0.0},
			expected: 1.0,
		},
		{
			name:     "unit vector k",
			v:        Vec3{0.0, 0.0, 1.0},
			expected: 1.0,
		},
		{
			name:     "zero vector",
			v:        Vec3{0.0, 0.0, 0.0},
			expected: 0.0,
		},
		{
			name:     "3-4-5 right triangle (pythagorean)",
			v:        Vec3{3.0, 4.0, 0.0},
			expected: 5.0,
		},
		{
			name:     "general vector",
			v:        Vec3{1.0, 2.0, 2.0},
			expected: 3.0, // sqrt(1 + 4 + 4) = sqrt(9) = 3
		},
		{
			name:     "negative components (magnitude always positive)",
			v:        Vec3{-1.0, -2.0, -2.0},
			expected: 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Length()
			if math.Abs(result-tt.expected) > epsilon {
				t.Errorf("Length() = %f, want %f", result, tt.expected)
			}
		})
	}
}

// TestVec3_Normalize tests vector normalization
func TestVec3_Normalize_ReturnsUnitVector(t *testing.T) {
	tests := []struct {
		name     string
		v        Vec3
		expected Vec3
	}{
		{
			name:     "already unit vector",
			v:        Vec3{1.0, 0.0, 0.0},
			expected: Vec3{1.0, 0.0, 0.0},
		},
		{
			name:     "scale down",
			v:        Vec3{2.0, 0.0, 0.0},
			expected: Vec3{1.0, 0.0, 0.0},
		},
		{
			name:     "general vector",
			v:        Vec3{3.0, 4.0, 0.0},
			expected: Vec3{0.6, 0.8, 0.0}, // 3/5, 4/5, 0
		},
		{
			name:     "3D vector",
			v:        Vec3{1.0, 2.0, 2.0},
			expected: Vec3{1.0 / 3.0, 2.0 / 3.0, 2.0 / 3.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v.Normalize()
			if !vec3Equals(result, tt.expected, epsilon) {
				t.Errorf("Normalize() = %v, want %v", result, tt.expected)
			}
			// Also verify the result has length 1
			length := result.Length()
			if math.Abs(length-1.0) > epsilon {
				t.Errorf("Normalize() result length = %f, want 1.0", length)
			}
		})
	}
}

// TestVec3_Normalize_ZeroVector tests normalization edge case with zero vector
func TestVec3_Normalize_ZeroVector_ReturnsZeroVector(t *testing.T) {
	v := Vec3{0.0, 0.0, 0.0}
	result := v.Normalize()
	expected := Vec3{0.0, 0.0, 0.0}

	if !vec3Equals(result, expected, epsilon) {
		t.Errorf("Normalize() on zero vector = %v, want %v", result, expected)
	}
}

// TestVec3_Distance tests distance calculation between two points
func TestVec3_Distance_ReturnsDistanceBetweenPoints(t *testing.T) {
	tests := []struct {
		name     string
		v1       Vec3
		v2       Vec3
		expected float64
	}{
		{
			name:     "same point",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{1.0, 2.0, 3.0},
			expected: 0.0,
		},
		{
			name:     "unit distance along x-axis",
			v1:       Vec3{0.0, 0.0, 0.0},
			v2:       Vec3{1.0, 0.0, 0.0},
			expected: 1.0,
		},
		{
			name:     "3-4-5 triangle",
			v1:       Vec3{0.0, 0.0, 0.0},
			v2:       Vec3{3.0, 4.0, 0.0},
			expected: 5.0,
		},
		{
			name:     "general 3D distance",
			v1:       Vec3{1.0, 2.0, 3.0},
			v2:       Vec3{4.0, 6.0, 8.0},
			expected: math.Sqrt(50.0), // sqrt((4-1)^2 + (6-2)^2 + (8-3)^2) = sqrt(9+16+25)
		},
		{
			name:     "negative coordinates",
			v1:       Vec3{-1.0, -2.0, -3.0},
			v2:       Vec3{2.0, 2.0, 1.0},
			expected: math.Sqrt(41.0), // sqrt(9 + 16 + 16)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Distance(tt.v2)
			if math.Abs(result-tt.expected) > epsilon {
				t.Errorf("Distance() = %f, want %f", result, tt.expected)
			}
		})
	}
}

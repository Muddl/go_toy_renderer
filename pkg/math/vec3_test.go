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

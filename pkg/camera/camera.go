// Package camera provides a virtual camera with view and projection matrix generation.
// The camera uses a right-handed coordinate system where +X is right, +Y is up,
// and +Z points toward the viewer (camera looks down the -Z axis).
// Matrices follow OpenGL column-major convention (matching pkg/math/Mat4x4).
package camera

import (
	stdmath "math"

	"github.com/muddl/go_toy_renderer/pkg/math"
)

// Camera represents a virtual camera in 3D space with perspective projection.
type Camera struct {
	Position math.Vec3 // Camera location in world space
	Target   math.Vec3 // Point the camera looks at in world space
	Up       math.Vec3 // World up vector, typically {0, 1, 0}
	FOV      float64   // Vertical field of view in degrees (e.g. 45, 60, 90)
	Aspect   float64   // Viewport width / height ratio (e.g. 16/9, 800/600)
	Near     float64   // Near clipping plane distance, must be > 0
	Far      float64   // Far clipping plane distance, must be > Near
}

// New creates a Camera with the given parameters.
// position: camera location in world space.
// target: point the camera looks at.
// up: world up vector (usually {0,1,0}).
// fov: vertical field of view in degrees.
// aspect: viewport width/height.
// near, far: near and far clipping plane distances (both must be > 0, far > near).
func New(position, target, up math.Vec3, fov, aspect, near, far float64) Camera {
	return Camera{
		Position: position,
		Target:   target,
		Up:       up,
		FOV:      fov,
		Aspect:   aspect,
		Near:     near,
		Far:      far,
	}
}

// ViewMatrix returns the view matrix using the LookAt algorithm.
// Transforms world space into camera space where:
//   - Camera is at the origin
//   - Camera looks down the -Z axis
//   - +Y is up, +X is right
//
// Formula: rotation_transpose * translation(-eye)
// right   = normalize(cross(up, forward))
// newUp   = cross(forward, right)
// forward = normalize(eye - target)  [points from target toward camera, i.e. +Z in camera space].
func (c Camera) ViewMatrix() math.Mat4x4 {
	// Compute right-handed camera basis vectors
	forward := c.Position.Subtract(c.Target).Normalize() // +Z in camera space
	right := c.Up.Cross(forward).Normalize()             // +X in camera space
	newUp := forward.Cross(right)                        // +Y in camera space (recomputed for orthogonality)

	// Translation components: project negative eye position onto each basis axis
	tx := -right.Dot(c.Position)
	ty := -newUp.Dot(c.Position)
	tz := -forward.Dot(c.Position)

	// Column-major layout (each group of 4 values is one column):
	//   Column 0: right.X,   newUp.X,  forward.X,  0
	//   Column 1: right.Y,   newUp.Y,  forward.Y,  0
	//   Column 2: right.Z,   newUp.Z,  forward.Z,  0
	//   Column 3: tx,        ty,       tz,         1
	return math.Mat4x4{
		right.X, newUp.X, forward.X, 0,
		right.Y, newUp.Y, forward.Y, 0,
		right.Z, newUp.Z, forward.Z, 0,
		tx, ty, tz, 1,
	}
}

// ProjectionMatrix returns a perspective projection matrix.
// Converts camera space to clip space. After perspective divide (dividing xyz by w),
// visible geometry maps into NDC: x,y in [-1,1], z in [-1,1].
// Points in front of the camera have positive W in clip space.
// Near plane maps to NDC z = -1; far plane maps to NDC z = +1.
//
// Formula uses vertical FOV (converted to radians internally).
// f = 1 / tan(fov / 2), rangeInv = 1 / (near - far).
func (c Camera) ProjectionMatrix() math.Mat4x4 {
	fovRad := c.FOV * stdmath.Pi / 180.0
	f := 1.0 / stdmath.Tan(fovRad/2.0)
	rangeInv := 1.0 / (c.Near - c.Far) // negative, since Near < Far

	// Column-major layout:
	//   Column 0: f/aspect, 0, 0, 0
	//   Column 1: 0, f, 0, 0
	//   Column 2: 0, 0, (Near+Far)*rangeInv, -1
	//   Column 3: 0, 0, 2*Near*Far*rangeInv, 0
	return math.Mat4x4{
		f / c.Aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (c.Near + c.Far) * rangeInv, -1,
		0, 0, 2 * c.Near * c.Far * rangeInv, 0,
	}
}

// ViewProjectionMatrix returns the combined view-projection matrix: Projection × View.
// Apply to world-space vertices to obtain clip-space coordinates.
// After perspective divide (clip.xyz / clip.w), you get NDC coordinates.
func (c Camera) ViewProjectionMatrix() math.Mat4x4 {
	return c.ProjectionMatrix().Multiply(c.ViewMatrix())
}

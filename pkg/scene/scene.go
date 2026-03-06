// Package scene provides a scene graph for the GPU rendering path.
// Each Node pairs a geometry.Mesh with a Transform that produces
// a model matrix for per-mesh positioning.
package scene

import (
	"github.com/muddl/go_toy_renderer/pkg/geometry"
	"github.com/muddl/go_toy_renderer/pkg/math"
)

// Transform holds position, rotation (Euler angles in radians), and scale.
// ModelMatrix computes the model matrix as T * Rz * Ry * Rx * S.
type Transform struct {
	Position math.Vec3 // World-space position.
	Rotation math.Vec3 // Euler angles (X, Y, Z) in radians.
	Scale    math.Vec3 // Per-axis scale factors.
}

// NewTransform returns an identity transform (origin, no rotation, uniform scale 1).
func NewTransform() Transform {
	return Transform{
		Scale: math.Vec3{X: 1, Y: 1, Z: 1},
	}
}

// ModelMatrix computes the 4x4 model matrix: Translation * Rz * Ry * Rx * Scale.
func (tr Transform) ModelMatrix() math.Mat4x4 {
	s := math.NewScale(tr.Scale.X, tr.Scale.Y, tr.Scale.Z)
	rx := math.NewRotationX(tr.Rotation.X)
	ry := math.NewRotationY(tr.Rotation.Y)
	rz := math.NewRotationZ(tr.Rotation.Z)
	tl := math.NewTranslation(tr.Position.X, tr.Position.Y, tr.Position.Z)
	return tl.Multiply(rz.Multiply(ry.Multiply(rx.Multiply(s))))
}

// Node is a scene graph entry: a mesh with an independent transform.
type Node struct {
	Mesh      *geometry.Mesh
	Transform Transform
	Static    bool // If true, the renderer skips per-frame animation for this node.
}

// Scene holds nodes to render and camera parameters.
type Scene struct {
	Nodes []Node
}

// NewScene creates an empty scene.
func NewScene() *Scene {
	return &Scene{Nodes: make([]Node, 0)}
}

// AddNode appends a node to the scene.
func (s *Scene) AddNode(n Node) {
	s.Nodes = append(s.Nodes, n)
}

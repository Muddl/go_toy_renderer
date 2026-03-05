// Package shaders embeds WGSL shader source files for use by pkg/gpu.
package shaders

import _ "embed"

// CubeWGSLTemplate is the combined vertex+fragment WGSL shader for the cube.
// The MVP matrix constant contains %f format verbs that are substituted at
// runtime by pkg/gpu with an aspect-correct MVP matrix (camera at (3,2,5)
// looking at origin, fov=60 deg). Uniform buffers replace this in Phase 14.
//
//go:embed cube.wgsl
var CubeWGSLTemplate string

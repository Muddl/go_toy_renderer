// Package shaders embeds WGSL shader source files for use by pkg/gpu.
package shaders

import _ "embed"

// CubeWGSL is the combined vertex+fragment WGSL shader for the cube.
// The vertex shader reads camera (viewProj) and mesh (model) uniforms
// via @group(0) @binding(0) and @group(0) @binding(1) respectively.
//
//go:embed cube.wgsl
var CubeWGSL string

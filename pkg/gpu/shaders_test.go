//go:build !headless

package gpu

import (
	"strings"
	"testing"
)

// TestMakeCubeShaderWGSL_HasRequiredEntryPoints verifies that the shader string
// produced from the embedded template contains the expected WGSL entry points.
func TestMakeCubeShaderWGSL_HasRequiredEntryPoints(t *testing.T) {
	wgsl := makeCubeShaderWGSL(640, 480)
	for _, want := range []string{
		"fn vs_main",
		"fn fs_main",
		"struct VertexInput",
		"struct VertexOutput",
	} {
		if !strings.Contains(wgsl, want) {
			t.Fatalf("makeCubeShaderWGSL: output missing %q", want)
		}
	}
}

// TestMakeCubeShaderWGSL_MVPSubstituted verifies that the %f format verbs are
// replaced with float values and do not appear in the generated WGSL.
func TestMakeCubeShaderWGSL_MVPSubstituted(t *testing.T) {
	const verb = "%f"
	wgsl := makeCubeShaderWGSL(640, 480)
	if strings.Contains(wgsl, verb) {
		t.Fatal("makeCubeShaderWGSL: output still contains unsubstituted format verbs")
	}
}

// TestMakeCubeShaderWGSL_AspectAffectsMVP verifies that different window dimensions
// produce different MVP values (aspect ratio is baked in until Phase 14 uniforms).
func TestMakeCubeShaderWGSL_AspectAffectsMVP(t *testing.T) {
	wgsl1 := makeCubeShaderWGSL(640, 480)
	wgsl2 := makeCubeShaderWGSL(1280, 720)
	if wgsl1 == wgsl2 {
		t.Fatal("makeCubeShaderWGSL: expected different WGSL for different window dimensions")
	}
}

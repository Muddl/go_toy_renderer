//go:build !headless

package gpu

import (
	"strings"
	"testing"

	"github.com/muddl/go_toy_renderer/assets/shaders"
)

func TestCubeShaderWGSL_HasRequiredEntryPoints(t *testing.T) {
	for _, want := range []string{
		"fn vs_main",
		"fn fs_main",
		"struct VertexInput",
		"struct VertexOutput",
	} {
		if !strings.Contains(shaders.CubeWGSL, want) {
			t.Fatalf("CubeWGSL: missing %q", want)
		}
	}
}

func TestCubeShaderWGSL_HasUniformBindings(t *testing.T) {
	for _, want := range []string{
		"@group(0) @binding(0)",
		"@group(0) @binding(1)",
		"var<uniform> camera",
		"var<uniform> mesh",
		"camera.viewProj",
		"mesh.model",
	} {
		if !strings.Contains(shaders.CubeWGSL, want) {
			t.Fatalf("CubeWGSL: missing %q", want)
		}
	}
}

func TestCubeShaderWGSL_NoFormatVerbs(t *testing.T) {
	if strings.Contains(shaders.CubeWGSL, "%f") {
		t.Fatal("CubeWGSL: should not contain format verbs (uniforms replaced template)")
	}
}

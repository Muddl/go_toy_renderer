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
		"@group(0) @binding(2)",
		"var<uniform> camera",
		"var<uniform> mesh",
		"var<uniform> light",
		"camera.viewProj",
		"mesh.model",
		"mesh.normalMatrix",
	} {
		if !strings.Contains(shaders.CubeWGSL, want) {
			t.Fatalf("CubeWGSL: missing %q", want)
		}
	}
}

func TestCubeShaderWGSL_HasLightingStructs(t *testing.T) {
	for _, want := range []string{
		"struct LightUniforms",
		"direction",
		"ambient",
	} {
		if !strings.Contains(shaders.CubeWGSL, want) {
			t.Fatalf("CubeWGSL: missing %q", want)
		}
	}
}

func TestCubeShaderWGSL_HasNormalInput(t *testing.T) {
	if !strings.Contains(shaders.CubeWGSL, "@location(2) normal") {
		t.Fatal("CubeWGSL: missing normal input at location 2")
	}
}

func TestCubeShaderWGSL_NoFormatVerbs(t *testing.T) {
	if strings.Contains(shaders.CubeWGSL, "%f") {
		t.Fatal("CubeWGSL: should not contain format verbs (uniforms replaced template)")
	}
}

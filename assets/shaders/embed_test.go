package shaders_test

import (
	"strings"
	"testing"

	"github.com/muddl/go_toy_renderer/assets/shaders"
)

func TestCubeWGSL_IsNonEmpty(t *testing.T) {
	if len(shaders.CubeWGSL) == 0 {
		t.Fatal("CubeWGSL must not be empty")
	}
}

func TestCubeWGSL_HasVertexEntryPoint(t *testing.T) {
	if !strings.Contains(shaders.CubeWGSL, "fn vs_main") {
		t.Fatal("CubeWGSL: missing vs_main vertex entry point")
	}
}

func TestCubeWGSL_HasFragmentEntryPoint(t *testing.T) {
	if !strings.Contains(shaders.CubeWGSL, "fn fs_main") {
		t.Fatal("CubeWGSL: missing fs_main fragment entry point")
	}
}

func TestCubeWGSL_HasUniformBindings(t *testing.T) {
	for _, want := range []string{
		"@group(0) @binding(0)",
		"@group(0) @binding(1)",
		"var<uniform> camera",
		"var<uniform> mesh",
	} {
		if !strings.Contains(shaders.CubeWGSL, want) {
			t.Fatalf("CubeWGSL: missing %q", want)
		}
	}
}

package shaders_test

import (
	"strings"
	"testing"

	"github.com/muddl/go_toy_renderer/assets/shaders"
)

// TestCubeWGSLTemplate_IsNonEmpty verifies the embedded shader file is present.
func TestCubeWGSLTemplate_IsNonEmpty(t *testing.T) {
	if len(shaders.CubeWGSLTemplate) == 0 {
		t.Fatal("CubeWGSLTemplate must not be empty")
	}
}

// TestCubeWGSLTemplate_HasVertexEntryPoint verifies the template declares vs_main.
func TestCubeWGSLTemplate_HasVertexEntryPoint(t *testing.T) {
	if !strings.Contains(shaders.CubeWGSLTemplate, "fn vs_main") {
		t.Fatal("CubeWGSLTemplate: missing vs_main vertex entry point")
	}
}

// TestCubeWGSLTemplate_HasFragmentEntryPoint verifies the template declares fs_main.
func TestCubeWGSLTemplate_HasFragmentEntryPoint(t *testing.T) {
	if !strings.Contains(shaders.CubeWGSLTemplate, "fn fs_main") {
		t.Fatal("CubeWGSLTemplate: missing fs_main fragment entry point")
	}
}

// TestCubeWGSLTemplate_HasMVPFormatVerbs verifies the template contains %f verbs
// for runtime MVP substitution (Phase 14 replaces these with uniform buffers).
func TestCubeWGSLTemplate_HasMVPFormatVerbs(t *testing.T) {
	const verb = "%f"
	if !strings.Contains(shaders.CubeWGSLTemplate, verb) {
		t.Fatal("CubeWGSLTemplate: must contain format verbs for runtime MVP substitution")
	}
}

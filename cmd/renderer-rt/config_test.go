package main

import (
	"testing"
)

func TestParseConfig_Defaults(t *testing.T) {
	cfg, err := parseConfig([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Width != 1280 {
		t.Errorf("Width = %d, want 1280", cfg.Width)
	}
	if cfg.Height != 720 {
		t.Errorf("Height = %d, want 720", cfg.Height)
	}
	if cfg.Backend != "cpu" {
		t.Errorf("Backend = %q, want %q", cfg.Backend, "cpu")
	}
}

func TestParseConfig_CustomFlags(t *testing.T) {
	cfg, err := parseConfig([]string{"--width", "800", "--height", "600", "--backend", "gpu"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Width != 800 {
		t.Errorf("Width = %d, want 800", cfg.Width)
	}
	if cfg.Height != 600 {
		t.Errorf("Height = %d, want 600", cfg.Height)
	}
	if cfg.Backend != "gpu" {
		t.Errorf("Backend = %q, want %q", cfg.Backend, "gpu")
	}
}

func TestParseConfig_UnknownFlag(t *testing.T) {
	_, err := parseConfig([]string{"--bogus", "value"})
	if err == nil {
		t.Error("expected error for unknown flag, got nil")
	}
}

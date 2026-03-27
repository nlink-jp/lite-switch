package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
}

const validConfig = `
[api]
base_url = "http://localhost:1234/v1"
api_key  = "test-key"

[model]
name = "test-model"
`

const validSwitches = `
switches:
  - tag: foo
    description: Foo things
  - tag: bar
    description: Bar things
`

func TestLoad_Valid(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	swPath := filepath.Join(dir, "switches.yaml")
	writeFile(t, cfgPath, validConfig)
	writeFile(t, swPath, validSwitches)

	cfg, err := Load(cfgPath, swPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.API.BaseURL != "http://localhost:1234/v1" {
		t.Errorf("BaseURL = %q", cfg.API.BaseURL)
	}
	if cfg.API.APIKey != "test-key" {
		t.Errorf("APIKey = %q", cfg.API.APIKey)
	}
	if cfg.Model.Name != "test-model" {
		t.Errorf("Model.Name = %q", cfg.Model.Name)
	}
	if len(cfg.Switches) != 2 {
		t.Fatalf("len(Switches) = %d, want 2", len(cfg.Switches))
	}
	if cfg.Switches[0].Tag != "foo" {
		t.Errorf("Switches[0].Tag = %q", cfg.Switches[0].Tag)
	}
}

func TestLoad_Defaults(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	swPath := filepath.Join(dir, "switches.yaml")
	writeFile(t, cfgPath, validConfig)
	writeFile(t, swPath, validSwitches)

	cfg, err := Load(cfgPath, swPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.API.TimeoutSeconds != 30 {
		t.Errorf("TimeoutSeconds = %d, want 30", cfg.API.TimeoutSeconds)
	}
	if cfg.API.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", cfg.API.MaxRetries)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	swPath := filepath.Join(dir, "switches.yaml")
	writeFile(t, cfgPath, validConfig)
	writeFile(t, swPath, validSwitches)

	t.Setenv("LITE_SWITCH_BASE_URL", "http://env-host")
	t.Setenv("LITE_SWITCH_API_KEY", "env-key")
	t.Setenv("LITE_SWITCH_MODEL", "env-model")

	cfg, err := Load(cfgPath, swPath)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.API.BaseURL != "http://env-host" {
		t.Errorf("BaseURL = %q, want env override", cfg.API.BaseURL)
	}
	if cfg.API.APIKey != "env-key" {
		t.Errorf("APIKey = %q, want env override", cfg.API.APIKey)
	}
	if cfg.Model.Name != "env-model" {
		t.Errorf("Model.Name = %q, want env override", cfg.Model.Name)
	}
}

func TestLoad_MissingConfigFile(t *testing.T) {
	_, err := Load("/nonexistent/config.toml", "/nonexistent/switches.yaml")
	if err == nil {
		t.Error("expected error for missing config file")
	}
}

func TestLoad_InvalidSwitches(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	swPath := filepath.Join(dir, "switches.yaml")
	writeFile(t, cfgPath, validConfig)
	writeFile(t, swPath, "not: valid: yaml: :")

	_, err := Load(cfgPath, swPath)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoad_EmptySwitches(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	swPath := filepath.Join(dir, "switches.yaml")
	writeFile(t, cfgPath, validConfig)
	writeFile(t, swPath, "switches: []")

	_, err := Load(cfgPath, swPath)
	if err == nil {
		t.Error("expected validation error for empty switches")
	}
}

func TestLoad_DuplicateTags(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.toml")
	swPath := filepath.Join(dir, "switches.yaml")
	writeFile(t, cfgPath, validConfig)
	writeFile(t, swPath, `
switches:
  - tag: foo
    description: First foo
  - tag: foo
    description: Duplicate foo
`)

	_, err := Load(cfgPath, swPath)
	if err == nil {
		t.Error("expected validation error for duplicate tags")
	}
}

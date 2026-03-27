// Package config handles loading of lite-switch configuration.
//
// Configuration is split across two files:
//   - System config (TOML): LLM API settings, loaded from ~/.config/lite-switch/config.toml
//   - Switches file (YAML): classification definitions, typically co-located with the project
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Stderr is the writer used for warning messages; can be overridden in tests.
var Stderr io.Writer = os.Stderr

// Switch represents a single classification option.
type Switch struct {
	Tag         string `yaml:"tag"`
	Description string `yaml:"description"`
}

// Config holds the merged runtime configuration.
type Config struct {
	API      APIConfig
	Model    ModelConfig
	Switches []Switch
}

// APIConfig holds LLM API connection settings.
type APIConfig struct {
	BaseURL        string `toml:"base_url"`
	APIKey         string `toml:"api_key"`
	TimeoutSeconds int    `toml:"timeout_seconds"`
	MaxRetries     int    `toml:"max_retries"`
}

// ModelConfig specifies which model to use.
type ModelConfig struct {
	Name string `toml:"name"`
}

// Timeout returns the configured request timeout as a time.Duration.
func (c *Config) Timeout() time.Duration {
	if c.API.TimeoutSeconds <= 0 {
		return 30 * time.Second
	}
	return time.Duration(c.API.TimeoutSeconds) * time.Second
}

// systemFile mirrors the TOML structure of the system config file.
type systemFile struct {
	API   APIConfig   `toml:"api"`
	Model ModelConfig `toml:"model"`
}

// switchesFile mirrors the YAML structure of the switches file.
type switchesFile struct {
	Switches []Switch `yaml:"switches"`
}

// DefaultConfigPath returns the default system config path.
func DefaultConfigPath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("getting home directory: %w", err)
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "lite-switch", "config.toml"), nil
}

// Load reads the system config from configPath and the switches from switchesPath,
// applies environment variable overrides, and validates the result.
func Load(configPath, switchesPath string) (*Config, error) {
	sys, err := loadSystem(configPath)
	if err != nil {
		return nil, err
	}

	sw, err := loadSwitches(switchesPath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		API:      sys.API,
		Model:    sys.Model,
		Switches: sw.Switches,
	}

	applyDefaults(cfg)
	applyEnv(cfg)

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}

func loadSystem(path string) (systemFile, error) {
	info, err := os.Stat(path)
	if err != nil {
		return systemFile{}, fmt.Errorf("reading config %q: %w", path, err)
	}
	checkPermissions(path, info)

	var sys systemFile
	if _, err := toml.DecodeFile(path, &sys); err != nil {
		return systemFile{}, fmt.Errorf("reading config %q: %w", path, err)
	}
	return sys, nil
}

func checkPermissions(path string, info os.FileInfo) {
	perm := info.Mode().Perm()
	if perm&0077 != 0 {
		_, _ = fmt.Fprintf(Stderr,
			"Warning: config file %s has permissions %04o; expected 0600.\n"+
				"  The file may contain an API key. Run: chmod 600 %s\n",
			path, perm, path,
		)
	}
}

func loadSwitches(path string) (switchesFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return switchesFile{}, fmt.Errorf("reading switches %q: %w", path, err)
	}
	var sw switchesFile
	if err := yaml.Unmarshal(data, &sw); err != nil {
		return switchesFile{}, fmt.Errorf("parsing switches %q: %w", path, err)
	}
	return sw, nil
}

func applyDefaults(cfg *Config) {
	if cfg.API.TimeoutSeconds <= 0 {
		cfg.API.TimeoutSeconds = 30
	}
	if cfg.API.MaxRetries <= 0 {
		cfg.API.MaxRetries = 3
	}
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("LITE_SWITCH_BASE_URL"); v != "" {
		cfg.API.BaseURL = v
	}
	if v := os.Getenv("LITE_SWITCH_API_KEY"); v != "" {
		cfg.API.APIKey = v
	}
	if v := os.Getenv("LITE_SWITCH_MODEL"); v != "" {
		cfg.Model.Name = v
	}
}

func validate(cfg *Config) error {
	var errs []string

	if strings.TrimSpace(cfg.API.BaseURL) == "" {
		errs = append(errs, "api.base_url must not be empty")
	}
	if strings.TrimSpace(cfg.API.APIKey) == "" {
		errs = append(errs, "api.api_key must not be empty (or set LITE_SWITCH_API_KEY)")
	}
	if strings.TrimSpace(cfg.Model.Name) == "" {
		errs = append(errs, "model.name must not be empty")
	}
	if len(cfg.Switches) == 0 {
		errs = append(errs, "switches file must define at least one switch")
	}

	seen := make(map[string]bool)
	for i, sw := range cfg.Switches {
		if strings.TrimSpace(sw.Tag) == "" {
			errs = append(errs, fmt.Sprintf("switches[%d]: tag must not be empty", i))
		}
		if strings.TrimSpace(sw.Description) == "" {
			errs = append(errs, fmt.Sprintf("switches[%d]: description must not be empty", i))
		}
		if seen[sw.Tag] {
			errs = append(errs, fmt.Sprintf("switches[%d]: duplicate tag %q", i, sw.Tag))
		}
		seen[sw.Tag] = true
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}

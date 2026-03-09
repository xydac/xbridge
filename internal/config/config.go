package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the server configuration.
type Config struct {
	Project       string            `yaml:"project"`
	Scheme        string            `yaml:"scheme"`
	Configuration string            `yaml:"configuration"`
	Simulator     SimulatorConfig   `yaml:"simulator"`
	Server        ServerConfig      `yaml:"server"`
	Hooks         HooksConfig       `yaml:"hooks"`
	Profiles      map[string]Profile `yaml:"profiles"`
}

// SimulatorConfig holds simulator settings.
type SimulatorConfig struct {
	Device  string `yaml:"device"`
	Runtime string `yaml:"runtime"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Port int    `yaml:"port"`
	Key  string `yaml:"key"`
}

// HooksConfig holds hook commands.
type HooksConfig struct {
	PreBuild string `yaml:"pre_build"`
}

// Profile holds a named build profile.
type Profile struct {
	Configuration string            `yaml:"configuration"`
	BuildArgs     []string          `yaml:"build_args"`
	Simulator     SimulatorConfig   `yaml:"simulator"`
	Env           map[string]string `yaml:"env"`
}

// DefaultConfig returns a config with sane defaults.
func DefaultConfig() *Config {
	return &Config{
		Configuration: "Debug",
		Server: ServerConfig{
			Port: 7900,
		},
	}
}

// Load reads config from xbridge.yaml in the given directory, if it exists.
// Returns default config if file doesn't exist.
func Load(dir string) (*Config, error) {
	cfg := DefaultConfig()

	path := filepath.Join(dir, "xbridge.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 7900
	}
	if cfg.Configuration == "" {
		cfg.Configuration = "Debug"
	}

	return cfg, nil
}

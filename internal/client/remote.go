package client

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ClientConfig is the client-side config stored in ~/.xbridge.yaml.
type ClientConfig struct {
	DefaultRemote string            `yaml:"default_remote"`
	Key           string            `yaml:"key"`
	Remotes       map[string]Remote `yaml:"remotes"`
}

// Remote represents a named remote server.
type Remote struct {
	Host string `yaml:"host"`
	Key  string `yaml:"key"`
}

// LoadClientConfig reads the client config from ~/.xbridge.yaml.
func LoadClientConfig() (*ClientConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".xbridge.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ClientConfig{Remotes: make(map[string]Remote)}, nil
		}
		return nil, err
	}

	var cfg ClientConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Remotes == nil {
		cfg.Remotes = make(map[string]Remote)
	}

	return &cfg, nil
}

// SaveClientConfig writes the client config to ~/.xbridge.yaml.
func SaveClientConfig(cfg *ClientConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	path := filepath.Join(home, ".xbridge.yaml")
	return os.WriteFile(path, data, 0600)
}

// ResolveRemote resolves a remote name to host and key.
func (c *ClientConfig) ResolveRemote(name string) (host, key string) {
	if name != "" {
		if r, ok := c.Remotes[name]; ok {
			host = r.Host
			key = r.Key
			if key == "" {
				key = c.Key
			}
			return
		}
	}

	return c.DefaultRemote, c.Key
}

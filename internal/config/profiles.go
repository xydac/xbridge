package config

import "fmt"

// ResolveProfile returns the effective configuration for a given profile name.
// If profileName is empty, returns the base config values.
func (c *Config) ResolveProfile(profileName string) (*ResolvedConfig, error) {
	resolved := &ResolvedConfig{
		Project:       c.Project,
		Scheme:        c.Scheme,
		Configuration: c.Configuration,
		Device:        c.Simulator.Device,
		Runtime:       c.Simulator.Runtime,
		PreBuild:      c.Hooks.PreBuild,
	}

	if profileName == "" {
		return resolved, nil
	}

	profile, ok := c.Profiles[profileName]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", profileName)
	}

	if profile.Configuration != "" {
		resolved.Configuration = profile.Configuration
	}
	if profile.Simulator.Device != "" {
		resolved.Device = profile.Simulator.Device
	}
	if profile.Simulator.Runtime != "" {
		resolved.Runtime = profile.Simulator.Runtime
	}
	resolved.BuildArgs = profile.BuildArgs
	resolved.Env = profile.Env

	return resolved, nil
}

// ResolvedConfig is the flattened configuration after profile resolution.
type ResolvedConfig struct {
	Project       string
	Scheme        string
	Configuration string
	Device        string
	Runtime       string
	PreBuild      string
	BuildArgs     []string
	Env           map[string]string
}

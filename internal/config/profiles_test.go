package config

import "testing"

func TestResolveProfileDefault(t *testing.T) {
	cfg := &Config{
		Project:       "./App.xcworkspace",
		Scheme:        "App",
		Configuration: "Debug",
		Simulator:     SimulatorConfig{Device: "iPhone 16"},
	}

	resolved, err := cfg.ResolveProfile("")
	if err != nil {
		t.Fatal(err)
	}

	if resolved.Project != "./App.xcworkspace" {
		t.Errorf("expected project ./App.xcworkspace, got %s", resolved.Project)
	}
	if resolved.Configuration != "Debug" {
		t.Errorf("expected Debug, got %s", resolved.Configuration)
	}
	if resolved.Device != "iPhone 16" {
		t.Errorf("expected iPhone 16, got %s", resolved.Device)
	}
}

func TestResolveProfileOverrides(t *testing.T) {
	cfg := &Config{
		Project:       "./App.xcworkspace",
		Scheme:        "App",
		Configuration: "Debug",
		Simulator:     SimulatorConfig{Device: "iPhone 16"},
		Profiles: map[string]Profile{
			"release": {
				Configuration: "Release",
				Simulator:     SimulatorConfig{Device: "iPad Pro (13-inch)"},
				BuildArgs:     []string{"EXTRA=1"},
				Env:           map[string]string{"ENV": "prod"},
			},
		},
	}

	resolved, err := cfg.ResolveProfile("release")
	if err != nil {
		t.Fatal(err)
	}

	if resolved.Configuration != "Release" {
		t.Errorf("expected Release, got %s", resolved.Configuration)
	}
	if resolved.Device != "iPad Pro (13-inch)" {
		t.Errorf("expected iPad Pro, got %s", resolved.Device)
	}
	if len(resolved.BuildArgs) != 1 || resolved.BuildArgs[0] != "EXTRA=1" {
		t.Errorf("expected BuildArgs [EXTRA=1], got %v", resolved.BuildArgs)
	}
	if resolved.Env["ENV"] != "prod" {
		t.Errorf("expected ENV=prod, got %s", resolved.Env["ENV"])
	}
}

func TestResolveProfileNotFound(t *testing.T) {
	cfg := &Config{}
	_, err := cfg.ResolveProfile("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent profile")
	}
}

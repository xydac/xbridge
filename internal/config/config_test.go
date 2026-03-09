package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Server.Port != 7900 {
		t.Errorf("expected default port 7900, got %d", cfg.Server.Port)
	}
	if cfg.Configuration != "Debug" {
		t.Errorf("expected default configuration Debug, got %s", cfg.Configuration)
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != 7900 {
		t.Errorf("expected default port, got %d", cfg.Server.Port)
	}
}

func TestLoadConfigFile(t *testing.T) {
	dir := t.TempDir()
	content := `
project: ./MyApp.xcworkspace
scheme: MyApp
configuration: Release
simulator:
  device: "iPhone 16 Pro"
  runtime: "iOS 18.2"
server:
  port: 8080
  key: "test-key"
hooks:
  pre_build: "pod install"
profiles:
  staging:
    configuration: Debug
    build_args:
      - "SWIFT_ACTIVE_COMPILATION_CONDITIONS=STAGING"
    simulator:
      device: "iPhone 16 Pro Max"
    env:
      API_URL: "https://staging.example.com"
`
	if err := os.WriteFile(filepath.Join(dir, "xbridge.yaml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Project != "./MyApp.xcworkspace" {
		t.Errorf("expected project ./MyApp.xcworkspace, got %s", cfg.Project)
	}
	if cfg.Scheme != "MyApp" {
		t.Errorf("expected scheme MyApp, got %s", cfg.Scheme)
	}
	if cfg.Configuration != "Release" {
		t.Errorf("expected configuration Release, got %s", cfg.Configuration)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.Key != "test-key" {
		t.Errorf("expected key test-key, got %s", cfg.Server.Key)
	}
	if cfg.Simulator.Device != "iPhone 16 Pro" {
		t.Errorf("expected device iPhone 16 Pro, got %s", cfg.Simulator.Device)
	}
	if cfg.Hooks.PreBuild != "pod install" {
		t.Errorf("expected pre_build pod install, got %s", cfg.Hooks.PreBuild)
	}
	if _, ok := cfg.Profiles["staging"]; !ok {
		t.Error("expected staging profile")
	}
}

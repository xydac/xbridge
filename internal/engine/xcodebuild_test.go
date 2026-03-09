package engine

import "testing"

func TestBuildArgsWorkspace(t *testing.T) {
	cfg := &XcodeBuildConfig{
		Project:       "MyApp.xcworkspace",
		Scheme:        "MyApp",
		Configuration: "Debug",
		Destination:   "platform=iOS Simulator,name=iPhone 16",
		DerivedData:   "/tmp/dd",
		ExtraArgs:     []string{"EXTRA=1"},
	}

	args := cfg.BuildArgs()

	expected := []string{
		"-workspace", "MyApp.xcworkspace",
		"-scheme", "MyApp",
		"-configuration", "Debug",
		"-destination", "platform=iOS Simulator,name=iPhone 16",
		"-derivedDataPath", "/tmp/dd",
		"EXTRA=1",
	}

	if len(args) != len(expected) {
		t.Fatalf("expected %d args, got %d: %v", len(expected), len(args), args)
	}

	for i, a := range args {
		if a != expected[i] {
			t.Errorf("arg[%d]: expected %q, got %q", i, expected[i], a)
		}
	}
}

func TestBuildArgsProject(t *testing.T) {
	cfg := &XcodeBuildConfig{
		Project: "MyApp.xcodeproj",
		Scheme:  "MyApp",
	}

	args := cfg.BuildArgs()
	if args[0] != "-project" {
		t.Errorf("expected -project, got %s", args[0])
	}
}

func TestBuildArgsEmpty(t *testing.T) {
	cfg := &XcodeBuildConfig{}
	args := cfg.BuildArgs()
	if len(args) != 0 {
		t.Errorf("expected empty args, got %v", args)
	}
}

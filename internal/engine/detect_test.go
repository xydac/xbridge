package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProjectWorkspace(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "MyApp.xcworkspace"), 0755)
	os.Mkdir(filepath.Join(dir, "MyApp.xcodeproj"), 0755)

	result, err := DetectProject(dir)
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(dir, "MyApp.xcworkspace")
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestDetectProjectXcodeproj(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "MyApp.xcodeproj"), 0755)

	result, err := DetectProject(dir)
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(dir, "MyApp.xcodeproj")
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestDetectProjectNone(t *testing.T) {
	dir := t.TempDir()

	result, err := DetectProject(dir)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Errorf("expected empty, got %s", result)
	}
}

func TestIsWorkspace(t *testing.T) {
	if !IsWorkspace("MyApp.xcworkspace") {
		t.Error("expected true for .xcworkspace")
	}
	if IsWorkspace("MyApp.xcodeproj") {
		t.Error("expected false for .xcodeproj")
	}
}

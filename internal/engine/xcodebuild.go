package engine

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// XcodeBuildConfig holds parameters for an xcodebuild invocation.
type XcodeBuildConfig struct {
	Project       string
	Scheme        string
	Configuration string
	Destination   string
	DerivedData   string
	ExtraArgs     []string
	Env           map[string]string
}

// BuildArgs constructs xcodebuild CLI arguments.
func (c *XcodeBuildConfig) BuildArgs() []string {
	var args []string

	if c.Project != "" {
		if IsWorkspace(c.Project) {
			args = append(args, "-workspace", c.Project)
		} else {
			args = append(args, "-project", c.Project)
		}
	}

	if c.Scheme != "" {
		args = append(args, "-scheme", c.Scheme)
	}
	if c.Configuration != "" {
		args = append(args, "-configuration", c.Configuration)
	}
	if c.Destination != "" {
		args = append(args, "-destination", c.Destination)
	}
	if c.DerivedData != "" {
		args = append(args, "-derivedDataPath", c.DerivedData)
	}

	args = append(args, c.ExtraArgs...)

	return args
}

// XcodeBuildCommand creates an exec.Cmd for xcodebuild with the given action.
func XcodeBuildCommand(ctx context.Context, cfg *XcodeBuildConfig, action string) *exec.Cmd {
	args := cfg.BuildArgs()
	args = append(args, action)

	cmd := exec.CommandContext(ctx, "xcodebuild", args...)
	for k, v := range cfg.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd
}

// XcodeVersion returns the current Xcode version string.
func XcodeVersion() (string, error) {
	out, err := exec.Command("xcodebuild", "-version").Output()
	if err != nil {
		return "", fmt.Errorf("xcodebuild -version: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// CleanBuild runs xcodebuild clean.
func CleanBuild(ctx context.Context, cfg *XcodeBuildConfig) error {
	cmd := XcodeBuildCommand(ctx, cfg, "clean")
	return cmd.Run()
}

// ListSchemes returns the list of schemes for a project.
func ListSchemes(project string) ([]string, error) {
	var args []string
	if IsWorkspace(project) {
		args = []string{"-workspace", project, "-list"}
	} else {
		args = []string{"-project", project, "-list"}
	}

	out, err := exec.Command("xcodebuild", args...).Output()
	if err != nil {
		return nil, err
	}

	var schemes []string
	inSchemes := false
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "Schemes:" {
			inSchemes = true
			continue
		}
		if inSchemes {
			if line == "" {
				break
			}
			schemes = append(schemes, line)
		}
	}
	return schemes, nil
}

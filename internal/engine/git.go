package engine

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitStatus holds the current git state.
type GitStatus struct {
	Branch string `json:"branch"`
	Commit string `json:"commit"`
	Dirty  bool   `json:"dirty"`
}

// GitGetStatus returns the current git status for the given directory.
func GitGetStatus(dir string) (*GitStatus, error) {
	branch, err := gitCmd(dir, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("git branch: %w", err)
	}

	commit, err := gitCmd(dir, "rev-parse", "--short", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("git commit: %w", err)
	}

	status, err := gitCmd(dir, "status", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("git status: %w", err)
	}

	return &GitStatus{
		Branch: branch,
		Commit: commit,
		Dirty:  status != "",
	}, nil
}

// GitPull runs git pull in the given directory.
func GitPull(dir string) (string, error) {
	return gitCmd(dir, "pull")
}

// GitCheckout checks out a branch in the given directory.
func GitCheckout(dir, branch string) (string, error) {
	return gitCmd(dir, "checkout", branch)
}

func gitCmd(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s: %s", strings.Join(args, " "), string(exitErr.Stderr))
		}
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// IsGitRepo returns true if the directory is a git repository.
func IsGitRepo(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	return cmd.Run() == nil
}

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

func TestHandleGitStatus(t *testing.T) {
	dir := t.TempDir()
	// Init a git repo in the temp dir
	cmds := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@test.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "commit", "--allow-empty", "-m", "init"},
	}
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("cmd %v: %s: %s", args, err, out)
		}
	}

	gh := &GitHandlers{WorkDir: dir}
	req := httptest.NewRequest("GET", "/git/status", nil)
	rec := httptest.NewRecorder()
	gh.HandleGitStatus().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var status map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatal(err)
	}

	if status["branch"] == "" {
		t.Error("expected branch to be set")
	}
	if status["commit"] == "" {
		t.Error("expected commit to be set")
	}
}

func TestHandleGitCheckout_MissingBranch(t *testing.T) {
	gh := &GitHandlers{WorkDir: t.TempDir()}

	req := httptest.NewRequest("POST", "/git/checkout", strings.NewReader(`{}`))
	rec := httptest.NewRecorder()
	gh.HandleGitCheckout().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

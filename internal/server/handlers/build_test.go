package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/xydac/xbridge/internal/build"
	"github.com/xydac/xbridge/internal/config"
)

func TestHandleBuildStatus_NotFound(t *testing.T) {
	mgr := build.NewManager(t.TempDir(), config.DefaultConfig())
	bh := &BuildHandlers{Manager: mgr}

	r := chi.NewRouter()
	r.Get("/build/{id}", bh.HandleBuildStatus())

	req := httptest.NewRequest("GET", "/build/nonexistent", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestHandleBuildLogs_NotFound(t *testing.T) {
	mgr := build.NewManager(t.TempDir(), config.DefaultConfig())
	bh := &BuildHandlers{Manager: mgr}

	r := chi.NewRouter()
	r.Get("/build/{id}/logs", bh.HandleBuildLogs())

	req := httptest.NewRequest("GET", "/build/nonexistent/logs", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestHandleBuildArtifact_NotFound(t *testing.T) {
	mgr := build.NewManager(t.TempDir(), config.DefaultConfig())
	bh := &BuildHandlers{Manager: mgr}

	r := chi.NewRouter()
	r.Get("/build/{id}/artifact", bh.HandleBuildArtifact())

	req := httptest.NewRequest("GET", "/build/nonexistent/artifact", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestHandleCleanBuild(t *testing.T) {
	// This test will fail without xcodebuild, but we verify the handler doesn't panic
	mgr := build.NewManager(t.TempDir(), config.DefaultConfig())
	bh := &BuildHandlers{Manager: mgr}

	req := httptest.NewRequest("POST", "/build/clean", nil)
	rec := httptest.NewRecorder()
	bh.HandleCleanBuild().ServeHTTP(rec, req)

	// Expect 500 since xcodebuild won't work in test
	// but the handler should not panic
	if rec.Code != http.StatusOK && rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 200 or 500, got %d", rec.Code)
	}
}

func TestHandleBuild_InvalidProfile(t *testing.T) {
	cfg := config.DefaultConfig()
	mgr := build.NewManager(t.TempDir(), cfg)
	bh := &BuildHandlers{Manager: mgr}

	body := `{"profile": "nonexistent"}`
	req := httptest.NewRequest("POST", "/build", strings.NewReader(body))
	rec := httptest.NewRecorder()
	bh.HandleBuild().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandleBuild_Success(t *testing.T) {
	cfg := config.DefaultConfig()
	mgr := build.NewManager(t.TempDir(), cfg)
	bh := &BuildHandlers{Manager: mgr}

	req := httptest.NewRequest("POST", "/build", nil)
	rec := httptest.NewRecorder()
	bh.HandleBuild().ServeHTTP(rec, req)

	// Build will start (and likely fail without xcodebuild), but should return 202
	if rec.Code != http.StatusAccepted {
		t.Errorf("expected 202, got %d", rec.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result["id"] == "" {
		t.Error("expected build ID in response")
	}
	if result["status"] != "running" {
		t.Errorf("expected status running, got %s", result["status"])
	}
}


package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xydac/xbridge/internal/build"
)

// BuildHandlers holds dependencies for build endpoints.
type BuildHandlers struct {
	Manager *build.Manager
}

// BuildRequest is the request body for POST /build.
type BuildRequest struct {
	Profile string `json:"profile"`
}

// HandleBuild triggers a new build.
func (h *BuildHandlers) HandleBuild() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BuildRequest
		if r.Body != nil {
			json.NewDecoder(r.Body).Decode(&req)
		}

		job, err := h.Manager.StartBuild(r.Context(), req.Profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"id":     job.ID,
			"status": string(job.Status),
		})
	}
}

// HandleCleanBuild runs xcodebuild clean.
func (h *BuildHandlers) HandleCleanBuild() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.Manager.CleanBuild(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "cleaned"})
	}
}

// HandleBuildStatus returns the status of a build job.
func (h *BuildHandlers) HandleBuildStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		job, ok := h.Manager.GetJob(id)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	}
}

// HandleBuildLogs streams build logs via SSE.
func (h *BuildHandlers) HandleBuildLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		job, ok := h.Manager.GetJob(id)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		WriteSSE(w, r, job.Logs)
	}
}

// HandleBuildArtifact downloads the build artifact.
func (h *BuildHandlers) HandleBuildArtifact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement artifact download (tar.gz of .app bundle from DerivedData)
		id := chi.URLParam(r, "id")
		_, ok := h.Manager.GetJob(id)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		http.Error(w, "artifact download not yet implemented", http.StatusNotImplemented)
	}
}

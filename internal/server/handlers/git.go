package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/xydac/xbridge/internal/engine"
)

// GitHandlers holds dependencies for git endpoints.
type GitHandlers struct {
	WorkDir string
}

// CheckoutRequest is the body for POST /git/checkout.
type CheckoutRequest struct {
	Branch string `json:"branch"`
}

// HandleGitStatus returns the current git status.
func (h *GitHandlers) HandleGitStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := engine.GitGetStatus(h.WorkDir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

// HandleGitPull triggers a git pull.
func (h *GitHandlers) HandleGitPull() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := engine.GitPull(h.WorkDir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"output": output,
		})
	}
}

// HandleGitCheckout switches to a branch.
func (h *GitHandlers) HandleGitCheckout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CheckoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.Branch == "" {
			http.Error(w, "branch is required", http.StatusBadRequest)
			return
		}

		output, err := engine.GitCheckout(h.WorkDir, req.Branch)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"output": output,
		})
	}
}

package handlers

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"

	"github.com/xydac/xbridge/internal/engine"
)

// HealthResponse is returned by GET /health.
type HealthResponse struct {
	Status       string `json:"status"`
	XcodeVersion string `json:"xcode_version,omitempty"`
	DiskFree     string `json:"disk_free,omitempty"`
}

// HandleHealth returns server health info.
func HandleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{Status: "ok"}

		if ver, err := engine.XcodeVersion(); err == nil {
			resp.XcodeVersion = ver
		}

		if out, err := exec.Command("df", "-h", "/").Output(); err == nil {
			lines := strings.Split(string(out), "\n")
			if len(lines) >= 2 {
				fields := strings.Fields(lines[1])
				if len(fields) >= 4 {
					resp.DiskFree = fields[3]
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

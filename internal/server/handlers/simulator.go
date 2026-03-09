package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/xydac/xbridge/internal/engine"
)

// HandleListSimulators returns available simulators.
func HandleListSimulators() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := engine.ListDevices()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(devices)
	}
}

// BootRequest is the body for POST /simulators/boot.
type BootRequest struct {
	UDID   string `json:"udid"`
	Device string `json:"device"`
}

// HandleBootSimulator boots a simulator.
func HandleBootSimulator() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BootRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		udid := req.UDID
		if udid == "" && req.Device != "" {
			var err error
			udid, err = engine.FindDeviceByName(req.Device)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}

		if udid == "" {
			http.Error(w, "udid or device name required", http.StatusBadRequest)
			return
		}

		if err := engine.BootDevice(udid); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "booted", "udid": udid})
	}
}

// ShutdownRequest is the body for POST /simulators/shutdown.
type ShutdownRequest struct {
	UDID string `json:"udid"`
}

// HandleShutdownSimulator shuts down a simulator.
func HandleShutdownSimulator() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ShutdownRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		udid := req.UDID
		if udid == "" {
			udid = "BOOTED"
		}

		resolved, err := engine.ResolveUDID(udid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := engine.ShutdownDevice(resolved); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "shutdown"})
	}
}

// HandleScreenshot takes a screenshot of a simulator.
func HandleScreenshot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		udid := chi.URLParam(r, "udid")

		data, err := engine.TakeScreenshot(udid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Write(data)
	}
}

// InstallRequest is the body for POST /simulators/:udid/install.
type InstallRequest struct {
	AppPath string `json:"app_path"`
}

// HandleInstallApp installs an app on a simulator.
func HandleInstallApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		udid := chi.URLParam(r, "udid")

		var req InstallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := engine.InstallApp(udid, req.AppPath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "installed"})
	}
}

// LaunchRequest is the body for POST /simulators/:udid/launch.
type LaunchRequest struct {
	BundleID string `json:"bundle_id"`
}

// HandleLaunchApp launches an app on a simulator.
func HandleLaunchApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		udid := chi.URLParam(r, "udid")

		var req LaunchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := engine.LaunchApp(udid, req.BundleID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "launched"})
	}
}

// OpenURLRequest is the body for POST /simulators/:udid/openurl.
type OpenURLRequest struct {
	URL string `json:"url"`
}

// HandleOpenURL opens a URL on a simulator.
func HandleOpenURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		udid := chi.URLParam(r, "udid")

		var req OpenURLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := engine.OpenURL(udid, req.URL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "opened"})
	}
}

// TapRequest is the body for POST /simulators/:udid/tap.
type TapRequest struct {
	X        int     `json:"x"`
	Y        int     `json:"y"`
	Duration float64 `json:"duration,omitempty"`
}

// HandleTap sends a tap event on a simulator.
func HandleTap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		udid := chi.URLParam(r, "udid")

		var req TapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := engine.Tap(udid, req.X, req.Y, req.Duration); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "tapped"})
	}
}

// SwipeRequest is the body for POST /simulators/:udid/swipe.
type SwipeRequest struct {
	X1       int     `json:"x1"`
	Y1       int     `json:"y1"`
	X2       int     `json:"x2"`
	Y2       int     `json:"y2"`
	Duration float64 `json:"duration,omitempty"`
}

// HandleSwipe sends a swipe event on a simulator.
func HandleSwipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		udid := chi.URLParam(r, "udid")

		var req SwipeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := engine.Swipe(udid, req.X1, req.Y1, req.X2, req.Y2, req.Duration); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "swiped"})
	}
}

// TextRequest is the body for POST /simulators/:udid/text.
type TextRequest struct {
	Text string `json:"text"`
}

// HandleInputText types text into the simulator.
func HandleInputText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		udid := chi.URLParam(r, "udid")

		var req TextRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := engine.InputText(udid, req.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "typed"})
	}
}

// KeyRequest is the body for POST /simulators/:udid/key.
type KeyRequest struct {
	Key      string  `json:"key"`
	Duration float64 `json:"duration,omitempty"`
}

// HandleKeyPress sends a key press event on a simulator.
func HandleKeyPress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		udid := chi.URLParam(r, "udid")

		var req KeyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := engine.KeyPress(udid, req.Key, req.Duration); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "pressed"})
	}
}

// HandleSimulatorLogs streams simulator logs via SSE.
func HandleSimulatorLogs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement simulator log streaming
		http.Error(w, "simulator log streaming not yet implemented", http.StatusNotImplemented)
	}
}

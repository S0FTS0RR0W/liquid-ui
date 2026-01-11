package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"strconv"

	"github.com/s0fts0rr0w/liquid-ui/backend/internal/liquidctl"
	"github.com/s0fts0rr0w/liquid-ui/backend/internal/profiles"
)

// writeJSON is a helper for consistent JSON responses.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// registerRoutes attaches all API endpoints to the mux.
func registerRoutes(mux *http.ServeMux, frontendFS fs.FS) {

	// ------------------------------------------------------------
	// Health check
	// ------------------------------------------------------------
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})

	// ------------------------------------------------------------
	// List devices
	// ------------------------------------------------------------
	mux.HandleFunc("/devices", func(w http.ResponseWriter, r *http.Request) {
		devs, err := liquidctl.ListDevices()
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, devs)
	})

	// ------------------------------------------------------------
	// Device status
	// ------------------------------------------------------------
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		idxStr := r.URL.Query().Get("device")
		if idxStr == "" {
			writeJSON(w, 400, map[string]string{"error": "missing device index"})
			return
		}

		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid device index"})
			return
		}

		status, err := liquidctl.GetStatus(idx)
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, 200, status)
	})

	// ------------------------------------------------------------
	// List profiles
	// ------------------------------------------------------------
	mux.HandleFunc("/profiles", func(w http.ResponseWriter, r *http.Request) {
		profs, err := profiles.LoadProfiles()
		if err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 200, profs)
	})

	// ------------------------------------------------------------
	// Save or update a profile
	// ------------------------------------------------------------
	mux.HandleFunc("/profiles/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, 405, map[string]string{"error": "POST required"})
			return
		}

		var p profiles.Profile
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}

		profs, _ := profiles.LoadProfiles()
		profs = profiles.AddOrUpdateProfile(profs, p)

		if err := profiles.SaveProfiles(profs); err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, 200, map[string]string{"status": "saved"})
	})

	// ------------------------------------------------------------
	// Delete a profile
	// ------------------------------------------------------------
	mux.HandleFunc("/profiles/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, 405, map[string]string{"error": "POST required"})
			return
		}

		var body struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}

		profs, _ := profiles.LoadProfiles()
		profs = profiles.DeleteProfile(profs, body.Name)

		if err := profiles.SaveProfiles(profs); err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, 200, map[string]string{"status": "deleted"})
	})

	// ------------------------------------------------------------
	// Apply a profile to a device
	// ------------------------------------------------------------
	mux.HandleFunc("/profiles/apply", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, 405, map[string]string{"error": "POST required"})
			return
		}

		var body struct {
			DeviceIndex int    `json:"deviceIndex"`
			ProfileName string `json:"profileName"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}

		profs, _ := profiles.LoadProfiles()
		p, err := profiles.FindProfile(profs, body.ProfileName)
		if err != nil {
			writeJSON(w, 404, map[string]string{"error": "profile not found"})
			return
		}

		if err := liquidctl.ApplyFanCurve(body.DeviceIndex, p.Points); err != nil {
			writeJSON(w, 500, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, 200, map[string]string{"status": "applied"})
	})

	// ------------------------------------------------------------
	// Frontend Static Files
	// ------------------------------------------------------------
	fileServer := http.FileServer(http.FS(frontendFS))
	mux.Handle("/", fileServer)
}

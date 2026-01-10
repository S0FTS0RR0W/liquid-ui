package profiles

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// CurvePoint represents a single point on a fan curve.
type CurvePoint struct {
	Temperature float64 `json:"temperature"`
	Speed       int     `json:"speed"` // percent (0–100)
}

// Profile represents a named fan curve.
type Profile struct {
	Name   string       `json:"name"`
	Points []CurvePoint `json:"points"`
}

// ------------------------------------------------------------
// Storage paths
// ------------------------------------------------------------

// ConfigDir returns ~/.config/liquid-ui
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "liquid-ui"), nil
}

// ProfilesFile returns ~/.config/liquid-ui/profiles.json
func ProfilesFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "profiles.json"), nil
}

// ------------------------------------------------------------
// Load / Save
// ------------------------------------------------------------

// LoadProfiles reads all profiles from disk.
// If the file doesn't exist, it returns an empty slice.
func LoadProfiles() ([]Profile, error) {
	path, err := ProfilesFile()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Profile{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var profiles []Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

// SaveProfiles writes all profiles to disk.
func SaveProfiles(profiles []Profile) error {
	path, err := ProfilesFile()
	if err != nil {
		return err
	}

	dir, _ := ConfigDir()
	os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ------------------------------------------------------------
// Profile helpers
// ------------------------------------------------------------

// FindProfile returns a profile by name.
func FindProfile(profiles []Profile, name string) (*Profile, error) {
	for _, p := range profiles {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, errors.New("profile not found")
}

// AddOrUpdateProfile inserts or replaces a profile.
func AddOrUpdateProfile(profiles []Profile, p Profile) []Profile {
	for i, existing := range profiles {
		if existing.Name == p.Name {
			profiles[i] = p
			return profiles
		}
	}
	return append(profiles, p)
}

// DeleteProfile removes a profile by name.
func DeleteProfile(profiles []Profile, name string) []Profile {
	out := []Profile{}
	for _, p := range profiles {
		if p.Name != name {
			out = append(out, p)
		}
	}
	return out
}

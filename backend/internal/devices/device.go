package devices

// Device represents a physical liquidctl-compatible device.
// This struct is intentionally minimal for now, but designed to expand
// as you add features like multiple fan channels, RGB, pump control, etc.
type Device struct {
	Index int    `json:"index"` // liquidctl device index
	Name  string `json:"name"`  // human-readable name
	Type  string `json:"type"`  // e.g. "NZXT Smart Device V2"
}

// Status represents live telemetry from a device.
// These fields will be populated by your liquidctl wrapper.
type Status struct {
	Temperature float64 `json:"temperature"`
	FanRPM      int     `json:"fanRpm"`
	PumpRPM     int     `json:"pumpRpm"`
}

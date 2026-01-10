package liquidctl

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/s0fts0rr0w/liquid-ui/backend/internal/devices"
)

var (
	deviceRe = regexp.MustCompile(`Device #(\d+): (.+)`)
	floatRe  = regexp.MustCompile(`([-+]?[\d.]+)`)
	intRe    = regexp.MustCompile(`(\d+)`)
)

// run executes a liquidctl command and returns stdout as a string.
func run(args ...string) (string, error) {
	cmd := exec.Command("liquidctl", args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("liquidctl error: %v (%s)", err, stderr.String())
	}

	return out.String(), nil
}

// ------------------------------------------------------------
// Device Discovery
// ------------------------------------------------------------

// ListDevices parses `liquidctl list` output into []Device.
func ListDevices() ([]devices.Device, error) {
	out, err := run("list")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var result []devices.Device

	// Example line: "Device #0: NZXT Smart Device V2"
	for _, line := range lines {
		matches := deviceRe.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}

		index, _ := strconv.Atoi(matches[1])
		name := matches[2]

		result = append(result, devices.Device{
			Index: index,
			Name:  name,
			Type:  name, // refine later if needed
		})
	}

	return result, nil
}

// ------------------------------------------------------------
// Status Parsing
// ------------------------------------------------------------

// GetStatus returns temperature, fan RPM, pump RPM for a device.
func GetStatus(index int) (devices.Status, error) {
	out, err := run("status", "--device", fmt.Sprint(index))
	if err != nil {
		return devices.Status{}, err
	}

	// Example output:
	// "Temperature: 32.0 °C\nFan speed: 1200 rpm\nPump speed: 2000 rpm"
	status := devices.Status{}

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Temperature:") {
			status.Temperature = parseFloat(line)
		}
		if strings.HasPrefix(line, "Fan speed:") {
			status.FanRPM = parseInt(line)
		}
		if strings.HasPrefix(line, "Pump speed:") {
			status.PumpRPM = parseInt(line)
		}
	}

	return status, nil
}

// ------------------------------------------------------------
// Fan Curve Application
// ------------------------------------------------------------

// ApplyFanCurve sends a temperature/speed curve to the device.
func ApplyFanCurve(index int, points []devices.CurvePoint) error {
	if len(points) == 0 {
		return errors.New("fan curve is empty")
	}

	// Build args: temp1 speed1 temp2 speed2 ...
	var args []string
	args = append(args, "set", "fan", "speed", "--device", fmt.Sprint(index))

	for _, p := range points {
		args = append(args, fmt.Sprint(int(p.Temperature)))
		args = append(args, fmt.Sprint(p.Speed))
	}

	_, err := run(args...)
	return err
}

// ------------------------------------------------------------
// Helpers
// ------------------------------------------------------------

func parseFloat(line string) float64 {
	match := floatRe.FindString(line)
	f, _ := strconv.ParseFloat(match, 64)
	return f
}

func parseInt(line string) int {
	match := intRe.FindString(line)
	i, _ := strconv.Atoi(match)
	return i
}

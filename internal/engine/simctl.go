package engine

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Device represents a simulator device.
type Device struct {
	UDID       string `json:"udid"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Runtime    string `json:"runtime"`
	IsAvail    bool   `json:"isAvailable"`
	DeviceType string `json:"deviceTypeIdentifier,omitempty"`
}

// simctlOutput is the raw output from `xcrun simctl list devices -j`.
type simctlOutput struct {
	Devices map[string][]simctlDevice `json:"devices"`
}

type simctlDevice struct {
	UDID                   string `json:"udid"`
	Name                   string `json:"name"`
	State                  string `json:"state"`
	IsAvailable            bool   `json:"isAvailable"`
	DeviceTypeIdentifier   string `json:"deviceTypeIdentifier"`
}

// ListDevices returns all simulator devices.
func ListDevices() ([]Device, error) {
	out, err := exec.Command("xcrun", "simctl", "list", "devices", "-j").Output()
	if err != nil {
		return nil, fmt.Errorf("simctl list devices: %w", err)
	}

	var raw simctlOutput
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("parse simctl output: %w", err)
	}

	var devices []Device
	for runtime, devs := range raw.Devices {
		// Extract short runtime name like "iOS 18.2"
		rt := parseRuntime(runtime)
		for _, d := range devs {
			devices = append(devices, Device{
				UDID:       d.UDID,
				Name:       d.Name,
				State:      d.State,
				Runtime:    rt,
				IsAvail:    d.IsAvailable,
				DeviceType: d.DeviceTypeIdentifier,
			})
		}
	}

	return devices, nil
}

// parseRuntime converts "com.apple.CoreSimulator.SimRuntime.iOS-18-2" to "iOS 18.2".
func parseRuntime(raw string) string {
	// Remove prefix
	raw = strings.TrimPrefix(raw, "com.apple.CoreSimulator.SimRuntime.")
	// Replace hyphens: iOS-18-2 -> iOS 18.2
	parts := strings.SplitN(raw, "-", 2)
	if len(parts) == 2 {
		version := strings.ReplaceAll(parts[1], "-", ".")
		return parts[0] + " " + version
	}
	return raw
}

// ResolveUDID resolves "BOOTED" to the actual UDID of the booted simulator.
// If udid is not "BOOTED", it returns it unchanged.
func ResolveUDID(udid string) (string, error) {
	if !strings.EqualFold(udid, "BOOTED") {
		return udid, nil
	}

	devices, err := ListDevices()
	if err != nil {
		return "", err
	}

	for _, d := range devices {
		if d.State == "Booted" {
			return d.UDID, nil
		}
	}
	return "", fmt.Errorf("no booted simulator found")
}

// BootDevice boots a simulator by UDID.
func BootDevice(udid string) error {
	return exec.Command("xcrun", "simctl", "boot", udid).Run()
}

// ShutdownDevice shuts down a simulator by UDID.
func ShutdownDevice(udid string) error {
	return exec.Command("xcrun", "simctl", "shutdown", udid).Run()
}

// TakeScreenshot captures a screenshot and returns the PNG data.
func TakeScreenshot(udid string) ([]byte, error) {
	resolved, err := ResolveUDID(udid)
	if err != nil {
		return nil, err
	}

	out, err := exec.Command("xcrun", "simctl", "io", resolved, "screenshot", "--type=png", "-").Output()
	if err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}
	return out, nil
}

// InstallApp installs an .app bundle on a simulator.
func InstallApp(udid, appPath string) error {
	resolved, err := ResolveUDID(udid)
	if err != nil {
		return err
	}
	return exec.Command("xcrun", "simctl", "install", resolved, appPath).Run()
}

// LaunchApp launches an app by bundle ID on a simulator.
func LaunchApp(udid, bundleID string) error {
	resolved, err := ResolveUDID(udid)
	if err != nil {
		return err
	}
	return exec.Command("xcrun", "simctl", "launch", resolved, bundleID).Run()
}

// OpenURL opens a URL/deep link on a simulator.
func OpenURL(udid, url string) error {
	resolved, err := ResolveUDID(udid)
	if err != nil {
		return err
	}
	return exec.Command("xcrun", "simctl", "openurl", resolved, url).Run()
}

// FindDeviceByName finds a device UDID by name, preferring available devices.
func FindDeviceByName(name string) (string, error) {
	devices, err := ListDevices()
	if err != nil {
		return "", err
	}

	for _, d := range devices {
		if strings.EqualFold(d.Name, name) && d.IsAvail {
			return d.UDID, nil
		}
	}

	return "", fmt.Errorf("simulator %q not found", name)
}

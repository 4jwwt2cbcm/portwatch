package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// Interval is how often to scan ports.
	Interval time.Duration `json:"interval_seconds"`
	// Ports is the list of ports to monitor. If empty, all open ports are monitored.
	Ports []uint16 `json:"ports"`
	// AlertOutput is the file path for alert output. Defaults to stdout if empty.
	AlertOutput string `json:"alert_output"`
	// Protocols is the list of protocols to scan ("tcp", "udp").
	Protocols []string `json:"protocols"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Interval:  30 * time.Second,
		Ports:     []uint16{},
		AlertOutput: "",
		Protocols: []string{"tcp"},
	}
}

// Load reads a JSON config file from the given path and returns a Config.
// Fields not present in the file retain their default values.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Use a raw struct to handle the seconds-as-int representation.
	var raw struct {
		IntervalSeconds int      `json:"interval_seconds"`
		Ports           []uint16 `json:"ports"`
		AlertOutput     string   `json:"alert_output"`
		Protocols       []string `json:"protocols"`
	}

	if err := json.NewDecoder(f).Decode(&raw); err != nil {
		return nil, err
	}

	if raw.IntervalSeconds > 0 {
		cfg.Interval = time.Duration(raw.IntervalSeconds) * time.Second
	}
	if len(raw.Ports) > 0 {
		cfg.Ports = raw.Ports
	}
	if raw.AlertOutput != "" {
		cfg.AlertOutput = raw.AlertOutput
	}
	if len(raw.Protocols) > 0 {
		cfg.Protocols = raw.Protocols
	}

	return cfg, nil
}

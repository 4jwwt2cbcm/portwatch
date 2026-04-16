package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	Interval  time.Duration `yaml:"interval"`
	StateDir  string        `yaml:"state_dir"`
	LogFormat string        `yaml:"log_format"`
	Filter    FilterConfig  `yaml:"filter"`
}

// FilterConfig restricts which ports are monitored.
type FilterConfig struct {
	Protocols []string `yaml:"protocols"`
	PortMin   int      `yaml:"port_min"`
	PortMax   int      `yaml:"port_max"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:  15 * time.Second,
		StateDir:  "/var/lib/portwatch",
		LogFormat: "text",
		Filter: FilterConfig{
			Protocols: []string{"tcp", "udp"},
			PortMin:   1,
			PortMax:   65535,
		},
	}
}

// Validate checks that the configuration values are logically consistent.
func (c Config) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("interval must be positive, got %s", c.Interval)
	}
	if c.Filter.PortMin < 1 || c.Filter.PortMin > 65535 {
		return fmt.Errorf("port_min must be between 1 and 65535, got %d", c.Filter.PortMin)
	}
	if c.Filter.PortMax < 1 || c.Filter.PortMax > 65535 {
		return fmt.Errorf("port_max must be between 1 and 65535, got %d", c.Filter.PortMax)
	}
	if c.Filter.PortMin > c.Filter.PortMax {
		return fmt.Errorf("port_min (%d) must not exceed port_max (%d)", c.Filter.PortMin, c.Filter.PortMax)
	}
	if c.LogFormat != "text" && c.LogFormat != "json" {
		return fmt.Errorf("log_format must be \"text\" or \"json\", got %q", c.LogFormat)
	}
	return nil
}

// Load reads a YAML config file from path and merges it over the defaults.
// If the file does not exist the defaults are returned without error.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	var partial Config
	if err := yaml.Unmarshal(data, &partial); err != nil {
		return cfg, err
	}

	if partial.Interval != 0 {
		cfg.Interval = partial.Interval
	}
	if partial.StateDir != "" {
		cfg.StateDir = partial.StateDir
	}
	if partial.LogFormat != "" {
		cfg.LogFormat = partial.LogFormat
	}
	if len(partial.Filter.Protocols) > 0 {
		cfg.Filter.Protocols = partial.Filter.Protocols
	}
	if partial.Filter.PortMin != 0 {
		cfg.Filter.PortMin = partial.Filter.PortMin
	}
	if partial.Filter.PortMax != 0 {
		cfg.Filter.PortMax = partial.Filter.PortMax
	}

	return cfg, nil
}

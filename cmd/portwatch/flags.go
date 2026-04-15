package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.1.0"

// CLIFlags holds parsed command-line options.
type CLIFlags struct {
	ConfigPath string
	Version    bool
}

// ParseFlags parses os.Args and returns CLIFlags.
// It calls os.Exit(2) on parse errors (standard flag behaviour).
func ParseFlags(args []string) CLIFlags {
	fs := flag.NewFlagSet("portwatch", flag.ExitOnError)

	var f CLIFlags
	fs.StringVar(&f.ConfigPath, "config", "portwatch.yaml", "path to configuration file")
	fs.BoolVar(&f.Version, "version", false, "print version and exit")

	_ = fs.Parse(args)
	return f
}

// MaybeVersion prints the version string and exits if f.Version is true.
func MaybeVersion(f CLIFlags) {
	if f.Version {
		fmt.Fprintf(os.Stdout, "portwatch %s\n", version)
		os.Exit(0)
	}
}

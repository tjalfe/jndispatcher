// Package arguments provides the command line arguments for the application.
package arguments

import (
	"github.com/spf13/pflag"
)

type Arguments struct {
	Debug   bool
	Encrypt bool
}

// ParseArguments parses the command line arguments and returns an Arguments struct.
func ParseArguments() (Arguments, error) {
	var args Arguments

	// Define command line flags
	encrypt := pflag.BoolP("encrypt", "e", false, "Encrypt a certificate store passphrase")
	debug := pflag.Bool("debug", false, "Run in debug mode")

	pflag.Parse()
	if *encrypt {
		args.Encrypt = true
		return args, nil
	}
	args.Debug = *debug
	return args, nil
}

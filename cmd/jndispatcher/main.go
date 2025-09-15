package main

import (
	"fmt"
	"os"

	"github.com/tjalfe/jndispatcher/internal/arguments"
	"github.com/tjalfe/jndispatcher/internal/config"
	"github.com/tjalfe/jndispatcher/internal/types"
	"github.com/tjalfe/pcrypt"
)

// PromptForPassword prompts the user for a password and encrypts it using pcrypt.
// The encrypted password is used in config yaml file so that the real password is not stored in plaintext.
func PromptForPassword() error {
	crypt, err := pcrypt.Init(types.PcryptInit)
	if err != nil {
		return err
	}
	err = pcrypt.Prompt(crypt)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	args, err := arguments.ParseArguments()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	if args.Encrypt {
		err := PromptForPassword()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encrypting password: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Normal execution of the application starts here
	_, err = config.ReadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

}

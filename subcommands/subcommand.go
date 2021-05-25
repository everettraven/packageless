package subcommands

import (
	"errors"
	"fmt"
	"os"
)

//Runner - Interface to enable easy interactions with the different subcommand objects
type Runner interface {
	Init(args []string) error
	Run() error
	Name() string
}

//SubCommand - Helper function that handles setting up and running subcommands
func SubCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("A subcommand must be passed")
	}

	cmds := []Runner{
		NewInstallCommand(),
		NewUpgradeCommand(),
		NewRunCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			err := cmd.Init(os.Args[2:])

			if err != nil {
				return err
			}

			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand %s", subcommand)
}

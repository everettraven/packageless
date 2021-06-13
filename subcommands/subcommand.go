package subcommands

import (
	"errors"
	"fmt"
)

//Runner - Interface to enable easy interactions with the different subcommand objects
type Runner interface {
	Init(args []string) error
	Run() error
	Name() string
}

//SubCommand - Helper function that handles setting up and running subcommands
func SubCommand(args []string, scmds []Runner) error {
	if len(args) < 1 {
		return errors.New("A subcommand must be passed")
	}

	subcommand := args[0]

	for _, cmd := range scmds {
		if cmd.Name() == subcommand {
			err := cmd.Init(args[1:])

			if err != nil {
				return err
			}

			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand %s", subcommand)
}

package subcommands

import (
	"errors"
	"fmt"
	"os"

	"github.com/everettraven/packageless/utils"
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

	tool := utils.NewUtility()

	cp := &utils.CopyTool{}

	cmds := []Runner{
		NewInstallCommand(tool, cp),
		NewUpgradeCommand(tool, cp),
		NewRunCommand(tool),
		NewUninstallCommand(tool),
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

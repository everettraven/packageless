package subcommands

import (
	"errors"
	"fmt"

	"github.com/everettraven/packageless/utils"
)

//Runner - Interface to enable easy interactions with the different subcommand objects
type Runner interface {
	Init(args []string) error
	Run() error
	Name() string
}

//SubCommand - Helper function that handles setting up and running subcommands
func SubCommand(args []string, scmds []Runner) error {
	utils.NewUtility().RenderInfoMarkdown("# packageless")
	if len(args) < 1 {
		return errors.New("A subcommand must be passed")
	}

	subcommand := args[0]

	args = args[1:]

	for _, cmd := range scmds {
		if cmd.Name() == subcommand {
			//Subcommands that take multiple pims as arguments
			if subcommand == "install" || subcommand == "upgrade" || subcommand == "uninstall" {
				//Execute subcommand for each pim in arguments
				for _, pim := range args {
					p := []string{pim}
					err := cmd.Init(p)

					if err != nil {
						return err
					}

					return cmd.Run()
				}
			} else {
				err := cmd.Init(args)

				if err != nil {
					return err
				}

				return cmd.Run()
			}

			return nil
		}
	}

	return fmt.Errorf("Unknown subcommand %s", subcommand)
}

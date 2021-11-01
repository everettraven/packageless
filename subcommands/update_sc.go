package subcommands

import (
	"errors"
	"flag"
	"fmt"

	"github.com/everettraven/packageless/utils"
)

type UpdateCommand struct {
	//FlagSet for the Update command
	fs *flag.FlagSet
	//Name of the pim to be updated
	name string

	tools utils.Tools

	config utils.Config
}

//Instantiation method for a new UpdateCommand
func NewUpdateCommand(tools utils.Tools, config utils.Config) *UpdateCommand {
	//Create a new InstallCommand and set the FlagSet
	uc := &UpdateCommand{
		fs:     flag.NewFlagSet("update", flag.ContinueOnError),
		tools:  tools,
		config: config,
	}

	return uc
}

//Name - Gets the name of the Sub-Command
func (uc *UpdateCommand) Name() string {
	return uc.fs.Name()
}

//Initialize the command, for this particular subcommand we should just do nothing
func (uc *UpdateCommand) Init(args []string) error {
	if len(args) <= 0 {
		fmt.Println("No pim specified, updating all currently installed pim configurations.")
	} else {
		uc.name = args[0]
	}
	return nil
}

//Run the command, this command should fetch the pim config for
//either the specified package or all currently installed packages
func (uc *UpdateCommand) Run() error {
	pimConfigDir := uc.config.BaseDir + uc.config.PimsConfigDir
	//Get list of installed pims
	pims, err := uc.tools.GetListOfInstalledPimConfigs(pimConfigDir)

	if err != nil {
		return errors.New("Encountered an error while trying to fetch list of installed pim configuration files: " + err.Error())
	}

	//Loop and download most recent pim configuration for pims
	for _, pim := range pims {

		//If a specific package name is specified, skip over all the installed pims that are not the specified one
		if uc.name != "" {
			if pim != uc.name {
				continue
			}
		}

		fmt.Println("Updating pim: " + pim)
		err = uc.tools.FetchPimConfig(uc.config.RepositoryHost, pim, pimConfigDir)
		if err != nil {
			return errors.New("Encountered an error while trying to fetch the latest pim configuration file for pim '" + pim + "': " + err.Error())
		}
	}

	return nil
}

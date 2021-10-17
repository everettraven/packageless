package subcommands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

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
		fmt.Println("No pim specified, upgrading all currently installed pims.")
	} else {
		uc.name = args[0]
	}
	return nil
}

//Run the command, this command should fetch the pim config for
//either the specified package or all currently installed packages
func (uc *UpdateCommand) Run() error {
	pimDir := uc.config.BaseDir + uc.config.PimsConfigDir
	//Get list of installed pims
	var pims []string
	fileInfo, err := ioutil.ReadDir(pimDir)

	if err != nil {
		return err
	}

	for _, file := range fileInfo {
		//if a package name was specified, lets only update the one package
		if uc.name != "" {
			filename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			if filename == uc.name {
				pims = append(pims, filename)
				break
			}
		} else {
			pims = append(pims, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		}
	}

	//Loop and download most recent pim configuration for pims
	for _, pim := range pims {
		fmt.Println("Updating pim: " + pim)
		err = uc.tools.FetchPimConfig(uc.config.RepositoryHost, pim, pimDir)
		if err != nil {
			return err
		}
	}

	return nil
}

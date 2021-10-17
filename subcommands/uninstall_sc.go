package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/docker/docker/client"
	"github.com/everettraven/packageless/utils"
)

//Uninstall Sub-Command Object
type UninstallCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the pim to Uninstall
	name string

	tools utils.Tools

	config utils.Config
}

//Instantiation method for a new UninstallCommand
func NewUninstallCommand(tools utils.Tools, config utils.Config) *UninstallCommand {
	//Create a new UninstallCommand and set the FlagSet
	uc := &UninstallCommand{
		fs:     flag.NewFlagSet("uninstall", flag.ContinueOnError),
		tools:  tools,
		config: config,
	}

	return uc
}

//Name - Gets the name of the Sub-Command
func (uc *UninstallCommand) Name() string {
	return uc.fs.Name()
}

//Init - Parses and Populates values of the Uninstall subcommand
func (uc *UninstallCommand) Init(args []string) error {

	if len(args) <= 0 {
		return errors.New("No pim name was found. You must include the name of the pim you wish to uninstall.")
	}

	uc.name = args[0]

	return nil
}

//Uninstall - Uninstalls the Uninstall subcommand
func (uc *UninstallCommand) Run() error {
	//Create variables to use later
	var found bool
	var pim utils.PackageImage
	var version utils.Version

	var pimName string
	var pimVersion string

	if strings.Contains(uc.name, ":") {
		split := strings.Split(uc.name, ":")
		pimName = split[0]
		pimVersion = split[1]
	} else {
		pimName = uc.name
		pimVersion = "latest"
	}

	pimConfigDir := uc.config.BaseDir + uc.config.PimsConfigDir
	pimPath := pimConfigDir + pimName + ".hcl"

	pimDir := uc.config.BaseDir + uc.config.PimsDir

	//Check if pim config already exists
	if !uc.tools.FileExists(pimPath) {
		return errors.New("configuration for pim: " + pimName + " could not be found. Have you installed " + pimName + "?")
	}

	//Create the Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	pimListBody, err := uc.tools.GetHCLBody(pimPath)

	if err != nil {
		return err
	}

	//Parse the pim list
	parseOut, err := uc.tools.ParseBody(pimListBody, utils.PimHCLUtil{})

	//Check for errors
	if err != nil {
		return err
	}

	pims := parseOut.(utils.PimHCLUtil)

	//Check for errors
	if err != nil {
		return err
	}

	//Look for the pim we want in the pim list
	for _, pimItem := range pims.Pims {
		//If we find it, set some variables and break
		if pimItem.Name == pimName {
			pim = pimItem

			for _, ver := range pim.Versions {
				if ver.Version == pimVersion {
					found = true
					version = ver
					break
				}
			}
		}
	}

	//Make sure we have found the pim in the pim list
	if !found {
		return errors.New("Could not find pim " + pimName + " with version '" + pimVersion + "' in the pim list")
	}

	//Check if the corresponding pim is already Uninstalled
	imgExist, err := uc.tools.ImageExists(version.Image, cli)

	//Check for errors
	if err != nil {
		return err
	}

	//If the image doesn't exist it can't be uninstalled
	if !imgExist {
		return errors.New("pim " + pim.Name + " with version '" + version.Version + "' is not installed.")
	}

	fmt.Println("Removing", pim.Name+":"+version.Version)

	//Check for the directories that correspond to this pims volumes
	fmt.Println("Removing pim directories")

	//Check the volumes and remove the directories if they exist
	for _, vol := range version.Volumes {
		//Make sure that a path is given.
		if vol.Path != "" {
			err = uc.tools.RemoveDir(pimDir + vol.Path)

			if err != nil {
				return err
			}
		}
	}

	//Remove the base directory for the pim
	err = uc.tools.RemoveDir(pimDir + pim.BaseDir)

	if err != nil {
		return err
	}

	fmt.Println("Removing Image")

	//Remove the image
	err = uc.tools.RemoveImage(version.Image, cli)

	//Check for errors
	if err != nil {
		return err
	}

	if uc.config.Alias {
		//Remove aliases
		fmt.Println("Removing Alias")

		if runtime.GOOS == "windows" {
			if version.Version != "latest" {
				err = uc.tools.RemoveAliasWin(pim.Name+":"+version.Version, pimDir)
			} else {
				err = uc.tools.RemoveAliasWin(pim.Name, pimDir)
			}
		} else {
			if version.Version != "latest" {
				err = uc.tools.RemoveAliasUnix(pim.Name+":"+version.Version, pimDir)
			} else {
				err = uc.tools.RemoveAliasUnix(pim.Name, pimDir)
			}
		}

		if err != nil {
			return err
		}

		fmt.Println("Removing pim configuration")
		err = uc.tools.RemoveFile(pimPath)

		if err != nil {
			return err
		}
	}

	return nil
}

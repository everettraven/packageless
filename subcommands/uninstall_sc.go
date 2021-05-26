package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/everettraven/packageless/utils"
)

//Uninstall Sub-Command Object
type UninstallCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the package to Uninstall
	name string
}

//Instantiation method for a new UninstallCommand
func NewUninstallCommand() *UninstallCommand {
	//Create a new UninstallCommand and set the FlagSet
	uc := &UninstallCommand{
		fs: flag.NewFlagSet("uninstall", flag.ContinueOnError),
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
		return errors.New("No package name was found. You must include the name of the package you wish to uninstall.")
	}

	uc.name = args[0]

	return nil
}

//Uninstall - Uninstalls the Uninstall subcommand
func (uc *UninstallCommand) Run() error {
	//Create variables to use later
	var found bool
	var pack utils.Package

	//Default location of the package list
	packageList := "./package_list.hcl"

	//Parse the package list
	parseOut, err := utils.Parse(packageList, utils.PackageHCLUtil{})

	//Check for errors
	if err != nil {
		return err
	}

	packages := parseOut.(utils.PackageHCLUtil)

	//Check for errors
	if err != nil {
		return err
	}

	//Look for the package we want in the package list
	for _, packs := range packages.Packages {
		//If we find it, set some variables and break
		if packs.Name == uc.name {
			found = true
			pack = packs
			break
		}
	}

	//Make sure we have found the package in the package list
	if !found {
		return errors.New("Could not find package " + uc.name + " in the package list")
	}

	//Check if the corresponding package image is already Uninstalled
	imgExist, err := utils.ImageExists(pack.Image)

	//Check for errors
	if err != nil {
		return err
	}

	//If the image doesn't exist it can't be uninstalled
	if !imgExist {
		return errors.New("Package " + pack.Name + " is not installed.")
	}

	fmt.Println("Removing package", pack.Name)

	//Check for the directories that correspond to this packages volumes
	fmt.Println("Removing package directories")

	//Check the volumes and remove the directories if they exist
	for _, vol := range pack.Volumes {
		//Make sure that a path is given.
		if vol.Path != "" {
			err = RemoveDir(vol.Path)

			if err != nil {
				return err
			}
		}
	}

	//Remove the base directory for the package
	err = RemoveDir(pack.BaseDir)

	if err != nil {
		return err
	}

	//Remove the image
	err = utils.RemoveImage(pack.Image)

	//Check for errors
	if err != nil {
		return err
	}

	return nil
}

//RemoveDir removes a given directory if it exists
func RemoveDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		err = os.RemoveAll(path)
	}

	return nil
}

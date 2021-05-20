package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

//Install Sub-Command Object
type InstallCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the package to install
	name string
}

//Instantiation method for a new InstallCommand
func NewInstallCommand() *InstallCommand {
	//Create a new InstallCommand and set the FlagSet
	ic := &InstallCommand{
		fs: flag.NewFlagSet("install", flag.ContinueOnError),
	}

	return ic
}

//Name - Gets the name of the Sub-Command
func (ic *InstallCommand) Name() string {
	return ic.fs.Name()
}

//Init - Parses and Populates values of the Install subcommand
func (ic *InstallCommand) Init(args []string) error {
	ic.name = args[0]
	return nil
}

//Run - Runs the install subcommand
func (ic *InstallCommand) Run() error {
	//Create variables to use later
	var found bool
	var pack Package

	//Default location of the package list
	packageList := "./package_list.hcl"

	//Parse the package list
	packages, err := Parse(packageList)

	//Check for errors
	if err != nil {
		return err
	}

	//Look for the package we want in the package list
	for _, packs := range packages.Packages {
		//If we find it, set some variables and break
		if packs.Name == ic.name {
			found = true
			pack = packs
			break
		}
	}

	//Make sure we have found the package in the package list
	if !found {
		return errors.New("Could not find package " + ic.name + " in the package list")
	}

	fmt.Println("Installing", ic.name)
	//Pull the image down from Docker Hub
	err = PullImage(pack.Image)

	if err != nil {
		return err
	}

	fmt.Println("Creating package directories")
	//Check the copies and make the destination directory if it doesn't exist
	for _, copy := range pack.Copies {
		if _, err := os.Stat(copy.Dest); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(copy.Dest, 0755)

				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	fmt.Println("Copying necessary files 1/3")
	//Create the container so that we can copy the files over to the right places
	containerID, err := CreateContainer(pack.Image)

	if err != nil {
		return err
	}

	fmt.Println("Copying necessary files 2/3")
	//Copy the files from the container to the locations
	for _, copy := range pack.Copies {
		err = CopyFromContainer(copy.Source, copy.Dest, containerID)

		if err != nil {
			return err
		}
	}

	fmt.Println("Copying necessary files 3/3")
	//Remove the Container
	err = RemoveContainer(containerID)

	if err != nil {
		return err
	}

	fmt.Println(ic.name, "successfully installed")

	return nil
}

//Runner - Interface to enable easy interactions with the different subcommand objects
type Runner interface {
	Init(args []string) error
	Run() error
	Name() string
}

//SubCommand - Helper function that handles setting up and running subcommands
func SubCommand(args []string) error {
	if len(args) < 1 {
		return errors.New("A subcommand must be passed ")
	}

	cmds := []Runner{
		NewInstallCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			cmd.Init(os.Args[2:])
			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand %s", subcommand)
}

package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/everettraven/packageless/utils"
)

//Install Sub-Command Object
type UpgradeCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the package to install
	name string
}

//Instantiation method for a new InstallCommand
func NewUpgradeCommand() *UpgradeCommand {
	//Create a new InstallCommand and set the FlagSet
	ic := &UpgradeCommand{
		fs: flag.NewFlagSet("upgrade", flag.ContinueOnError),
	}

	return ic
}

//Name - Gets the name of the Sub-Command
func (ic *UpgradeCommand) Name() string {
	return ic.fs.Name()
}

//Init - Parses and Populates values of the Install subcommand
func (ic *UpgradeCommand) Init(args []string) error {
	if len(args) <= 0 {
		fmt.Println("No package specified, upgrading all currently installed packages.")
	} else {
		ic.name = args[0]
		fmt.Println("Upgrading", ic.name)
	}
	return nil
}

//Run - Runs the install subcommand
func (ic *UpgradeCommand) Run() error {
	//Create variables to use later
	var found bool
	var pack utils.Package

	//Default location of the package list
	packageList := "./package_list.hcl"

	//Parse the package list
	parseOut, err := utils.Parse(packageList, utils.PackageHCLUtil{})

	packages := parseOut.(utils.PackageHCLUtil)

	//Check for errors
	if err != nil {
		return err
	}

	if ic.name != "" {
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

		//Check if the corresponding package image is already installed
		imgExist, err := utils.ImageExists(pack.Image)

		//Check for errors
		if err != nil {
			return err
		}

		//If the image exists the package is already installed
		if !imgExist {
			return errors.New("Package: " + pack.Name + " is not installed. It must be installed before it can be upgraded.")
		}

		fmt.Println("Upgrading", pack.Name)
		//Pull the image down from Docker Hub
		err = utils.PullImage(pack.Image)

		if err != nil {
			return err
		}

		fmt.Println("Updating package directories")

		//Check the volumes and create the directories for them if they don't already exist
		for _, vol := range pack.Volumes {
			//Make sure that a path is given. If not we already assume that the working directory will be mounted
			if vol.Path != "" {
				if _, err := os.Stat(vol.Path); err != nil {
					if os.IsNotExist(err) {
						err = os.MkdirAll(vol.Path, 0755)

						if err != nil {
							return err
						}
					} else {
						return err
					}
				} else {
					//Remove the directory if it already exists
					err = os.RemoveAll(vol.Path)

					if err != nil {
						return err
					}

					//Recreate the directory
					err = os.MkdirAll(vol.Path, 0755)

					if err != nil {
						return err
					}
				}
			}
		}

		//Check and see if any files need to be copied from the container to one of the volumes on the host.
		if len(pack.Copies) > 0 {

			fmt.Println("Copying necessary files 1/3")
			//Create the container so that we can copy the files over to the right places
			containerID, err := utils.CreateContainer(pack.Image)

			if err != nil {
				return err
			}

			fmt.Println("Copying necessary files 2/3")
			//Copy the files from the container to the locations
			for _, copy := range pack.Copies {
				err = utils.CopyFromContainer(copy.Source, copy.Dest, containerID)

				if err != nil {
					return err
				}
			}

			fmt.Println("Copying necessary files 3/3")
			//Remove the Container
			err = utils.RemoveContainer(containerID)

			if err != nil {
				return err
			}

			fmt.Println(pack.Name, "successfully upgraded")

		}
	} else {
		//Loop through the packages in the package list
		for _, pack := range packages.Packages {
			//Check if the corresponding package image is already installed
			imgExist, err := utils.ImageExists(pack.Image)

			//Check for errors
			if err != nil {
				return err
			}

			//If the image exists the package is already installed
			if !imgExist {
				continue
			}

			fmt.Println("Upgrading", pack.Name)
			//Pull the image down from Docker Hub
			err = utils.PullImage(pack.Image)

			if err != nil {
				return err
			}

			fmt.Println("Updating package directories")

			//Check the volumes and create the directories for them if they don't already exist
			for _, vol := range pack.Volumes {
				//Make sure that a path is given. If not we already assume that the working directory will be mounted
				if vol.Path != "" {
					if _, err := os.Stat(vol.Path); err != nil {
						if os.IsNotExist(err) {
							err = os.MkdirAll(vol.Path, 0755)

							if err != nil {
								return err
							}
						} else {
							return err
						}
					} else {
						//Remove the directory if it already exists
						err = os.RemoveAll(vol.Path)

						if err != nil {
							return err
						}

						//Recreate the directory
						err = os.MkdirAll(vol.Path, 0755)

						if err != nil {
							return err
						}
					}
				}
			}

			//Check and see if any files need to be copied from the container to one of the volumes on the host.
			if len(pack.Copies) > 0 {

				fmt.Println("Copying necessary files 1/3")
				//Create the container so that we can copy the files over to the right places
				containerID, err := utils.CreateContainer(pack.Image)

				if err != nil {
					return err
				}

				fmt.Println("Copying necessary files 2/3")
				//Copy the files from the container to the locations
				for _, copy := range pack.Copies {
					err = utils.CopyFromContainer(copy.Source, copy.Dest, containerID)

					if err != nil {
						return err
					}
				}

				fmt.Println("Copying necessary files 3/3")
				//Remove the Container
				err = utils.RemoveContainer(containerID)

				if err != nil {
					return err
				}

				fmt.Println(pack.Name, "successfully upgraded")
			}

		}

	}

	return nil
}

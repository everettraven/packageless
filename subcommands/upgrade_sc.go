package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/everettraven/packageless/utils"
)

//Upgrade Sub-Command Object
type UpgradeCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the pim to upgrade
	name string

	tools utils.Tools

	cp utils.Copier

	config utils.Config
}

//Instantiation method for a new UpgradeCommand
func NewUpgradeCommand(tools utils.Tools, cp utils.Copier, config utils.Config) *UpgradeCommand {
	//Create a new UpgradeCommand and set the FlagSet
	ic := &UpgradeCommand{
		fs:     flag.NewFlagSet("upgrade", flag.ContinueOnError),
		tools:  tools,
		cp:     cp,
		config: config,
	}

	return ic
}

//Name - Gets the name of the Sub-Command
func (ic *UpgradeCommand) Name() string {
	return ic.fs.Name()
}

//Init - Parses and Populates values of the Upgrade subcommand
func (ic *UpgradeCommand) Init(args []string) error {
	if len(args) <= 0 {
		fmt.Println("No pim specified, upgrading all currently installed pims.")
	} else {
		ic.name = args[0]
	}
	return nil
}

//Run - Runs the Upgrade subcommand
func (ic *UpgradeCommand) Run() error {
	//Create variables to use later
	var found bool
	var pim utils.PackageImage
	var version utils.Version

	var pimName string
	var pimVersion string

	if strings.Contains(ic.name, ":") {
		split := strings.Split(ic.name, ":")
		pimName = split[0]
		pimVersion = split[1]
	} else {
		pimName = ic.name
		pimVersion = "latest"
	}

	//Create the Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	pimConfigDir := ic.config.BaseDir + ic.config.PimsConfigDir
	pimDir := ic.config.BaseDir + ic.config.PimsDir

	if pimName != "" {

		pimPath := pimConfigDir + pimName + ".hcl"

		//Check if pim config already exists
		if !ic.tools.FileExists(pimPath) {
			return errors.New("Could not find pim configuration for: " + pimName + " has it been installed?")
		}

		pimListBody, err := ic.tools.GetHCLBody(pimPath)

		if err != nil {
			return err
		}

		//Parse the pim list
		parseOut, err := ic.tools.ParseBody(pimListBody, utils.PimHCLUtil{})

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

		//Check if the corresponding pim image is already installed
		imgExist, err := ic.tools.ImageExists(version.Image, cli)

		//Check for errors
		if err != nil {
			return err
		}

		//If the image exists the pim is already installed
		if !imgExist {
			return errors.New("pim: " + pim.Name + " with version '" + version.Version + "' is not installed. It must be installed before it can be upgraded.")
		}

		fmt.Println("Upgrading", pim.Name+":"+version.Version)
		//Pull the image down from Docker Hub
		err = ic.tools.PullImage(version.Image, cli)

		if err != nil {
			return err
		}

		fmt.Println("Updating pim directories")

		//Check the volumes and create the directories for them if they don't already exist
		for _, vol := range version.Volumes {
			//Make sure that a path is given. If not we already assume that the working directory will be mounted
			if vol.Path != "" {
				err = ic.tools.UpgradeDir(pimDir + vol.Path)

				if err != nil {
					return err
				}
			}
		}

		//Check and see if any files need to be copied from the container to one of the volumes on the host.
		if len(version.Copies) > 0 {

			fmt.Println("Copying necessary files 1/3")
			//Create the container so that we can copy the files over to the right places
			containerID, err := ic.tools.CreateContainer(version.Image, cli)

			if err != nil {
				return err
			}

			fmt.Println("Copying necessary files 2/3")
			//Copy the files from the container to the locations
			for _, copy := range version.Copies {
				err = ic.tools.CopyFromContainer(copy.Source, pimDir+copy.Dest, containerID, cli, ic.cp)

				if err != nil {
					return err
				}
			}

			fmt.Println("Copying necessary files 3/3")
			//Remove the Container
			err = ic.tools.RemoveContainer(containerID, cli)

			if err != nil {
				return err
			}

			fmt.Println(pim.Name, "successfully upgraded")

		}
	} else {

		//Get list of installed pims
		pimNames, err := ic.tools.GetListOfInstalledPimConfigs(pimConfigDir)

		if err != nil {
			return errors.New("Encountered an error while trying to fetch list of installed pim configuration files: " + err.Error())
		}

		for _, pimName := range pimNames {
			pimPath := pimConfigDir + pimName + ".hcl"

			pimListBody, err := ic.tools.GetHCLBody(pimPath)

			if err != nil {
				return err
			}

			//Parse the pim list
			parseOut, err := ic.tools.ParseBody(pimListBody, utils.PimHCLUtil{})

			pims := parseOut.(utils.PimHCLUtil)

			if err != nil {
				return err
			}

			//Loop through the pims in the pim list
			for _, pim := range pims.Pims {

				for _, ver := range pim.Versions {
					//Check if the corresponding pim image is already installed
					imgExist, err := ic.tools.ImageExists(ver.Image, cli)

					//Check for errors
					if err != nil {
						return err
					}

					//If the image does not exist, then this version for the pim is not installed
					//and therefore does not need to be upgraded
					if !imgExist {
						continue
					}

					fmt.Println("Upgrading", pim.Name+":"+ver.Version)
					//Pull the image down from Docker Hub
					err = ic.tools.PullImage(ver.Image, cli)

					if err != nil {
						return err
					}

					fmt.Println("Updating pim directories")

					//Check the volumes and create the directories for them if they don't already exist
					for _, vol := range ver.Volumes {
						//Make sure that a path is given. If not we already assume that the working directory will be mounted
						if vol.Path != "" {
							err = ic.tools.UpgradeDir(pimDir + vol.Path)

							if err != nil {
								return err
							}
						}
					}

					//Check and see if any files need to be copied from the container to one of the volumes on the host.
					if len(ver.Copies) > 0 {

						fmt.Println("Copying necessary files 1/3")
						//Create the container so that we can copy the files over to the right places
						containerID, err := ic.tools.CreateContainer(ver.Image, cli)

						if err != nil {
							return err
						}

						fmt.Println("Copying necessary files 2/3")
						//Copy the files from the container to the locations
						for _, copy := range ver.Copies {
							err = ic.tools.CopyFromContainer(copy.Source, pimDir+copy.Dest, containerID, cli, ic.cp)

							if err != nil {
								return err
							}
						}

						fmt.Println("Copying necessary files 3/3")
						//Remove the Container
						err = ic.tools.RemoveContainer(containerID, cli)

						if err != nil {
							return err
						}

						fmt.Println(pim.Name, "successfully upgraded")
					}

				}

			}
		}

	}

	return nil
}

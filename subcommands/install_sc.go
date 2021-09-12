package subcommands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/client"
	"github.com/everettraven/packageless/utils"
)

//Install Sub-Command Object
type InstallCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the package to install
	name string

	//Tools that can be used by the command
	tools utils.Tools

	cp utils.Copier

	config utils.Config
}

//Instantiation method for a new InstallCommand
func NewInstallCommand(tools utils.Tools, cp utils.Copier, config utils.Config) *InstallCommand {
	//Create a new InstallCommand and set the FlagSet
	ic := &InstallCommand{
		fs:     flag.NewFlagSet("install", flag.ContinueOnError),
		tools:  tools,
		cp:     cp,
		config: config,
	}

	return ic
}

//Name - Gets the name of the Sub-Command
func (ic *InstallCommand) Name() string {
	return ic.fs.Name()
}

//Init - Parses and Populates values of the Install subcommand
func (ic *InstallCommand) Init(args []string) error {

	if len(args) <= 0 {
		return errors.New("No package name was found. You must include the name of the package you wish to install.")
	}

	ic.name = args[0]
	return nil
}

//Run - Runs the install subcommand
func (ic *InstallCommand) Run() error {
	//Create variables to use later
	var found bool
	var pack utils.Package
	var version utils.Version

	var packName string
	var packVersion string

	if strings.Contains(ic.name, ":") {
		split := strings.Split(ic.name, ":")
		packName = split[0]
		packVersion = split[1]
	} else {
		packName = ic.name
		packVersion = "latest"
	}

	//Create the Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	//Create a variable for the executable directory
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	ed := filepath.Dir(ex)

	//Default location of the package list
	packageList := ed + "/package_list.hcl"

	packageListBody, err := ic.tools.GetHCLBody(packageList)

	if err != nil {
		return err
	}

	//Parse the package list
	parseOut, err := ic.tools.ParseBody(packageListBody, utils.PackageHCLUtil{})

	packages := parseOut.(utils.PackageHCLUtil)

	//Check for errors
	if err != nil {
		return err
	}

	//Look for the package we want in the package list
	for _, packs := range packages.Packages {
		//If we find it, set some variables and break
		if packs.Name == packName {
			pack = packs

			for _, ver := range pack.Versions {
				if ver.Version == packVersion {
					found = true
					version = ver
					break
				}
			}
		}
	}

	//Make sure we have found the package in the package list
	if !found {
		return errors.New("Could not find package " + packName + " with version '" + packVersion + "' in the package list")
	}

	//Check if the corresponding package image is already installed
	imgExist, err := ic.tools.ImageExists(version.Image, cli)

	//Check for errors
	if err != nil {
		return err
	}

	//If the image exists the package is already installed
	if imgExist {
		return errors.New("Package " + pack.Name + " is already installed")
	}

	fmt.Println("Installing", pack.Name)
	//Pull the image down from Docker Hub
	err = ic.tools.PullImage(version.Image, cli)

	if err != nil {
		return err
	}

	fmt.Println("Creating package directories")

	//Create the base directory for the package
	err = ic.tools.MakeDir(ed + pack.BaseDir)

	if err != nil {
		return err
	}

	//Check the volumes and create the directories for them if they don't already exist
	for _, vol := range version.Volumes {
		//Make sure that a path is given. If not we already assume that the working directory will be mounted
		if vol.Path != "" {
			err = ic.tools.MakeDir(ed + vol.Path)

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
			err = ic.tools.CopyFromContainer(copy.Source, ed+copy.Dest, containerID, cli, ic.cp)

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

	}

	if ic.config.Alias {
		//Set the alias for the command
		fmt.Println("Setting Alias")

		if runtime.GOOS == "windows" {
			err = ic.tools.AddAliasWin(pack.Name, ed)
		} else {
			err = ic.tools.AddAliasUnix(pack.Name, ed)
		}

		if err != nil {
			return err
		}
	}

	fmt.Println(pack.Name, "successfully installed")

	return nil
}

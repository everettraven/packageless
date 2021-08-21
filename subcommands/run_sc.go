package subcommands

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/everettraven/packageless/utils"
)

//Run Sub-Command Object
type RunCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the package to run
	name string

	args []string

	tools utils.Tools

	config utils.Config
}

//Instantiation method for a new RunCommand
func NewRunCommand(tools utils.Tools, config utils.Config) *RunCommand {
	//Create a new RunCommand and set the FlagSet
	rc := &RunCommand{
		fs:     flag.NewFlagSet("run", flag.ContinueOnError),
		tools:  tools,
		config: config,
	}

	return rc
}

//Name - Gets the name of the Sub-Command
func (rc *RunCommand) Name() string {
	return rc.fs.Name()
}

//Init - Parses and Populates values of the Run subcommand
func (rc *RunCommand) Init(args []string) error {

	if len(args) <= 0 {
		return errors.New("No package name was found. You must include the name of the package you wish to run.")
	}

	rc.name = args[0]

	rc.args = args[1:]

	return nil
}

//Run - Runs the Run subcommand
func (rc *RunCommand) Run() error {
	//Create variables to use later
	var found bool
	var pack utils.Package

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

	packageListBody, err := rc.tools.GetHCLBody(packageList)

	if err != nil {
		return err
	}

	//Parse the package list
	parseOut, err := rc.tools.ParseBody(packageListBody, utils.PackageHCLUtil{})

	//Check for errors
	if err != nil {
		return err
	}

	packages := parseOut.(utils.PackageHCLUtil)

	//Look for the package we want in the package list
	for _, packs := range packages.Packages {
		//If we find it, set some variables and break
		if packs.Name == rc.name {
			found = true
			pack = packs
			break
		}
	}

	//Make sure we have found the package in the package list
	if !found {
		return errors.New("Could not find package " + rc.name + " in the package list")
	}

	//Check if the corresponding package image is already installed
	imgExist, err := rc.tools.ImageExists(pack.Image, cli)

	//Check for errors
	if err != nil {
		return err
	}

	//If the image exists the package is already installed
	if !imgExist {
		return errors.New("Package " + pack.Name + " is not installed. You must install the package before running it.")
	}

	//Create the variables to use when running the container
	var ports []string
	var volumes []string

	ports = append(ports, strconv.Itoa(rc.config.StartPort)+":"+pack.Port)

	for _, vol := range pack.Volumes {
		if vol.Path != "" {
			volumes = append(volumes, ed+vol.Path+":"+vol.Mount)
		} else {
			sourcePath, err := os.Getwd()

			if err != nil {
				return err
			}

			volumes = append(volumes, sourcePath+":"+vol.Mount)
		}
	}

	//Run the container
	_, err = rc.tools.RunContainer(pack.Image, ports, volumes, pack.Name, rc.args)

	if err != nil {
		return err
	}

	return nil
}

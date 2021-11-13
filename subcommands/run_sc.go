package subcommands

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	"github.com/everettraven/packageless/utils"
)

//Run Sub-Command Object
type RunCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the pim to run
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
		return errors.New("No pim name was found. You must include the name of the pim you wish to run.")
	}

	rc.name = args[0]

	rc.args = args[1:]

	return nil
}

//Run - Runs the Run subcommand
func (rc *RunCommand) Run() error {
	//Create variables to use later
	var found bool
	var pim utils.PackageImage
	var version utils.Version

	var pimName string
	var pimVersion string

	if strings.Contains(rc.name, ":") {
		split := strings.Split(rc.name, ":")
		pimName = split[0]
		pimVersion = split[1]
	} else {
		pimName = rc.name
		pimVersion = "latest"
	}

	//Create the Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	pimConfigDir := rc.config.BaseDir + rc.config.PimsConfigDir

	//Default location of the pim list
	pimList := pimConfigDir + pimName + ".hcl"

	if !rc.tools.FileExists(pimList) {
		return errors.New("Could not find a configuration file for '" + pimName + "' has it been installed?")
	}

	pimListBody, err := rc.tools.GetHCLBody(pimList)

	if err != nil {
		return err
	}

	//Parse the pim list
	parseOut, err := rc.tools.ParseBody(pimListBody, utils.PimHCLUtil{})

	//Check for errors
	if err != nil {
		return err
	}

	pims := parseOut.(utils.PimHCLUtil)

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
		return errors.New("Could not find pim " + pimName + " with version '" + pimVersion + "' in the pim configuration")
	}

	//Check if the corresponding pim image is already installed
	imgExist, err := rc.tools.ImageExists(version.Image, cli)

	//Check for errors
	if err != nil {
		return err
	}

	//If the image exists the pim is already installed
	if !imgExist {
		return errors.New("pim " + pim.Name + " with version '" + version.Version + "' is not installed. You must install the pim before running it.")
	}

	//Create the variables to use when running the container
	var ports []string
	var volumes []string

	ports = append(ports, strconv.Itoa(rc.config.StartPort)+":"+version.Port)

	pimDir := rc.config.BaseDir + rc.config.PimsDir

	for _, vol := range version.Volumes {
		if vol.Path != "" {
			volumes = append(volumes, pimDir+vol.Path+":"+vol.Mount)
		} else {
			sourcePath, err := os.Getwd()

			if err != nil {
				return err
			}

			volumes = append(volumes, sourcePath+":"+vol.Mount)
		}
	}

	//Run the container
	_, err = rc.tools.RunContainer(version.Image, ports, volumes, pim.Name, rc.args)

	if err != nil {
		return err
	}

	return nil
}

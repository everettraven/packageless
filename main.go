package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/everettraven/packageless/subcommands"
	"github.com/everettraven/packageless/utils"
)

func main() {
	exitCode, exitErr := wrappedMain()

	if exitErr != nil {
		fmt.Println(exitErr)
	}

	os.Exit(exitCode)
}

func wrappedMain() (int, error) {
	//Create the utils for the subcommands
	util := utils.NewUtility()

	//Create the copier for the subcommands
	cp := &utils.CopyTool{}

	//Create a variable for the executable directory
	ex, err := os.Executable()
	if err != nil {
		return 1, err
	}
	ed := filepath.Dir(ex)

	//Config file location
	configLoc := ed + "/config.hcl"

	configBody, err := util.GetHCLBody(configLoc)

	if err != nil {
		return 1, err
	}

	//Parse the config file
	parseOut, err := util.ParseBody(configBody, utils.Config{})

	if err != nil {
		return 1, err
	}

	config := parseOut.(utils.Config)

	//Create the list of subcommands
	scmds := []subcommands.Runner{
		subcommands.NewInstallCommand(util, cp, config),
		subcommands.NewUninstallCommand(util, config),
		subcommands.NewUpgradeCommand(util, cp),
		subcommands.NewRunCommand(util, config),
		subcommands.NewVersionCommand(),
	}

	//Run the subcommands
	if err := subcommands.SubCommand(os.Args[1:], scmds); err != nil {
		return 1, err
	}

	return 0, nil
}

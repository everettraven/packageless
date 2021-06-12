package main

import (
	"fmt"
	"os"

	"github.com/everettraven/packageless/subcommands"
	"github.com/everettraven/packageless/utils"
)

func main() {
	//Create the utils for the subcommands
	util := utils.NewUtility()

	//Create the copier for the subcommands
	cp := &utils.CopyTool{}

	//Create the list of subcommands
	scmds := []subcommands.Runner{
		subcommands.NewInstallCommand(util, cp),
		subcommands.NewUninstallCommand(util),
		subcommands.NewUpgradeCommand(util, cp),
		subcommands.NewRunCommand(util),
	}

	//Run the subcommands
	if err := subcommands.SubCommand(os.Args[1:], scmds); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

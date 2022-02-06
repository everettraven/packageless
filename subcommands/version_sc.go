package subcommands

import (
	"flag"
	"fmt"

	"github.com/everettraven/packageless/utils"
)

var version = "v0.0.0"

type VersionCommand struct {
	//FlagSet for the version command
	fs *flag.FlagSet

	tools utils.Tools
}

//Instantiation method for a new VersionCommand
func NewVersionCommand(tools utils.Tools) *VersionCommand {
	//Create a new InstallCommand and set the FlagSet
	vc := &VersionCommand{
		fs:    flag.NewFlagSet("version", flag.ContinueOnError),
		tools: tools,
	}

	return vc
}

//Name - Gets the name of the Sub-Command
func (vc *VersionCommand) Name() string {
	return vc.fs.Name()
}

//Initialize the command, for this particular subcommand we should just do nothing
func (vc *VersionCommand) Init(args []string) error {
	return nil
}

//Run the command, this particular command should be a
//simple print of the value of the version variable
func (vc *VersionCommand) Run() error {
	vc.tools.RenderInfoMarkdown(fmt.Sprintf("**Packageless Version**: *%s*", version))
	return nil
}

package subcommands

import (
	"errors"
	"reflect"
	"testing"
)

//Create a mock subcommand struct
type MockSC struct {
	//Create a variable to hold the init args
	Args []string

	//Create a variable to set the command Name
	CmdName string

	//Create a variable to see if the command was ran
	Ran bool

	//Create a variable to set where to error
	ErrorAt string

	//Create a variable to set the error message that should be returned
	ErrorMsg string
}

//Function to create a new MockSC
func NewMockSC() *MockSC {
	msc := &MockSC{}
	return msc
}

//MockSC Init function
func (msc *MockSC) Init(args []string) error {
	if msc.ErrorAt == "Init" {
		return errors.New(msc.ErrorMsg)
	}

	msc.Args = args

	return nil
}

//MockSC Name function
func (msc *MockSC) Name() string {
	return msc.CmdName
}

//MockSC Run function
func (msc *MockSC) Run() error {
	if msc.ErrorAt == "Run" {
		return errors.New(msc.ErrorMsg)
	}

	msc.Ran = true

	return nil
}

//Test the subcommands SubCommand Function
func TestSubCommand(t *testing.T) {
	//Create a mock subcommand
	msc := NewMockSC()

	//Set MockSC Values
	msc.CmdName = "test"

	//Create an argument array
	args := []string{msc.CmdName, "that"}

	//Create an array of Runner interface containing the mock subcommand
	scmds := []Runner{
		msc,
	}

	//Run the SubCommand function
	err := SubCommand(args, scmds)

	//Shouldn't have an error
	if err != nil {
		t.Fatal(err)
	}

	//Make sure that the init args match the args list we created
	if !reflect.DeepEqual(args[1:], msc.Args) {
		t.Fatalf("SubCommand: Expected Init Args: %v | Received: %v", args[1:], msc.Args)
	}

	//Make sure the command ran as no errors should have occurred
	if !msc.Ran {
		t.Fatalf("SubCommand: Expected the subcommand to have been ran but it was not")
	}
}

//Test the SubCommand function if an error happens at the Init function
func TestSubCommandErrorAtInit(t *testing.T) {
	//Create a mock subcommand
	msc := NewMockSC()

	//Tell it when to error
	msc.ErrorAt = "Init"

	//Set the error message
	msc.ErrorMsg = "Testing error at Init()"

	//Set MockSC Values
	msc.CmdName = "test"

	//Create an argument array
	args := []string{msc.CmdName, "that"}

	//Create an array of Runner interface containing the mock subcommand
	scmds := []Runner{
		msc,
	}

	//Run the SubCommand function
	err := SubCommand(args, scmds)

	//Should have an error
	if err == nil {
		t.Fatalf("SubCommand: Expected to have error: %s | Received No Error", msc.ErrorMsg)
	}

	if err != nil {
		if err.Error() != msc.ErrorMsg {
			t.Fatalf("SubCommand: Expected to have error: %s | Received: %s", msc.ErrorMsg, err.Error())
		}
	}
}

//Test the SubCommand function if an error happens at the Run function
func TestSubCommandErrorAtRun(t *testing.T) {
	//Create a mock subcommand
	msc := NewMockSC()

	//Tell it when to error
	msc.ErrorAt = "Run"

	//Set the error message
	msc.ErrorMsg = "Testing error at Run()"

	//Set MockSC Values
	msc.CmdName = "test"

	//Create an argument array
	args := []string{msc.CmdName, "that"}

	//Create an array of Runner interface containing the mock subcommand
	scmds := []Runner{
		msc,
	}

	//Run the SubCommand function
	err := SubCommand(args, scmds)

	//Should have an error
	if err == nil {
		t.Fatalf("SubCommand: Expected to have error: %s | Received No Error", msc.ErrorMsg)
	}

	if err != nil {
		if err.Error() != msc.ErrorMsg {
			t.Fatalf("SubCommand: Expected to have error: %s | Received: %s", msc.ErrorMsg, err.Error())
		}
	}
}

//Test the SubCommand function if an unknown subcommand is passed
func TestSubCommandUnknownSC(t *testing.T) {
	//Create a mock subcommand
	msc := NewMockSC()

	//Set MockSC Values
	msc.CmdName = "test"

	sc := "unknown"

	//Create an argument array
	args := []string{sc}

	exErr := "Unknown subcommand " + sc

	//Create an array of Runner interface containing the mock subcommand
	scmds := []Runner{
		msc,
	}

	//Run the SubCommand function
	err := SubCommand(args, scmds)

	//Should have an error
	if err == nil {
		t.Fatalf("SubCommand: Expected to have error: %s | Received No Error", exErr)
	}

	if err != nil {
		if err.Error() != exErr {
			t.Fatalf("SubCommand: Expected to have error: %s | Received: %s", exErr, err.Error())
		}
	}
}

//Test the SubCommand function if multiple pims passed to install subcommand
func TestSubCommandInstallMultiple(t *testing.T) {
	//Create a mock subcommand
	msc := NewMockSC()

	//Set MockSC Values
	msc.CmdName = "install"

	sc := "install"

	//Create an argument array
	args := []string{sc, "python", "git"}

	//Create an array of Runner interface containing the mock subcommand
	scmds := []Runner{
		msc,
	}

	//Run the SubCommand function
	err := SubCommand(args, scmds)

	//There shouldn't be an error
	if err != nil {
		t.Fatalf("SubCommand: Unexpected error: %s", err)
	}
}

//Test the SubCommand function if multiple pims passed to uninstall subcommand
func TestSubCommandUninstallMultiple(t *testing.T) {
	//Create a mock subcommand
	msc := NewMockSC()

	//Set MockSC Values
	msc.CmdName = "uninstall"

	sc := "uninstall"

	//Create an argument array
	args := []string{sc, "python", "git"}

	//Create an array of Runner interface containing the mock subcommand
	scmds := []Runner{
		msc,
	}

	//Run the SubCommand function
	err := SubCommand(args, scmds)

	//There shouldn't be an error
	if err != nil {
		t.Fatalf("SubCommand: Unexpected error: %s", err)
	}
}

//Test the SubCommand function if multiple pims passed to upgrade subcommand
func TestSubCommandUpgradeMultiple(t *testing.T) {
	//Create a mock subcommand
	msc := NewMockSC()

	//Set MockSC Values
	msc.CmdName = "upgrade"

	sc := "upgrade"

	//Create an argument array
	args := []string{sc, "python", "git"}

	//Create an array of Runner interface containing the mock subcommand
	scmds := []Runner{
		msc,
	}

	//Run the SubCommand function
	err := SubCommand(args, scmds)

	//There shouldn't be an error
	if err != nil {
		t.Fatalf("SubCommand: Unexpected error: %s", err)
	}
}

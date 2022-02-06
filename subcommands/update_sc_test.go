package subcommands

import (
	"reflect"
	"testing"

	"github.com/everettraven/packageless/utils"
)

//Test to make sure the update subcommand has the proper name upon creation
func TestUpdateName(t *testing.T) {
	mu := utils.NewMockUtility()

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	updateCommand := NewUpdateCommand(mu, config)

	if updateCommand.Name() != "update" {
		t.Fatal("The update subcommand's name should be: update | Subcommand Name: " + updateCommand.Name())
	}
}

//Test to make sure the update subcommand initializes correctly
func TestUpdateInit(t *testing.T) {
	mu := utils.NewMockUtility()

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	updateCommand := NewUpdateCommand(mu, config)

	args := []string{"python"}

	err := updateCommand.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	if updateCommand.name != args[0] {
		t.Fatal("pim Name should have been initialized as: " + args[0] + " but is: " + updateCommand.name)
	}
}

func TestUpdateNoArgsFlow(t *testing.T) {
	mu := utils.NewMockUtility()

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	mu.InstalledPims = []string{"python", "another"}

	updateCommand := NewUpdateCommand(mu, config)

	args := []string{}

	err := updateCommand.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = updateCommand.Run()

	if err != nil {
		t.Fatal(err)
	}

	callStack := []string{
		"RenderInfoMarkdown",
		"GetListOfInstalledPimConfigs",
		"RenderInfoMarkdown",
		"FetchPimConfig",
		"RenderInfoMarkdown",
		"FetchPimConfig",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	if !reflect.DeepEqual(mu.FetchedPims, mu.InstalledPims) {
		t.Fatalf("The fetched pims should be the same as the installed pims, but it is not. Fetched Pims: %v | Installed Pims: %v", mu.FetchedPims, mu.InstalledPims)
	}
}

func TestUpdateArgsFlow(t *testing.T) {
	mu := utils.NewMockUtility()

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	mu.InstalledPims = []string{"python", "another"}

	updateCommand := NewUpdateCommand(mu, config)

	args := []string{"python"}
	expectedFetchedPims := []string{"python"}

	err := updateCommand.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = updateCommand.Run()

	if err != nil {
		t.Fatal(err)
	}

	callStack := []string{
		"GetListOfInstalledPimConfigs",
		"RenderInfoMarkdown",
		"FetchPimConfig",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	if !reflect.DeepEqual(mu.FetchedPims, expectedFetchedPims) {
		t.Fatalf("The fetched pims does not match the expected. Fetched Pims: %v | Expected: %v", mu.FetchedPims, expectedFetchedPims)
	}
}

func TestUpdateErrorAtGetListOfInstalledPimConfigs(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "GetListOfInstalledPimConfigs"

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	updateCommand := NewUpdateCommand(mu, config)

	expectedErr := "Encountered an error while trying to fetch list of installed pim configuration files: " + mu.ErrorMsg

	args := []string{"python"}

	err := updateCommand.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = updateCommand.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"GetListOfInstalledPimConfigs",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

func TestUpdateErrorAtFetchPimConfig(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "FetchPimConfig"

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	mu.InstalledPims = []string{"python"}

	updateCommand := NewUpdateCommand(mu, config)

	expectedErr := "Encountered an error while trying to fetch the latest pim configuration file for pim 'python': " + mu.ErrorMsg

	args := []string{"python"}

	err := updateCommand.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = updateCommand.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"GetListOfInstalledPimConfigs",
		"RenderInfoMarkdown",
		"FetchPimConfig",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

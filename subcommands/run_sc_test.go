package subcommands

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/everettraven/packageless/utils"
)

//Test to make sure the run subcommand has the proper name upon creation
func TestRunName(t *testing.T) {
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

	rc := NewRunCommand(mu, config)

	if rc.Name() != "run" {
		t.Fatal("The run subcommand's name should be: run | Subcommand Name: " + rc.Name())
	}
}

//Test to make sure the Run subcommand initializes correctly
func TestRunInit(t *testing.T) {
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

	rc := NewRunCommand(mu, config)

	args := []string{"python"}

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	if rc.name != args[0] {
		t.Fatal("pim Name should have been initialized as: " + args[0] + " but is: " + rc.name)
	}
}

//Test the Run subcommand with no pim specified
func TestRunNoPackage(t *testing.T) {
	mu := utils.NewMockUtility()

	expectedErr := "No pim name was found. You must include the name of the pim you wish to run."

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	rc := NewRunCommand(mu, config)

	args := []string{}

	err := rc.Init(args)

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}
}

//Test the Run subcommand with a non existent pim specified
func TestRunNonExistPackage(t *testing.T) {
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

	rc := NewRunCommand(mu, config)

	args := []string{"nonexistent"}

	expectedErr := "Could not find pim nonexistent with version 'latest' in the pim configuration"

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Tests the Run subcommand if the image does not exist
func TestRunImageNotExist(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = false

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	rc := NewRunCommand(mu, config)

	args := []string{"python"}

	expectedErr := "pim python with version 'latest' is not installed. You must install the pim before running it."

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Tests the flow of a correctly ran Run subcommand with no run args
func TestRunFlowNoRunArgs(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	rc := NewRunCommand(mu, config)

	args := []string{"python"}

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RunContainer",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	//Make sure the image that was ran matches the pim image
	if mu.RunImage != mu.Pim.Pims[0].Versions[0].Image {
		t.Fatalf("RunContainer: Image does not match the expected Image. Received Image: %s | Expected Image: %s", mu.RunImage, mu.Pim.Pims[0].Versions[0].Image)
	}

	pimDir := config.BaseDir + config.PimsDir

	port := []string{strconv.Itoa(mu.Conf.StartPort) + ":" + mu.Pim.Pims[0].Versions[0].Port}
	volume := []string{pimDir + mu.Pim.Pims[0].Versions[0].Volumes[0].Path + ":" + mu.Pim.Pims[0].Versions[0].Volumes[0].Mount}

	//Make sure the ports passed in matches
	if !reflect.DeepEqual(mu.RunPorts, port) {
		t.Fatalf("RunContainer: Ports do not match the expected Ports. Received Ports: %v | Expected Ports: %v", mu.RunPorts, port)
	}

	//Make sure the volumes passed in matches
	if !reflect.DeepEqual(mu.RunVolumes, volume) {
		t.Fatalf("RunContainer: Volumes do not match the expected Volumes. Received Volumes: %v | Expected Volumes: %v", mu.RunVolumes, volume)
	}

	//Make sure there are no run args
	if len(mu.RunArgs) > 0 {
		t.Fatalf("RunContainer: RunArgs were received but no RunArgs were expected. Received RunArgs: %v", mu.RunArgs)
	}
}

//Tests the flow of a correctly ran Run subcommand with run args
func TestRunFlowRunArgs(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	rc := NewRunCommand(mu, config)

	args := []string{"python"}

	runArgs := []string{"-m", "pip", "install", "flask"}

	args = append(args, runArgs...)

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RunContainer",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	//Make sure the image that was ran matches the pim image
	if mu.RunImage != mu.Pim.Pims[0].Versions[0].Image {
		t.Fatalf("RunContainer: Image does not match the expected Image. Received Image: %s | Expected Image: %s", mu.RunImage, mu.Pim.Pims[0].Versions[0].Image)
	}

	pimDir := config.BaseDir + config.PimsDir

	port := []string{strconv.Itoa(mu.Conf.StartPort) + ":" + mu.Pim.Pims[0].Versions[0].Port}
	volume := []string{pimDir + mu.Pim.Pims[0].Versions[0].Volumes[0].Path + ":" + mu.Pim.Pims[0].Versions[0].Volumes[0].Mount}

	//Make sure the ports passed in matches
	if !reflect.DeepEqual(mu.RunPorts, port) {
		t.Fatalf("RunContainer: Ports do not match the expected Ports. Received Ports: %v | Expected Ports: %v", mu.RunPorts, port)
	}

	//Make sure the volumes passed in matches
	if !reflect.DeepEqual(mu.RunVolumes, volume) {
		t.Fatalf("RunContainer: Volumes do not match the expected Volumes. Received Volumes: %v | Expected Volumes: %v", mu.RunVolumes, volume)
	}

	//Make sure there are no run args
	if !reflect.DeepEqual(mu.RunArgs, runArgs) {
		t.Fatalf("RunContainer: RunArgs do not match the expected RunArgs. Received RunArgs: %v | Expected RunArgs: %v", mu.RunArgs, runArgs)
	}
}

//Test if an error happens at the GetHCLBody function
func TestRunErrorAtGetHCLBody(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "GetHCLBody"

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	rc := NewRunCommand(mu, config)

	args := []string{"python"}

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test if there is an error from the ParseBody function
func TestRunErrorAtParseBody(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "ParseBody"

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	rc := NewRunCommand(mu, config)

	args := []string{"python"}

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test if there is an error from the ImageExists function
func TestRunErrorAtImageExists(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "ImageExists"

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	rc := NewRunCommand(mu, config)

	args := []string{"python"}

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test if there is an error from the RunContainer function
func TestRunErrorAtRunContainer(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "RunContainer"

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	rc := NewRunCommand(mu, config)

	args := []string{"python"}

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RunContainer",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test the Run subcommand with a pim with a nonexistent version specified
func TestRunNonExistVersion(t *testing.T) {
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

	rc := NewRunCommand(mu, config)

	args := []string{"python:idontexist"}

	expectedErr := "Could not find pim python with version 'idontexist' in the pim configuration"

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

func TestRunFileNotExists(t *testing.T) {
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

	mu.PimConfigShouldExist = false

	rc := NewRunCommand(mu, config)

	args := []string{"python:idontexist"}

	expectedErr := "Could not find a configuration file for 'python' has it been installed?"

	err := rc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = rc.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

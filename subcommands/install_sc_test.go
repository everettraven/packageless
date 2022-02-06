package subcommands

import (
	"reflect"
	"testing"

	"github.com/everettraven/packageless/utils"
)

//Test to make sure the install subcommand has the proper name upon creation
func TestInstallName(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	if ic.Name() != "install" {
		t.Fatal("The install subcommand's name should be: install | Subcommand Name: " + ic.Name())
	}
}

//Test to make sure the install subcommand initializes correctly
func TestInstallInit(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	if ic.name != args[0] {
		t.Fatal("pim Name should have been initialized as: " + args[0] + " but is: " + ic.name)
	}
}

//Tests the flow of a correctly ran install subcommand
func TestInstallFlow(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
		"PullImage",
		"RenderInfoMarkdown",
		"MakeDir",
		"MakeDir",
		"RenderInfoMarkdown",
		"CreateContainer",
		"RenderInfoMarkdown",
		"CopyFromContainer",
		"RenderInfoMarkdown",
		"RemoveContainer",
		"RenderInfoMarkdown",
		"AddAlias",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	//Make a list of images that should have been pulled and make sure it matches from the MockUtility
	var images []string

	//Lists of copy data
	var copySources []string
	var copyDests []string

	//directories to be created
	var mkdirs []string

	//commands to have aliases created
	var aliasCmds []string

	pimDir := config.BaseDir + config.PimsDir
	pimConfigDir := config.BaseDir + config.PimsConfigDir

	mkdirs = append(mkdirs, pimConfigDir)
	mkdirs = append(mkdirs, pimDir)
	//Fill lists
	for _, pim := range mu.Pim.Pims {
		//Just use the first version
		version := pim.Versions[0]
		images = append(images, version.Image)
		mkdirs = append(mkdirs, pimDir+pim.BaseDir)
		aliasCmds = append(aliasCmds, pim.Name)

		//Loop through volumes in the pim
		for _, vol := range version.Volumes {
			mkdirs = append(mkdirs, pimDir+vol.Path)
		}

		//Loop through the copies in the pim
		for _, copy := range version.Copies {
			copySources = append(copySources, copy.Source)
			copyDests = append(copyDests, pimDir+copy.Dest)
		}

		//Just use the first pim
		break
	}

	//If the pulled images doesn't match the test fails
	if !reflect.DeepEqual(images, mu.PulledImgs) {
		t.Fatalf("Pulled Images does not match the expected Pulled Images. Pulled Images: %v | Expected Pulled Images: %v", mu.PulledImgs, images)
	}

	//If the directories made don't match, the test fails
	if !reflect.DeepEqual(mkdirs, mu.MadeDirs) {
		t.Fatalf("Made directories does not match the expected directories. Made Directories: %v | Expected Made Directories: %v", mu.MadeDirs, mkdirs)
	}

	//Make sure that the image passed into the CreateContainer function is correct
	if !reflect.DeepEqual(mu.CreateImages, images) {
		t.Fatalf("CreateContainer images does not match the expected images. Images: %v | Expected Images: %v", mu.CreateImages, images)
	}

	//Make sure the proper ContainerID is being passed into the CopyFromContainer function
	if mu.CopyContainerID != mu.ContainerID {
		t.Fatalf("CopyFromContainer ContainerID does not match the expected ContainerID. ContainerID: %s | Expected ContainerID: %s", mu.CopyContainerID, mu.ContainerID)
	}

	//Ensure that the Copy sources are correct
	if !reflect.DeepEqual(mu.CopySources, copySources) {
		t.Fatalf("CopyFromContainer Copy Sources does not match the expected Copy Sources. Copy Sources: %v | Expected Copy Sources: %v", mu.CopySources, copySources)
	}

	//Ensure that the Copy destinations are correct
	if !reflect.DeepEqual(mu.CopyDests, copyDests) {
		t.Fatalf("CopyFromContainer Copy Destinations does not match the expected Copy Destinations. Copy Destinations: %v | Expected Copy Destinations: %v", mu.CopyDests, copyDests)
	}

	//Ensure that the ContainerID is passed correctly to the RemoveContainer function
	if mu.RemoveContainerID != mu.ContainerID {
		t.Fatalf("RemoveContainer ContainerID does not match the expected ContainerID. ContainerID: %s | Expected ContainerID: %s", mu.RemoveContainerID, mu.ContainerID)
	}

	//Make sure that the commands being passed to the alias functions is correct
	if !reflect.DeepEqual(mu.CmdToAlias, aliasCmds) {
		t.Fatalf("AddAlias Alias Commands does not match the expected Alias Commands. Alias Commands: %v | Expected Alias Commands: %v", mu.CmdToAlias, aliasCmds)
	}

}

//Test the install subcommand getting an error after calling the GetHCLBody function
func TestInstallErrorAtGetHCLBody(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "GetHCLBody"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand getting an error after calling the ParseBody function
func TestInstallErrorAtParseBody(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "ParseBody"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand getting an error after calling the ImageExists function
func TestInstallErrorAtImageExists(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "ImageExists"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand getting an error after calling the PullImage function
func TestInstallErrorAtPullImage(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "PullImage"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
		"PullImage",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand getting an error after calling the MakeDir function
func TestInstallErrorAtMakeDir(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "MakeDir"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand getting an error after calling the CreateContainer function
func TestInstallErrorAtCreateContainer(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "CreateContainer"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
		"PullImage",
		"RenderInfoMarkdown",
		"MakeDir",
		"MakeDir",
		"RenderInfoMarkdown",
		"CreateContainer",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand getting an error after calling the CopyFromContainer function
func TestInstallErrorAtCopyFromContainer(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "CopyFromContainer"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
		"PullImage",
		"RenderInfoMarkdown",
		"MakeDir",
		"MakeDir",
		"RenderInfoMarkdown",
		"CreateContainer",
		"RenderInfoMarkdown",
		"CopyFromContainer",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand getting an error after calling the RemoveContainer function
func TestInstallErrorAtRemoveContainer(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "RemoveContainer"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
		"PullImage",
		"RenderInfoMarkdown",
		"MakeDir",
		"MakeDir",
		"RenderInfoMarkdown",
		"CreateContainer",
		"RenderInfoMarkdown",
		"CopyFromContainer",
		"RenderInfoMarkdown",
		"RemoveContainer",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand getting an error after calling the corresponding AddAlias function
func TestInstallErrorAtAddAlias(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ErrorAt = "AddAlias"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + " but did not receive an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatal("Expected the following error: " + mu.ErrorMsg + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
		"PullImage",
		"RenderInfoMarkdown",
		"MakeDir",
		"MakeDir",
		"RenderInfoMarkdown",
		"CreateContainer",
		"RenderInfoMarkdown",
		"CopyFromContainer",
		"RenderInfoMarkdown",
		"RemoveContainer",
		"RenderInfoMarkdown",
		"AddAlias",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the install subcommand when ImageExists function returns true
func TestInstallImageExists(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	expectedErr := "pim python is already installed"

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test the install subcommand with no arguments passed
func TestInstallNoPackage(t *testing.T) {
	mu := utils.NewMockUtility()

	expectedErr := "No pim name was found. You must include the name of the pim you wish to install."

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{}

	err := ic.Init(args)

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}
}

//Test the install subcommand if the passed in pim does not exist
func TestInstallNonExistPackage(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"nonexistent"}

	expectedErr := "Could not find pim nonexistent with version 'latest' in the pim list"

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test if the config has the Alias property set to false
func TestInstallAliasFalse(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          false,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
		"PullImage",
		"RenderInfoMarkdown",
		"MakeDir",
		"MakeDir",
		"RenderInfoMarkdown",
		"CreateContainer",
		"RenderInfoMarkdown",
		"CopyFromContainer",
		"RenderInfoMarkdown",
		"RemoveContainer",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	//Make a list of images that should have been pulled and make sure it matches from the MockUtility
	var images []string

	//Lists of copy data
	var copySources []string
	var copyDests []string

	//directories to be created
	var mkdirs []string

	pimDir := config.BaseDir + config.PimsDir
	pimConfigDir := config.BaseDir + config.PimsConfigDir

	mkdirs = append(mkdirs, pimConfigDir)
	mkdirs = append(mkdirs, pimDir)

	//Fill lists
	for _, pim := range mu.Pim.Pims {
		//Just use the first version
		version := pim.Versions[0]
		images = append(images, version.Image)
		mkdirs = append(mkdirs, pimDir+pim.BaseDir)

		//Loop through volumes in the pim
		for _, vol := range version.Volumes {
			mkdirs = append(mkdirs, pimDir+vol.Path)
		}

		//Loop through the copies in the pim
		for _, copy := range version.Copies {
			copySources = append(copySources, copy.Source)
			copyDests = append(copyDests, pimDir+copy.Dest)
		}

		//Just use the first pim
		break
	}

	//If the pulled images doesn't match the test fails
	if !reflect.DeepEqual(images, mu.PulledImgs) {
		t.Fatalf("Pulled Images does not match the expected Pulled Images. Pulled Images: %v | Expected Pulled Images: %v", mu.PulledImgs, images)
	}

	//If the directories made don't match, the test fails
	if !reflect.DeepEqual(mkdirs, mu.MadeDirs) {
		t.Fatalf("Made directories does not match the expected directories. Made Directories: %v | Expected Made Directories: %v", mu.MadeDirs, mkdirs)
	}

	//Make sure that the image passed into the CreateContainer function is correct
	if !reflect.DeepEqual(mu.CreateImages, images) {
		t.Fatalf("CreateContainer images does not match the expected images. Images: %v | Expected Images: %v", mu.CreateImages, images)
	}

	//Make sure the proper ContainerID is being passed into the CopyFromContainer function
	if mu.CopyContainerID != mu.ContainerID {
		t.Fatalf("CopyFromContainer ContainerID does not match the expected ContainerID. ContainerID: %s | Expected ContainerID: %s", mu.CopyContainerID, mu.ContainerID)
	}

	//Ensure that the Copy sources are correct
	if !reflect.DeepEqual(mu.CopySources, copySources) {
		t.Fatalf("CopyFromContainer Copy Sources does not match the expected Copy Sources. Copy Sources: %v | Expected Copy Sources: %v", mu.CopySources, copySources)
	}

	//Ensure that the Copy destinations are correct
	if !reflect.DeepEqual(mu.CopyDests, copyDests) {
		t.Fatalf("CopyFromContainer Copy Destinations does not match the expected Copy Destinations. Copy Destinations: %v | Expected Copy Destinations: %v", mu.CopyDests, copyDests)
	}

	//Ensure that the ContainerID is passed correctly to the RemoveContainer function
	if mu.RemoveContainerID != mu.ContainerID {
		t.Fatalf("RemoveContainer ContainerID does not match the expected ContainerID. ContainerID: %s | Expected ContainerID: %s", mu.RemoveContainerID, mu.ContainerID)
	}

}

//Test the install function with a nonexistent pim version
func TestInstallNonExistVersion(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python:idontexist"}

	expectedErr := "Could not find pim python with version 'idontexist' in the pim list"

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

func TestInstallFlowPimConfigNotExist(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

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

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"FetchPimConfig",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
		"PullImage",
		"RenderInfoMarkdown",
		"MakeDir",
		"MakeDir",
		"RenderInfoMarkdown",
		"CreateContainer",
		"RenderInfoMarkdown",
		"CopyFromContainer",
		"RenderInfoMarkdown",
		"RemoveContainer",
		"RenderInfoMarkdown",
		"AddAlias",
		"RenderInfoMarkdown",
		"RenderInfoMarkdown",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	//Make a list of images that should have been pulled and make sure it matches from the MockUtility
	var images []string

	//Lists of copy data
	var copySources []string
	var copyDests []string

	//directories to be created
	var mkdirs []string

	//commands to have aliases created
	var aliasCmds []string

	pimDir := config.BaseDir + config.PimsDir
	pimConfigDir := config.BaseDir + config.PimsConfigDir

	mkdirs = append(mkdirs, pimConfigDir)
	mkdirs = append(mkdirs, pimDir)
	//Fill lists
	for _, pim := range mu.Pim.Pims {
		//Just use the first version
		version := pim.Versions[0]
		images = append(images, version.Image)
		mkdirs = append(mkdirs, pimDir+pim.BaseDir)
		aliasCmds = append(aliasCmds, pim.Name)

		//Loop through volumes in the pim
		for _, vol := range version.Volumes {
			mkdirs = append(mkdirs, pimDir+vol.Path)
		}

		//Loop through the copies in the pim
		for _, copy := range version.Copies {
			copySources = append(copySources, copy.Source)
			copyDests = append(copyDests, pimDir+copy.Dest)
		}

		//Just use the first pim
		break
	}

	//If the pulled images doesn't match the test fails
	if !reflect.DeepEqual(images, mu.PulledImgs) {
		t.Fatalf("Pulled Images does not match the expected Pulled Images. Pulled Images: %v | Expected Pulled Images: %v", mu.PulledImgs, images)
	}

	//If the directories made don't match, the test fails
	if !reflect.DeepEqual(mkdirs, mu.MadeDirs) {
		t.Fatalf("Made directories does not match the expected directories. Made Directories: %v | Expected Made Directories: %v", mu.MadeDirs, mkdirs)
	}

	//Make sure that the image passed into the CreateContainer function is correct
	if !reflect.DeepEqual(mu.CreateImages, images) {
		t.Fatalf("CreateContainer images does not match the expected images. Images: %v | Expected Images: %v", mu.CreateImages, images)
	}

	//Make sure the proper ContainerID is being passed into the CopyFromContainer function
	if mu.CopyContainerID != mu.ContainerID {
		t.Fatalf("CopyFromContainer ContainerID does not match the expected ContainerID. ContainerID: %s | Expected ContainerID: %s", mu.CopyContainerID, mu.ContainerID)
	}

	//Ensure that the Copy sources are correct
	if !reflect.DeepEqual(mu.CopySources, copySources) {
		t.Fatalf("CopyFromContainer Copy Sources does not match the expected Copy Sources. Copy Sources: %v | Expected Copy Sources: %v", mu.CopySources, copySources)
	}

	//Ensure that the Copy destinations are correct
	if !reflect.DeepEqual(mu.CopyDests, copyDests) {
		t.Fatalf("CopyFromContainer Copy Destinations does not match the expected Copy Destinations. Copy Destinations: %v | Expected Copy Destinations: %v", mu.CopyDests, copyDests)
	}

	//Ensure that the ContainerID is passed correctly to the RemoveContainer function
	if mu.RemoveContainerID != mu.ContainerID {
		t.Fatalf("RemoveContainer ContainerID does not match the expected ContainerID. ContainerID: %s | Expected ContainerID: %s", mu.RemoveContainerID, mu.ContainerID)
	}

	//Make sure that the commands being passed to the alias functions is correct
	if !reflect.DeepEqual(mu.CmdToAlias, aliasCmds) {
		t.Fatalf("AddAlias Alias Commands does not match the expected Alias Commands. Alias Commands: %v | Expected Alias Commands: %v", mu.CmdToAlias, aliasCmds)
	}
}

func TestInstallErrorAtFetchPimConfig(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

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

	mu.ErrorAt = "FetchPimConfig"
	mu.ErrorMsg = "Test error message"

	ic := NewInstallCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err == nil {
		t.Fatal("Should have thrown an error")
	}

	if err.Error() != mu.ErrorMsg {
		t.Fatalf("Expected error: %s | Received: %s", mu.ErrorMsg, err.Error())
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"MakeDir",
		"MakeDir",
		"FileExists",
		"FetchPimConfig",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

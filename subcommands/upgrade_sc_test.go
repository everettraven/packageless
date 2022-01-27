package subcommands

import (
	"reflect"
	"testing"

	"github.com/everettraven/packageless/utils"
)

//Test to make sure the Upgrade subcommand has the proper name upon creation
func TestUpgradeName(t *testing.T) {
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

	ic := NewUpgradeCommand(mu, mcp, config)

	if ic.Name() != "upgrade" {
		t.Fatal("The Upgrade subcommand's name should be: upgrade | Subcommand Name: " + ic.Name())
	}
}

//Test to make sure the Upgrade subcommand initializes correctly
func TestUpgradeInit(t *testing.T) {
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

	ic := NewUpgradeCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	if ic.name != args[0] {
		t.Fatal("pim Name should have been initialized as: " + args[0] + " but is: " + ic.name)
	}
}

//Tests the flow of a correctly ran Upgrade subcommand
func TestUpgradeFlow(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
		"CreateContainer",
		"CopyFromContainer",
		"RemoveContainer",
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
	var updirs []string

	pimDir := config.BaseDir + config.PimsDir

	//Fill lists
	for _, pim := range mu.Pim.Pims {
		//Just use the first version
		version := pim.Versions[0]
		images = append(images, version.Image)

		//Loop through volumes in the pim
		for _, vol := range version.Volumes {
			updirs = append(updirs, pimDir+vol.Path)
		}

		//Loop through the copies in the pim
		for _, copy := range version.Copies {
			copySources = append(copySources, copy.Source)
			copyDests = append(copyDests, pimDir+copy.Dest)
		}

	}

	//If the pulled images doesn't match the test fails
	if !reflect.DeepEqual(images, mu.PulledImgs) {
		t.Fatalf("Pulled Images does not match the expected Pulled Images. Pulled Images: %v | Expected Pulled Images: %v", mu.PulledImgs, images)
	}

	//If the directories made don't match, the test fails
	if !reflect.DeepEqual(updirs, mu.UpgradedDirs) {
		t.Fatalf("Upgraded directories does not match the expected directories. Upgraded Directories: %v | Expected Upgraded Directories: %v", mu.MadeDirs, updirs)
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

//Test the Upgrade subcommand getting an error after calling the GetHCLBody function
func TestUpgradeErrorAtGetHCLBody(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the Upgrade subcommand getting an error after calling the ParseBody function
func TestUpgradeErrorAtParseBody(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the Upgrade subcommand getting an error after calling the ImageExists function
func TestUpgradeErrorAtImageExists(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the Upgrade subcommand getting an error after calling the PullImage function
func TestUpgradeErrorAtPullImage(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the Upgrade subcommand getting an error after calling the MakeDir function
func TestUpgradeErrorAtUpgradeDir(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "UpgradeDir"

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the Upgrade subcommand getting an error after calling the CreateContainer function
func TestUpgradeErrorAtCreateContainer(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
		"CreateContainer",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the Upgrade subcommand getting an error after calling the CopyFromContainer function
func TestUpgradeErrorAtCopyFromContainer(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
		"CreateContainer",
		"CopyFromContainer",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the Upgrade subcommand getting an error after calling the RemoveContainer function
func TestUpgradeErrorAtRemoveContainer(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
		"CreateContainer",
		"CopyFromContainer",
		"RemoveContainer",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

}

//Test the Upgrade subcommand when ImageExists function returns true
func TestUpgradeImageNotExists(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = false

	args := []string{"python"}
	expectedErr := "pim: python with version 'latest' is not installed. It must be installed before it can be upgraded."

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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test the Upgrade subcommand with no arguments passed and 2 packages in the pim list
func TestUpgradeNoPackageWithTwoPacks(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

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

	mu.InstalledPims = []string{"python", "second"}

	ic := NewUpgradeCommand(mu, mcp, config)

	args := []string{}

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
		"GetListOfInstalledPimConfigs",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
		"CreateContainer",
		"CopyFromContainer",
		"RemoveContainer",
		//Repeat the cycle from GetHCLBody since we should be reading a new pim file
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
		"CreateContainer",
		"CopyFromContainer",
		"RemoveContainer",
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
	var updirs []string

	//Fill lists
	for _, pim := range mu.Pim.Pims {
		//Since we are doing two packages we need to repeat this loop a second time
		for i := 0; i < 2; i++ {
			//Just get the first version
			version := pim.Versions[0]
			images = append(images, version.Image)

			//Loop through volumes in the pim
			for _, vol := range version.Volumes {
				updirs = append(updirs, config.BaseDir+config.PimsDir+vol.Path)
			}

			//Loop through the copies in the pim
			for _, copy := range version.Copies {
				copySources = append(copySources, copy.Source)
				copyDests = append(copyDests, config.BaseDir+config.PimsDir+copy.Dest)
			}
		}
	}

	//If the pulled images doesn't match the test fails
	if !reflect.DeepEqual(images, mu.PulledImgs) {
		t.Fatalf("Pulled Images does not match the expected Pulled Images. Pulled Images: %v | Expected Pulled Images: %v", mu.PulledImgs, images)
	}

	//If the directories made don't match, the test fails
	if !reflect.DeepEqual(updirs, mu.UpgradedDirs) {
		t.Fatalf("Upgraded directories does not match the expected directories. Upgraded Directories: %v | Expected Upgraded Directories: %v", mu.MadeDirs, updirs)
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

	pimConfigDir := config.BaseDir + config.PimsConfigDir
	//Make sure we are getting the correct pim config dir passed in
	if mu.PimConfigDir != pimConfigDir {
		t.Fatalf("The pim configuration directory was: %s | Expected: %s", mu.PimConfigDir, pimConfigDir)
	}

}

//Test the Upgrade subcommand if the passed in pim does not exist
func TestUpgradeNonExistPackage(t *testing.T) {
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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test the Upgrade subcommand if the passed in pim version does not exist
func TestUpgradeNonExistVersion(t *testing.T) {
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

	ic := NewUpgradeCommand(mu, mcp, config)

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

func TestUpgradeErrorAtGetListOfInstalledPimConfigs(t *testing.T) {
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

	mu.ErrorAt = "GetListOfInstalledPimConfigs"
	mu.ErrorMsg = "error message"

	ic := NewUpgradeCommand(mu, mcp, config)

	expectedErr := "Encountered an error while trying to fetch list of installed pim configuration files: " + mu.ErrorMsg

	err := ic.Init([]string{})

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
		"GetListOfInstalledPimConfigs",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

func TestUpgradeErrorAtPimConfigurationFileNotFound(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.PimConfigShouldExist = false

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

	ic := NewUpgradeCommand(mu, mcp, config)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	expectedErr := "Could not find pim configuration for: " + ic.name + " has it been installed?"

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

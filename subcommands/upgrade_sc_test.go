package subcommands

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/everettraven/packageless/utils"
)

//Test to make sure the Upgrade subcommand has the proper name upon creation
func TestUpgradeName(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	ic := NewUpgradeCommand(mu, mcp)

	if ic.Name() != "upgrade" {
		t.Fatal("The Upgrade subcommand's name should be: upgrade | Subcommand Name: " + ic.Name())
	}
}

//Test to make sure the Upgrade subcommand initializes correctly
func TestUpgradeInit(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	ic := NewUpgradeCommand(mu, mcp)

	args := []string{"python"}

	err := ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	if ic.name != args[0] {
		t.Fatal("Package Name should have been initialized as: " + args[0] + " but is: " + ic.name)
	}
}

//Tests the flow of a correctly ran Upgrade subcommand
func TestUpgradeFlow(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	//Get the executable directory
	ex, err := os.Executable()

	if err != nil {
		t.Fatal(err)
	}

	ed := filepath.Dir(ex)

	mcp := &utils.MockCopyTool{}

	ic := NewUpgradeCommand(mu, mcp)

	args := []string{"python"}

	err = ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
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
	for _, pack := range mu.Pack.Packages {
		//Just use the first version
		version := pack.Versions[0]
		images = append(images, version.Image)

		//Loop through volumes in the package
		for _, vol := range version.Volumes {
			updirs = append(updirs, ed+vol.Path)
		}

		//Loop through the copies in the package
		for _, copy := range version.Copies {
			copySources = append(copySources, copy.Source)
			copyDests = append(copyDests, ed+copy.Dest)
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

	ic := NewUpgradeCommand(mu, mcp)

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

	ic := NewUpgradeCommand(mu, mcp)

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

	ic := NewUpgradeCommand(mu, mcp)

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

	ic := NewUpgradeCommand(mu, mcp)

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

	ic := NewUpgradeCommand(mu, mcp)

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

	ic := NewUpgradeCommand(mu, mcp)

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

	ic := NewUpgradeCommand(mu, mcp)

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

	ic := NewUpgradeCommand(mu, mcp)

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
	expectedErr := "Package: python is not installed. It must be installed before it can be upgraded."

	mcp := &utils.MockCopyTool{}

	ic := NewUpgradeCommand(mu, mcp)

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
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test the Upgrade subcommand with no arguments passed and 2 packages in the package list
func TestUpgradeNoPackageWithTwoPacks(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.Pack.Packages = append(mu.Pack.Packages, utils.Package{
		Name:    "package",
		BaseDir: "/package",
		Versions: []utils.Version{
			{
				Version: "latest",
				Image:   "packageless/package",
				Volumes: []utils.Volume{
					{
						Path:  "/package/config/",
						Mount: "/package/config_data/",
					},
				},
				Port: "4000",
			},
		},
	})

	mu.ImgExist = true

	//Get the executable directory
	ex, err := os.Executable()

	if err != nil {
		t.Fatal(err)
	}

	ed := filepath.Dir(ex)

	mcp := &utils.MockCopyTool{}

	ic := NewUpgradeCommand(mu, mcp)

	args := []string{}

	err = ic.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = ic.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
		"CreateContainer",
		"CopyFromContainer",
		"RemoveContainer",
		"ImageExists",
		"PullImage",
		"UpgradeDir",
		//Second Package has no Copy fields so it should end at UpgradeDir
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
	for _, pack := range mu.Pack.Packages {
		//Just get the first version
		version := pack.Versions[0]
		images = append(images, version.Image)

		//Loop through volumes in the package
		for _, vol := range version.Volumes {
			updirs = append(updirs, ed+vol.Path)
		}

		//Loop through the copies in the package
		for _, copy := range version.Copies {
			copySources = append(copySources, copy.Source)
			copyDests = append(copyDests, ed+copy.Dest)
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
	if !reflect.DeepEqual(mu.CreateImages, images[:1]) {
		t.Fatalf("CreateContainer images does not match the expected images. Images: %v | Expected Images: %v", mu.CreateImages, images[:1])
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

//Test the Upgrade subcommand if the passed in package does not exist
func TestUpgradeNonExistPackage(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	ic := NewUpgradeCommand(mu, mcp)

	args := []string{"nonexistent"}

	expectedErr := "Could not find package nonexistent with version 'latest' in the package list"

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
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test the Upgrade subcommand if the passed in package version does not exist
func TestUpgradeNonExistVersion(t *testing.T) {
	mu := utils.NewMockUtility()

	mcp := &utils.MockCopyTool{}

	ic := NewUpgradeCommand(mu, mcp)

	args := []string{"python:idontexist"}

	expectedErr := "Could not find package python with version 'idontexist' in the package list"

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
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

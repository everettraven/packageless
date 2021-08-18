package subcommands

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/everettraven/packageless/utils"
)

//Test to make sure the uninstall subcommand has the proper name upon creation
func TestUninstallName(t *testing.T) {
	mu := utils.NewMockUtility()

	uc := NewUninstallCommand(mu)

	if uc.Name() != "uninstall" {
		t.Fatal("The uninstall subcommand's name should be: uninstall | Subcommand Name: " + uc.Name())
	}
}

//Test to make sure the uninstall subcommand initializes correctly
func TestUninstallInit(t *testing.T) {
	mu := utils.NewMockUtility()

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	if uc.name != args[0] {
		t.Fatal("Package Name should have been initialized as: " + args[0] + " but is: " + uc.name)
	}
}

//Test the uninstall subcommand with no package specified
func TestUninstallNoPackage(t *testing.T) {
	mu := utils.NewMockUtility()

	expectedErr := "No package name was found. You must include the name of the package you wish to uninstall."

	uc := NewUninstallCommand(mu)

	args := []string{}

	err := uc.Init(args)

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}
}

//Test the uninstall subcommand with a non existent package specified
func TestUninstallNonExistPackage(t *testing.T) {
	mu := utils.NewMockUtility()

	uc := NewUninstallCommand(mu)

	args := []string{"nonexistent"}

	expectedErr := "Could not find package nonexistent in the package list"

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

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

//Tests the uninstall subcommand if the image does not exist
func TestUninstallImageNotExist(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = false

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	expectedErr := "Package python is not installed."

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

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

//Tests the flow of a correctly ran uninstall subcommand
func TestUninstallFlow(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	//Get the executable directory
	ex, err := os.Executable()

	if err != nil {
		t.Fatal(err)
	}

	ed := filepath.Dir(ex)

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	err = uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RemoveDir",
		"RemoveDir",
		"RemoveImage",
		"RemoveAlias",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	//Make a list of images that should have been removed and make sure it matches from the MockUtility
	var images []string

	//directories to be removed
	var rmdirs []string

	//commands that should have had their aliases removed
	var aliasCmds []string

	//Fill lists
	for _, pack := range mu.Pack.Packages {
		images = append(images, pack.Image)
		aliasCmds = append(aliasCmds, pack.Name)

		//Loop through volumes in the package
		for _, vol := range pack.Volumes {
			rmdirs = append(rmdirs, ed+vol.Path)
		}

		rmdirs = append(rmdirs, ed+pack.BaseDir)
	}

	//If the pulled images doesn't match the test fails
	if !reflect.DeepEqual(images, mu.RemovedImgs) {
		t.Fatalf("Removed Images does not match the expected Removed Images. Removed Images: %v | Expected Removed Images: %v", mu.RemovedImgs, images)
	}

	//If the directories made don't match, the test fails
	if !reflect.DeepEqual(rmdirs, mu.RemovedDirs) {
		t.Fatalf("Removed directories does not match the expected directories. Removed Directories: %v | Expected Removed Directories: %v", mu.RemovedDirs, rmdirs)
	}

	//Make sure that the commands being passed to the alias functions is correct
	if !reflect.DeepEqual(mu.CmdToAlias, aliasCmds) {
		t.Fatalf("AddAlias Alias Commands does not match the expected Alias Commands. Alias Commands: %v | Expected Alias Commands: %v", mu.CmdToAlias, aliasCmds)
	}
}

//Test if an error happens at the GetHCLBody function
func TestUninstallErrorAtGetHCLBody(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "GetHCLBody"

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

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

//Test if there is an error from the ParseBody function
func TestUninstallErrorAtParseBody(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "ParseBody"

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

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

//Test if there is an error from the ImageExists function
func TestUninstallErrorAtImageExists(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "ImageExists"

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

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

//Test if there is an error from the RemoveDir function
func TestUninstallErrorAtRemoveDir(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "RemoveDir"

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

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
		"RemoveDir",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test if there is an error from the RemoveAlias function
func TestUninstallErrorAtRemoveAlias(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "RemoveAlias"

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

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
		"RemoveDir",
		"RemoveDir",
		"RemoveImage",
		"RemoveAlias",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

//Test if there is an error from the RemoveImage function
func TestUninstallErrorAtRemoveImage(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	mu.ErrorAt = "RemoveImage"

	uc := NewUninstallCommand(mu)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

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
		"RemoveDir",
		"RemoveDir",
		"RemoveImage",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

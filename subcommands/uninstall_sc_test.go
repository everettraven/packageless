package subcommands

import (
	"reflect"
	"testing"

	"github.com/everettraven/packageless/utils"
)

//Test to make sure the uninstall subcommand has the proper name upon creation
func TestUninstallName(t *testing.T) {
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

	uc := NewUninstallCommand(mu, config)

	if uc.Name() != "uninstall" {
		t.Fatal("The uninstall subcommand's name should be: uninstall | Subcommand Name: " + uc.Name())
	}
}

//Test to make sure the uninstall subcommand initializes correctly
func TestUninstallInit(t *testing.T) {
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

	uc := NewUninstallCommand(mu, config)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	if uc.name != args[0] {
		t.Fatal("pim Name should have been initialized as: " + args[0] + " but is: " + uc.name)
	}
}

//Test the uninstall subcommand with no pim specified
func TestUninstallNoPackage(t *testing.T) {
	mu := utils.NewMockUtility()

	expectedErr := "No pim name was found. You must include the name of the pim you wish to uninstall."

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	uc := NewUninstallCommand(mu, config)

	args := []string{}

	err := uc.Init(args)

	if err == nil {
		t.Fatal("Expected the following error: '" + expectedErr + "' but did not receive an error")
	}

	if err.Error() != expectedErr {
		t.Fatal("Expected the following error: " + expectedErr + "| Received: " + err.Error())
	}
}

//Test the uninstall subcommand with a non existent pim specified
func TestUninstallNonExistPackage(t *testing.T) {
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

	uc := NewUninstallCommand(mu, config)

	args := []string{"nonexistent"}

	expectedErr := "Could not find pim nonexistent with version 'latest' in the pim configuration"

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
		"FileExists",
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

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	uc := NewUninstallCommand(mu, config)

	args := []string{"python"}

	expectedErr := "pim python with version 'latest' is not installed."

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
		"FileExists",
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

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	uc := NewUninstallCommand(mu, config)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RemoveDir",
		"RemoveDir",
		"RemoveImage",
		"RemoveAlias",
		"RemoveFile",
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

	pimDir := config.BaseDir + config.PimsDir

	//Fill lists
	for _, pim := range mu.Pim.Pims {
		//Just use the first version
		version := pim.Versions[0]
		images = append(images, version.Image)
		aliasCmds = append(aliasCmds, pim.Name)

		//Loop through volumes in the pim
		for _, vol := range version.Volumes {
			rmdirs = append(rmdirs, pimDir+vol.Path)
		}

		rmdirs = append(rmdirs, pimDir+pim.BaseDir)
		//Just use the first pim for the test
		break
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

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	uc := NewUninstallCommand(mu, config)

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
		"FileExists",
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

	config := utils.Config{
		BaseDir:   "./",
		PortInc:   1,
		StartPort: 5000,
		Alias:     false,
	}

	uc := NewUninstallCommand(mu, config)

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
		"FileExists",
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

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	uc := NewUninstallCommand(mu, config)

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
		"FileExists",
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

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	uc := NewUninstallCommand(mu, config)

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
		"FileExists",
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

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	uc := NewUninstallCommand(mu, config)

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
		"FileExists",
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

	config := utils.Config{
		BaseDir:        "~/.packageless/",
		StartPort:      3000,
		PortInc:        1,
		Alias:          true,
		RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
		PimsConfigDir:  "pims_config/",
		PimsDir:        "pims/",
	}

	uc := NewUninstallCommand(mu, config)

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
		"FileExists",
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

//Test uninstall when config Alias attribute is set to false
func TestUninstallAliasFalse(t *testing.T) {
	mu := utils.NewMockUtility()

	mu.ImgExist = true

	config := utils.Config{
		BaseDir:   "./",
		PortInc:   1,
		StartPort: 5000,
		Alias:     false,
	}

	uc := NewUninstallCommand(mu, config)

	args := []string{"python"}

	err := uc.Init(args)

	if err != nil {
		t.Fatal(err)
	}

	err = uc.Run()

	if err != nil {
		t.Fatal(err)
	}

	//Set a variable with the proper call stack and see if the call stack matches
	callStack := []string{
		"FileExists",
		"GetHCLBody",
		"ParseBody",
		"ImageExists",
		"RemoveDir",
		"RemoveDir",
		"RemoveImage",
		"RemoveFile",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}

	//Make a list of images that should have been removed and make sure it matches from the MockUtility
	var images []string

	//directories to be removed
	var rmdirs []string

	pimDir := config.BaseDir + config.PimsDir

	//Fill lists
	for _, pim := range mu.Pim.Pims {
		//Just use the first version
		version := pim.Versions[0]
		images = append(images, version.Image)

		//Loop through volumes in the pim
		for _, vol := range version.Volumes {
			rmdirs = append(rmdirs, pimDir+vol.Path)
		}

		rmdirs = append(rmdirs, pimDir+pim.BaseDir)

		//Just use the first pim
		break
	}

	//If the pulled images doesn't match the test fails
	if !reflect.DeepEqual(images, mu.RemovedImgs) {
		t.Fatalf("Removed Images does not match the expected Removed Images. Removed Images: %v | Expected Removed Images: %v", mu.RemovedImgs, images)
	}

	//If the directories made don't match, the test fails
	if !reflect.DeepEqual(rmdirs, mu.RemovedDirs) {
		t.Fatalf("Removed directories does not match the expected directories. Removed Directories: %v | Expected Removed Directories: %v", mu.RemovedDirs, rmdirs)
	}
}

//Test the uninstall subcommand with a pim with a nonexistent version specified
func TestUninstallNonExistVersion(t *testing.T) {
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

	uc := NewUninstallCommand(mu, config)

	args := []string{"python:idontexist"}

	expectedErr := "Could not find pim python with version 'idontexist' in the pim configuration"

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
		"FileExists",
		"GetHCLBody",
		"ParseBody",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

func TestUninstallPimConfigFileNotExist(t *testing.T) {
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

	uc := NewUninstallCommand(mu, config)

	args := []string{"python:idontexist"}

	expectedErr := "configuration for pim: python could not be found. Have you installed python?"

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
		"FileExists",
	}

	if !reflect.DeepEqual(callStack, mu.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mu.Calls, callStack)
	}
}

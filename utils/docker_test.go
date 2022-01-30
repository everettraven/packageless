package utils

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

//Unit Tests
//------------------------------------------------------------------------

//Test PullImage Function
func TestPullImage(t *testing.T) {
	//Create a new Docker Client Mock
	dm := NewDockMock()

	//Create a new utility object
	util := NewUtility()

	//Set the image name to test
	img := "image"

	//Pull the image
	err := util.PullImage(img, dm)

	//If error occurs the test fails
	if err != nil {
		t.Fatal(err)
	}

	//Check that the proper image was passed to the docker client
	if dm.IPRefStr != img {
		t.Fatal("PullImage: Image passed into the Docker SDK ImagePull Function should be '" + img + "' | Received: " + dm.IPRefStr)
	}

}

//Test PullImage Function when it returns an error
func TestPullImageError(t *testing.T) {
	//Create a new Docker Client Mock
	dm := NewDockMock()

	//Create a new utility object
	util := NewUtility()

	//Set the image name to test
	img := "image"

	//Set the error at and error message
	dm.ErrorAt = "ImagePull"
	dm.ErrorMsg = "Testing error at ImagePull()"

	//Pull the image
	err := util.PullImage(img, dm)

	//Error should occur
	if err == nil {
		t.Fatal("PullImage: Expected to receive an error, but did not receive one.")
	}

	if err != nil {
		if err.Error() != dm.ErrorMsg {
			t.Fatal("PullImage: Expected Error: " + dm.ErrorMsg + " | Received Error: " + err.Error())
		}
	}

}

//Test ImageExists Function when the image does exist
func TestImageExistsDoesExist(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the image for testing
	img := "image:faketag"

	//Set the return images array in the Mock Docker Client
	dm.ILRet = []types.ImageSummary{
		{
			RepoTags: []string{"image:faketag"},
		},
	}

	//Create a new utility
	util := NewUtility()

	//Check for the image
	exists, err := util.ImageExists(img, dm)

	//If an error occurs, the test fails
	if err != nil {
		t.Fatal(err)
	}

	//The image should exist in this case
	if !exists {
		t.Fatal("ImageExists: Image should exist, but it does not.")
	}

}

//Test ImageExists Function when the image does not exist
func TestImageExistsDoesNotExist(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the image for testing
	img := "image"

	//Set the return images array in the Mock Docker Client
	dm.ILRet = []types.ImageSummary{}

	//Create a new utility
	util := NewUtility()

	//Check for the image
	exists, err := util.ImageExists(img, dm)

	//If an error occurs, the test fails
	if err != nil {
		t.Fatal(err)
	}

	//The image should not exist in this case
	if exists {
		t.Fatal("ImageExists: Image should not exist, but it does.")
	}
}

//Test the ImageExists Function when an error occurs
func TestImageExistError(t *testing.T) {
	//Create a new Docker Client Mock
	dm := NewDockMock()

	//Create a new utility object
	util := NewUtility()

	//Set the image name to test
	img := "image"

	//Set the error at and error message
	dm.ErrorAt = "ImageList"
	dm.ErrorMsg = "Testing error at ImageList()"

	//Set the return images array in the Mock Docker Client
	dm.ILRet = []types.ImageSummary{}

	//Check for the image
	_, err := util.ImageExists(img, dm)

	//Error should occur
	if err == nil {
		t.Fatal("ImageExists: Expected to receive an error, but did not receive one.")
	}

	if err != nil {
		if err.Error() != dm.ErrorMsg {
			t.Fatal("ImageExists: Expected Error: " + dm.ErrorMsg + " | Received Error: " + err.Error())
		}
	}
}

//Test ImageExists function when RepoTags is not populated
func TestImageExistContinueOnEmptyRepoTag(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the image for testing
	img := "image:faketag"

	//Set the return images array in the Mock Docker Client
	dm.ILRet = []types.ImageSummary{
		// First one should be an empty RepoTags
		{
			RepoTags: []string{},
		},
		// When image loop continues it should find this one.
		{
			RepoTags: []string{"image:faketag"},
		},
	}

	//Create a new utility
	util := NewUtility()

	//Check for the image
	exists, err := util.ImageExists(img, dm)

	//If an error occurs, the test fails
	if err != nil {
		t.Fatal(err)
	}

	//The image should exist in this case
	if !exists {
		t.Fatal("ImageExists: Image should exist, but it does not.")
	}

}

//Test CreateContainer Function
func TestCreateContainer(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the image that should be used
	img := "image"

	//Set what the containerID should be
	containerID := "testcontainer"

	//Set what the create container cmd should be
	cmd := []string{"bash"}

	dm.CCRet = container.ContainerCreateCreatedBody{
		ID: containerID,
	}

	//Create the util
	util := NewUtility()

	//Test creating the container
	cID, err := util.CreateContainer(img, dm)

	//If there is an error then the test fails
	if err != nil {
		t.Fatal(err)
	}

	//Make sure the containerID matches
	if cID != containerID {
		t.Fatal("CreateContainer: Expected ContainerID: " + containerID + " | Received: " + cID)
	}

	//Make sure the proper config settings were set when running the container
	if dm.CCConfig.Image != img {
		t.Fatal("CreateContainer: Expected Container Config Image: " + img + " | Received: " + cID)
	}

	if dm.CCConfig.Cmd[0] != cmd[0] {
		t.Fatalf("CreateContainer: Expected Container Config Cmd: %v | Received: %v", cmd, dm.CCConfig.Cmd)
	}
}

//Test CreateContainer Function with an error
func TestCreateContainerError(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the image that should be used
	img := "image"

	//Set what the containerID should be
	containerID := "testcontainer"

	//Set the error at and error message
	dm.ErrorAt = "ContainerCreate"
	dm.ErrorMsg = "Testing error at ContainerCreate()"

	dm.CCRet = container.ContainerCreateCreatedBody{
		ID: containerID,
	}

	//Create the util
	util := NewUtility()

	//Test creating the container
	_, err := util.CreateContainer(img, dm)

	//Error should occur
	if err == nil {
		t.Fatal("CreateContainer: Expected to receive an error, but did not receive one.")
	}

	if err != nil {
		if err.Error() != dm.ErrorMsg {
			t.Fatal("CreateContainer: Expected Error: " + dm.ErrorMsg + " | Received Error: " + err.Error())
		}
	}
}

//Test CopyFromContainer Function
func TestCopyFromContainer(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the containerID to be used
	cID := "fake"

	//Set what the source should be
	source := "/fake/source"

	//Set what the destination should be
	dest := "/fake/dest"

	//Create the util
	util := NewUtility()

	//Create the mock copy tool
	mcp := &MockCopyTool{}

	//Test creating the container
	err := util.CopyFromContainer(source, dest, cID, dm, mcp)

	//If error occurs the test fails
	if err != nil {
		t.Fatal(err)
	}

	//Make sure the containerID was passed in successfully
	if dm.CFCID != cID {
		t.Fatalf("CopyFromContainer: Expected ContainerID: %s | Received: %s", cID, dm.CFCID)
	}

	//Make sure the source was passed in correctly
	if dm.CFCSource != source {
		t.Fatalf("CopyFromContainer: Expected Source: %s | Received: %s", source, dm.CFCSource)
	}

	//Make sure the destination gets passed to the copy tool correctly
	if mcp.Dest != dest {
		t.Fatalf("CopyFromContainer -> CopyFiles: Expected Dest: %s | Received: %s", dest, mcp.Dest)
	}

}

//Test CopyFromContainer Function with an error
func TestCopyFromContainerError(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the containerID to be used
	cID := "fake"

	//Set what the source should be
	source := "/fake/source"

	//Set what the destination should be
	dest := "/fake/destination"

	//Set the error at and error message
	dm.ErrorAt = "CopyFromContainer"
	dm.ErrorMsg = "Testing error at CopyFromContainer()"

	//Create the util
	util := NewUtility()

	//Create the mock copy tool
	mcp := &MockCopyTool{}

	//Test creating the container
	err := util.CopyFromContainer(source, dest, cID, dm, mcp)

	//If error occurs the test fails
	if err == nil {
		t.Fatal("CopyFromContainer: Expected to receive an error, but did not receive one.")
	}

	if err != nil {
		if err.Error() != dm.ErrorMsg {
			t.Fatal("CopyFromContainer: Expected Error: " + dm.ErrorMsg + " | Received Error: " + err.Error())
		}
	}
}

//Test CopyFromContainer Function with an error in the CopyFiles Function
func TestCopyFromContainerErrorCopyFiles(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the containerID to be used
	cID := "fake"

	//Set what the source should be
	source := "/fake/source"

	//Set what the destination should be
	dest := "/fake/destination"

	dm.ErrorMsg = "Testing error at CopyFromContainer()"

	//Create the util
	util := NewUtility()

	//Create the mock copy tool with error sets
	mcp := &MockCopyTool{
		Error:    true,
		ErrorMsg: "Testing error at CopyFiles()",
	}

	//Test creating the container
	err := util.CopyFromContainer(source, dest, cID, dm, mcp)

	//If error occurs the test fails
	if err == nil {
		t.Fatal("CopyFromContainer: Expected to receive an error, but did not receive one.")
	}

	if err != nil {
		if err.Error() != mcp.ErrorMsg {
			t.Fatal("CopyFromContainer: Expected Error: " + mcp.ErrorMsg + " | Received Error: " + err.Error())
		}
	}
}

//Test RemoveContainer Function
func TestRemoveContainer(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the containerID to be used
	cID := "fake"

	//Create the util
	util := NewUtility()

	//Test creating the container
	err := util.RemoveContainer(cID, dm)

	//If error occurs the test fails
	if err != nil {
		t.Fatal(err)
	}

	//Make sure the RemoveContainer options are correct
	if !dm.CROptions.Force {
		t.Fatal("RemoveContainer: Expected the ContainerRemoveOptions Force field to be set to true but it is set to false")
	}

	//Make sure the containerID is correct
	if dm.CRContainer != cID {
		t.Fatalf("RemoveContainer: Expected Container ID: %s | Received: %s", cID, dm.CRContainer)
	}
}

//Test RemoveContainer Function with an error
func TestRemoveContainerError(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the containerID to be used
	cID := "fake"

	//Set the error at and error message
	dm.ErrorAt = "ContainerRemove"
	dm.ErrorMsg = "Testing error at ContainerRemove()"

	//Create the util
	util := NewUtility()

	//Test creating the container
	err := util.RemoveContainer(cID, dm)

	//Error should occur
	if err == nil {
		t.Fatal("RemoveContainer: Expected to receive an error, but did not receive one.")
	}

	if err != nil {
		if err.Error() != dm.ErrorMsg {
			t.Fatal("RemoveContainer: Expected Error: " + dm.ErrorMsg + " | Received Error: " + err.Error())
		}
	}
}

//Test RunContainer Function without arguments
func TestRunContainerNoArgs(t *testing.T) {
	//Ser the image to be run
	image := "image"

	//Set the ports
	ports := []string{"3000:3000"}

	//Get absolute path for beginning of volume
	absPath, err := filepath.Abs("/a/path/")

	//Should not be an error here
	if err != nil {
		t.Fatal(err)
	}

	//Set the volumes
	volumes := []string{absPath + ":/another/path"}

	//Set the container name
	cName := "test"

	//Set the empty args
	args := []string{}

	//Set the expected command
	exCmd := "docker run -it --rm --name " + cName + " -p " + ports[0] + " -v " + volumes[0] + " " + image + " "

	//Create the util tool
	util := NewUtility()

	//Run the RunContainer function and ignore any errors since we just want to make sure the cmd is built properly
	cmd, _ := util.RunContainer(image, ports, volumes, cName, args)

	//Returned cmd should equal the expected one
	if cmd != exCmd {
		t.Fatalf("RunContainer: Expected CMD: %s | Received CMD: %s", exCmd, cmd)
	}

}

//Test RunContainer Function with arguments
func TestRunContainerWithArgs(t *testing.T) {
	//Ser the image to be run
	image := "image"

	//Set the ports
	ports := []string{"3000:3000"}

	//Get absolute path for beginning of volume
	absPath, err := filepath.Abs("/a/path/")

	//Should not be an error here
	if err != nil {
		t.Fatal(err)
	}

	//Set the volumes
	volumes := []string{absPath + ":/another/path"}

	//Set the container name
	cName := "test"

	//Set the empty args
	args := []string{"some", "arguments"}

	argStr := strings.Join(args, " ")

	//Set the expected command
	exCmd := "docker run -it --rm --name " + cName + " -p " + ports[0] + " -v " + volumes[0] + " " + image + " " + argStr

	//Create the util tool
	util := NewUtility()

	//Run the RunContainer function and ignore any errors since we just want to make sure the cmd is built properly
	cmd, _ := util.RunContainer(image, ports, volumes, cName, args)

	//Returned cmd should equal the expected one
	if cmd != exCmd {
		t.Fatalf("RunContainer: Expected CMD: %s | Received CMD: %s", exCmd, cmd)
	}

}

func TestRunContainer(t *testing.T) {
	t.Run("ShouldReturnError_WhenSplitVolumeIs1", func(t *testing.T) {
		//Set the image to be run
		image := "image"

		//Set the ports
		ports := []string{"3000:3000"}

		//Set the volumes
		volumes := []string{"/path1"}

		//Set the container name
		cName := "test"

		//Set the empty args
		args := []string{}

		//Create the util tool
		util := NewUtility()

		//Run the RunContainer function and assert the error
		exErr := errors.New("utils: Invalid split volume 1")
		_, err := util.RunContainer(image, ports, volumes, cName, args)
		if err.Error() != exErr.Error() {
			t.Fatalf("RunContainer: Expected err: %s | Received err: %s", exErr.Error(), err.Error())
		}
	})

	t.Run("ShouldReturnCorrectCmdStr_WhenSplitVolumeIs2", func(t *testing.T) {
		//Set the image to be run
		image := "image"

		//Set the ports
		ports := []string{"3000:3000"}

		//Set the volumes
		volumes := []string{"/path1:/path2"}

		//Set the container name
		cName := "test"

		//Set the empty args
		args := []string{}

		//Set the expected command
		exCmd := "docker run -it --rm --name " + cName + " -p " + ports[0] + " -v " + volumes[0] + " " + image + " "

		//Create the util tool
		util := NewUtility()

		//Run the RunContainer function and ignore any errors since we just want to make sure the cmd is built properly
		cmd, _ := util.RunContainer(image, ports, volumes, cName, args)
		if cmd != exCmd {
			t.Fatalf("RunContainer: Expected CMD: %s | Received CMD: %s", exCmd, cmd)
		}
	})

	t.Run("ShouldReturnCorrectCmdStr_WhenSplitVolumeIs3", func(t *testing.T) {
		//Set the image to be run
		image := "image"

		//Set the ports
		ports := []string{"3000:3000"}

		//Set the volumes
		volumes := []string{"/path1:/path2:/path3"}

		//Set the container name
		cName := "test"

		//Set the empty args
		args := []string{}

		//Set the expected command
		exCmd := "docker run -it --rm --name " + cName + " -p " + ports[0] + " -v " + volumes[0] + " " + image + " "

		//Create the util tool
		util := NewUtility()

		//Run the RunContainer function and ignore any errors since we just want to make sure the cmd is built properly
		cmd, _ := util.RunContainer(image, ports, volumes, cName, args)
		if cmd != exCmd {
			t.Fatalf("RunContainer: Expected CMD: %s | Received CMD: %s", exCmd, cmd)
		}
	})

	t.Run("ShouldReturnError_WhenSplitVolumeIs4", func(t *testing.T) {
		//Set the image to be run
		image := "image"

		//Set the ports
		ports := []string{"3000:3000"}

		//Set the volumes
		volumes := []string{"/path1:/path2:/path3:/path4"}

		//Set the container name
		cName := "test"

		//Set the empty args
		args := []string{}

		//Create the util tool
		util := NewUtility()

		//Run the RunContainer function and assert the error
		exErr := errors.New("utils: Invalid split volume 4")
		_, err := util.RunContainer(image, ports, volumes, cName, args)
		if err.Error() != exErr.Error() {
			t.Fatalf("RunContainer: Expected err: %s | Received err: %s", exErr.Error(), err.Error())
		}
	})
}

//Test RemoveImage Function
func TestRemoveImage(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the image to be used
	img := "image"

	//Set the image ID
	imgID := "fakeImgID"

	dm.ILRet = []types.ImageSummary{
		{
			ID:       imgID,
			RepoTags: []string{img + ":faketag"},
		},
	}

	//Create the util
	util := NewUtility()

	//Test creating the container
	err := util.RemoveImage(img, dm)

	//Error shouldn't occur
	if err != nil {
		t.Fatal(err)
	}

	//Check and make sure the image passed in is correct
	if dm.IRImgID != img {
		t.Fatalf("RemoveImage: Expected ImageID: %s | Received: %s", img, dm.IRImgID)
	}

	//Check and make sure the image removal options passed in is correct
	if !dm.IROptions.Force {
		t.Fatal("RemoveImage: Expected the ImageRemovalOptions Force field to be set to 'true' but it is not")
	}
}

//Test RemoveImage Function when it encounters an error
func TestRemoveImageError(t *testing.T) {
	//Create the Mock Docker Client
	dm := NewDockMock()

	//Set the image to be used
	img := "image"

	//Set the error at and error message
	dm.ErrorAt = "ImageRemove"
	dm.ErrorMsg = "Testing error at ImageRemove()"

	//Create the util
	util := NewUtility()

	//Test creating the container
	err := util.RemoveImage(img, dm)

	//Error should occur
	if err == nil {
		t.Fatal("RemoveContainer: Expected to receive an error, but did not receive one.")
	}

	if err != nil {
		if err.Error() != dm.ErrorMsg {
			t.Fatal("RemoveContainer: Expected Error: " + dm.ErrorMsg + " | Received Error: " + err.Error())
		}
	}
}

//Integration Tests
//------------------------------------------------------------------------

//Integration test for the pull image function
func TestDocker_Integration(t *testing.T) {
	//if we want to run short tests then skip this
	if testing.Short() {
		t.Skip("skipping test, short tests specified")
	}

	//Create the util object
	util := NewUtility()

	//Create the Docker CLI
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	//Shouldn't be any errors
	if err != nil {
		t.Fatal(err)
	}

	//set the image to pull
	img := "bpalmer/alpine-base-ssh"

	//pull the image
	err = util.PullImage(img, cli)

	//shouldn't have any errors
	if err != nil {
		t.Fatal(err)
	}

	//make sure the image exists after being pulled
	imgExist, err := util.ImageExists(img+":latest", cli)

	if err != nil {
		t.Fatal(err)
	}

	if !imgExist {
		t.Fatalf("Docker Integration: Expected Image %s to exist but it did not", img)
	}

	//Try creating the container
	cont, err := util.CreateContainer(img, cli)

	if err != nil {
		t.Fatal(err)
	}

	//Make sure the container exists
	contExist := false
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		t.Fatal(err)
	}

	for _, c := range containers {
		if c.ID == cont {
			contExist = true
		}
	}

	if !contExist {
		t.Fatalf("Docker Integration: Expected Container %s to exist but it did not", cont)
	}

	//Test getting a file from the container
	file := "hostname"
	source := "/etc/" + file

	//get the executing directory

	//set the destination
	dest, err := filepath.Abs("../testing/")

	if err != nil {
		t.Fatal(err)
	}

	// runtime.Breakpoint()

	err = util.CopyFromContainer(source, dest, cont, cli, &CopyTool{})

	if err != nil {
		t.Fatal(err)
	}

	//See if the file exists in the destination directory
	if _, err = os.Stat(dest + "/" + file); err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("Docker Integration: File %s was not copied from the container to host destination %s", file, dest)
		}
	}

	//Remove the file now
	err = os.Remove(dest + "/" + file)

	if err != nil {
		t.Fatal(err)
	}

	//Test removing the container
	err = util.RemoveContainer(cont, cli)

	if err != nil {
		t.Fatal(err)
	}

	//Make sure the container no longer exists

	contExist = false
	containers, err = cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})

	if err != nil {
		t.Fatal(err)
	}

	for _, c := range containers {
		if c.ID == cont {
			contExist = true
		}
	}

	if contExist {
		t.Fatalf("Docker Integration: Expected Container %s to not exist but it does exist", cont)
	}

	//Now test removing the image
	err = util.RemoveImage(img, cli)

	if err != nil {
		t.Fatal(err)
	}

	//Check if it exists again
	imgExist, err = util.ImageExists(img, cli)

	if imgExist {
		t.Fatalf("Docker: Expected Image %s to not exist but it did", img)
	}
}

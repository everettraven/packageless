package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

//PullImage - This function pulls a Docker Image from the packageless organization in Docker Hub
func (u *Utility) PullImage(name string, cli Client) error {
	//Set the context
	ctx := context.Background()

	//Begin pulling the image
	out, err := cli.ImagePull(ctx, name, types.ImagePullOptions{})

	//Check for errors
	if err != nil {
		return err
	}

	//Close the output buffer after the function exits
	defer out.Close()

	//Copy the output to the screen
	io.Copy(os.Stdout, out)

	//No errors
	return nil
}

//ImageExists - Function to check and see if Docker has the image downloaded
func (u *Utility) ImageExists(imageID string, cli Client) (bool, error) {
	//Create a context and get a list of images on the system
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	//Check for errors
	if err != nil {
		return false, err
	}

	//Loop through all the images and check if a match is found
	for _, image := range images {

		// If RepoTags returned isnt populated then skip to the next image
		if len(image.RepoTags) < 1 {
			continue
		}

		if image.RepoTags[0] == imageID {
			return true, nil
		}
	}

	//No match found
	return false, nil
}

//CreateContainer - Create a Docker Container from a Docker Image. Returns the containerID and any errors
func (u *Utility) CreateContainer(image string, cli Client) (string, error) {
	//Create the context and create the container
	ctx := context.Background()
	container, err := cli.ContainerCreate(ctx, &container.Config{Image: image, Cmd: []string{"bash"}}, nil, nil, nil, "")

	//Check for errors
	if err != nil {
		return "", err
	}

	//No errors
	return container.ID, err
}

//CopyFromContainer will copy files from within a Docker Container to the source location on the host
func (u *Utility) CopyFromContainer(source string, dest string, containerID string, cli Client, cp Copier) error {
	//Set the context and begin copying from the container
	ctx := context.Background()
	reader, _, err := cli.CopyFromContainer(ctx, containerID, source)

	//Check for errors
	if err != nil {
		return err
	}

	//Close the reader after the function ends
	defer reader.Close()

	//Copy the files over
	err = cp.CopyFiles(reader, dest, source)

	if err != nil {
		return err
	}

	return nil
}

//RemoveContainer is used to remove a container Docker given the container ID
func (u *Utility) RemoveContainer(containerID string, cli Client) error {

	//Create the context and remove the container
	ctx := context.Background()
	err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})

	//Check for errors
	if err != nil {
		return err
	}

	//No Errors
	return nil
}

//RunContainer - Runs a container for the specified package
func (u *Utility) RunContainer(image string, ports []string, volumes []string, containerName string, args []string) (string, error) {
	// Build the command to run the docker container
	var cmdStr string
	var cmd *exec.Cmd

	cmdStr += "docker "

	// add the base docker command details
	cmdStr += "run -it --rm --name " + containerName + " "

	// add the ports to the command
	for _, port := range ports {
		cmdStr += "-p " + port + " "
	}

	// add the volumes to the command
	for _, vol := range volumes {
		if vErr := u.validateRunContainerVolume(vol); vErr != nil {
			return "", vErr
		}

		splitVol := strings.Split(vol, ":")
		var source string
		var target string

		if len(splitVol) == 3 {
			source = strings.Join(splitVol[:2], ":")
			target = splitVol[2]
		} else {
			source = splitVol[0]
			target = splitVol[1]
		}

		source, err := filepath.Abs(source)

		if err != nil {
			return "", err
		}

		cmdStr += "-v " + source + ":" + target + " "
	}

	// add the image name and the arguments
	cmdStr += image + " "

	//Combine the arguments into one string
	argStr := strings.Join(args, " ")

	//add the arguments
	cmdStr += argStr

	//Instantiate the command based on OS
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", cmdStr)
	} else {
		cmd = exec.Command("bash", "-c", cmdStr)
	}

	//Connect the command stderr, stdout, and stdin to the OS stderr, stdout, stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	//Run the command
	err := cmd.Run()

	//Check for errors
	if err != nil {
		return cmdStr, err
	}

	return cmdStr, nil

}

func (u *Utility) validateRunContainerVolume(volume string) error {
	splitVolume := strings.Split(volume, ":")

	if len(splitVolume) != 2 && len(splitVolume) != 3 {
		return fmt.Errorf("utils: Invalid split volume %d", len(splitVolume))
	}

	return nil
}

//RemoveImage removes the image with the given name from local Docker
func (u *Utility) RemoveImage(image string, cli Client) error {
	//Create the context and search for the image in the list of images
	ctx := context.Background()

	//Remove the image
	_, err := cli.ImageRemove(ctx, image, types.ImageRemoveOptions{Force: true})

	//Check for errors
	if err != nil {
		return err
	}

	return nil
}

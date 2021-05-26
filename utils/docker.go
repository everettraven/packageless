package utils

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

//PullImage - This function pulls a Docker Image from the packageless organization in Docker Hub
func PullImage(name string) error {
	//Set up a Docker API client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	//Check for errors
	if err != nil {
		return err
	}

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
func ImageExists(imageID string) (bool, error) {
	//Create a client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	//Check for errors
	if err != nil {
		return false, err
	}

	//Create a context and get a list of images on the system
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	//Check for errors
	if err != nil {
		return false, err
	}

	//Loop through all the images and check if a match is found
	for _, image := range images {
		if strings.Split(image.RepoTags[0], ":")[0] == imageID {
			return true, nil
		}
	}

	//No match found
	return false, nil
}

//CreateContainer - Create a Docker Container from a Docker Image. Returns the containerID and any errors
func CreateContainer(image string) (string, error) {
	//Create the client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	//Check for errors
	if err != nil {
		return "", err
	}

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
func CopyFromContainer(source string, dest string, containerID string) error {
	//Create the Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	//Check for errors
	if err != nil {
		return err
	}

	//Set the context and begin copying from the container
	ctx := context.Background()
	reader, _, err := cli.CopyFromContainer(ctx, containerID, source)

	//Check for errors
	if err != nil {
		return err
	}

	//Close the reader after the function ends
	defer reader.Close()

	//Create a tar Reader
	tarReader := tar.NewReader(reader)

	//Skip the first header as it is the source folder name
	tarReader.Next()

	//Loop through the reader and write the files
	for {
		//Get the tar header
		header, err := tarReader.Next()
		//Make sure we havent reached the end of the tar
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		newHeaderPath := strings.Split(header.Name, "/")[1:]
		joinPath := strings.Join(newHeaderPath[:], "/")

		//Create the destination file path on the host
		path := filepath.Join(dest, joinPath)
		//Get the file info from the header
		info := header.FileInfo()

		//Check if the current file is a directory
		if info.IsDir() {

			//Check if the directory exists
			if _, err = os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					//Make the directory
					err = os.MkdirAll(path, info.Mode())
				} else {
					return err
				}
			}

		} else {
			//Create the file and open it in the destination path on the host
			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

			//Check for errors
			if err != nil {
				return err
			}

			//Copy the contents of the tar reader to the file
			_, err = io.Copy(file, tarReader)

			//Check for errors
			if err != nil {
				return err
			}

			//Close the file when all the writing is finished
			file.Close()
		}

	}

	return nil
}

//RemoveContainer is used to remove a container Docker given the container ID
func RemoveContainer(containerID string) error {
	//Create the Docker API client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	//Check for errors
	if err != nil {
		return err
	}

	//Create the context and remove the container
	ctx := context.Background()
	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})

	//Check for errors
	if err != nil {
		return err
	}

	//No Errors
	return nil
}

//RunContainer - Runs a container for the specified package
func RunContainer(image string, ports []string, volumes []string, containerName string, args []string) error {
	// Build the command to run the docker container
	var cmdStr string
	cmdBase := "docker"

	//If it is running on windows run the command through powershell
	if runtime.GOOS == "windows" {
		cmdBase = "powershell"
		cmdStr += "docker "
	}

	// add the base docker command details
	cmdStr += "run -it --rm --name " + containerName + " "

	// add the ports to the command
	for _, port := range ports {
		cmdStr += "-p " + port + " "
	}

	// add the volumes to the command
	for _, vol := range volumes {
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
			return err
		}

		cmdStr += "-v " + source + ":" + target + " "
	}

	// add the image name and the arguments
	cmdStr += image + " "

	//Combine the arguments into one string
	argStr := strings.Join(args, " ")

	//add the arguments
	cmdStr += argStr

	//Instantiate the command
	cmd := exec.Command(cmdBase, cmdStr)

	//Connect the command stderr, stdout, and stdin to the OS stderr, stdout, stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	//Run the command
	err := cmd.Run()

	//Check for errors
	if err != nil {
		return err
	}

	return nil

}

//RemoveImage removes the image with the given name from local Docker
func RemoveImage(image string) error {
	//Create the client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	//Check for errors
	if err != nil {
		return err
	}

	//Create the context and search for the image in the list of images
	ctx := context.Background()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	//Check for errors
	if err != nil {
		return err
	}

	//Create a variable to hold the ID of the image we want to remove
	var imageID string

	//Loop through all the images and check if a match is found
	for _, img := range images {
		if strings.Split(img.RepoTags[0], ":")[0] == image {
			imageID = img.ID
		}
	}

	//Remove the image
	_, err = cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{Force: true})

	//Check for errors
	if err != nil {
		return err
	}

	return nil
}

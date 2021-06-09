package utils

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/hashicorp/hcl2/hcl"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

//Client interface so that we can create a mock of the docker SDK interactions in our unit tests
type Client interface {
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
	CopyFromContainer(ctx context.Context, containerID string, srcPath string) (io.ReadCloser, types.ContainerPathStat, error)
	ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error
	ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error)
}

//Tools interface so that we can create a mock of our utility functions in our unit tests
type Tools interface {
	MakeDir(path string) error
	OpenFile(path string) (*os.File, error)
	RemoveDir(path string) error
	UpgradeDir(path string) error
	ParseBody(body hcl.Body, out interface{}) (interface{}, error)
	GetHCLBody(filepath string) (hcl.Body, error)
	PullImage(name string, cli Client) error
	ImageExists(imageID string, cli Client) (bool, error)
	CreateContainer(image string, cli Client) (string, error)
	CopyFromContainer(source string, dest string, containerID string, cli Client) error
	RemoveContainer(containerID string, cli Client) error
	RunContainer(image string, ports []string, volumes []string, containerName string, args []string) (string, error)
	RemoveImage(image string, cli Client) error
	CopyFiles(reader io.ReadCloser, dest string) error
	AddAliasWin(name string, ed string) error
	RemoveAliasWin(name string, ed string) error
	AddAliasUnix(name string, ed string) error
	RemoveAliasUnix(name string, ed string) error
}

//Utility Tool struct with its functions
type Utility struct{}

func NewUtility() *Utility {
	util := &Utility{}
	return util
}

//MakeDir makes a directory if it does not exist
func (u *Utility) MakeDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0765)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

//OpenFile opens the specified file, creating it if it does not exist
func (u *Utility) OpenFile(path string) (*os.File, error) {
	var file *os.File
	//Check if the path exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			//Create the file
			file, err = os.Create(path)

			if err != nil {
				return nil, err
			}
		}
	} else {
		//Open the file
		file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0755)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

//RemoveDir removes the specified directory
func (u *Utility) RemoveDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	} else {
		err = os.RemoveAll(path)
	}

	return nil
}

//UpgradeDir resets the directory by removing it if it exists and then recreating it
func (u *Utility) UpgradeDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0765)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		//Remove the directory if it already exists
		err = os.RemoveAll(path)

		if err != nil {
			return err
		}

		//Recreate the directory
		err = os.MkdirAll(path, 0765)

		if err != nil {
			return err
		}
	}
	return nil
}

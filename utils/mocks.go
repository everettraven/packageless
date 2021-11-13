package utils

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/hashicorp/hcl2/hcl"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

//Mock utility and its functions
type MockUtility struct {
	//Keeps track of the function calls that are made during a test
	Calls []string

	//Package object that can be changed for different tests
	Pim PimHCLUtil

	//Config Object that can be changed for different tests
	Conf Config

	//HCL Body that can be changed for different tests
	HCLBody hcl.Body

	//Set if the image should exist or not for testing
	ImgExist bool

	//Set the function name that should throw an error when called
	ErrorAt string

	//Set an error message for the tests
	ErrorMsg string

	//Keep track of the directories that are created
	MadeDirs []string

	//Keep track of the files that are opened/created
	OpenedFiles []string

	//Keep track of the directories that are removed
	RemovedDirs []string

	//Keep track of the upgraded directories
	UpgradedDirs []string

	//Keep track of the HCLFiles read
	HCLFiles []string

	//Keep track of the images that are pulled
	PulledImgs []string

	//Use a fake containerID to make sure it is being used correctly
	ContainerID string

	//CreateContainer data
	CreateImages []string

	//Keep track of the CopyFromContainer data
	CopySources     []string
	CopyDests       []string
	CopyContainerID string

	//Keep track of the RemoveContainer data
	RemoveContainerID string

	//Keep track of the RunContainer data
	RunImage         string
	RunPorts         []string
	RunVolumes       []string
	RunContainerName string
	RunArgs          []string

	//Keep track of the RemoveImage data
	RemovedImgs []string

	//Keep track of the alias data
	CmdToAlias []string

	//Should the Pim Configuration file exist
	PimConfigShouldExist bool

	//Pim Config Directory passed in
	PimConfigDir string

	//List of pim names to return
	InstalledPims []string

	//List of pims fetched using FetchPimConfigs
	FetchedPims []string
}

//Create a new Mock Utility and set any default variables
func NewMockUtility() *MockUtility {
	mu := &MockUtility{
		Pim: PimHCLUtil{
			Pims: []PackageImage{
				{
					Name:    "python",
					BaseDir: "/base",
					Versions: []Version{
						{
							Version: "latest",
							Image:   "packageless/python",

							Volumes: []Volume{
								{
									Path:  "a/path",
									Mount: "/another/one",
								},
							},
							Copies: []*Copy{
								{
									Source: "/source/path",
									Dest:   "destination",
								},
							},
							Port: "3000",
						},
					},
				},
			},
		},
		Conf: Config{
			BaseDir:        "~/.packageless/",
			StartPort:      3000,
			PortInc:        1,
			Alias:          true,
			RepositoryHost: "https://raw.githubusercontent.com/everettraven/packageless-pims/main/pims/",
			PimsConfigDir:  "pims_config/",
			PimsDir:        "pims/",
		},
		HCLBody:              hcl.EmptyBody(),
		ErrorMsg:             "Testing for error handling",
		ContainerID:          "FakeContainer123",
		PimConfigShouldExist: true,
	}

	return mu
}

//Mock of the MakeDir Utility function
func (mu *MockUtility) MakeDir(path string) error {
	mu.Calls = append(mu.Calls, "MakeDir")
	mu.MadeDirs = append(mu.MadeDirs, path)

	if mu.ErrorAt == "MakeDir" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the OpenFile Utility function
func (mu *MockUtility) OpenFile(path string) (*os.File, error) {
	mu.Calls = append(mu.Calls, "OpenFile")
	mu.OpenedFiles = append(mu.OpenedFiles, path)

	if mu.ErrorAt == "OpenFile" {
		return nil, errors.New(mu.ErrorMsg)
	}

	return &os.File{}, nil
}

//Mock of the RemoveDir Utility function
func (mu *MockUtility) RemoveDir(path string) error {
	mu.Calls = append(mu.Calls, "RemoveDir")
	mu.RemovedDirs = append(mu.RemovedDirs, path)

	if mu.ErrorAt == "RemoveDir" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the UpgradeDir Utility function
func (mu *MockUtility) UpgradeDir(path string) error {
	mu.Calls = append(mu.Calls, "UpgradeDir")
	mu.UpgradedDirs = append(mu.UpgradedDirs, path)

	if mu.ErrorAt == "UpgradeDir" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the ParseBody Utility function
func (mu *MockUtility) ParseBody(body hcl.Body, out interface{}) (interface{}, error) {
	mu.Calls = append(mu.Calls, "ParseBody")

	if mu.ErrorAt == "ParseBody" {
		return out, errors.New(mu.ErrorMsg)
	}

	switch out.(type) {
	default:
		return nil, errors.New("Unexpected type in parse")
	case PimHCLUtil:
		return mu.Pim, nil
	case Config:
		return mu.Conf, nil
	}
}

//Mock of the GetHCLBody Utility function
func (mu *MockUtility) GetHCLBody(filepath string) (hcl.Body, error) {
	mu.Calls = append(mu.Calls, "GetHCLBody")
	mu.HCLFiles = append(mu.HCLFiles, filepath)

	if mu.ErrorAt == "GetHCLBody" {
		return mu.HCLBody, errors.New(mu.ErrorMsg)
	}

	return mu.HCLBody, nil
}

//Mock of the PullImage Utility function
func (mu *MockUtility) PullImage(name string, cli Client) error {
	mu.Calls = append(mu.Calls, "PullImage")
	mu.PulledImgs = append(mu.PulledImgs, name)

	if mu.ErrorAt == "PullImage" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the ImageExists Utility function
func (mu *MockUtility) ImageExists(imageID string, cli Client) (bool, error) {
	mu.Calls = append(mu.Calls, "ImageExists")

	if mu.ErrorAt == "ImageExists" {
		return mu.ImgExist, errors.New(mu.ErrorMsg)
	}

	return mu.ImgExist, nil
}

//Mock of the CreateContainer Utility function
func (mu *MockUtility) CreateContainer(image string, cli Client) (string, error) {
	mu.Calls = append(mu.Calls, "CreateContainer")
	mu.CreateImages = append(mu.CreateImages, image)

	if mu.ErrorAt == "CreateContainer" {
		return "", errors.New(mu.ErrorMsg)
	}

	return mu.ContainerID, nil
}

//Mock of the CopyFromContainer Utility function
func (mu *MockUtility) CopyFromContainer(source string, dest string, containerID string, cli Client, cp Copier) error {
	mu.Calls = append(mu.Calls, "CopyFromContainer")
	mu.CopySources = append(mu.CopySources, source)
	mu.CopyDests = append(mu.CopyDests, dest)
	mu.CopyContainerID = containerID

	if mu.ErrorAt == "CopyFromContainer" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the RemoveContainer Utility function
func (mu *MockUtility) RemoveContainer(containerID string, cli Client) error {
	mu.Calls = append(mu.Calls, "RemoveContainer")
	mu.RemoveContainerID = containerID

	if mu.ErrorAt == "RemoveContainer" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the RunContainer Utility function
func (mu *MockUtility) RunContainer(image string, ports []string, volumes []string, containerName string, args []string) (string, error) {
	mu.Calls = append(mu.Calls, "RunContainer")
	mu.RunImage = image
	mu.RunPorts = ports
	mu.RunVolumes = volumes
	mu.RunContainerName = containerName
	mu.RunArgs = args

	if mu.ErrorAt == "RunContainer" {
		return "", errors.New(mu.ErrorMsg)
	}

	return "", nil
}

//Mock of the RemoveImage Utility function
func (mu *MockUtility) RemoveImage(image string, cli Client) error {
	mu.Calls = append(mu.Calls, "RemoveImage")
	mu.RemovedImgs = append(mu.RemovedImgs, image)

	if mu.ErrorAt == "RemoveImage" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the AddAliasWin Utility function
func (mu *MockUtility) AddAliasWin(name string, ed string) error {
	mu.Calls = append(mu.Calls, "AddAlias")
	mu.CmdToAlias = append(mu.CmdToAlias, name)

	if mu.ErrorAt == "AddAlias" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the RemoveAliasWin Utility function
func (mu *MockUtility) RemoveAliasWin(name string, ed string) error {
	mu.Calls = append(mu.Calls, "RemoveAlias")
	mu.CmdToAlias = append(mu.CmdToAlias, name)

	if mu.ErrorAt == "RemoveAlias" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the AddAliasUnix Utility function
func (mu *MockUtility) AddAliasUnix(name string, ed string) error {
	mu.Calls = append(mu.Calls, "AddAlias")
	mu.CmdToAlias = append(mu.CmdToAlias, name)

	if mu.ErrorAt == "AddAlias" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

//Mock of the RemoveAliasUnix Utility function
func (mu *MockUtility) RemoveAliasUnix(name string, ed string) error {
	mu.Calls = append(mu.Calls, "RemoveAlias")
	mu.CmdToAlias = append(mu.CmdToAlias, name)

	if mu.ErrorAt == "RemoveAlias" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

func (mu *MockUtility) FetchPimConfig(baseUrl string, pimName string, savePath string) error {
	mu.Calls = append(mu.Calls, "FetchPimConfig")

	if mu.ErrorAt == "FetchPimConfig" {
		return errors.New(mu.ErrorMsg)
	}

	mu.FetchedPims = append(mu.FetchedPims, pimName)

	return nil
}

func (mu *MockUtility) FileExists(path string) bool {
	mu.Calls = append(mu.Calls, "FileExists")

	return mu.PimConfigShouldExist
}

func (mu *MockUtility) RemoveFile(path string) error {
	mu.Calls = append(mu.Calls, "RemoveFile")

	if mu.ErrorAt == "RemoveFile" {
		return errors.New(mu.ErrorMsg)
	}

	return nil
}

func (mu *MockUtility) GetListOfInstalledPimConfigs(pimConfigDir string) ([]string, error) {
	mu.Calls = append(mu.Calls, "GetListOfInstalledPimConfigs")

	if mu.ErrorAt == "GetListOfInstalledPimConfigs" {
		return nil, errors.New(mu.ErrorMsg)
	}

	mu.PimConfigDir = pimConfigDir

	return mu.InstalledPims, nil
}

//Create a Mock for the Docker client
type DockMock struct {
	//Variable to know what function to return an error from
	ErrorAt string

	//Variable to store the error message
	ErrorMsg string

	//Keep track of the values from ImagePull Function
	IPRefStr string

	//Keep track of the values from the ContainerCreate Function
	CCConfig *container.Config
	CCName   string
	CCRet    container.ContainerCreateCreatedBody

	//Keep track of the values from the CopyFromContainer Function
	CFCID     string
	CFCSource string

	//Keep track of the values from the ContainerRemove Function
	CRContainer string
	CROptions   types.ContainerRemoveOptions

	//Keep track of the values from the ImageRemove Function
	IRImgID   string
	IROptions types.ImageRemoveOptions

	//ImageList return value
	ILRet []types.ImageSummary
}

//Function to create a new DockMock
func NewDockMock() *DockMock {
	dm := &DockMock{}

	return dm
}

//Mock function of the Docker SDK ImagePull function
func (dm *DockMock) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	if dm.ErrorAt == "ImagePull" {
		return nil, errors.New(dm.ErrorMsg)
	}

	dm.IPRefStr = refStr

	//Create the ReadCloser
	rc := io.NopCloser(bytes.NewReader([]byte("ImagePull")))

	return rc, nil
}

//Mock function of the Docker SDK ImagePull function
func (dm *DockMock) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	if dm.ErrorAt == "ImageList" {
		return nil, errors.New(dm.ErrorMsg)
	}

	return dm.ILRet, nil
}

//Mock function of the Docker SDK ContainerCreate function
func (dm *DockMock) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	if dm.ErrorAt == "ContainerCreate" {
		return container.ContainerCreateCreatedBody{}, errors.New(dm.ErrorMsg)
	}

	dm.CCConfig = config
	dm.CCName = containerName
	return dm.CCRet, nil
}

//Mock function of the Docker SDK CopyFromContainer function
func (dm *DockMock) CopyFromContainer(ctx context.Context, containerID string, srcPath string) (io.ReadCloser, types.ContainerPathStat, error) {
	if dm.ErrorAt == "CopyFromContainer" {
		return nil, types.ContainerPathStat{}, errors.New(dm.ErrorMsg)
	}

	dm.CFCID = containerID
	dm.CFCSource = srcPath

	//Create the ReadCloser
	rc := io.NopCloser(bytes.NewReader([]byte("")))

	return rc, types.ContainerPathStat{}, nil
}

//Mock function of the Docker SDK ContainerRemove function
func (dm *DockMock) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	if dm.ErrorAt == "ContainerRemove" {
		return errors.New(dm.ErrorMsg)
	}

	dm.CRContainer = container
	dm.CROptions = options
	return nil
}

//Mock function of the Docker SDK ImageRemove function
func (dm *DockMock) ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	if dm.ErrorAt == "ImageRemove" {
		return nil, errors.New(dm.ErrorMsg)
	}

	dm.IRImgID = imageID
	dm.IROptions = options
	return nil, nil
}

//CopyTool Mock
type MockCopyTool struct {
	Error    bool
	ErrorMsg string
	Dest     string
}

//Mock of the CopyFiles Utility function
func (mcp *MockCopyTool) CopyFiles(reader io.ReadCloser, dest string, source string) error {
	if mcp.Error {
		return errors.New(mcp.ErrorMsg)
	}

	mcp.Dest = dest

	return nil
}

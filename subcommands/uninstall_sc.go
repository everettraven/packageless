package subcommands

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/everettraven/packageless/utils"
)

//Uninstall Sub-Command Object
type UninstallCommand struct {
	//FlagSet so that we can create a custom flag
	fs *flag.FlagSet

	//String for the name of the package to Uninstall
	name string
}

//Instantiation method for a new UninstallCommand
func NewUninstallCommand() *UninstallCommand {
	//Create a new UninstallCommand and set the FlagSet
	uc := &UninstallCommand{
		fs: flag.NewFlagSet("uninstall", flag.ContinueOnError),
	}

	return uc
}

//Name - Gets the name of the Sub-Command
func (uc *UninstallCommand) Name() string {
	return uc.fs.Name()
}

//Init - Parses and Populates values of the Uninstall subcommand
func (uc *UninstallCommand) Init(args []string) error {

	if len(args) <= 0 {
		return errors.New("No package name was found. You must include the name of the package you wish to uninstall.")
	}

	uc.name = args[0]

	return nil
}

//Uninstall - Uninstalls the Uninstall subcommand
func (uc *UninstallCommand) Run() error {
	//Create variables to use later
	var found bool
	var pack utils.Package

	//Create a variable for the executable directory
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	ed := filepath.Dir(ex)

	//Default location of the package list
	packageList := ed + "/package_list.hcl"

	//Parse the package list
	parseOut, err := utils.Parse(packageList, utils.PackageHCLUtil{})

	//Check for errors
	if err != nil {
		return err
	}

	packages := parseOut.(utils.PackageHCLUtil)

	//Check for errors
	if err != nil {
		return err
	}

	//Look for the package we want in the package list
	for _, packs := range packages.Packages {
		//If we find it, set some variables and break
		if packs.Name == uc.name {
			found = true
			pack = packs
			break
		}
	}

	//Make sure we have found the package in the package list
	if !found {
		return errors.New("Could not find package " + uc.name + " in the package list")
	}

	//Check if the corresponding package image is already Uninstalled
	imgExist, err := utils.ImageExists(pack.Image)

	//Check for errors
	if err != nil {
		return err
	}

	//If the image doesn't exist it can't be uninstalled
	if !imgExist {
		return errors.New("Package " + pack.Name + " is not installed.")
	}

	fmt.Println("Removing package", pack.Name)

	//Check for the directories that correspond to this packages volumes
	fmt.Println("Removing package directories")

	//Check the volumes and remove the directories if they exist
	for _, vol := range pack.Volumes {
		//Make sure that a path is given.
		if vol.Path != "" {
			err = RemoveDir(ed + vol.Path)

			if err != nil {
				return err
			}
		}
	}

	//Remove the base directory for the package
	err = RemoveDir(ed + pack.BaseDir)

	if err != nil {
		return err
	}

	//Remove aliases
	err = RemoveAlias(pack.Name, ed)

	if err != nil {
		return err
	}

	//Remove the image
	err = utils.RemoveImage(pack.Image)

	//Check for errors
	if err != nil {
		return err
	}

	return nil
}

//RemoveDir removes a given directory if it exists
func RemoveDir(path string) error {
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

//Remove Alias will remove the alias for the specified package name from the corresponding files
func RemoveAlias(name string, ed string) error {
	//If runtime is on windows remove from macros.doskey and the powershell profile files
	if runtime.GOOS == "windows" {

		//Command Prompt
		//------------------------------------------------
		//Open the doskey file
		file, err := os.OpenFile(ed+"/macros.doskey", os.O_RDWR, 0755)

		var newOut []string

		if err != nil {
			return err
		}

		reader := bufio.NewReader(file)

		//Read the file line by line
		for {
			line, err := reader.ReadString('\n')

			//Check for EOF
			if err != nil && err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			dos := name + "=" + ed + "/packageless run " + name + "\n"

			//if the line is the doskey for this package dont include it in the new file
			if line != dos {
				newOut = append(newOut, line)
			}
		}

		//Close the file
		file.Close()

		//Recreate the doskey file
		newFile, err := os.OpenFile(ed+"/macros.doskey", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

		if err != nil {
			return err
		}

		//Write the contents back to the doskey file
		for _, line := range newOut {
			_, err = newFile.WriteString(line)

			if err != nil {
				return err
			}
		}

		//Close the file
		newFile.Close()

		//------------------------------------------------

		//PowerShell
		//------------------------------------------------

		pwshPath := os.Getenv("USERPROFILE") + "/Documents/WindowsPowerShell/"

		//Open the powershell profile file
		file, err = os.OpenFile(pwshPath+"Microsoft.PowerShell_profile.ps1", os.O_RDWR, 0755)

		//Reset the newOut array
		newOut = nil

		if err != nil {
			return err
		}

		reader = bufio.NewReader(file)

		//Read the file line by line
		for {
			line, err := reader.ReadString('\n')

			//Check for EOF
			if err != nil && err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			alias := "function " + name + "(){ " + ed + "\\packageless.exe run " + name + " }\n"

			//if the line is the alias for this package dont include it in the new file
			if line != alias {
				newOut = append(newOut, line)
			}
		}

		//Close the file
		file.Close()

		//Recreate the powershell profile file
		newFile, err = os.OpenFile(pwshPath+"Microsoft.PowerShell_profile.ps1", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

		if err != nil {
			return err
		}

		//Write the contents back to the powershell profile file
		for _, line := range newOut {
			_, err = newFile.WriteString(line)

			if err != nil {
				return err
			}
		}

		//Close the file
		newFile.Close()

	} else {
		//If it isnt windows, remove it from the bash aliases file
		file, err := os.OpenFile("~/.bash_aliases", os.O_RDWR, 0755)

		var newOut []string

		if err != nil {
			return err
		}

		reader := bufio.NewReader(file)

		//Read the file line by line
		for {
			line, err := reader.ReadString('\n')

			//Check for EOF
			if err != nil && err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			alias := "alias " + name + "=" + "\"" + ed + "/packageless run " + name + "\"" + "\n"

			//if the line is the alias for this package dont include it in the new file
			if line != alias {
				newOut = append(newOut, line)
			}
		}

		//Close the file
		file.Close()

		//Recreate the bash aliases file
		newFile, err := os.OpenFile("~/.bash_aliases", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

		if err != nil {
			return err
		}

		//Write the contents back to the bash aliases file
		for _, line := range newOut {
			_, err = newFile.WriteString(line)

			if err != nil {
				return err
			}
		}

		//Close the file
		newFile.Close()

	}
	return nil
}

package subcommands

import (
	"bufio"
	"io"
	"os"
)

//AddAlias will add the alias for the package name specified
func AddAliasWin(name string, ed string) error {
	//Set the alias for Powershell
	//--------------------------------
	pwshPath := os.Getenv("USERPROFILE") + "/Documents/WindowsPowerShell/"

	//Create the powershell directory if it doesnt exist
	err := MakeDir(pwshPath)

	if err != nil {
		return err
	}

	//Open the powershell alias file
	file, err := OpenFile(pwshPath + "Microsoft.PowerShell_profile.ps1")

	if err != nil {
		return err
	}

	alias := "function " + name + "(){ " + ed + "\\packageless.exe run " + name + " $args }\n"

	_, err = file.WriteString(alias)

	if err != nil {
		return err
	}

	file.Close()

	return nil
}

//Remove Alias will remove the alias for the specified package name from the corresponding files
func RemoveAliasWin(name string, ed string) error {
	//PowerShell
	//------------------------------------------------

	pwshPath := os.Getenv("USERPROFILE") + "/Documents/WindowsPowerShell/"

	//Open the powershell profile file
	file, err := os.OpenFile(pwshPath+"Microsoft.PowerShell_profile.ps1", os.O_RDWR, 0755)

	//Create the newOut array
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

		alias := "function " + name + "(){ " + ed + "\\packageless.exe run " + name + " }\n"

		//if the line is the alias for this package dont include it in the new file
		if line != alias {
			newOut = append(newOut, line)
		}
	}

	//Close the file
	file.Close()

	//Recreate the powershell profile file
	newFile, err := os.OpenFile(pwshPath+"Microsoft.PowerShell_profile.ps1", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

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

	return nil
}

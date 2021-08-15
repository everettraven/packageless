package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

//AddAlias will add the alias for the package name specified
func (u *Utility) AddAliasUnix(name string, ed string) error {

	//get the home directory file path
	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	//Get the shell path
	shellPath := os.Getenv("SHELL")

	var path string

	shellSplit := strings.Split(shellPath, "/")

	shell := shellSplit[len(shellSplit)-1]

	fmt.Println("Shell is " + shell)

	//Get the filepath for the correct shell rc file
	if shell == "bash" {
		path = home + "/.bashrc"
	} else if shell == "zsh" {
		path = home + "/.zshrc"
	} else {
		return errors.New("Shell: " + shell + " is currently unsupported.")
	}

	//If run on linux lets modify the shell rc file to include the new aliases
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)

	if err != nil {
		return err
	}

	//Create the alias and write it to the file
	alias := "alias " + name + "=" + "\"" + ed + "/packageless run " + name + "\"" + "\n"

	_, err = file.WriteString(alias)

	if err != nil {
		return err
	}

	file.Close()

	return nil
}

//Remove Alias will remove the alias for the specified package name from the corresponding files
func (u *Utility) RemoveAliasUnix(name string, ed string) error {

	//get the home directory file path
	home, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	//Get the shell path
	shellPath := os.Getenv("SHELL")

	var path string

	shellSplit := strings.Split(shellPath, "/")

	shell := shellSplit[len(shellSplit)-1]

	//Get the filepath for the correct shell rc file
	if shell == "bash" {
		path = home + "/.bashrc"
	} else if shell == "zsh" {
		path = home + "/.zshrc"
	} else {
		return errors.New("Shell: " + shell + " is currently unsupported.")
	}

	//Open the shell rc file
	file, err := os.OpenFile(path, os.O_RDWR, 0755)

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

	//Recreate the shell rc file
	newFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		return err
	}

	//Write the contents back to the shell rc file
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

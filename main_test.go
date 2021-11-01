package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/everettraven/packageless/utils"
)

type TestCase struct {
	Name        string
	Args        []string
	Err         bool
	ExpectedErr string
}

func TestMain_OsArgs(t *testing.T) {

	if testing.Short() {
		t.Skip("short testing specified, skipping this test")
	}

	//Create a list of structs to hold data info for testing cases for windows
	casesWindows := []TestCase{
		{
			"Install Test",
			[]string{"packageless", "install", "python"},
			false,
			"",
		},
		{
			//We know this one will fail because go test doesn't include tty which the docker command we run uses
			//This test failing lets us know the command is attempting to run properly
			"Run Test",
			[]string{"packageless", "run", "python"},
			true,
			"",
		},
		{
			"Update Test",
			[]string{"packageless", "update", "python"},
			false,
			"",
		},
		{
			"Upgrade Test",
			[]string{"packageless", "upgrade", "python"},
			false,
			"",
		},
		{
			"Uninstall Test",
			[]string{"packageless", "uninstall", "python"},
			false,
			"",
		},
	}

	//Create a list of structs to hold data info for testing cases for windows
	casesUnix := []TestCase{
		{
			"Install Test",
			[]string{"packageless", "install", "python"},
			true,
			"Shell: go is currently unsupported.",
		},
		{
			//We know this one will fail because go test doesn't include tty which the docker command we run uses
			//This test failing lets us know the command is attempting to run properly
			"Run Test",
			[]string{"packageless", "run", "python"},
			true,
			"",
		},
		{
			"Update Test",
			[]string{"packageless", "update", "python"},
			false,
			"",
		},
		{
			"Upgrade Test",
			[]string{"packageless", "upgrade", "python"},
			false,
			"",
		},
		{
			"Uninstall Test",
			[]string{"packageless", "uninstall", "python"},
			true,
			"Shell: go is currently unsupported.",
		},
	}

	//Copy necessary files to the test build location
	ex, err := os.Executable()

	if err != nil {
		t.Fatal(err)
	}

	ed := filepath.Dir(ex)

	if err != nil {
		t.Fatal(err)
	}

	//Create the .packageless directory and necessary subdirectories in the executable directory
	err = utils.NewUtility().MakeDir(ed + "/.packageless")

	if err != nil {
		t.Fatal(err)
	}

	err = utils.NewUtility().MakeDir(ed + "/.packageless/pims_config")

	if err != nil {
		t.Fatal(err)
	}

	err = utils.NewUtility().MakeDir(ed + "/.packageless/pims")

	if err != nil {
		t.Fatal(err)
	}

	err = utils.NewUtility().MakeDir(ed + "/Documents/WindowsPowerShell")

	if err != nil {
		t.Fatal(err)
	}

	err = Copy("./config.hcl", ed+"/.packageless/config.hcl")

	if err != nil {
		t.Fatal(err)
	}

	//Change the HOME/USERPROFILE environment variable to point to the executable directory

	//First we need to save the old HOME/USERPROFILE environment variable so we can change it back later
	var oldHomeEnv string

	//If the user is on windows we need to change the USERPROFILE value, otherwise it is the HOME value
	if runtime.GOOS == "windows" {
		oldHomeEnv = os.Getenv("USERPROFILE")

		//set a new HOME ENV variable
		err = os.Setenv("USERPROFILE", ed)

		if err != nil {
			t.Fatalf("Error trying to set new USERPROFILE env value: %s", err.Error())
		}
	} else {
		oldHomeEnv = os.Getenv("HOME")

		//set a new HOME ENV variable
		err = os.Setenv("HOME", ed)

		if err != nil {
			t.Fatalf("Error trying to set new HOME env value: %s", err.Error())
		}
	}

	var cases []TestCase

	if runtime.GOOS == "windows" {
		cases = casesWindows
	} else {
		cases = casesUnix
	}

	//Loop through the test cases
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			fmt.Println(i)
			//Set the args when running
			os.Args = tc.Args

			//Run the main function
			exit, err := wrappedMain()

			if (exit != 0) != tc.Err {
				t.Fatalf("Fail - Exit code did not match the expected | Received Error: %s", err)
			} else {
				if tc.ExpectedErr != "" {
					if err.Error() != tc.ExpectedErr {
						t.Fatalf("Fail - Expected Error: %s | Received: %s", tc.ExpectedErr, err.Error())
					}
				}
			}

		})
	}

	//Before we exit the tests lets set the HOME/USERPROFILE env var value back
	if runtime.GOOS == "windows" {
		//set a new HOME ENV variable
		err = os.Setenv("USERPROFILE", oldHomeEnv)

		if err != nil {
			t.Fatalf("Failed to set the USERPROFILE env value back to the original value of '%s': %s", oldHomeEnv, err.Error())
		}
	} else {

		err = os.Setenv("HOME", oldHomeEnv)

		if err != nil {
			t.Fatalf("Failed to set the HOME env value back to the original value of '%s': %s", oldHomeEnv, err.Error())
		}
	}

}

//Function to copy files
func Copy(src string, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(dest)

	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)

	if err != nil {
		return err
	}

	return out.Close()
}

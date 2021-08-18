package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
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

	err = Copy("./package_list.hcl", ed+"/package_list.hcl")

	if err != nil {
		t.Fatal(err)
	}

	err = Copy("./config.hcl", ed+"/config.hcl")

	if err != nil {
		t.Fatal(err)
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
				t.Fatalf("Fail - Exit code did not match the expected")
			} else {
				if tc.ExpectedErr != "" {
					if err.Error() != tc.ExpectedErr {
						t.Fatalf("Fail - Expected Error: %s | Received: %s", tc.ExpectedErr, err.Error())
					}
				}
			}

		})
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

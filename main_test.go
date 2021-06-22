package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestMain_OsArgs(t *testing.T) {

	if testing.Short() {
		t.Skip("short testing specified, skipping this test")
	}

	//Create a list of structs to hold data info for testing cases
	cases := []struct {
		Name string
		Args []string
		Err  bool
	}{
		{
			"Install Test",
			[]string{"packageless", "install", "python"},
			false,
		},
		{
			//We know this one will fail because go test doesn't include tty which the docker command we run uses
			//This test failing lets us know the command is attempting to run properly
			"Run Test",
			[]string{"packageless", "run", "python"},
			true,
		},
		{
			"Upgrade Test",
			[]string{"packageless", "upgrade", "python"},
			false,
		},
		{
			"Uninstall Test",
			[]string{"packageless", "uninstall", "python"},
			false,
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

	//Loop through the test cases
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Name), func(t *testing.T) {
			fmt.Println(i)
			//Set the args when running
			os.Args = tc.Args

			//Run the main function
			exit := wrappedMain()

			if (exit != 0) != tc.Err {
				t.Fatalf("Fail: %d", exit)
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

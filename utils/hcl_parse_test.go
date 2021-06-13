package utils

import (
	"strconv"
	"testing"

	"github.com/hashicorp/hcl2/hclparse"
)

//Test the ParseBody function with a Config Object HCL Body
func TestParseBodyConfig(t *testing.T) {
	//Create the HCL byte array
	hcl := []byte(`base_dir="./"
	start_port=3000
	port_increment=1`)

	//Create the parser
	parser := hclparse.NewParser()

	//Parse the byte array
	f, diags := parser.ParseHCL(hcl, "config_test")

	//If error it fails
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}

	//Create a new utility
	util := NewUtility()

	//Parse the HCL Body
	parseOut, err := util.ParseBody(f.Body, Config{})

	//If error then fails
	if err != nil {
		t.Fatal(err)
	}

	//Get the config object
	config := parseOut.(Config)

	//Test the config values
	if config.BaseDir != "./" {
		t.Fatal("Config Base Directory should be './' | Received: " + config.BaseDir)
	}

	if config.PortInc != 1 {
		t.Fatal("Config Port Increment should be '1' | Received: " + strconv.Itoa(config.PortInc))
	}

	if config.StartPort != 3000 {
		t.Fatal("Config Start Port should be '3000' | Received: " + strconv.Itoa(config.StartPort))
	}
}

//Test the parse body function with a package object that contains the optional copy fields
func TestParseBodyPackageWithCopy(t *testing.T) {
	//Create the HCL byte array
	hcl := []byte(`package "test_pack" {
		image="test"
		base_dir="/base"
		
		volume {
			path="/test/path"
			mount="/test/"
		}
	
		copy {
			source="/test_source/"
			dest="/test_dest/"
		}
	
		port="3000"
	}`)

	//Create the parser
	parser := hclparse.NewParser()

	//Parse the byte array
	f, diags := parser.ParseHCL(hcl, "config_test")

	//If error it fails
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}

	//Create a new utility
	util := NewUtility()

	//Parse the HCL Body
	parseOut, err := util.ParseBody(f.Body, PackageHCLUtil{})

	if err != nil {
		t.Fatal(err)
	}

	//Get the package object
	packs := parseOut.(PackageHCLUtil)

	//If there is more or less than one package, the test should fail
	if len(packs.Packages) != 1 {
		t.Fatal("The # of packages expected is '1' | Received: " + strconv.Itoa(len(packs.Packages)))
	}

	pack := packs.Packages[0]
	//Make sure the package name is correct
	if pack.Name != "test_pack" {
		t.Fatal("Package name should be 'test_pack' | Received: " + pack.Name)
	}

	//Make sure the package base directory is correct
	if pack.BaseDir != "/base" {
		t.Fatal("Package base directory should be '/base' | Received: " + pack.BaseDir)
	}

	//Make sure the package port is correct
	if pack.Port != "3000" {
		t.Fatal("Package port should be '3000' | Received: " + pack.Port)
	}

	//Make sure the package image is correct
	if pack.Image != "test" {
		t.Fatal("Package image should be 'test' | Received: " + pack.Image)
	}

	//Make sure the volumes array is of length 1
	if len(pack.Volumes) != 1 {
		t.Fatal("Package # of volumes should be '1' | Received: " + strconv.Itoa(len(pack.Volumes)))
	}

	vol := pack.Volumes[0]

	//Make sure the volumes array path is correct
	if vol.Path != "/test/path" {
		t.Fatal("Package volume host path should be '/python/packages/' | Received: " + vol.Path)
	}

	//Make sure the volumes mount path is correct
	if vol.Mount != "/test/" {
		t.Fatal("Package volume mount should be '/test/' | Received: " + vol.Mount)
	}

	//Make sure the copies array is of length 1
	if len(pack.Copies) != 1 {
		t.Fatal("Package # of copies should be '1' | Received: " + strconv.Itoa(len(pack.Copies)))
	}

	cp := pack.Copies[0]

	//Make sure the copy source is correct
	if cp.Source != "/test_source/" {
		t.Fatal("Package copy source should be '/test_source/' | Received: " + cp.Source)
	}

	//Make sure the copy destination is correct
	if cp.Dest != "/test_dest/" {
		t.Fatal("Package copy dest should be '/test_dest/' | Received: " + cp.Dest)
	}
}

//Test
//Test the parse body function with a package object that does not container the optional copy fields
func TestParseBodyPackageNoCopy(t *testing.T) {
	//Create the HCL byte array
	hcl := []byte(`package "test_pack" {
		image="test"
		base_dir="/base"
		
		volume {
			path="/test/path"
			mount="/test/"
		}
	
		port="3000"
	}`)

	//Create the parser
	parser := hclparse.NewParser()

	//Parse the byte array
	f, diags := parser.ParseHCL(hcl, "config_test")

	//If error it fails
	if diags.HasErrors() {
		t.Fatal(diags.Error())
	}

	//Create a new utility
	util := NewUtility()

	//Parse the HCL Body
	parseOut, err := util.ParseBody(f.Body, PackageHCLUtil{})

	if err != nil {
		t.Fatal(err)
	}

	//Get the package object
	packs := parseOut.(PackageHCLUtil)

	//If there is more or less than one package, the test should fail
	if len(packs.Packages) != 1 {
		t.Fatal("The # of packages expected is '1' | Received: " + strconv.Itoa(len(packs.Packages)))
	}

	pack := packs.Packages[0]
	//Make sure the package name is correct
	if pack.Name != "test_pack" {
		t.Fatal("Package name should be 'test_pack' | Received: " + pack.Name)
	}

	//Make sure the package base directory is correct
	if pack.BaseDir != "/base" {
		t.Fatal("Package base directory should be '/base' | Received: " + pack.BaseDir)
	}

	//Make sure the package port is correct
	if pack.Port != "3000" {
		t.Fatal("Package port should be '3000' | Received: " + pack.Port)
	}

	//Make sure the package image is correct
	if pack.Image != "test" {
		t.Fatal("Package image should be 'test' | Received: " + pack.Image)
	}

	//Make sure the volumes array is of length 1
	if len(pack.Volumes) != 1 {
		t.Fatal("Package # of volumes should be '1' | Received: " + strconv.Itoa(len(pack.Volumes)))
	}

	vol := pack.Volumes[0]

	//Make sure the volumes array path is correct
	if vol.Path != "/test/path" {
		t.Fatal("Package volume host path should be '/python/packages/' | Received: " + vol.Path)
	}

	//Make sure the volumes mount path is correct
	if vol.Mount != "/test/" {
		t.Fatal("Package volume mount should be '/test/' | Received: " + vol.Mount)
	}

	//Make sure the copies array is empty
	if len(pack.Copies) > 0 {
		t.Fatal("Package # of copies should be '0' | Received: " + strconv.Itoa(len(pack.Copies)))
	}
}

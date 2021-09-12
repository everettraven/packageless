package utils

import (
	"strconv"
	"testing"

	"github.com/hashicorp/hcl2/hclparse"
)

// Unit Tests For HCL Parsing
//-----------------------------------------------------------------------------------------

//Test the ParseBody function with a Config Object HCL Body
func TestParseBodyConfig(t *testing.T) {
	//Create the HCL byte array
	hcl := []byte(`base_dir="./"
	start_port=3000
	port_increment=1
	alias=true`)

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

	if !config.Alias {
		t.Fatal("Config attribute 'alias' should be set to true. Instead it was set to false.")
	}
}

//Test the parse body function with a package object that contains the optional copy fields
func TestParseBodyPackageWithCopy(t *testing.T) {
	//Create the HCL byte array
	hcl := []byte(`package "test_pack" {
		base_dir="/base"
		version "latest" {
			image="test"
			
			volume {
				path="/test/path"
				mount="/test/"
			}
		
			copy {
				source="/test_source/"
				dest="/test_dest/"
			}
		
			port="3000"
		}
		
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

	//Make sure there is only one version

	if len(pack.Versions) != 1 {
		t.Fatal("The # of versions expected is '1' | Received: " + strconv.Itoa(len(pack.Versions)))
	}

	//Get the package version and make sure the fields are correct
	version := pack.Versions[0]

	if version.Version != "latest" {
		t.Fatal("Package version should be 'latest' | Received: " + version.Version)
	}

	//Make sure the package port is correct
	if version.Port != "3000" {
		t.Fatal("Package port should be '3000' | Received: " + version.Port)
	}

	//Make sure the package image is correct
	if version.Image != "test" {
		t.Fatal("Package image should be 'test' | Received: " + version.Image)
	}

	//Make sure the volumes array is of length 1
	if len(version.Volumes) != 1 {
		t.Fatal("Package # of volumes should be '1' | Received: " + strconv.Itoa(len(version.Volumes)))
	}

	vol := version.Volumes[0]

	//Make sure the volumes array path is correct
	if vol.Path != "/test/path" {
		t.Fatal("Package volume host path should be '/python/packages/' | Received: " + vol.Path)
	}

	//Make sure the volumes mount path is correct
	if vol.Mount != "/test/" {
		t.Fatal("Package volume mount should be '/test/' | Received: " + vol.Mount)
	}

	//Make sure the copies array is of length 1
	if len(version.Copies) != 1 {
		t.Fatal("Package # of copies should be '1' | Received: " + strconv.Itoa(len(version.Copies)))
	}

	cp := version.Copies[0]

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
		base_dir="/base"
		version "latest" {
			image="test"
			
			volume {
				path="/test/path"
				mount="/test/"
			}
		
			port="3000"
		}
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

	//Make sure the number of versions is correct
	if len(pack.Versions) != 1 {
		t.Fatal("The # of versions expected is '1' | Received: " + strconv.Itoa(len(pack.Versions)))
	}

	//Get and check the version values
	version := pack.Versions[0]

	if version.Version != "latest" {
		t.Fatal("Package version should be 'latest' | Received: " + version.Version)
	}

	//Make sure the package port is correct
	if version.Port != "3000" {
		t.Fatal("Package port should be '3000' | Received: " + version.Port)
	}

	//Make sure the package image is correct
	if version.Image != "test" {
		t.Fatal("Package image should be 'test' | Received: " + version.Image)
	}

	//Make sure the volumes array is of length 1
	if len(version.Volumes) != 1 {
		t.Fatal("Package # of volumes should be '1' | Received: " + strconv.Itoa(len(version.Volumes)))
	}

	vol := version.Volumes[0]

	//Make sure the volumes array path is correct
	if vol.Path != "/test/path" {
		t.Fatal("Package volume host path should be '/python/packages/' | Received: " + vol.Path)
	}

	//Make sure the volumes mount path is correct
	if vol.Mount != "/test/" {
		t.Fatal("Package volume mount should be '/test/' | Received: " + vol.Mount)
	}

	//Make sure the copies array is empty
	if len(version.Copies) > 0 {
		t.Fatal("Package # of copies should be '0' | Received: " + strconv.Itoa(len(version.Copies)))
	}
}

// Integration Tests For HCL Parsing
//-----------------------------------------------------------------------------------------

//Integration test for reading the test Config file
func TestHCLParse_Integration_Config(t *testing.T) {
	//Create the util tool
	util := NewUtility()

	//Read the Test HCL Config file
	body, err := util.GetHCLBody("../testing/test_config.hcl")

	//Shouldn't throw an error
	if err != nil {
		t.Fatal(err)
	}

	//Parse the HCL Body into an object
	parseOut, err := util.ParseBody(body, Config{})

	//Shouldn't throw an error
	if err != nil {
		t.Fatal(err)
	}

	//Get the parsed object
	config := parseOut.(Config)

	//Set the expected variables
	cBD := "./test"
	cSP := 5000
	cPI := 100
	cA := true

	//Ensure the base dir is correct
	if config.BaseDir != cBD {
		t.Fatalf("HCL Parse Integration: Expected BaseDir: %s | Received: %s", cBD, config.BaseDir)
	}

	//Ensure the start port is correct
	if config.StartPort != cSP {
		t.Fatalf("HCL Parse Integration: Expected StartPort: %d | Received: %d", cSP, config.StartPort)
	}

	//Ensure the port increment is correct
	if config.PortInc != cPI {
		t.Fatalf("HCL Parse Integration: Expected PortInc: %d | Received: %d", cPI, config.PortInc)
	}

	if config.Alias != cA {
		t.Fatalf("HCL Parse Integration: Expected Alias: %s | Received: %s", strconv.FormatBool(cA), strconv.FormatBool(config.Alias))
	}

}

//Integration test for reading the test package list file
func TestHCLParse_Integration_PackageList(t *testing.T) {
	//Create the util tool
	util := NewUtility()

	//Read the Test HCL Config file
	body, err := util.GetHCLBody("../testing/test_package_list.hcl")

	//Shouldn't throw an error
	if err != nil {
		t.Fatal(err)
	}

	//Parse the HCL Body into an object
	parseOut, err := util.ParseBody(body, PackageHCLUtil{})

	//Shouldn't throw an error
	if err != nil {
		t.Fatal(err)
	}

	//Get the parsed object
	packs := parseOut.(PackageHCLUtil)

	//Set expected variables
	pLen := 1
	pName := "test"
	pImage := "packageless/testing"
	pBD := "/base"
	pPort := "3000"
	pVersion := "latest"
	vLen := 2
	v1Path := "/a/path"
	v1Mount := "/mount/path"
	v2Path := ""
	v2Mount := "/run/"
	cpLen := 1
	cpSource := "/a/source"
	cpDest := "/a/dest"

	//Ensure packages length is correct
	if len(packs.Packages) != pLen {
		t.Fatalf("Parse HCL Integration: Expected Packages Length: %d | Received: %d", pLen, len(packs.Packages))
	}

	p := packs.Packages[0]

	version := p.Versions[0]

	if version.Version != pVersion {
		t.Fatalf("Parse HCL Intergration: Expected Package Version: %s | Received: %s", pVersion, version.Version)
	}

	//Ensure the package name is correct
	if p.Name != pName {
		t.Fatalf("Parse HCL Integration: Expected Package Name: %s | Received: %s", pName, p.Name)
	}

	//Ensure the package image is correct
	if version.Image != pImage {
		t.Fatalf("Parse HCL Integration: Expected Package Image: %s | Received: %s", pImage, version.Image)
	}

	//Ensure the package base directory is correct
	if p.BaseDir != pBD {
		t.Fatalf("ParseHCL Integration: Expected Package BaseDir: %s | Received: %s", pBD, p.BaseDir)
	}

	//Ensure the package port is correct
	if version.Port != pPort {
		t.Fatalf("ParseHCL Integration: Expected Package Port: %s | Received: %s", pPort, version.Port)
	}

	//Ensure the volumes length matches
	if len(version.Volumes) != vLen {
		t.Fatalf("ParseHCL Integration: Expected Package Volumes Len: %d | Received: %d", vLen, len(version.Volumes))
	}

	vols := version.Volumes

	//Ensure the first volume path matches
	if vols[0].Path != v1Path {
		t.Fatalf("ParseHCL Integration: Expected Package Volume 1 Path: %s | Received: %s", v1Path, vols[0].Path)
	}

	//Ensure the first volume mount path matches
	if vols[0].Mount != v1Mount {
		t.Fatalf("ParseHCL Integration: Expected Package Volume 1 Mount: %s | Received: %s", v1Mount, vols[0].Mount)
	}

	//Ensure the second volume path matches
	if vols[1].Path != v2Path {
		t.Fatalf("ParseHCL Integration: Expected Package Volume 2 Path: %s | Received: %s", v2Path, vols[1].Path)
	}

	//Ensure the second volume mount matches
	if vols[1].Mount != v2Mount {
		t.Fatalf("ParseHCL Integration: Expected Package Volume 2 Mount: %s | Received: %s", v2Mount, vols[1].Mount)
	}

	//Ensure the copies length matches
	if len(version.Copies) != cpLen {
		t.Fatalf("ParseHCL Integration: Expected Package Copies Len: %d | Received: %d", cpLen, len(version.Copies))
	}

	cp := version.Copies[0]

	//Ensure the copy source matches
	if cp.Source != cpSource {
		t.Fatalf("ParseHCL Integration: Expected Package Copy Source: %s | Received: %s", cpSource, cp.Source)
	}

	//Ensure the copy dest matches
	if cp.Dest != cpDest {
		t.Fatalf("ParseHCL Integration: Expected Package Copy Dest: %s | Received: %s", cpDest, cp.Dest)
	}

}

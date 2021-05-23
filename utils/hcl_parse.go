package utils

import (
	"errors"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
)

//Copy object to parse the copy block in the package list
type Copy struct {
	Source string `hcl:"source,attr"`
	Dest   string `hcl:"dest,attr"`
}

//Volume object to parse the volume block in the package list
type Volume struct {
	Path  string `hcl:"path,optional"`
	Mount string `hcl:"mount,attr"`
}

//Package object to parse the package block in the package list
type Package struct {
	Name    string   `hcl:"name,label"`
	Image   string   `hcl:"image,attr"`
	Volumes []Volume `hcl:"volume,block"`
	Copies  []*Copy  `hcl:"copy,block"`
}

//PackageHCLUtil object to contain a list of packages and all their attributes after the parsing of the package list
type PackageHCLUtil struct {
	Packages []Package `hcl:"package,block"`
}

//Parse function to parse the HCL file given in the filepath
func Parse(filepath string) (PackageHCLUtil, error) {
	//Create a parser
	parser := hclparse.NewParser()

	//Create the object to be decoded to
	var packages PackageHCLUtil

	//Parse the data
	parseData, parseDiags := parser.ParseHCLFile(filepath)

	//Check for errors
	if parseDiags.HasErrors() {
		return packages, errors.New("ParseDiags: " + parseDiags.Error())
	}

	//Decode the parsed HCL to the Object
	decodeDiags := gohcl.DecodeBody(parseData.Body, nil, &packages)

	//Check for errors
	if decodeDiags.HasErrors() {
		return packages, errors.New("DecodeDiags: " + decodeDiags.Error())
	}

	return packages, nil
}

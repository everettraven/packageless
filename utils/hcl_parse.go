package utils

import (
	"errors"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
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
	BaseDir string   `hcl:"base_dir,attr"`
	Volumes []Volume `hcl:"volume,block"`
	Copies  []*Copy  `hcl:"copy,block"`
	Port    string   `hcl:"port,optional"`
}

//PackageHCLUtil object to contain a list of packages and all their attributes after the parsing of the package list
type PackageHCLUtil struct {
	Packages []Package `hcl:"package,block"`
}

//Config object to contain the configuration details
type Config struct {
	BaseDir   string `hcl:"base_dir,attr"`
	StartPort int    `hcl:"start_port,attr"`
	PortInc   int    `hcl:"port_increment,attr"`
	Alias     bool   `hcl:"alias,attr"`
}

//Parse function to parse the HCL body given
func (u *Utility) ParseBody(body hcl.Body, out interface{}) (interface{}, error) {

	switch out.(type) {
	default:
		return nil, errors.New("Unexpected type passed into the HCL parse function")

	case PackageHCLUtil:
		//Create the object to be decoded to
		var packages PackageHCLUtil

		//Decode the parsed HCL to the Object
		decodeDiags := gohcl.DecodeBody(body, nil, &packages)

		//Check for errors
		if decodeDiags.HasErrors() {
			return packages, errors.New("DecodeDiags: " + decodeDiags.Error())
		}

		return packages, nil

	case Config:
		//Create the object to be decoded to
		var config Config

		//Decode the parsed HCL to the Object
		decodeDiags := gohcl.DecodeBody(body, nil, &config)

		//Check for errors
		if decodeDiags.HasErrors() {
			return config, errors.New("DecodeDiags: " + decodeDiags.Error())
		}

		return config, nil
	}

}

//GetHCLBody gets the HCL Body from a given filepath
func (u *Utility) GetHCLBody(filepath string) (hcl.Body, error) {
	//create a parser
	parser := hclparse.NewParser()

	file, diags := parser.ParseHCLFile(filepath)

	if diags.HasErrors() {
		return nil, errors.New("Parse Error: " + diags.Error())
	}

	return file.Body, nil
}

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
type PackageImage struct {
	Name     string    `hcl:"name,label"`
	BaseDir  string    `hcl:"base_dir,attr"`
	Versions []Version `hcl:"version,block"`
}

type Version struct {
	Version string   `hcl:"version,label"`
	Image   string   `hcl:"image,attr"`
	Volumes []Volume `hcl:"volume,block"`
	Copies  []*Copy  `hcl:"copy,block"`
	Port    string   `hcl:"port,optional"`
}

//PackageHCLUtil object to contain a list of packages and all their attributes after the parsing of the package list
type PimHCLUtil struct {
	Pims []PackageImage `hcl:"pim,block"`
}

//Config object to contain the configuration details
type Config struct {
	BaseDir        string `hcl:"base_dir,attr"`
	StartPort      int    `hcl:"start_port,attr"`
	PortInc        int    `hcl:"port_increment,attr"`
	Alias          bool   `hcl:"alias,attr"`
	RepositoryHost string `hcl:"repository_host,attr"`
	PimsConfigDir  string `hcl:"pims_config_dir,attr"`
	PimsDir        string `hcl:"pims_dir,attr"`
}

//Parse function to parse the HCL body given
func (u *Utility) ParseBody(body hcl.Body, out interface{}) (interface{}, error) {

	switch out.(type) {
	default:
		return nil, errors.New("Unexpected type passed into the HCL parse function")

	case PimHCLUtil:
		//Create the object to be decoded to
		var pims PimHCLUtil

		//Decode the parsed HCL to the Object
		decodeDiags := gohcl.DecodeBody(body, nil, &pims)

		//Check for errors
		if decodeDiags.HasErrors() {
			return pims, errors.New("DecodeDiags: " + decodeDiags.Error())
		}

		return pims, nil

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

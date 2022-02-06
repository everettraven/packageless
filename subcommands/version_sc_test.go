package subcommands

import (
	"testing"

	"github.com/everettraven/packageless/utils"
)

func TestVersionName(t *testing.T) {
	mockUtility := utils.NewMockUtility()

	expected := "version"
	vc := NewVersionCommand(mockUtility)

	if vc.Name() != expected {
		t.Fatalf("The version subcommand's name should be: '%s' but was '%s'", expected, vc.Name())
	}
}

func TestVersionInit(t *testing.T) {
	mockUtility := utils.NewMockUtility()

	vc := NewVersionCommand(mockUtility)

	err := vc.Init([]string{})

	if err != nil {
		t.Fatalf("This method should do nothing except return nil | Received: %s", err)
	}
}

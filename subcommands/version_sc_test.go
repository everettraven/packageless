package subcommands

import (
	"reflect"
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

func TestVersionRun(t *testing.T) {
	mockUtility := utils.NewMockUtility()

	vc := NewVersionCommand(mockUtility)

	err := vc.Init([]string{})

	if err != nil {
		t.Fatalf("This method should do nothing except return nil | Received: %s", err)
	}

	err = vc.Run()

	if err != nil {
		t.Fatalf("This subcommand should not return an error | Received: %s", err)
	}

	callStack := []string{
		"RenderInfoMarkdown",
	}

	//If the call stack doesn't match the test fails
	if !reflect.DeepEqual(callStack, mockUtility.Calls) {
		t.Fatalf("Call Stack does not match the expected call stack. Call Stack: %v | Expected Call Stack: %v", mockUtility.Calls, callStack)
	}
}

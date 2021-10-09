package subcommands

import "testing"

func TestVersionName(t *testing.T) {
	expected := "version"
	vc := NewVersionCommand()

	if vc.Name() != expected {
		t.Fatalf("The version subcommand's name should be: '%s' but was '%s'", expected, vc.Name())
	}
}

func TestVersionInit(t *testing.T) {
	vc := NewVersionCommand()

	err := vc.Init([]string{})

	if err != nil {
		t.Fatalf("This method should do nothing except return nil | Received: %s", err)
	}
}

func ExampleVersion() {
	vc := NewVersionCommand()

	vc.Run()

	// Output:
	// Packageless Version: v0.0.0

}

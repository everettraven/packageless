package main

import (
	"fmt"
	"os"

	"github.com/everettraven/packageless/subcommands"
)

func main() {
	if err := subcommands.SubCommand(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

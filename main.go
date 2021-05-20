package main

import (
	"fmt"
	"os"

	"github.com/everettraven/packageless/utils"
)

func main() {
	if err := utils.SubCommand(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

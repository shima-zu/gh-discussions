package main

import (
	"fmt"
	"os"

	"github.com/shima-zu/gh-discussions/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

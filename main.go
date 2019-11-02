package main

import (
	"os"

	"github.com/uitml/frink/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

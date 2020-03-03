package main

import (
	"os"

	"github.com/uitml/frink/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}

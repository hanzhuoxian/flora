// floractl is the command line tool for Flora.
package main

import (
	"os"

	"github.com/hanzhuoxian/flora/internal/floractl/cmd"
)

func main() {
	command := cmd.NewDefaultCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

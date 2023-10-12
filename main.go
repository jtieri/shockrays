package main

import (
	"github.com/jtieri/shockrays/cmd"
	"os"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cmd.Execute(homeDir)
}

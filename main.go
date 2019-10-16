package main

import (
	"os"

	"github.com/marian-craciunescu/symantecbeat/cmd"

	_ "github.com/marian-craciunescu/symantecbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

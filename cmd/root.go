package cmd

import (
	"github.com/marian-craciunescu/symantecbeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
	"github.com/elastic/beats/libbeat/cmd/instance"
)

// Name of this beat
var Name = "symantecbeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, instance.Settings{Name: Name})

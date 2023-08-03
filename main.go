package main

import (
	"github.com/sailhouse/sailhouse/cmd"
)

var configFile string
var app string
var topic string

var (
	version string
)

func main() {
	if version == "" {
		version = "v0.0.0"
	}

	cmd.Execute(version)
}

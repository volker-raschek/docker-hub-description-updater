package main

import (
	"github.com/volker-raschek/docker-hub-description-updater/cmd"
)

var (
	version string
)

func main() {
	_ = cmd.Execute(version)
}

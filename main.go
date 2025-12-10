package main

import (
	"gobi/cmd"
	"gobi/components/logger"
	"gobi/components/project"
)

func main() {
	// steps:
	// 1. setup project
	// 2. go through all config files and files and cache all required stuff -> will also be used by incremental build
	// 3. setup builder
	// 4. run build
	project.Setup()

	file := logger.Setup()
	defer file.Close()

	cmd.Execute()
}

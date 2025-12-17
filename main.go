package main

import (
	"gobi/cmd"
	"gobi/components/logger"
	"gobi/components/project"
)

// @todo work on archi, as data shall be kept and trasnferred as optimally as possible

func main() {
	// steps:
	// 1. setup project
	// 2. go through all config files and files and cache all required stuff -> will also be used by incremental build
	// 3. setup builder
	// 4. run build

	// load project config
	project.LoadProjectConfiguration()

	file := logger.Setup()
	defer file.Close()

	// load so we have incremental build
	project.LoadBuildCache()
	project.Setup()

	cmd.Execute()

	// @todo check how pacgo does this
	project.CacheBuildData()
}

package main

import (
	"gobi/cmd"
	"gobi/components/logger"
	"gobi/components/project"
)

func main() {
	project.Setup()

	file := logger.Setup()
	defer file.Close()

	cmd.Execute()
}

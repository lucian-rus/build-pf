package main

import (
	"gobi/cmd"
	"gobi/components/logger"
)

func main() {
	file := logger.Setup()
	defer file.Close()

	cmd.Execute()
}

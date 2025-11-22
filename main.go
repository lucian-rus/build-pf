package main

import (
	"gobi/components/builder"
	"gobi/components/crawler"
	"os"
)

func main() {
	var sourceFilesList []string
	crawler.ScanDirectoryForSourceFiles(os.Args[1], &sourceFilesList, false)

	os.Chdir(os.Args[1])
	builder.PrepareBuildCommands(sourceFilesList)
	builder.RunBuilder()
}

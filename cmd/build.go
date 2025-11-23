package cmd

import (
	"gobi/components/builder"
	"gobi/components/crawler"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	buildCmd = &cobra.Command{
		Use: "build",
		Run: func(cmd *cobra.Command, args []string) {

			runBuildCmd(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuildCmd(args []string) {
	if err := validateNumberOfArguments(len(args), 1, 1); err != nil {
		log.Fatal(err)
	}
	var sourceFilesList []string
	crawler.ScanDirectoryForSourceFiles(args[0], &sourceFilesList, false)
	os.Chdir(args[0])
	builder.PrepareBuildCommands(sourceFilesList)
	builder.RunBuilder()
}

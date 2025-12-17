package cmd

import (
	"fmt"
	"gobi/modules/builder"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		runBuildCmd(args)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuildCmd(args []string) {
	// @todo for now only support current dir
	if err := validateNumberOfArguments(len(args), 0, 0); err != nil {
		log.Fatal(err)
	}

	startNow := time.Now()

	builder.Build()

	// step 6 -> check benchmark
	fmt.Println("this took", time.Since(startNow))
}

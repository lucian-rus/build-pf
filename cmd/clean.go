package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	cleanCmd = &cobra.Command{
		Use: "clean",
		Run: func(cmd *cobra.Command, args []string) {
			runCleanCmd(args)
		},
	}

	// @todo handle this once proper init is done. for now, init is bad so no need to implement this

	cleanBuildCmd = &cobra.Command{
		Use: "build",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("NOT YET IMPLEMENTED")
		},
	}

	cleanLogCmd = &cobra.Command{
		Use: "log",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("NOT YET IMPLEMENTED")
		},
	}
)

func init() {
	rootCmd.AddCommand(cleanCmd)

	cleanCmd.AddCommand(cleanBuildCmd)
	cleanCmd.AddCommand(cleanLogCmd)
}

func runCleanCmd(args []string) {
	// @todo for now only support current dir
	if err := validateNumberOfArguments(len(args), 0, 0); err != nil {
		log.Fatal(err)
	}

	// project.CleanProject()
}

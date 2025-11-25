/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gobi",
	Short: "gobi is a c/c++ build tool",
	Long:  `gobi is a c/c++ build tool`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pacgo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func validateNumberOfArguments(num, minNum, maxNum int) error {
	if err := checkIfTooManyArguments(num, maxNum); err != nil {
		return err
	}

	return checkIfNotEnoughArguments(num, minNum)
}

func checkIfTooManyArguments(num, maxNum int) error {
	if num > maxNum {
		return errors.New("too many arguments")
	}

	return nil
}

func checkIfNotEnoughArguments(num, minNum int) error {
	if num < minNum {
		return errors.New("not enough arguments")
	}

	return nil
}

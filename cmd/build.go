package cmd

import (
	"gobi/components/builder"
	"gobi/components/crawler"
	"gobi/components/env"
	"log"
	"os"
	"path/filepath"

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

	cwd, _ := os.Getwd()
	projectWorkingDirectory := filepath.Join(cwd, args[0])
	projectConfigFileName := filepath.Join(projectWorkingDirectory, "gobi.json")
	if err := crawler.ReadProjectConfigFileContent(projectConfigFileName, &env.ProjectConfiguration); err != nil {
		log.Fatal(err)
	}

	// @todo start handling errors
	for _, subdir := range env.ProjectConfiguration.Subdirectories {
		buildLibrary(args, subdir, projectWorkingDirectory)
	}

	// @todo resolve build and squash stuff together
	// global include paths have been resolved
	// dependency output global path -> yet to be resolved
	// make build generic -> yet to be resolved
	// build project -> yet to be resolved
}

func buildLibrary(args []string, subdir, projectWorkingDirectory string) {
	libraryWorkingDirectory := filepath.Join(projectWorkingDirectory, subdir)
	os.Chdir(libraryWorkingDirectory)
	log.Println("Changing directory to ", libraryWorkingDirectory)

	libraryConfigFileName := filepath.Join(libraryWorkingDirectory, "lib.json")
	log.Println("Reading ", libraryConfigFileName)

	var localLibraryConfiguration env.LibraryProperties
	crawler.ReadLibraryConfigFileContent(libraryConfigFileName, &localLibraryConfiguration)

	if 0 == len(localLibraryConfiguration.Sources) {
		log.Fatal("THIS IS NOT YET SUPPORTED")
		crawler.ScanDirectoryForSourceFiles(args[0], &localLibraryConfiguration.Sources, false)
	}

	// since libraries do not contain the main function, use `-c` flag
	localLibraryConfiguration.Flags = append(localLibraryConfiguration.Flags, "-c")
	builder.PrepareBuildCommands(localLibraryConfiguration)
	builder.RunBuilder()

	resolveGlobalDependencies(localLibraryConfiguration, projectWorkingDirectory)
	// finally map the index and append the new config
	env.LibrariesMap[localLibraryConfiguration.Name] = len(env.LibrariesConfiguration)
	env.LibrariesConfiguration = append(env.LibrariesConfiguration, localLibraryConfiguration)

	os.Chdir(projectWorkingDirectory)
	log.Println("Changing directory back to ", projectWorkingDirectory)
}

func resolveGlobalDependencies(libraryConfiguration env.LibraryProperties, projectWorkingDirectory string) {
	// @todo this may not be needed
	// for index, path := range libraryConfiguration.Sources {
	// 	libraryConfiguration.Sources[index] = filepath.Join(projectWorkingDirectory, path)
	// }

	for index, path := range libraryConfiguration.Includes.Public {
		libraryConfiguration.Includes.Public[index] = filepath.Join(projectWorkingDirectory, path)
	}
}

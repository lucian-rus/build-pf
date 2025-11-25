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
	if err := validateNumberOfArguments(len(args), 0, 1); err != nil {
		log.Fatal(err)
	}

	// hack to not use more memory here
	if 0 == len(args) {
		args = append(args, ".")
	}

	cwd, _ := os.Getwd()
	projectWorkingDirectory := filepath.Join(cwd, args[0])
	projectConfigFileName := filepath.Join(projectWorkingDirectory, "gobi.json")

	// @todo handle case in which configuration file is not present. this shall default to a specific config
	// set based on machine
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

	// @todo properly extract generic functions here
	for _, dependency := range env.ProjectConfiguration.Dependencies.Public {
		dependencyIndex := env.LibrariesMap[dependency]
		env.ProjectConfiguration.Includes.Public = append(env.ProjectConfiguration.Includes.Public,
			env.LibrariesConfiguration[dependencyIndex].Includes.Public...)
	}

	for _, dependency := range env.ProjectConfiguration.Dependencies.Private {
		dependencyIndex := env.LibrariesMap[dependency]
		env.ProjectConfiguration.Includes.Private = append(env.ProjectConfiguration.Includes.Private,
			env.LibrariesConfiguration[dependencyIndex].Includes.Public...)
	}

	builder.PrepareBuildCommands(env.ProjectConfiguration.LibraryProperties)
	builder.RunBuilder()

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
		// @todo fix this
		crawler.ScanDirectoryForSourceFiles(args[0], &localLibraryConfiguration.Sources, false)
	}

	// since libraries do not contain the main function, use `-c` flag
	localLibraryConfiguration.Flags = append(localLibraryConfiguration.Flags, "-c")
	builder.PrepareBuildCommands(localLibraryConfiguration)
	builder.RunBuilder()

	resolveGlobalDependencies(localLibraryConfiguration, libraryWorkingDirectory)
	// finally map the index and append the new config
	env.LibrariesMap[localLibraryConfiguration.Name] = len(env.LibrariesConfiguration)
	env.LibrariesConfiguration = append(env.LibrariesConfiguration, localLibraryConfiguration)

	// @todo fix this -> this is a terrible way of doing this
	env.ProjectConfiguration.Dependencies.Libraries = append(env.ProjectConfiguration.Dependencies.Libraries,
		filepath.Join(libraryWorkingDirectory, localLibraryConfiguration.Name))

	os.Chdir(projectWorkingDirectory)
	log.Println("Changing directory back to ", projectWorkingDirectory)
}

func resolveGlobalDependencies(libraryConfiguration env.LibraryProperties, libraryWorkingDirectory string) {
	// @todo this may not be needed
	// for index, path := range libraryConfiguration.Sources {
	// 	libraryConfiguration.Sources[index] = filepath.Join(libraryWorkingDirectory, path)
	// }

	for index, path := range libraryConfiguration.Includes.Public {
		libraryConfiguration.Includes.Public[index] = filepath.Join(libraryWorkingDirectory, path)
	}
}

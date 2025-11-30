package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"gobi/components/builder"
	"gobi/components/filesystem"
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
	if err := validateNumberOfArguments(len(args), 0, 1); err != nil {
		log.Fatal(err)
	}

	// hack to not use more memory here
	if len(args) == 0 {
		args = append(args, ".")
	}

	startNow := time.Now()

	cwd, _ := os.Getwd()
	projWorkingDir := filepath.Join(cwd, args[0])
	projConfigFileName := filepath.Join(projWorkingDir, "gobi.json")

	// @todo handle case in which configuration file is not present. this shall default to a specific config
	// set based on machine
	if err := filesystem.ReadProjectConfigFileContent(projConfigFileName, &builder.ProjectConfiguration); err != nil {
		log.Fatal(err)
	}

	builder.ProjectConfiguration.ResolveSubdirPaths(projWorkingDir)
	builder.ProjectConfiguration.ResolveOutputPath(projWorkingDir)

	// step 1 -> setup filesystem for build to happen
	filesystem.SetupFilesystem(builder.ProjectConfiguration)

	// step 2 -> create lib entry for all dependencies (@todo use goroutines)
	createLibEntries(projWorkingDir)
	// step 3 -> sort based on dependency tree

	// step 4 -> build libs
	buildLibs()

	// step 5 -> build proj (building proj requires running steps 3 beforehand)
	builder.ProjectConfiguration.ResolvePrivateIncludesGlobalPaths(projWorkingDir)
	builder.ProjectConfiguration.ResolvePublicIncludesGlobalPaths(projWorkingDir)
	builder.ProjectConfiguration.ResolveSourcesGlobalPaths(projWorkingDir)
	builder.ProjectConfiguration.ResolvePrivateDependencies()
	builder.ProjectConfiguration.ResolvePublicDependencies()

	builder.ProjectConfiguration.Build()

	// step 6 -> check benchmark
	fmt.Println("this took", time.Since(startNow))
}

func createLibEntries(projWorkingDir string) {
	for _, subdir := range builder.ProjectConfiguration.Subdirectories {
		libWorkingDir := subdir
		libConfigFileName := filepath.Join(libWorkingDir, "lib.json")

		var localLibConfig builder.LibraryProperties
		filesystem.ReadLibraryConfigFile(libConfigFileName, &localLibConfig)

		// since libraries do not contain the main function, use `-c` flag
		localLibConfig.SpecifyNoMain()
		localLibConfig.ResolvePrivateIncludesGlobalPaths(libWorkingDir)
		localLibConfig.ResolvePublicIncludesGlobalPaths(libWorkingDir)
		localLibConfig.ResolveSourcesGlobalPaths(libWorkingDir)
		localLibConfig.ResolvePrivateDependencies()
		localLibConfig.ResolvePublicDependencies()

		builder.LibrariesMap[localLibConfig.Name] = len(builder.LibConfigurations)
		builder.LibConfigurations = append(builder.LibConfigurations, localLibConfig)
	}
}

func buildLibs() {
	for _, lib := range builder.LibConfigurations {
		lib.Build()
	}
}

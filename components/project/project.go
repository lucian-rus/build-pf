package project

import (
	"gobi/components/builder"
	"gobi/components/filesystem"
	"log"
	"os"
	"path/filepath"
)

var (
	projectDir string
)

func Setup() {
	//@todo this has to be properly done outside here

	projectDir, _ := os.Getwd()
	projConfigFileName := filepath.Join(projectDir, "gobi.json")

	// @todo handle case in which configuration file is not present. this shall default to a specific config
	// set based on machine
	if err := filesystem.ReadProjectConfigFileContent(projConfigFileName, &builder.ProjectConfiguration); err != nil {
		log.Fatal(err)
	}

	// @todo paths shall be resolved properly. there are limitations of how paths are resolved when using relative paths
	builder.ProjectConfiguration.ResolveSubdirPaths(projectDir)
	builder.ProjectConfiguration.ResolveOutputPath(projectDir)

	// step 1 -> setup filesystem for build to happen
	filesystem.SetupFilesystem(builder.ProjectConfiguration)

	// step 2 -> create lib entry for all dependencies (@todo use goroutines)
	createLibEntries(projectDir)
	// step 3 -> sort based on dependency tree

}

func BuildLibraries() {
	// step 4 -> build libs
	buildLibs()
}

func BuildProject() {
	// step 5 -> build proj (building proj requires running steps 3 beforehand)
	builder.ProjectConfiguration.ResolvePrivateIncludesGlobalPaths(projectDir)
	builder.ProjectConfiguration.ResolvePublicIncludesGlobalPaths(projectDir)
	builder.ProjectConfiguration.ResolveSourcesGlobalPaths(projectDir)
	builder.ProjectConfiguration.ResolvePrivateDependencies()
	builder.ProjectConfiguration.ResolvePublicDependencies()

	builder.ProjectConfiguration.Build()

}

func CleanProject() {
	os.RemoveAll(builder.ProjectConfiguration.OutputPath)
}

package project

import (
	"encoding/json"
	"fmt"
	"gobi/components/builder"
	"gobi/components/filesystem"
	"gobi/components/project/crawler"
	"log"
	"os"
	"path/filepath"
)

var (
	dependencyTree = make(map[string][]string) // this maps the "build level"

	// a build level is defined as the relative order in which the lib shall be built in order
	// to have all it's dependencies fullfilled
	buildLevelMap = make(map[string]int)
)

func LoadProjectConfiguration() {
	// @todo maybe add support for targeted build in the future. e.g `gobi build <path>`
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

	for _, item := range builder.ProjectConfiguration.Subdirectories {
		fmt.Println(item)
	}
}

func Setup() {
	crawler.ScanBuildDirectoryForLibraries(builder.ProjectConfiguration.OutputPath)
	// @todo treat defines added at project level shall be treated as project-wide

	// step 1 -> setup filesystem for build to happen
	filesystem.SetupFilesystem(builder.ProjectConfiguration)

	// step 2 -> create lib entry for all dependencies (@todo use goroutines)
	createLibEntries()
	// step 3 -> sort based on dependency tree

}

func BuildLibraries() {
	prepareLibrariesforBuild()

	// step 4 -> build libs
	for _, lib := range builder.LibConfigurations {
		// @todo update skip of library build. this is a dumb method to do it, but works for now
		if len(lib.Sources) == 0 {
			continue
		}

		lib.Build()
	}
}

func BuildProject() {
	// step 5 -> build proj (building proj requires running steps 3 beforehand)

	// @todo maybe concat this to only iterate once through the libs
	// while it loses configurability, gains a lot of speed
	// maybe allow both methods
	builder.ProjectConfiguration.ResolvePrivateIncludesGlobalPaths()
	builder.ProjectConfiguration.ResolvePublicIncludesGlobalPaths()
	builder.ProjectConfiguration.ResolveSourcesGlobalPaths()
	builder.ProjectConfiguration.ResolvePrivateDependencies()
	builder.ProjectConfiguration.ResolvePublicDependencies()

	builder.ProjectConfiguration.Build()

}

func CleanProject() {
	os.RemoveAll(builder.ProjectConfiguration.OutputPath)
}

func LoadBuildCache() error {
	reportFile := filepath.Join(builder.ProjectConfiguration.OutputPath, "report.json")
	fileContent, err := os.ReadFile(reportFile)
	if err != nil {
		log.Println("Error when trying to open", reportFile)
		return err
	}

	if err := json.Unmarshal(fileContent, &builder.BuildData); err != nil {
		log.Println("Error when unmarshalling file")
		return err
	}

	for key, value := range builder.BuildData {
		fmt.Println(key, value)
	}

	return nil
}

func CacheBuildData() error {
	reportFile := filepath.Join(builder.ProjectConfiguration.OutputPath, "report.json")
	file, err := os.Create(reportFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // pretty print
	return encoder.Encode(builder.BuildData)
}

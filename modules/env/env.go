// this package is a placeholder for now, as most functionality was moved over to the builder package
package env

import (
	"encoding/json"
	"fmt"
	"gobi/modules/cache"
	"gobi/modules/filesystem"
	"gobi/modules/library"
	"gobi/modules/project"
	"path/filepath"
)

const (
	ProjectConfigFileName = "gobi.json"
	LibConfigFileName     = "lib.json"

	// cache files
	BuildCacheFileName  = "cache/build-cache.json"
	SourceCacheFileName = "cache/source-cache.json"
	HeaderCacheFileName = "cache/header-cache.json"

	// enable debugging/printing of data
	EnableDebugData = true
)

var (
	// internal
	projectConfiguration project.ProjectProperties
	libConfigurations    = make(map[string]library.LibraryProperties)

	buildCacheMap  = make(map[string]cache.BuildCache)
	sourceCacheMap = make(map[string]cache.SourceCache)
	headerCacheMap = make(map[string]cache.HeaderCache)

	sourceFilesMap = make(map[string]cache.SourceCache)
	headerFilesMap = make(map[string]cache.HeaderCache)

	skipBuildPhase = false
)

func Setup() {
	loadprojectConfiguration()
	loadLibraryConfigurations()

	loadBuildCache()
	loadSourceCache()
	loadHeaderCache()

	// @todo get a way to fix files that have the same name.
	// some projects may have multiple files having the same name
	// this will cause key conflicts. check how to fix
	scanEnvForSourceFiles()
	scanEnvForHeaderFiles()

	parseLibraryConfigurations()

	// handle eveything required for build
	prepareLibrariesforBuild()
	prepareProjectForBuild()

	runIncrementalBuildChecks()
	// @todo this is not done ideally and should be modified. once the dependency system is up and running, replace this.
	if skipBuildPhase {
		fmt.Println("Nothing to be done. Skipping...")
	} else {
		runCommandCreator()
	}

	// after loading is done, start creating required directories
	filesystem.CreateDirectory(projectConfiguration.OutputPath)
	filesystem.CreateDirectory(filepath.Join(projectConfiguration.OutputPath, "cache"))
	filesystem.CreateDirectory(filepath.Join(projectConfiguration.OutputPath, "libs"))
}

// @todo check if this actually works as intended
func CacheData() error {
	data, err := json.MarshalIndent(buildCacheMap, "", "  ")
	if err != nil {
		return err
	}
	cacheFilePath := filepath.Join(projectConfiguration.OutputPath, BuildCacheFileName)
	filesystem.WriteDataToJson(data, cacheFilePath)

	data, err = json.MarshalIndent(sourceFilesMap, "", "  ")
	if err != nil {
		return err
	}
	cacheFilePath = filepath.Join(projectConfiguration.OutputPath, SourceCacheFileName)
	filesystem.WriteDataToJson(data, cacheFilePath)

	data, err = json.MarshalIndent(headerFilesMap, "", "  ")
	if err != nil {
		return err
	}
	cacheFilePath = filepath.Join(projectConfiguration.OutputPath, HeaderCacheFileName)
	filesystem.WriteDataToJson(data, cacheFilePath)

	return nil
}

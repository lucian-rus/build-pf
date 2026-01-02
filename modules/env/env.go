// this package is a placeholder for now, as most functionality was moved over to the builder package
package env

import (
	"encoding/json"
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
	// public
	ProjectConfiguration project.ProjectProperties

	LibConfigurations = make(map[string]library.LibraryProperties)

	// internal
	BuildCacheMap  = make(map[string]cache.BuildCache)
	SourceCacheMap = make(map[string]cache.SourceCache)
	HeaderCacheMap = make(map[string]cache.HeaderCache)
)

func Setup() {
	loadProjectConfiguration()
	loadBuildCache()

	// @todo get a way to fix files that have the same name.
	// some projects may have multiple files having the same name
	// this will cause key conflicts. check how to fix
	scanEnvForSourceFiles()
	scanEnvForHeaderFiles()

	// load library configuration AFTER the cache in order to avoid unnecessary crawling
	loadLibraryConfigurations()
	parseLibraryConfigurations()

	// handle eveything required for build
	prepareLibrariesforBuild()
	prepareProjectForBuild()

	// after loading is done, start creating required directories
	filesystem.CreateDirectory(ProjectConfiguration.OutputPath)
	filesystem.CreateDirectory(filepath.Join(ProjectConfiguration.OutputPath, "cache"))
	filesystem.CreateDirectory(filepath.Join(ProjectConfiguration.OutputPath, "libs"))
}

// @todo check if this actually works as intended
func CacheData() error {
	data, err := json.MarshalIndent(BuildCacheMap, "", "  ")
	if err != nil {
		return err
	}
	cacheFilePath := filepath.Join(ProjectConfiguration.OutputPath, BuildCacheFileName)
	filesystem.WriteDataToJson(data, cacheFilePath)

	data, err = json.MarshalIndent(SourceCacheMap, "", "  ")
	if err != nil {
		return err
	}
	cacheFilePath = filepath.Join(ProjectConfiguration.OutputPath, SourceCacheFileName)
	filesystem.WriteDataToJson(data, cacheFilePath)

	data, err = json.MarshalIndent(HeaderCacheMap, "", "  ")
	if err != nil {
		return err
	}
	cacheFilePath = filepath.Join(ProjectConfiguration.OutputPath, HeaderCacheFileName)
	filesystem.WriteDataToJson(data, cacheFilePath)

	return nil
}

// this package is a placeholder for now, as most functionality was moved over to the builder package
package env

import (
	"encoding/json"
	"gobi/modules/cache"
	"gobi/modules/filesystem"
	"gobi/modules/library"
	"gobi/modules/project"
	"os"
)

const (
	ProjectConfigFileName = "gobi.json"
	LibConfigFileName     = "lib.json"
	CacheConfigFileName   = "cache.json"

	// enable debugging/printing of data
	EnableDebugData = true
)

var (
	// public
	ProjectConfiguration project.ProjectProperties

	LibConfigurations = make(map[string]library.LibraryProperties)
	BuildCacheMap     = make(map[string]cache.BuildCache)

// internal
)

func Setup() {
	loadProjectConfiguration()
	loadBuildCache()

	// load library configuration AFTER the cache in order to avoid unnecessary crawling
	loadLibraryConfigurations()

	// handle eveything required for build
	prepareLibrariesforBuild()
	prepareProjectForBuild()

	// after loading is done, start creating required directories
	filesystem.CreateDirectory(ProjectConfiguration.OutputPath)
}

// @todo check if this actually works as intended
func CacheBuildData() {
	data, err := json.MarshalIndent(BuildCacheMap, "", "  ")
	if err != nil {
		// Optionally log or handle the error
		return
	}
	f, err := os.Create(CacheConfigFileName)
	if err != nil {
		// Optionally log or handle the error
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		// Optionally log or handle the error
		return
	}
}

// this package is a placeholder for now, as most functionality was moved over to the builder package
package env

import (
	"encoding/json"
	"gobi/modules/cache"
	"gobi/modules/filesystem"
	"gobi/modules/library"
	"gobi/modules/project"
	"os"
	"path/filepath"
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
func CacheBuildData() error {
	data, err := json.MarshalIndent(BuildCacheMap, "", "  ")
	if err != nil {
		return err
	}

	cacheFilePath := filepath.Join(ProjectConfiguration.OutputPath, CacheConfigFileName)
	file, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

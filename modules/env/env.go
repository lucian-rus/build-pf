// this package is a placeholder for now, as most functionality was moved over to the builder package
package env

import (
	"fmt"
	"gobi/modules/cache"
	"gobi/modules/env/crawler"
	"gobi/modules/filesystem"
	"gobi/modules/library"
	"gobi/modules/project"
	"os"
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
	loadProjectConfiguration()

	// scanEnvForSourceFiles()
	// scanEnvForHeaderFiles()

	if err := crawler.ScanDirectoryForFiles(".", &projectConfiguration.Sources, ".c"); err != nil {
		fmt.Println("could not read sources")
	}

	var headerFiles []string
	if err := crawler.ScanDirectoryForFiles(".", &headerFiles, ".h"); err != nil {
		fmt.Println("could not read headers")
	}

	headerDirs := make(map[string]bool)
	projectConfiguration.Includes.Private = nil
	for _, header := range headerFiles {
		dir := filepath.Dir(header)
		if !headerDirs[dir] {
			headerDirs[dir] = true
			projectConfiguration.Includes.Private = append(projectConfiguration.Includes.Private, dir)
		}
	}

	prepareProjectForBuild()
	runCommandCreator()

	filesystem.CreateDirectory(projectConfiguration.OutputPath)
}

func Cleanup() {
	os.RemoveAll(projectConfiguration.OutputPath)
	// os.RemoveAll(projectConfiguration.LogPath)
}

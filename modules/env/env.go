// this package is a placeholder for now, as most functionality was moved over to the builder package
package env

import (
	"gobi/modules/library"
	"gobi/modules/project"
)

const (
	ProjectConfigFileName = "gobi.json"
	LibConfigFileName     = "lib.json"

	// enable debugging/printing of data
	EnableDebugData = true
)

var (
	ProjectConfiguration project.ProjectProperties // this HAS to be unique

	LibConfigurations = make(
		map[string]library.LibraryProperties,
	) // meant to hold the configurations instead of the array

// internal
)

func Setup() {
	loadProjectConfiguration()
	loadBuildCache()

	// load library configuration AFTER the cache in order to avoid unnecessary crawling
	loadLibraryConfigurations()
	prepareLibrariesforBuild()
	prepareProjectForBuild()

	// after loading is done, start creating required directories

}

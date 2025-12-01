package project

import (
	"gobi/components/builder"
	"gobi/components/filesystem"
	"path/filepath"
)

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

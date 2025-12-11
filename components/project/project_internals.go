package project

import (
	"fmt"
	"gobi/components/builder"
	"gobi/components/crawler"
	"gobi/components/env"
	"gobi/components/filesystem"
	"gobi/components/parser"
	"path/filepath"
)

func createLibEntries() {
	for _, subdir := range builder.ProjectConfiguration.Subdirectories {
		libConfigFileName := filepath.Join(subdir, "lib.json")

		var localLibConfig builder.LibraryProperties
		localLibConfig.Root, _ = filepath.Abs(subdir)
		filesystem.ReadLibraryConfigFile(libConfigFileName, &localLibConfig)

		if len(localLibConfig.Sources) == 0 {
			crawler.ScanDirectoryForSources(localLibConfig.Root, &localLibConfig.Sources)
		} else {
			localLibConfig.ResolveSourcesGlobalPaths()
		}

		// @todo update this to properly handle includes -> extract required files and analyse the libs
		// for the available header files
		for _, source := range localLibConfig.Sources {
			parser.GetIncludesList(source)
		}
		crawler.ScanDirectoryForHeaders(localLibConfig.Root, &localLibConfig.Headers)

		builder.LibConfigurations[localLibConfig.Name] = localLibConfig
	}
}

func resolveLibraries() {
	for _, lib := range builder.LibConfigurations {
		// since libraries do not contain the main function, use `-c` flag
		lib.SpecifyNoMain()
		lib.ResolvePrivateIncludesGlobalPaths()
		lib.ResolvePublicIncludesGlobalPaths()
		lib.ResolvePrivateDependencies()
		lib.ResolvePublicDependencies()

		dependencyTree[lib.Name] = append(dependencyTree[lib.Name], lib.Dependencies.Public...)
		dependencyTree[lib.Name] = append(dependencyTree[lib.Name], lib.Dependencies.Private...)

		builder.LibConfigurations[lib.Name] = lib // update the map
	}

	// return here to not debug anymore
	if !env.EnableDebugData {
		return
	}

	for _, lib := range builder.LibConfigurations {
		fmt.Println(lib.Name)
		fmt.Println("	* private")
		for _, item := range lib.Includes.Private {
			fmt.Println("	- ", item)
		}

		fmt.Println("	* public")
		for _, item := range lib.Includes.Public {
			fmt.Println("	- ", item)
		}
	}

	fmt.Println("dependency list:")

	for key, depList := range dependencyTree {
		fmt.Printf("	* %s :", key)
		fmt.Println(depList)
	}
}

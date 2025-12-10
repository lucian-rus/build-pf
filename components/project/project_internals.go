package project

import (
	"fmt"
	"gobi/components/builder"
	"gobi/components/filesystem"
	"path/filepath"
)

func createLibEntries(projWorkingDir string) {
	for _, subdir := range builder.ProjectConfiguration.Subdirectories {
		libConfigFileName := filepath.Join(subdir, "lib.json")

		var localLibConfig builder.LibraryProperties
		filesystem.ReadLibraryConfigFile(libConfigFileName, &localLibConfig)

		// since libraries do not contain the main function, use `-c` flag
		localLibConfig.SpecifyNoMain()
		localLibConfig.ResolvePrivateIncludesGlobalPaths(subdir)
		localLibConfig.ResolvePublicIncludesGlobalPaths(subdir)
		localLibConfig.ResolveSourcesGlobalPaths(subdir)
		localLibConfig.ResolvePrivateDependencies()
		localLibConfig.ResolvePublicDependencies()

		builder.LibConfigurations[localLibConfig.Name] = localLibConfig

		dependencyTree[localLibConfig.Name] = append(dependencyTree[localLibConfig.Name], localLibConfig.Dependencies.Public...)
		dependencyTree[localLibConfig.Name] = append(dependencyTree[localLibConfig.Name], localLibConfig.Dependencies.Private...)
	}

	for _, lib := range builder.LibConfigurations {
		fmt.Println(lib.Name)

		fmt.Println("	* private")
		for _, item := range lib.Includes.Private {
			fmt.Println(item)
		}

		fmt.Println("	* public")
		for _, item := range lib.Includes.Public {
			fmt.Println(item)
		}
	}

	fmt.Println("---------------------")
	fmt.Println("dependency list:")

	for key, depList := range dependencyTree {
		fmt.Printf("	* %s :", key)
		fmt.Println(depList)
	}
}

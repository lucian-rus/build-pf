package env

import (
	"encoding/json"
	"fmt"
	"gobi/modules/builder"
	"gobi/modules/cache"
	"gobi/modules/env/crawler"
	"gobi/modules/filesystem"
	"gobi/modules/library"
	"log"
	"os"
	"path/filepath"
)

func loadProjectConfiguration() error {
	fmt.Println("--------------- loading project --------------------")
	projectDir, _ := os.Getwd()
	projConfigFileName := filepath.Join(projectDir, ProjectConfigFileName)

	fileContent, err := filesystem.ReadJsonConfigFile(projConfigFileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileContent, &ProjectConfiguration); err != nil {
		log.Println("Error when unmarshalling JSON file", projConfigFileName)
		return err
	}

	// @todo update the formatter
	if EnableDebugData {
		fmt.Println("name of project:		", ProjectConfiguration.Name)
		fmt.Println("list of private includes:	", ProjectConfiguration.Includes.Private)
		fmt.Println("list of public includes:	", ProjectConfiguration.Includes.Public)
		fmt.Println("list of private dependencies:	", ProjectConfiguration.Dependencies.Private)
		fmt.Println("list of public dependencies:	", ProjectConfiguration.Dependencies.Public)
	}

	ProjectConfiguration.ResolveSubdirPaths(projectDir)
	ProjectConfiguration.ResolveOutputPath(projectDir)

	return nil
}

func loadBuildCache() error {
	fmt.Println("---------------- loading cache ---------------------")
	cacheFilePath := filepath.Join(ProjectConfiguration.OutputPath, CacheConfigFileName)

	fileContent, err := filesystem.ReadJsonConfigFile(cacheFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileContent, &BuildCacheMap); err != nil {
		log.Println("Error when unmarshalling JSON file", cacheFilePath)
		return err
	}

	return nil
}

func loadLibraryConfigurations() error {
	fmt.Println("-------------- loading libraries -------------------")

	for _, subdir := range ProjectConfiguration.Subdirectories {
		libConfigFileName := filepath.Join(subdir, LibConfigFileName)

		var localLibConfig library.LibraryProperties
		fileContent, err := filesystem.ReadJsonConfigFile(libConfigFileName)
		if err := json.Unmarshal(fileContent, &localLibConfig); err != nil {
			log.Println("Error when unmarshalling JSON file", libConfigFileName)
			return err
		}

		localLibConfig.Root, _ = filepath.Abs(subdir)
		if err != nil {
			return err
		}

		// @todo first check if binary exists. if so, only then do the source check
		if len(localLibConfig.Sources) == 0 {
			crawler.ScanDirectoryForSources(localLibConfig.Root, &localLibConfig.Sources)
		} else {
			localLibConfig.ResolveSourcesGlobalPaths()
		}

		for _, source := range localLibConfig.Sources {
			var aux int

			crawler.GetTimestampForFile(source, &aux)
			BuildCacheMap[source] = cache.BuildCache{
				Name:      source,
				Timestamp: aux,
			}
		}
		// @todo update this to properly handle includes -> extract required files and analyse the libs
		// for the available header files
		// for _, source := range localLibConfig.Sources {
		// 	parser.GetIncludesList(source)
		// }
		crawler.ScanDirectoryForHeaders(localLibConfig.Root, &localLibConfig.Headers)

		LibConfigurations[localLibConfig.Name] = localLibConfig
	}

	return nil
}

func prepareLibrariesforBuild() {
	fmt.Println("-------------- baking libraries --------------------")

	for _, lib := range LibConfigurations {
		// since libraries do not contain the main function, use `-c` flag
		lib.SpecifyNoMain()
		lib.ResolvePrivateIncludesGlobalPaths()
		lib.ResolvePublicIncludesGlobalPaths()
		lib.ResolvePrivateDependencies(ProjectConfiguration.OutputPath, LibConfigurations)
		lib.ResolvePublicDependencies(ProjectConfiguration.OutputPath, LibConfigurations)

		LibConfigurations[lib.Name] = lib // update the map

		commandList := createCommandSequence(lib)
		builder.AddBuildSequence(commandList)
	}

	for _, lib := range LibConfigurations {
		printLibraryDebugData(lib)
	}

	// fmt.Println("dependency list:")
	// for key, depList := range dependencyTree {
	// 	fmt.Printf("	* %s :", key)
	// 	fmt.Println(depList)
	// }
}

func prepareProjectForBuild() {
	fmt.Println("--------------- baking project ---------------------")

	ProjectConfiguration.ResolvePrivateIncludesGlobalPaths()
	ProjectConfiguration.ResolvePublicIncludesGlobalPaths()
	ProjectConfiguration.ResolvePrivateDependencies(ProjectConfiguration.OutputPath, LibConfigurations)
	ProjectConfiguration.ResolvePublicDependencies(ProjectConfiguration.OutputPath, LibConfigurations)

	commandList := createCommandSequence(ProjectConfiguration.LibraryProperties)
	builder.AddBuildSequence(commandList)

	printLibraryDebugData(ProjectConfiguration.LibraryProperties)
}

func createCommandSequence(lib library.LibraryProperties) []string {
	var commandList []string
	// append compiler
	commandList = append(commandList, ProjectConfiguration.Compiler)
	commandList = append(commandList, lib.Flags...)

	// append definitions
	for _, item := range lib.Defines {
		parsedArgument := "-D" + item
		commandList = append(commandList, parsedArgument)
	}

	// append includes
	for _, item := range lib.Includes.Public {
		parsedArgument := "-I" + item
		commandList = append(commandList, parsedArgument)
	}

	for _, item := range lib.Includes.Private {
		parsedArgument := "-I" + item
		commandList = append(commandList, parsedArgument)
	}

	// append output
	commandList = append(commandList, "-o")
	// @todo need to resolve global dependency for output
	objOutputPath := filepath.Join(ProjectConfiguration.OutputPath, lib.Name)
	commandList = append(commandList, objOutputPath)

	// append sources and dependencies
	commandList = append(commandList, lib.Sources...)
	commandList = append(commandList, lib.LinkedObjects...)
	return commandList
}

func printLibraryDebugData(lib library.LibraryProperties) {
	// return here to not debug anymore
	if !EnableDebugData {
		return
	}

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

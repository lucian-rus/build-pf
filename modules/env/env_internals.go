package env

import (
	"fmt"
	"gobi/modules/builder"
	"gobi/modules/cache"
	"gobi/modules/env/crawler"
	"gobi/modules/library"
	"path/filepath"
	"time"
)

func prepareProjectForBuild() {
	fmt.Println("--------------- baking project ---------------------")
	fmt.Println(" - Baking", projectConfiguration.Name)

	projectConfiguration.ResolvePrivateIncludesGlobalPaths()
	projectConfiguration.ResolvePublicIncludesGlobalPaths()
	projectConfiguration.ResolvePrivateDependencies(projectConfiguration.OutputDirPath, libConfigurations)
	projectConfiguration.ResolvePublicDependencies(projectConfiguration.OutputDirPath, libConfigurations)
	// unlike libraries, do this here, as libraries are sent to `libs` dir
	// @todo check if there is a better way to do it
	projectConfiguration.OutputPath = filepath.Join(projectConfiguration.OutputDirPath, projectConfiguration.Name)
}

// @todo account for the project checks as well
func runCommandCreator() {
	for _, lib := range libConfigurations {
		commandList := createCommandSequence(lib)
		builder.AddBuildSequence(commandList)
	}

	for _, lib := range libConfigurations {
		printLibraryDebugData(lib)
	}

	commandList := createCommandSequence(projectConfiguration.LibraryProperties)
	builder.AddBuildSequence(commandList)

	printLibraryDebugData(projectConfiguration.LibraryProperties)
}

func createCommandSequence(lib library.LibraryProperties) []string {
	var commandList []string
	// append compiler
	commandList = append(commandList, projectConfiguration.Compiler)
	commandList = append(commandList, lib.Flags.Public...)
	commandList = append(commandList, lib.Flags.Private...)

	// append definitions
	for _, item := range lib.Defines.Public {
		parsedArgument := "-D" + item
		commandList = append(commandList, parsedArgument)
	}
	for _, item := range lib.Defines.Private {
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
	commandList = append(commandList, lib.OutputPath)

	// append sources and dependencies
	commandList = append(commandList, lib.Sources...)

	buildCacheMap[lib.Name] = cache.BuildCache{
		FileCache: cache.FileCache{
			Timestamp: int(time.Now().Unix()),
			Path:      lib.OutputPath,
		},
	}

	return commandList
}

func updatesourceFilesMap(lib library.LibraryProperties) {
	for _, source := range lib.Sources {
		var timestamp int
		crawler.GetTimestampForFile(source, &timestamp)

		// update source cache
		sourceFilesMap[filepath.Base(source)] = cache.SourceCache{
			FileCache: cache.FileCache{
				Timestamp: timestamp,
				Path:      source,
			},
			Library: lib.Name,
		}
	}
}

func doFileTimestampsMatch(lib library.LibraryProperties) bool {
	for _, source := range lib.Sources {
		var liveTimestamp int
		var cachedTimestamp int

		sourceBaseName := filepath.Base(source)
		if _, ok := sourceFilesMap[sourceBaseName]; ok {
			liveTimestamp = sourceFilesMap[sourceBaseName].Timestamp
		}

		if _, ok := sourceCacheMap[sourceBaseName]; ok {
			cachedTimestamp = sourceCacheMap[sourceBaseName].Timestamp
		}

		if liveTimestamp != cachedTimestamp {
			return false
		}
	}

	return true
}

func printLibraryDebugData(lib library.LibraryProperties) {
	// return here to not debug anymore
	if !EnableDebugData {
		return
	}

	fmt.Println(lib.Name)

	fmt.Println("	* sources")
	for _, item := range lib.Sources {
		fmt.Println("	- ", item)
	}

	fmt.Println("	* private")
	for _, item := range lib.Includes.Private {
		fmt.Println("	- ", item)
	}

	fmt.Println("	* public")
	for _, item := range lib.Includes.Public {
		fmt.Println("	- ", item)
	}
}

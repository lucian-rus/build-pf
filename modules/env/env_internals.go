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

func scanEnvForSourceFiles() {
	fmt.Println("-------------- scanning sources --------------------")

	var sourceList []string
	crawler.ScanDirectoryForFiles(".", &sourceList, ".c")

	for _, source := range sourceList {
		var lastEditTimestamp int
		crawler.GetTimestampForFile(source, &lastEditTimestamp)

		sourceFilesMap[filepath.Base(source)] = cache.SourceCache{
			FileCache: cache.FileCache{
				Path:      source,
				Timestamp: lastEditTimestamp,
			},
		}
	}

	// if debug data is disabled, do not print data
	if !EnableDebugData {
		return
	}

	fmt.Println("Source map:")
	for key, value := range sourceFilesMap {
		fmt.Println(key, value)
	}
}

func scanEnvForHeaderFiles() {
	fmt.Println("-------------- scanning headers --------------------")

	var headerList []string
	crawler.ScanDirectoryForFiles(".", &headerList, ".h")

	for _, header := range headerList {
		var lastEditTimestamp int
		crawler.GetTimestampForFile(header, &lastEditTimestamp)

		headerFilesMap[filepath.Base(header)] = cache.HeaderCache{
			FileCache: cache.FileCache{
				Path:      header,
				Timestamp: lastEditTimestamp,
			},
		}
	}

	// if debug data is disabled, do not print data
	if !EnableDebugData {
		return
	}

	fmt.Println("Header map:")
	for key, value := range headerFilesMap {
		fmt.Println(key, value)
	}
}

func parseLibraryConfigurations() {
	// @todo check if this can be optimized
	for _, lib := range libConfigurations {
		updatesourceFilesMap(lib)
		updateHeadersFilesMap(lib)
	}
}

func prepareLibrariesforBuild() {
	fmt.Println("-------------- baking libraries --------------------")

	for _, lib := range libConfigurations {
		fmt.Println(" - Baking", lib.Name)
		// since libraries do not contain the main function, use `-c` flag
		lib.SpecifyNoMain()
		lib.ResolvePrivateIncludesGlobalPaths()
		lib.ResolvePublicIncludesGlobalPaths()
		lib.ResolvePrivateDependencies(projectConfiguration.OutputPath, libConfigurations)
		lib.ResolvePublicDependencies(projectConfiguration.OutputPath, libConfigurations)
		lib.ResolveObjectPath(projectConfiguration.OutputPath)

		lib.InheritProjectDefines(projectConfiguration.LibraryProperties)
		lib.InheritProjectFlags(projectConfiguration.LibraryProperties)

		libConfigurations[lib.Name] = lib // update the map
	}
}

func prepareProjectForBuild() {
	fmt.Println("--------------- baking project ---------------------")
	fmt.Println(" - Baking", projectConfiguration.Name)

	projectConfiguration.ResolvePrivateIncludesGlobalPaths()
	projectConfiguration.ResolvePublicIncludesGlobalPaths()
	projectConfiguration.ResolvePrivateDependencies(projectConfiguration.OutputPath, libConfigurations)
	projectConfiguration.ResolvePublicDependencies(projectConfiguration.OutputPath, libConfigurations)
	// unlike libraries, do this here, as libraries are sent to `libs` dir
	// @todo check if there is a better way to do it
	projectConfiguration.ObjectPath = filepath.Join(projectConfiguration.OutputPath, projectConfiguration.Name)
}

// @todo account for the project checks as well
func runIncrementalBuildChecks() {
	var libsToBeSkipped []string

	for key, lib := range libConfigurations {
		var libPreviouslyBuilt bool
		if _, ok := buildCacheMap[lib.Name]; ok {
			libPreviouslyBuilt = true
		}

		cachedTimestampMatch := doFileTimestampsMatch(lib)
		if libPreviouslyBuilt && cachedTimestampMatch {
			libsToBeSkipped = append(libsToBeSkipped, key)
		}
	}
	// run the delete sequence
	for _, lib := range libsToBeSkipped {
		delete(libConfigurations, lib)
	}

	// do the same for the project configuration
	var libPreviouslyBuilt bool
	if _, ok := buildCacheMap[projectConfiguration.Name]; ok {
		libPreviouslyBuilt = true
	}

	cachedTimestampMatch := doFileTimestampsMatch(projectConfiguration.LibraryProperties)
	if libPreviouslyBuilt && cachedTimestampMatch && (len(libConfigurations) == 0) {
		skipBuildPhase = true
	}

}

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
	commandList = append(commandList, lib.ObjectPath)

	// append sources and dependencies
	commandList = append(commandList, lib.Sources...)
	commandList = append(commandList, lib.LinkedObjects...)

	buildCacheMap[lib.Name] = cache.BuildCache{
		FileCache: cache.FileCache{
			Timestamp: int(time.Now().Unix()),
			Path:      lib.ObjectPath,
		},
	}

	return commandList
}

func updateHeadersFilesMap(lib library.LibraryProperties) {
	for _, header := range lib.Headers {
		var timestamp int
		crawler.GetTimestampForFile(header, &timestamp)

		// update header cache
		headerFilesMap[filepath.Base(header)] = cache.HeaderCache{
			FileCache: cache.FileCache{
				Timestamp: timestamp,
				Path:      header,
			},
			Library: lib.Name,
		}
	}
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
	fmt.Println("	* private")
	for _, item := range lib.Includes.Private {
		fmt.Println("	- ", item)
	}

	fmt.Println("	* public")
	for _, item := range lib.Includes.Public {
		fmt.Println("	- ", item)
	}
}

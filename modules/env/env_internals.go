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

		SourceFilesMap[filepath.Base(source)] = cache.SourceCache{
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
	for key, value := range SourceFilesMap {
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

		HeaderFilesMap[filepath.Base(header)] = cache.HeaderCache{
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
	for key, value := range HeaderFilesMap {
		fmt.Println(key, value)
	}
}

func parseLibraryConfigurations() {
	// @todo check if this can be optimized
	for _, lib := range LibConfigurations {
		updateSourceFilesMap(lib)
		updateHeadersFilesMap(lib)
	}
}

func prepareLibrariesforBuild() {
	fmt.Println("-------------- baking libraries --------------------")

	if len(LibConfigurations) == 0 {
		fmt.Println("Nothing to do. Skip...")
	}

	for _, lib := range LibConfigurations {
		// since libraries do not contain the main function, use `-c` flag
		lib.SpecifyNoMain()
		lib.ResolvePrivateIncludesGlobalPaths()
		lib.ResolvePublicIncludesGlobalPaths()
		lib.ResolvePrivateDependencies(ProjectConfiguration.OutputPath, LibConfigurations)
		lib.ResolvePublicDependencies(ProjectConfiguration.OutputPath, LibConfigurations)
		lib.ResolveObjectPath(ProjectConfiguration.OutputPath)

		lib.InheritProjectDefines(ProjectConfiguration.LibraryProperties)
		lib.InheritProjectFlags(ProjectConfiguration.LibraryProperties)

		LibConfigurations[lib.Name] = lib // update the map
	}
}

func prepareProjectForBuild() {
	fmt.Println("--------------- baking project ---------------------")

	ProjectConfiguration.ResolvePrivateIncludesGlobalPaths()
	ProjectConfiguration.ResolvePublicIncludesGlobalPaths()
	ProjectConfiguration.ResolvePrivateDependencies(ProjectConfiguration.OutputPath, LibConfigurations)
	ProjectConfiguration.ResolvePublicDependencies(ProjectConfiguration.OutputPath, LibConfigurations)
	// unlike libraries, do this here, as libraries are sent to `libs` dir
	// @todo check if there is a better way to do it
	ProjectConfiguration.ObjectPath = filepath.Join(ProjectConfiguration.OutputPath, ProjectConfiguration.Name)
}

// @todo account for the project checks as well
func runIncrementalBuildChecks() {
	var libsToBeSkipped []string

	for key, lib := range LibConfigurations {
		var libPreviouslyBuilt bool
		if _, ok := BuildCacheMap[lib.Name]; ok {
			libPreviouslyBuilt = true
		}

		cachedTimestampMatch := doFileTimestampsMatch(lib)
		if libPreviouslyBuilt && cachedTimestampMatch {
			libsToBeSkipped = append(libsToBeSkipped, key)
		}
	}

	for _, lib := range libsToBeSkipped {
		delete(LibConfigurations, lib)
	}
}

func runCommandCreator() {
	for _, lib := range LibConfigurations {
		commandList := createCommandSequence(lib)
		builder.AddBuildSequence(commandList)
	}

	for _, lib := range LibConfigurations {
		printLibraryDebugData(lib)
	}

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
	commandList = append(commandList, lib.ObjectPath)

	// append sources and dependencies
	commandList = append(commandList, lib.Sources...)
	commandList = append(commandList, lib.LinkedObjects...)

	BuildCacheMap[lib.Name] = cache.BuildCache{
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
		HeaderFilesMap[filepath.Base(header)] = cache.HeaderCache{
			FileCache: cache.FileCache{
				Timestamp: timestamp,
				Path:      header,
			},
			Library: lib.Name,
		}
	}
}

func updateSourceFilesMap(lib library.LibraryProperties) {
	for _, source := range lib.Sources {
		var timestamp int
		crawler.GetTimestampForFile(source, &timestamp)

		// update source cache
		SourceFilesMap[filepath.Base(source)] = cache.SourceCache{
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
		if _, ok := SourceFilesMap[sourceBaseName]; ok {
			liveTimestamp = SourceFilesMap[sourceBaseName].Timestamp
		}

		if _, ok := SourceCacheMap[sourceBaseName]; ok {
			cachedTimestamp = SourceCacheMap[sourceBaseName].Timestamp
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

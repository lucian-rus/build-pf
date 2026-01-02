package env

import (
	"fmt"
	"gobi/modules/builder"
	"gobi/modules/cache"
	"gobi/modules/env/crawler"
	"gobi/modules/library"
	"path/filepath"
)

func scanEnvForSourceFiles() {
	fmt.Println("-------------- scanning sources --------------------")

	var sourceList []string
	crawler.ScanDirectoryForFiles(".", &sourceList, ".c")

	for _, source := range sourceList {
		var lastEditTimestamp int
		crawler.GetTimestampForFile(source, &lastEditTimestamp)

		SourceCacheMap[filepath.Base(source)] = cache.SourceCache{
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
	for key, value := range SourceCacheMap {
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

		HeaderCacheMap[filepath.Base(header)] = cache.HeaderCache{
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
	for key, value := range HeaderCacheMap {
		fmt.Println(key, value)
	}
}

func parseLibraryConfigurations() {
	for _, lib := range LibConfigurations {
		for _, source := range lib.Sources {
			var timestamp int
			crawler.GetTimestampForFile(source, &timestamp)

			// update source cache
			SourceCacheMap[filepath.Base(source)] = cache.SourceCache{
				FileCache: cache.FileCache{
					Timestamp: timestamp,
					Path:      source,
				},
				Library: lib.Name,
			}
		}

		for _, header := range lib.Headers {
			var timestamp int
			crawler.GetTimestampForFile(header, &timestamp)

			// update header cache
			HeaderCacheMap[filepath.Base(header)] = cache.HeaderCache{
				FileCache: cache.FileCache{
					Timestamp: timestamp,
					Path:      header,
				},
				Library: lib.Name,
			}
		}
	}
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

		lib.InheritProjectDefines(ProjectConfiguration.LibraryProperties)
		lib.InheritProjectFlags(ProjectConfiguration.LibraryProperties)

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

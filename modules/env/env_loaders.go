// this file shall contain the loaders -> functions that load cached data
package env

import (
	"encoding/json"
	"fmt"
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

func loadLibraryConfigurations() error {
	fmt.Println("-------------- loading libraries -------------------")

	for _, subdir := range ProjectConfiguration.Subdirectories {
		libConfigFileName := filepath.Join(subdir, LibConfigFileName)

		// declare and init default values
		var localLibConfig library.LibraryProperties
		localLibConfig.SetDefaultValues()

		fileContent, err := filesystem.ReadJsonConfigFile(libConfigFileName)
		if err := json.Unmarshal(fileContent, &localLibConfig); err != nil {
			log.Println("Error when unmarshalling JSON file", libConfigFileName)
			return err
		}

		localLibConfig.Root, _ = filepath.Abs(subdir)
		if err != nil {
			return err
		}

		// @todo first check if binary exists. if so, only then do the source check - this is done in a dumb way
		if len(localLibConfig.Sources) == 0 {
			crawler.ScanDirectoryForFiles(localLibConfig.Root, &localLibConfig.Sources, ".c")
		} else {
			localLibConfig.ResolveSourcesGlobalPaths()
		}

		crawler.ScanDirectoryForFiles(localLibConfig.Root, &localLibConfig.Headers, ".h")
		LibConfigurations[localLibConfig.Name] = localLibConfig
	}

	return nil
}

func loadBuildCache() error {
	cacheFilePath := filepath.Join(ProjectConfiguration.OutputPath, BuildCacheFileName)

	fileContent, err := filesystem.ReadJsonConfigFile(cacheFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileContent, &BuildCacheMap); err != nil {
		log.Println("Error when unmarshalling JSON file", cacheFilePath)
		return err
	}

	for key, value := range BuildCacheMap {
		fmt.Println(key, value)
	}

	return nil
}

func loadSourceCache() error {
	fmt.Println("------------ loading source cache ------------------")

	return nil
}

func loadHeaderCache() error {
	fmt.Println("------------ loading header cache ------------------")

	return nil
}

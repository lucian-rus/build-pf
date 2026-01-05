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

func loadprojectConfiguration() error {
	fmt.Println("--------------- loading project --------------------")
	projectDir, _ := os.Getwd()
	projConfigFileName := filepath.Join(projectDir, ProjectConfigFileName)

	fileContent, err := filesystem.ReadJsonConfigFile(projConfigFileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileContent, &projectConfiguration); err != nil {
		log.Println("Error when unmarshalling JSON file", projConfigFileName)
		return err
	}

	// @todo update the formatter
	if EnableDebugData {
		fmt.Println("name of project:		", projectConfiguration.Name)
		fmt.Println("list of private includes:	", projectConfiguration.Includes.Private)
		fmt.Println("list of public includes:	", projectConfiguration.Includes.Public)
		fmt.Println("list of private dependencies:	", projectConfiguration.Dependencies.Private)
		fmt.Println("list of public dependencies:	", projectConfiguration.Dependencies.Public)
	}

	projectConfiguration.ResolveSubdirPaths(projectDir)
	projectConfiguration.ResolveOutputPath(projectDir)

	return nil
}

func loadLibraryConfigurations() error {
	fmt.Println("-------------- loading libraries -------------------")

	for _, subdir := range projectConfiguration.Subdirectories {
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
		libConfigurations[localLibConfig.Name] = localLibConfig
	}

	return nil
}

// @todo refactor these functions to extract common code
func loadBuildCache() error {
	fmt.Println("------------- loading build cache ------------------")
	cacheFilePath := filepath.Join(projectConfiguration.OutputPath, BuildCacheFileName)

	fileContent, err := filesystem.ReadJsonConfigFile(cacheFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileContent, &buildCacheMap); err != nil {
		log.Println("Error when unmarshalling JSON file", cacheFilePath)
		return err
	}

	for key, value := range buildCacheMap {
		fmt.Println(key, value)
	}

	return nil
}

func loadSourceCache() error {
	fmt.Println("------------ loading source cache ------------------")
	cacheFilePath := filepath.Join(projectConfiguration.OutputPath, SourceCacheFileName)

	fileContent, err := filesystem.ReadJsonConfigFile(cacheFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileContent, &sourceCacheMap); err != nil {
		log.Println("Error when unmarshalling JSON file", cacheFilePath)
		return err
	}

	for key, value := range sourceCacheMap {
		fmt.Println(key, value)
	}

	return nil
}

func loadHeaderCache() error {
	fmt.Println("------------ loading header cache ------------------")
	cacheFilePath := filepath.Join(projectConfiguration.OutputPath, HeaderCacheFileName)

	fileContent, err := filesystem.ReadJsonConfigFile(cacheFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileContent, &headerCacheMap); err != nil {
		log.Println("Error when unmarshalling JSON file", cacheFilePath)
		return err
	}

	for key, value := range headerCacheMap {
		fmt.Println(key, value)
	}

	return nil
}

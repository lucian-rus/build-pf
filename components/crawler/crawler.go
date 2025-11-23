package crawler

import (
	"encoding/json"
	"fmt"
	"gobi/components/env"
	"log"
	"os"
	"path/filepath"
)

func ReadLibraryConfigFileContent(jsonFilePath string, libraryProperties *env.LibraryProperties) error {
	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Println("Error when trying to open", jsonFilePath)
		return err
	}

	if err := json.Unmarshal(fileContent, libraryProperties); err != nil {
		log.Println("Error when unmarshalling file")
		return err
	}

	if true == env.EnableDebugData {
		fmt.Println((*libraryProperties).Name)
		fmt.Println((*libraryProperties).Includes)
		fmt.Println((*libraryProperties).Dependencies)
	}
	return nil
}

func ReadProjectConfigFileContent(jsonFilePath string, projectProperties *env.ProjectProperties) error {
	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Println("Error when trying to open", jsonFilePath)
		return err
	}

	if err := json.Unmarshal(fileContent, projectProperties); err != nil {
		log.Println("Error when unmarshalling file")
		return err
	}

	if true == env.EnableDebugData {
		fmt.Println((*projectProperties).Name)
		fmt.Println((*projectProperties).Includes)
		fmt.Println((*projectProperties).Dependencies)
	}
	return nil
}

// at the moment, this functions support both absolute and relative paths. tbd what would be best
func ScanDirectoryForSourceFiles(directoryPath string, sourceFilesList *[]string, absolutePaths bool) error {
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if ".c" != filepath.Ext(entry.Name()) {
			continue
		}

		if false == absolutePaths {
			fmt.Println("* found file: ", entry.Name())
			*sourceFilesList = append(*sourceFilesList, entry.Name())
			continue
		}

		sourceFilePath, err := filepath.Abs(entry.Name())
		if err != nil {
			return nil
		}

		fmt.Println("* found file: ", sourceFilePath)
		*sourceFilesList = append(*sourceFilesList, sourceFilePath)
	}

	return nil
}

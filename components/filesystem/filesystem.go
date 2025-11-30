package filesystem

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gobi/components/builder"
	"gobi/components/env"
)

func SetupFilesystem(cwd string, project builder.ProjectProperties) error {
	buildDir := filepath.Join(cwd, project.OutputPath)
	logDir := filepath.Join(cwd, "log")

	if DoesEntityExist(buildDir) && DoesEntityExist(logDir) {
		return nil
	}

	if err := os.Mkdir(buildDir, 0777); err != nil {
		return err
	}

	if err := os.Mkdir(logDir, 0777); err != nil {
		return err
	}

	return nil
}

func ReadLibraryConfigFileContent(
	jsonFilePath string,
	libraryProperties *builder.LibraryProperties,
) error {
	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Println("Error when trying to open", jsonFilePath)
		return err
	}

	if err := json.Unmarshal(fileContent, libraryProperties); err != nil {
		log.Println("Error when unmarshalling file")
		return err
	}

	if env.EnableDebugData {
		fmt.Println((*libraryProperties).Name)
		fmt.Println((*libraryProperties).Includes)
		fmt.Println((*libraryProperties).Dependencies)
	}
	return nil
}

func ReadProjectConfigFileContent(
	jsonFilePath string,
	projectProperties *builder.ProjectProperties,
) error {
	fileContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Println("Error when trying to open", jsonFilePath)
		return err
	}

	if err := json.Unmarshal(fileContent, projectProperties); err != nil {
		log.Println("Error when unmarshalling file")
		return err
	}

	if env.EnableDebugData {
		fmt.Println((*projectProperties).Name)
		fmt.Println((*projectProperties).Includes)
		fmt.Println((*projectProperties).Dependencies)
	}
	return nil
}

// entity is defined as either dir/file
func DoesEntityExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

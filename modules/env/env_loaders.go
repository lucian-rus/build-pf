// this file shall contain the loaders -> functions that load cached data
package env

import (
	"encoding/json"
	"fmt"
	"gobi/modules/filesystem"
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

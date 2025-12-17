package filesystem

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateDirectory(dirPath string) error {
	dirsToCreate := strings.Split(dirPath, "/")

	lastBuiltDir := ""
	for _, dir := range dirsToCreate {
		// this appends the last built directory in  to provide proper pathing
		lastBuiltDir = filepath.Join(lastBuiltDir, dir)

		if DoesEntityExist(lastBuiltDir) {
			continue
		}

		if err := os.Mkdir(lastBuiltDir, 0777); err != nil {
			return err
		}
	}

	return nil
}

func ReadJsonConfigFile(jsonFilePath string) ([]byte, error) {
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Println("Error when trying to open JSON file", jsonFilePath)
		return nil, err
	}

	return data, nil
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

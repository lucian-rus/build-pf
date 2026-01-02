package filesystem

import (
	"log"
	"os"
	"path/filepath"
)

func CreateDirectory(dirPath string) error {
	// normalize path and create all parent directories in a cross-platform way.
	normalizedDirPath := filepath.Clean(dirPath)

	// if it already exists, nothing to do.
	if DoesEntityExist(normalizedDirPath) {
		return nil
	}

	// use MkdirAll which works on Windows and Unix and creates any necessary parents.
	if err := os.MkdirAll(normalizedDirPath, 0777); err != nil {
		return err
	}

	return nil
}

func ReadJsonConfigFile(jsonFilePath string) ([]byte, error) {
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Println("Error when trying to open JSON file", jsonFilePath, err)
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

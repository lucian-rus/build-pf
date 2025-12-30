package filesystem

import (
	"log"
	"os"
	"path/filepath"
)

func CreateDirectory(dirPath string) error {
	// Normalize path and create all parent directories in a cross-platform way.
	cleaned := filepath.Clean(dirPath)

	// If it already exists, nothing to do.
	if DoesEntityExist(cleaned) {
		return nil
	}

	// Use MkdirAll which works on Windows and Unix and creates any necessary parents.
	if err := os.MkdirAll(cleaned, 0755); err != nil {
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

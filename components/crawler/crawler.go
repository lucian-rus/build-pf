package crawler

import (
	"fmt"
	"os"
	"path/filepath"
)

// @todo crawler shall store last edited timestamp for each file, as this will be used by incremental build
func ScanDirectoryForSourceFiles(
	directoryPath string,
	sourceFilesList *[]string,
	absolutePaths bool,
) error {
	// at the moment, this functions support both absolute and relative paths. tbd what would be best
	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".c" {
			continue
		}

		if !absolutePaths {
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

func ScanDirectoryForConfigurationFiles() {
}

package crawler

import (
	"fmt"
	"os"
	"path/filepath"
)

func ScanDirectoryForFiles(dirPath string, fileList *[]string, fileType string) error {
	// at the moment, this functions support both absolute and relative paths. tbd what would be best
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, _ error) error {
		if filepath.Ext(info.Name()) != fileType {
			return nil
		}

		sourceFilePath, err := filepath.Abs(path)
		fmt.Println("Scanning directory for files -", dirPath)
		fmt.Println("	* found file: ", sourceFilePath)

		(*fileList) = append(*fileList, sourceFilePath)
		return err
	})

	return err
}

func ScanDirectoryForConfigurationFiles() {
}

func GetTimestampForFile(filePath string, timestamp *int) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	(*timestamp) = int(info.ModTime().Unix())

	return nil
}

func ScanBuildDirectoryForLibraries(dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, _ error) error {

		// if filepath.Ext(info.Name()) == "" && !info.IsDir() {
		// 	builder.LibsBuilt = append(builder.LibsBuilt, info.Name())
		// }
		return nil
	})

	return err
}

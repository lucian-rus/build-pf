package crawler

import (
	"os"
	"path/filepath"
)

// @todo crawler shall store last edited timestamp for each file, as this will be used by incremental build
func ScanDirectoryForSources(dirPath string, fileList *[]string) error {
	// @todo at the moment this only supports c files. should be able to support cpp as well
	return scanDirectoryForFiles(dirPath, fileList, ".c")
}

// @todo remove name as param
func ScanDirectoryForHeaders(dirPath string, fileList *[]string) error {
	return scanDirectoryForFiles(dirPath, fileList, ".h")
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

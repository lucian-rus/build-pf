package crawler

import (
	"gobi/components/builder"
	"os"
	"path/filepath"
)

// @todo crawler shall store last edited timestamp for each file, as this will be used by incremental build
func ScanDirectoryForSources(directoryPath string, fileList *[]string, name string) error {
	// @todo at the moment this only supports c files. should be able to support cpp as well
	return scanDirectoryForFiles(directoryPath, fileList, ".c", name)
}

// @todo remove name as param
func ScanDirectoryForHeaders(directoryPath string, fileList *[]string) error {
	return scanDirectoryForFiles(directoryPath, fileList, ".h", "")
}

func ScanDirectoryForConfigurationFiles() {
}

func ScanBuildDirectoryForLibraries(directoryPath string) error {
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, _ error) error {

		if filepath.Ext(info.Name()) == "" && !info.IsDir() {
			builder.LibsBuilt = append(builder.LibsBuilt, info.Name())
		}
		return nil
	})

	return err
}

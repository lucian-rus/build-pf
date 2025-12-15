package crawler

import (
	"fmt"
	"os"
	"path/filepath"
)

func scanDirectorForFiles(directoryPath string, fileList *[]string, fileType string) error {
	// at the moment, this functions support both absolute and relative paths. tbd what would be best
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, _ error) error {
		if filepath.Ext(info.Name()) != fileType {
			return nil
		}

		sourceFilePath, err := filepath.Abs(path)
		fmt.Println(info.ModTime().Unix())
		fmt.Println("* found file: ", sourceFilePath)
		*fileList = append(*fileList, sourceFilePath)
		return err
	})

	return err
}

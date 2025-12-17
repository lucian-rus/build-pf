package crawler

import (
	"fmt"
	"gobi/components/builder"
	"os"
	"path/filepath"
	"slices"
)

func scanDirectoryForFiles(directoryPath string, fileList *[]string, fileType string, name string) error {
	// at the moment, this functions support both absolute and relative paths. tbd what would be best
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, _ error) error {
		if filepath.Ext(info.Name()) != fileType {
			return nil
		}

		sourceFilePath, err := filepath.Abs(path)
		fmt.Println(info.ModTime().Unix())

		fmt.Println(builder.BuildData)
		// @todo update this. the name shall NOT Be passed as param
		timestamp, ok := builder.BuildData[sourceFilePath]
		if ok && timestamp == int(info.ModTime().Unix()) && slices.Contains(builder.LibsBuilt, name) {
			return nil
		}

		builder.BuildData[sourceFilePath] = int(info.ModTime().Unix())
		fmt.Println("* found file: ", sourceFilePath)
		*fileList = append(*fileList, sourceFilePath)
		return err
	})

	return err
}

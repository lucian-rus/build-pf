package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gobi/modules/filesystem"
)

func Setup(targetPath string) *os.File {
	logDirPath := filepath.Join(targetPath, "log")

	if err := filesystem.CreateDirectory(logDirPath); err != nil {
		log.Fatal("Could not setup logging directory. Error: ", err)
	}

	// @todo ensure generation of log folder is done properly
	currentTime := time.Now()
	fileName := fmt.Sprintf("%-2d%-d%d-%d:%d:%d.log", currentTime.Year(), int(currentTime.Month()),
		currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())

	filePath := filepath.Join(logDirPath, fileName)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Could not setup logging file. Error: ", err)
	}

	log.SetOutput(file)
	log.Println("Logger setup done.")
	return file
}

package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func Setup() *os.File {
	// @todo ensure generation of log folder is done properly
	currentTime := time.Now()
	fileName := fmt.Sprintf("%-2d%-d%d-%d:%d:%d.log", currentTime.Year(), int(currentTime.Month()),
		currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())

	filePath := filepath.Join("log", fileName)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Could not setup logging file. Error: ", err)
	}

	log.SetOutput(file)
	log.Println("Logger setup done.")
	return file
}

package crawler

// @todo crawler shall store last edited timestamp for each file, as this will be used by incremental build
func ScanDirectoryForSources(directoryPath string, fileList *[]string) error {
	return scanDirectorForFiles(directoryPath, fileList, ".c")
}

func ScanDirectoryForHeaders(directoryPath string, fileList *[]string) error {
	return scanDirectorForFiles(directoryPath, fileList, ".h")
}

func ScanDirectoryForConfigurationFiles() {
}

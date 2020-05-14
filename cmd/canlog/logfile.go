package main

import (
	"fmt"
	"os"
	"time"
)

const logsDirectory = "./logs"

// LogFile is a simple wrapper around an os.File, providing convenience
// functionality for file creation, file rotation, and path generation.
type LogFile struct {
	Directory    string
	Filename     string
	IsInProgress bool
	File         *os.File
}

// CreateLogFile opens a new log file, of the format <date>-<suffix>.part,
// and returns a LogFile for manipulation and management.
func CreateLogFile(fileSuffix string) *LogFile {
	logDateTime := time.Now().Format("2006-01-02-15-04-05")

	lf := &LogFile{
		Directory:    logsDirectory,
		Filename:     fmt.Sprintf("%s_%s", logDateTime, fileSuffix),
		IsInProgress: true,
	}

	logFile, err := os.OpenFile(lf.FilePath(), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	lf.File = logFile
	return lf
}

// FilePath returns a fully qualified path to the log file, including '.part'
// suffix in the case of concurrent processing.
func (lf *LogFile) FilePath() string {
	filePath := fmt.Sprintf("%s/%s", lf.Directory, lf.Filename)
	if lf.IsInProgress {
		filePath = filePath + ".part"
	}

	return filePath
}

// Finish closes the log file, and removes the `.part` suffix on the filename.
func (lf *LogFile) Finish() {
	if err := lf.File.Close(); err != nil {
		panic(err)
	}

	lf.IsInProgress = false

	logFileName := lf.FilePath()
	if err := os.Rename(logFileName+".part", logFileName); err != nil {
		panic(err)
	}
}

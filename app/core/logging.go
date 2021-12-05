package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// loggingDir : The logging directory to store the logs.
var loggingDir = filepath.Join(getConfDir(), "logs")

// setUpLogging : Set up the logger to log any useful information such as errors when running the application.
// The log file is stored in the configuration directory.
func (m *MangaDesk) setUpLogging() error {
	if err := os.MkdirAll(loggingDir, os.ModePerm); err != nil {
		return err
	}

	// Create file for current session logging
	formattedDate := time.Now().Format("2006-01-02 15-04-05")
	logFilePath := filepath.Join(loggingDir, fmt.Sprintf("%s.log", formattedDate))

	var err error
	if m.LogFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm); err != nil {
		return err
	}
	log.SetOutput(m.LogFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Printf("Sessions started at %s\n", formattedDate)

	return nil
}

// stopLogging : Closes the log file.
func (m *MangaDesk) stopLogging() error {
	return m.LogFile.Close()
}

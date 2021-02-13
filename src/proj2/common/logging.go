package common

import (
	"fmt"
	"log"
	"os"
	"time"
)

// CreateLogger returns a new logger instance for the client or server;
// logType should indicate who the log is for, e.g., "client" or "server"
func CreateLogger(logType string) (*log.Logger, *os.File) {
	curTime := time.Now().Format("20060102-150405")
	logFile, err := os.OpenFile(fmt.Sprintf("logs/%s_%s.log",
		logType, curTime), os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	// specifically use raw file (unbuffered) similar to stderr rather than
	// bufio writer so that it records even in panic condition
	logger := log.New(logFile, logType, log.LstdFlags)

	// write initial message
	logger.Printf("Beginning log for %s\n", logType)

	return logger, logFile
}

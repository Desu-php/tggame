
package logging

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	
	logDir := "logs"
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFilePath := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(file)
	log.SetLevel(logrus.InfoLevel)

	return log
}

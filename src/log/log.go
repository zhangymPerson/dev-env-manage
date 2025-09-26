package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	IsDebug     = false  // Global debug flag
	LogFilePath = ""     // Global log file path
	logFile     *os.File // File handle for log file
)

// Configure sets global logging configuration
func Configure(debug bool, logPath string) error {
	IsDebug = debug
	LogFilePath = logPath
	return InitLog()
}

// Debug logs debug messages if IsDebug is true
func Debug(format string, v ...interface{}) {
	if IsDebug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs informational messages
func Info(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// Warning logs warning messages
func Warning(format string, v ...interface{}) {
	log.Printf("[WARNING] "+format, v...)
}

// Error logs error messages
func Error(format string, v ...any) {
	log.Printf("[ERROR] "+format, v...)
}

// Fatal logs fatal messages and exits
func Fatal(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}

// SetOutput sets the output destination for the logger
func SetOutput(output io.Writer) {
	log.SetOutput(output)
}

// InitLog initializes the logger with the configured log file
func InitLog() error {
	if LogFilePath != "" {
		if err := os.MkdirAll(filepath.Dir(LogFilePath), 0755); err != nil {
			return err
		}
		file, err := os.OpenFile(LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		logFile = file

		// Set output to both console and file when in debug mode
		if IsDebug {
			multiWriter := io.MultiWriter(os.Stdout, file)
			log.SetOutput(multiWriter)
		} else {
			log.SetOutput(file)
		}
	} else {
		// If no log file path is specified, output only to console
		log.SetOutput(os.Stdout)
	}
	return nil
}

// Close closes the log file if it's open
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

func SetDebug() {
	IsDebug = true
}

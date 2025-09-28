package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	IsDebug     = false  // Global debug flag
	LogFilePath = ""     // Global log file path
	logFile     *os.File // File handle for log file
)

// getCallerInfo 获取调用者信息（文件名和行号）
func getCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0
	}
	// 获取文件名而不包含完整路径
	filename := filepath.Base(file)
	return filename, line
}

// Configure sets global logging configuration
func Configure(debug bool, logPath string) error {
	IsDebug = debug
	LogFilePath = logPath
	return InitLog()
}

// Debug logs debug messages if IsDebug is true
func Debug(format string, v ...interface{}) {
	if IsDebug {
		file, line := getCallerInfo(2) // skip Debug function and its caller
		log.Printf("[%s:%d] [DEBUG] "+format, append([]interface{}{file, line}, v...)...)
	}
}

// Info logs informational messages
func Info(format string, v ...interface{}) {
	file, line := getCallerInfo(2) // skip Info function and its caller
	log.Printf("[%s:%d] [INFO] "+format, append([]interface{}{file, line}, v...)...)
}

// Warning logs warning messages
func Warning(format string, v ...interface{}) {
	file, line := getCallerInfo(2) // skip Warning function and its caller
	log.Printf("[%s:%d] [WARNING] "+format, append([]interface{}{file, line}, v...)...)
}

// Error logs error messages
func Error(format string, v ...any) {
	file, line := getCallerInfo(2) // skip Error function and its caller
	log.Printf("[%s:%d] [ERROR] "+format, append([]any{file, line}, v...)...)
}

// Fatal logs fatal messages and exits
func Fatal(format string, v ...interface{}) {
	file, line := getCallerInfo(2) // skip Fatal function and its caller
	log.Fatalf("[%s:%d] [FATAL] "+format, append([]interface{}{file, line}, v...)...)
}

// Fatalf logs fatal messages and exits
func Fatalf(format string, v ...interface{}) {
	file, line := getCallerInfo(2) // skip Fatalf function and its caller
	log.Fatalf("[%s:%d] [FATAL] "+format, append([]interface{}{file, line}, v...)...)
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

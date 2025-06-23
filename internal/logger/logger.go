package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func AccessLog(r *http.Request, status int, responseSize int64) {
	// Get binary path and determine log directory
	binaryPath, err := os.Executable()
	if err != nil {
		log.Printf("ERROR: Failed to get executable path: %v\n", err)
		return
	}

	folderPath := filepath.Join(filepath.Dir(binaryPath), "logs")
	fileName := filepath.Join(folderPath, "access.log")

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Printf("ERROR: Failed to create log directory '%s': %v\n", folderPath, err)
		return
	}

	// Open file in append mode, create if it doesn't exist
	accessLog, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to open log file '%s': %v\n", fileName, err)
		return
	}
	defer accessLog.Close()

	// Get client IP (handling X-Forwarded-For)
	remoteAddr := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		remoteAddr = forwarded
	}

	// Combined Log Format:
	logMessage := fmt.Sprintf(
		"%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
		remoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL.Path,
		r.Proto,
		status,
		responseSize,
		r.Header.Get("Referer"),
		r.Header.Get("User-Agent"),
	)

	if _, err := accessLog.WriteString(logMessage); err != nil {
		log.Printf("ERROR: Failed to write to log file '%s': %v\n", fileName, err)
	}
}

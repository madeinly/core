package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// AccessLog appends a message to the access log file.
// It ensures the log directory and file exist with proper permissions.
func AccessLog(r *http.Request, status int, responseSize int64) {
	fileName := "logs/access.log"
	folderPath := filepath.Dir(fileName)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Printf("ERROR: Failed to create log directory '%s': %v\n", folderPath, err)
		return
	}

	// Open file in append mode
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to open log file '%s': %v\n", fileName, err)
		return
	}
	defer file.Close()

	// Get client IP (supports X-Forwarded-For if behind a proxy)
	remoteAddr := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		remoteAddr = forwarded
	}

	// Combined Log Format:
	// $remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"
	logMessage := fmt.Sprintf(
		"%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"\n",
		remoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"), // Apache-style timestamp
		r.Method,
		r.URL.Path,
		r.Proto,
		status,
		responseSize,
		r.Header.Get("Referer"),
		r.Header.Get("User-Agent"),
	)

	// Write to file
	if _, err := file.WriteString(logMessage); err != nil {
		log.Printf("ERROR: Failed to write to log file '%s': %v\n", fileName, err)
	}
}

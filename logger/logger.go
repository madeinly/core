package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

func ErrorLog() {

}

func Log(err error, msg string, data string) {
	binaryPath, err := os.Executable()
	if err != nil {
		log.Printf("ERROR: Failed to get executable path: %v\n", err)
		return
	}

	folderPath := filepath.Join(filepath.Dir(binaryPath), "logs")
	fileName := filepath.Join(folderPath, "error.log")

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		log.Printf("ERROR: Failed to create log directory '%s': %v\n", folderPath, err)
		return
	}

	errorLog, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to open log file '%s': %v\n", fileName, err)
		return
	}
	defer errorLog.Close()

	logMessage := fmt.Sprintf(
		"%s - %s - %s\n",
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		msg,
		data,
	)

	if _, err := errorLog.WriteString(logMessage); err != nil {
		log.Printf("ERROR: Failed to write to log file '%s': %v\n", fileName, err)
	}

}

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

func DebugLog(req *http.Request, res *httptest.ResponseRecorder) {

	var (
		green   = color.New(color.FgGreen).SprintFunc()
		red     = color.New(color.FgRed).SprintFunc()
		cyan    = color.New(color.FgCyan).SprintFunc()
		yellow  = color.New(color.FgYellow).SprintFunc()
		blue    = color.New(color.FgBlue).SprintFunc()
		magenta = color.New(color.FgMagenta).SprintFunc()
		gray    = color.New(color.FgHiBlack).SprintFunc()
	)

	start := time.Now()

	// get Body
	var requestBody string

	bodyBytes, err := io.ReadAll(res.Body)
	if err == nil {
		requestBody = string(bodyBytes)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// start template
	fmt.Printf("%s\n", green("───────────────────────────────────────────────────────────────"))

	fmt.Printf(
		"%s - - [%s] \"%s %s %s\" %s %s\n%s %v\n%s %v",
		gray(req.RemoteAddr),
		cyan(time.Now().Format("02/Jan/2006:15:04:05 -0700")),
		green(req.Method),
		yellow(req.URL.Path),
		magenta(req.Proto),
		blue(req.ContentLength),
		cyan(time.Since(start)),
		gray("Params:"), req.URL.Query(),
		gray("Headers:"), req.Header,
	)

	fmt.Println("Cookies:", req.Cookies()) // Explicitly log cookies

	if requestBody != "" {
		fmt.Printf("\n%s\n%s", gray("Request Body:"), formatJSON(requestBody))
	}

	if res != nil {
		fmt.Printf("\n\n%s", gray("Response: "))
		statusStr := fmt.Sprintf("\nStatus: %d", res.Code)
		if res.Code >= 400 {
			log.Print(red(statusStr))
		} else {
			log.Print(green(statusStr))
		}
		fmt.Println(gray("Headers:"))
		fmt.Printf("%v\n", res.Header())

		responseBody := res.Body.String()
		if responseBody != "" {
			fmt.Println(gray("Body:"))
			fmt.Printf("%s\n", formatJSON(responseBody))
		}

	}

	fmt.Printf("%s\n\n", red("───────────────────────────────────────────────────────────────"))

	AccessLog(req, res.Code, int64(res.Body.Len()))

}

func formatJSON(input string) string {
	var raw interface{}
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return input
	}
	pretty, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return input
	}
	return string(pretty)
}

package fatal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/madeinly/core/email"
	"github.com/madeinly/core/settings"
)

// OnErr logs, attempts to notify, and then exits the program.
// It never returns.
func OnErr(err error, msg string, args ...any) {
	if err == nil {
		return
	}

	formattedMsg := fmt.Sprintf(msg, args...)

	// 1. Log the fatal error to a file
	logToFile(err, formattedMsg)

	// 2. Attempt to put the frontend into maintenance mode
	mockAPIKey := "your-secret-api-key" // In a real scenario, this should come from a secure source
	maintenanceURL := settings.GetSettings().FrontDomain + "/api/maintenance"

	req, _ := http.NewRequest("POST", maintenanceURL, nil)
	req.Header.Set("Authorization", "Bearer "+mockAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, httpErr := client.Do(req)

	maintenanceStatus := "succeeded"
	if httpErr != nil || resp.StatusCode != 200 {
		maintenanceStatus = fmt.Sprintf("failed with status: %d and error: %v", resp.StatusCode, httpErr)
	}

	// 3. Send an email notification
	subject := "Fatal Error in Application"
	body := fmt.Sprintf(
		"A fatal error occurred:\n\nMessage: %s\nError: %v\n\nMaintenance mode activation: %s",
		formattedMsg,
		err,
		maintenanceStatus,
	)

	// This assumes the email settings are valid. If email fails, we still exit.
	_ = email.Send(settings.GetSettings().Email.User, subject, body)

	// 4. Exit the program
	os.Exit(1)
}

// logToFile handles writing the fatal error to the log file.
func logToFile(err error, formattedMsg string) {
	exe, exeErr := os.Executable()
	if exeErr != nil {
		log.Printf("FATAL: Could not get executable path: %v", exeErr)
		log.Printf("Original error: %s - %v", formattedMsg, err)
		return
	}
	logDir := filepath.Join(filepath.Dir(exe), "logs")
	_ = os.MkdirAll(logDir, 0o755)

	logPath := filepath.Join(logDir, "error.log")
	f, fileErr := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if fileErr != nil {
		log.Printf("FATAL: Could not open log file: %v", fileErr)
		log.Printf("Original error: %s - %v", formattedMsg, err)
		return
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("FATAL: %s", formattedMsg)
	logger.Printf("Error: %v", err)
	logger.Printf("goroutine %d\n", runtime.NumGoroutine())
}


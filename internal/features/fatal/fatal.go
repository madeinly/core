package fatal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/madeinly/core/internal/features/email"
	"github.com/madeinly/core/internal/features/settings"
)

// OnErr logs, attempts to notify, and then exits the program.
// It never returns.
func OnErr(err error, msg string, args ...any) {
	if err == nil {
		return
	}

	formattedMsg := fmt.Sprintf(msg, args...)

	// 1. Log the fatal error using the standard logger
	log.Printf("FATAL: %s", formattedMsg)
	log.Printf("Error: %v", err)
	log.Printf("goroutine %d\n", runtime.NumGoroutine())

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

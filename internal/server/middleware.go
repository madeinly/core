package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/fatih/color"
)

// Create colored printers

// type contextKey string

// const (
// 	userIDKey contextKey = "userID"
// )

// func LoggedRoutes(next http.Handler) http.Handler {

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		isLogged, claims, _ := auth.ValidateCookie(r)
// 		if !isLogged {
// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
// 			return
// 		}
// 		fmt.Println(claims.UserID)

// 		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func Logging(next http.Handler) http.Handler {
	const devMode = true // Set this to false in production

	// Define color functions using the faith/color package
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc() // HiBlack is often a more visible gray

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		start := time.Now()

		// In dev mode, read and log the request body
		var requestBody string
		if devMode {
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// Restore the body so handlers can read it
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Create a response recorder to capture the response
		var responseRecorder *httptest.ResponseRecorder
		var originalWriter http.ResponseWriter = w
		if devMode {
			responseRecorder = httptest.NewRecorder()
			w = responseRecorder
		}

		// Start log with green line
		if devMode {
			fmt.Printf("%s\n", green("───────────────────────────────────────────────────────────────"))
		}

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the basic request info
		// Using log.Printf might strip colors depending on the log output.
		// For consistent color output, fmt.Printf directly to os.Stdout might be better
		// if log's default output is stripping them.
		// However, sticking to log.Printf as per the original for now.
		fmt.Printf(
			"%s - - [%s] \"%s %s %s\" %s %s\n%s %v\n%s %v",
			gray(r.RemoteAddr),
			cyan(time.Now().Format("02/Jan/2006:15:04:05 -0700")),
			green(r.Method),
			yellow(r.URL.Path),
			magenta(r.Proto),
			blue(r.ContentLength),
			cyan(time.Since(start)),
			gray("Params:"), r.URL.Query(),
			gray("Headers:"), r.Header,
		)

		// Dev mode logging
		if devMode {
			// Log request body if present
			if requestBody != "" {
				fmt.Printf("\n%s\n%s", gray("Request Body:"), formatJSON(requestBody))
			}

			// Log response if we captured it
			if responseRecorder != nil {
				fmt.Printf("\n\n%s", gray("Response: "))
				statusStr := fmt.Sprintf("\nStatus: %d", responseRecorder.Code)
				if responseRecorder.Code >= 400 {
					log.Print(red(statusStr))
				} else {
					log.Print(green(statusStr))
				}
				fmt.Println(gray("Headers:"))
				fmt.Printf("%v\n", responseRecorder.Header())

				responseBody := responseRecorder.Body.String()
				if responseBody != "" {
					fmt.Println(gray("Body:"))
					fmt.Printf("%s\n", formatJSON(responseBody))
				}

				// Write the recorded response to the original writer
				for k, v := range responseRecorder.Header() {
					originalWriter.Header()[k] = v
				}
				originalWriter.WriteHeader(responseRecorder.Code)
				originalWriter.Write(responseRecorder.Body.Bytes())
			}

			// End log with red line
			fmt.Printf("%s\n\n", red("───────────────────────────────────────────────────────────────"))
		}
	})
}

// formatJSON attempts to pretty-print JSON strings
func formatJSON(input string) string {
	var raw interface{}
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return input // Return as-is if not JSON
	}
	pretty, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return input
	}
	return string(pretty)
}

// func AuthRoute(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		isLogged, _, err := auth.ValidateCookie(r)

// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		if !isLogged {
// 			http.Error(w, "Bad request", http.StatusForbidden)
// 			return

// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

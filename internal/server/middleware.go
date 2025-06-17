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
	"github.com/madeinly/core/internal/logger"
)

func Logging(next http.Handler) http.Handler {
	const devMode = false

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		start := time.Now()

		// get response with recorder
		var originalWriter http.ResponseWriter = w
		responseRecorder := httptest.NewRecorder()
		w = responseRecorder

		if devMode {

			//define colors
			var (
				green   = color.New(color.FgGreen).SprintFunc()
				red     = color.New(color.FgRed).SprintFunc()
				cyan    = color.New(color.FgCyan).SprintFunc()
				yellow  = color.New(color.FgYellow).SprintFunc()
				blue    = color.New(color.FgBlue).SprintFunc()
				magenta = color.New(color.FgMagenta).SprintFunc()
				gray    = color.New(color.FgHiBlack).SprintFunc()
			)

			// get Body
			var requestBody string
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			next.ServeHTTP(w, r)

			// start template
			fmt.Printf("%s\n", green("───────────────────────────────────────────────────────────────"))

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

			if requestBody != "" {
				fmt.Printf("\n%s\n%s", gray("Request Body:"), formatJSON(requestBody))
			}

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

				for k, v := range responseRecorder.Header() {
					originalWriter.Header()[k] = v
				}
				originalWriter.WriteHeader(responseRecorder.Code)
				originalWriter.Write(responseRecorder.Body.Bytes())
			}

			fmt.Printf("%s\n\n", red("───────────────────────────────────────────────────────────────"))
		} else {

			next.ServeHTTP(w, r)
		}

		logger.AccessLog(r, responseRecorder.Code, int64(responseRecorder.Body.Len()))

	})
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

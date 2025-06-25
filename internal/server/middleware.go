package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/madeinly/core/internal/logger"
	"github.com/madeinly/core/internal/settings"
)

var debugMode bool = settings.Settings.Debug

func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		enableCORS(w, r)

		fmt.Println("debug mode", debugMode)

		if debugMode {
			recorder := httptest.NewRecorder()
			next.ServeHTTP(recorder, r)

			for k, v := range recorder.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(recorder.Code)
			w.Write(recorder.Body.Bytes())

			logger.DebugLog(r, recorder)
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	// Get the requesting origin
	origin := r.Header.Get("Origin")

	// Allow local development ports
	if strings.HasPrefix(origin, "http://localhost:") ||
		strings.HasPrefix(origin, "http://192.168.11.") {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// For preflight requests
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600") // 1 hour cache
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

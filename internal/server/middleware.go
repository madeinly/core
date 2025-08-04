package server

import (
	"net/http"
	"net/http/httptest"

	"github.com/madeinly/core/logger"
	"github.com/madeinly/core/settings"
)

func Logging(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AppSettings := settings.GetSettings()
		debugMode := AppSettings.Debug

		enableCORS(w, r)

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

	// frontDomain := settings.GetSettings().FrontDomain
	origin := r.Header.Get("Origin")

	// fmt.Println(frontDomain, origin)

	// 1. CORS Configuration
	// if strings.HasPrefix(origin, fmt.Sprintf("http://%s", frontDomain)) {
	if true {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// 2. Preflight Handling
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Vary", "Origin")
}

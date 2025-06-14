package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MadeSimplest/core/internal/features"
)

func SetupRoutes(mux *http.ServeMux) {

	for _, feature := range features.Available {

		for _, route := range feature.Routes {

			mux.Handle(route.Type+" "+route.Pattern, route.Handler)
		}

	}

}

func Start(port string) {

	if port == "" {
		port = "8080"
	}

	router := http.NewServeMux()
	SetupRoutes(router)

	addr := ":" + port
	fmt.Printf("Server running on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, Logging(router)); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

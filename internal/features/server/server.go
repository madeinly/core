package server

import (
	"net/http"
	"os"

	"github.com/madeinly/core/internal/extensions"
)

func SetupRoutes(mux *http.ServeMux) {

	for _, feature := range extensions.Available {

		for _, route := range feature.Routes {

			mux.Handle(route.Type+" "+route.Pattern, route.Handler)
		}

	}

}

func Start(address string, port string) error {

	router := http.NewServeMux()
	SetupRoutes(router)

	fulladdress := address + ":" + port

	if err := http.ListenAndServe(fulladdress, Logging(router)); err != nil {
		os.Exit(1)
	}

	return nil
}

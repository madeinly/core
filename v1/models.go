package v1

import (
	"net/http"

	"github.com/madeinly/core/internal/features/validation"
)

type Route struct {
	Type    string
	Pattern string
	Handler http.Handler
}

type Migration struct {
	Name   string
	Schema string
}

type Features []FeaturePackage

type Arg struct {
	Name        string
	Default     string
	Required    bool
	Description string
}

type FeaturePackage struct {
	Name      string
	Routes    []Route
	Migration Migration
	Setup     func(map[string]string) error
	Cmd       func()
	Args      []Arg
}

type Error = validation.Error

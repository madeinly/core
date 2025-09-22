package extensions

import "net/http"

type FeaturePackage struct {
	Name      string
	Routes    []Route
	Migration Migration
	Setup     func(map[string]string) error
	Cmd       func()
	Args      []Arg
}

type Route struct {
	Type    string
	Pattern string
	Handler http.Handler
}

type Migration struct {
	Name   string
	Schema string
}

type Arg struct {
	Name        string
	Default     string
	Required    bool
	Description string
}

type Features []FeaturePackage

package models

import (
	"net/http"
)

type Route struct {
	Type    string
	Pattern string
	Handler http.HandlerFunc
}

type Migration struct {
	Name   string
	Schema string
}

type Features []FeaturePackage

type FeaturePackage struct {
	Name      string
	Routes    []Route
	Migration Migration
	Setup     func() error
	Cmd       func()
}

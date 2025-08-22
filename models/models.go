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

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

package extensions

import "net/http"

type Mod struct {
	Name        string
	Routes      []Route
	Migration   Migration
	Setup       func(map[string]string) error
	Cmd         func()
	InstallArgs []InstallArg
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

type InstallArg struct {
	Name        string
	Default     string
	Required    bool
	Description string
}

type Mods []Mod

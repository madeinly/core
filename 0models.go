package core

import (
	"github.com/madeinly/core/internal/extensions"
	"github.com/madeinly/core/internal/features/validation"
)

type Route = extensions.Route

type Migration = extensions.Migration

type Mods = []extensions.Mod

type InstallArg = extensions.InstallArg

type Mod = extensions.Mod

type Error = validation.Error

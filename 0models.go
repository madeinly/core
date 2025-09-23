package v1

import (
	"github.com/madeinly/core/internal/extensions"
	"github.com/madeinly/core/internal/features/validation"
)

type Route = extensions.Route

type Migration = extensions.Migration

type Features = []extensions.FeaturePackage

type Arg = extensions.Arg

type FeaturePackage = extensions.FeaturePackage

type Error = validation.Error

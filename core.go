package core

import (
	"github.com/madeinly/core/internal/extensions"
	"github.com/madeinly/core/internal/gateways/cmd"
	v1 "github.com/madeinly/core/v1"
)

func Start(featuresAvailable v1.Features) {

	extensions.RegisterFeatures(featuresAvailable)

	cmd.Execute()

}

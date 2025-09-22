package flows

import (
	"github.com/madeinly/core/internal/extensions"
	v1 "github.com/madeinly/core/v1"
)

func Start(featuresAvailable v1.Features) {

	extensions.RegisterFeatures(featuresAvailable)

	extensions.CmdRouter()

}

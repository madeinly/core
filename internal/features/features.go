package features

import (
	"github.com/madeinly/core/models"
)

var Available models.Features

func RegisterFeatures(features models.Features) {

	Available = features

}

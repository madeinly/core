package features

import (
	"github.com/MadeSimplest/core/models"
)

var Available models.Features

func RegisterFeatures(features models.Features) {

	Available = features

}

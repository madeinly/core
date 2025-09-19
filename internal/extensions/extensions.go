package extensions

import v1 "github.com/madeinly/core/v1"

var Available v1.Features

func RegisterFeatures(features v1.Features) {

	Available = features

}

package runners

import (
	"fmt"

	"github.com/madeinly/core/models"
)

type SetupFunc func() error

type SetupFeature struct {
	Name string
	Fn   SetupFunc
}

// RunAll executes all registered setup functions
func RunSetupPackages(features models.Features) error {

	for _, feature := range features {

		err := feature.Setup()
		if err != nil {
			return fmt.Errorf("failed to execute setup %s: %w", feature.Name, err)
		}

	}

	fmt.Println("todo salio chevere")

	return nil
}

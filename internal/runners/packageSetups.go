package runners

import (
	"fmt"

	"github.com/madeinly/core/internal/features"
)

type SetupFunc func() error

type SetupFeature struct {
	Name string
	Fn   SetupFunc
}

// RunAll executes all registered setup functions
func RunSetupPackages(setupArgs map[string]map[string]string) error {

	fmt.Println(setupArgs)

	for _, feature := range features.Available {

		fmt.Println(feature.Name)

		err := feature.Setup(setupArgs[feature.Name])
		fmt.Println("passed the user")
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("failed to execute setup %s: %w", feature.Name, err)
		}

	}

	fmt.Println("all modules were successfully setup")

	return nil
}

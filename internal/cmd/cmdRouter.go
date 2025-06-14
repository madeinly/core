package cmd

import (
	"fmt"
	"os"

	"github.com/MadeSimplest/core/internal/features"
)

func CmdRouter() {
	if len(os.Args) < 2 {
		fmt.Println("you need to pick a feature")
		os.Exit(1)
		return
	}

	featureName := os.Args[1]
	remainingArgs := os.Args[2:]

	// Handle core execution
	if featureName == "core" {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = append([]string{featureName}, remainingArgs...)

		Execute()
		return
	}

	// Feature routing
	var targetCmd func()
	for _, feature := range features.Available {
		if feature.Name == featureName {
			targetCmd = feature.Cmd
			break
		}
	}

	if targetCmd == nil {
		fmt.Printf("Unknown feature: %s\n", featureName)
		os.Exit(1)
	}

	// Prepare args for feature CLI
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = append([]string{featureName}, remainingArgs...)

	targetCmd()
}

package extensions

import (
	"fmt"
	"os"
)

func CmdRouter() {
	if len(os.Args) < 2 {
		fmt.Println("you need to pick a feature")
		os.Exit(1)
		return
	}

	featureName := os.Args[1]
	remainingArgs := os.Args[2:]

	/*
		will hold the feature routing (if exist) entry point
	*/
	var targetCmd func()

	for _, feature := range Available {
		if feature.Name == featureName {
			targetCmd = feature.Cmd
			break
		}
	}

	/*
		maybe in the future avoid sending data from here and instead adding a central channel for that
	*/
	if targetCmd == nil {
		fmt.Printf("Unknown feature: %s", featureName)
		os.Exit(1)
	}

	// Prepare args for feature CLI
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = append([]string{featureName}, remainingArgs...)

	targetCmd()
}

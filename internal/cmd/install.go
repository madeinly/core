package cmd

import (
	"fmt"

	"github.com/madeinly/core/internal/features"
	"github.com/madeinly/core/internal/runners"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Initialize and configure the application database",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {

		// crear un slice y popularlo con los datos de las flags

		setupArgs := make(map[string]map[string]string)

		var err error

		for _, feature := range features.Available {
			if feature.Args == nil {
				continue
			}

			// create inner map for this feature
			if setupArgs[feature.Name] == nil {
				setupArgs[feature.Name] = make(map[string]string)
			}

			for _, arg := range feature.Args {

				val, err := cmd.Flags().GetString(arg.Name)
				if err != nil {
					fmt.Printf("error getting %s: %v\n", arg.Name, err)
					return
				}
				setupArgs[feature.Name][arg.Name] = val
			}
		}

		err = runners.RunMigrations(features.Available)

		if err != nil {
			fmt.Println(err.Error())
		}

		runners.RunSetupPackages(setupArgs)
	},
}

func init() {

	rootCmd.AddCommand(InstallCmd)

	fmt.Println("feature", features.Available)

	for _, feature := range features.Available {

		if feature.Args == nil {
			continue
		}

		for _, arg := range feature.Args {

			InstallCmd.Flags().String(arg.Name, "", arg.Description)

			if arg.Required {
				InstallCmd.MarkFlagRequired(feature.Name + "_" + arg.Name)
			}

		}

	}

}

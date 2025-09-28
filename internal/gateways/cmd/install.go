package cmd

import (
	"fmt"

	"github.com/madeinly/core/internal/extensions"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Initialise and configure the application",

	Run: runCmd,
}

func runCmd(cmd *cobra.Command, args []string) {
	// setupArgs := make(map[string]map[string]string)
	// for _, feat := range extensions.Available {
	// 	if feat.InstallArgs == nil {
	// 		continue
	// 	}
	// 	if setupArgs[feat.Name] == nil {
	// 		setupArgs[feat.Name] = make(map[string]string)
	// 	}
	// 	for _, arg := range feat.InstallArgs {
	// 		val, _ := cmd.Flags().GetString(arg.Name) // err already checked
	// 		setupArgs[feat.Name][arg.Name] = val
	// 	}
	// }

	// if err := extensions.RunMigrations(extensions.Available); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// extensions.RunSetupPackages(setupArgs)
}

func init() {
	var installArgs []extensions.InstallArg

	fmt.Println(extensions.Available)
	fmt.Println("testing")

	for _, mod := range extensions.Available {

		fmt.Println(mod.Name)

		if mod.InstallArgs == nil {
			continue
		}

		installArgs = append(installArgs, mod.InstallArgs...)
	}

	fmt.Println("this are the install arguments", installArgs)

	for _, arg := range installArgs {

		if arg.Required {
			installCmd.Flags().String(arg.Name, "", arg.Description)
			installCmd.MarkFlagRequired(arg.Name)
			continue
		}

		installCmd.Flags().String(arg.Name, arg.Default, arg.Description)
	}

	rootCmd.AddCommand(installCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/madeinly/core/internal/extensions"
	"github.com/spf13/cobra"
)

// installCmd is created **without** any dynamic flags
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Initialise and configure the application database",
	Long:  ``,
	// 1. Cobra calls this **after** parsing the raw CLI tokens
	//    but **before** normal flag parsing and **before** Run.
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// 2. Dynamically add the flags
		for _, feat := range extensions.Available {
			if feat.Args == nil {
				continue
			}
			for _, arg := range feat.Args {
				if arg.Required {
					cmd.Flags().String(arg.Name, "", arg.Description)
					_ = cmd.MarkFlagRequired(arg.Name)
				} else {
					cmd.Flags().String(arg.Name, arg.Default, arg.Description)
				}
			}
		}
		// 3. Tell Cobra to re-parse the command line with the new flags
		return cmd.ParseFlags(os.Args[1:])
	},
	Run: func(cmd *cobra.Command, args []string) {
		// normal body unchanged
		setupArgs := make(map[string]map[string]string)
		for _, feat := range extensions.Available {
			if feat.Args == nil {
				continue
			}
			if setupArgs[feat.Name] == nil {
				setupArgs[feat.Name] = make(map[string]string)
			}
			for _, arg := range feat.Args {
				val, _ := cmd.Flags().GetString(arg.Name) // err already checked
				setupArgs[feat.Name][arg.Name] = val
			}
		}
		if err := extensions.RunMigrations(extensions.Available); err != nil {
			fmt.Println(err)
			return
		}
		extensions.RunSetupPackages(setupArgs)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

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
	Long: `The install command performs complete database setup for the application.
    
This command will:
1. Create all necessary database tables and structures
2. Run all available schema migrations
3. Populate the database with initial seed data

This is typically run:
- When first setting up the application
- After cloning the project to a new environment
- When you need to reset the database to a fresh state

Examples:
  # Basic installation
  myapp install
  
  # Installation with verbose output
  myapp install --verbose

Warning: Running this command on an existing database will:
- Apply any pending migrations
- Recreate seed data
- Potentially overwrite existing records in seed tables

The command is idempotent - it can be safely run multiple times as it will
only make changes when necessary.`,

	Run: func(cmd *cobra.Command, args []string) {
		runInstall()
	},
}

func init() {
	rootCmd.AddCommand(InstallCmd)
}

func runInstall() {

	err := runners.RunMigrations(features.Available)

	if err != nil {
		fmt.Println(err.Error())
	}

	runners.RunSetupPackages(features.Available)

}

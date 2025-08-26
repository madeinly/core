package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "simplest cms commands",
	Long: `simplest - A Fast and Flexible Static Site Generator built with Go.

Complete documentation is available at http://hugo.spf13.com

Key Features:
- Blazing fast performance
- Simple markdown content authoring
- Flexible templating system
- Beautiful themes and easy customization`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is core entry")
	},
}

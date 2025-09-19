package cmd

import (
	"fmt"
	"os"

	"github.com/madeinly/core/internal/extensions"
	"github.com/spf13/cobra"
)

// This is simply the root command

/*
Not sure if it has any use beyond being the start point
[TODO]
if it has no use find the best way to disable its access to users
*/
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "simplest cms commands",
	Run: func(cmd *cobra.Command, args []string) {
		extensions.CmdRouter()

	},
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

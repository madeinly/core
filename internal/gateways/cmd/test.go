package cmd

import (
	"fmt"

	"github.com/madeinly/core/internal/flows"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

// installCmd represents the install command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Initialize and configure the application database",

	Run: func(cmd *cobra.Command, args []string) {

		tempFile := flows.TempFile()

		tempFile.WriteString("yenny es lo mas bonito que hay")

		fmt.Println(tempFile.Name())

		defer tempFile.Close()

	},
}

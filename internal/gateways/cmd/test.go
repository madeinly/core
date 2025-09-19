package cmd

import (
	"fmt"

	"github.com/madeinly/core/internal/features/email"
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
		err := email.Send("soyrbto@gmail.com", "testing", "hola, esta es una prueba")

		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

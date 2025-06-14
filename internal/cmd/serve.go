package cmd

import (
	"github.com/MadeSimplest/core/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long: `Start the HTTP server on the specified port.
If no port is provided, it defaults to 8080.

Examples:
  serve             # starts server on port 8080
  serve --port 3000 # starts server on port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		server.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

package cmd

import (
	"github.com/madeinly/core/internal/server"
	"github.com/madeinly/core/internal/settings"
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
		settings.InitSettings()

		server.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")

}

package cmd

import (
	"fmt"
	"sync"

	"github.com/madeinly/core/internal"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the HTTP server",
	Long: `Starts the program listening to address and port set by flags: --address --port
It takes priority from the flags you pass, if none are pass it uses the one set up on settings.toml
if there is no settings it will automatically generates a settings file with address as Slocalhost and port as 1234`,
	Run: func(cmd *cobra.Command, args []string) {

		address, err := cmd.Flags().GetString("address")
		if err != nil {
			fmt.Printf("error getting the address: %v", err)
		}

		port, err := cmd.Flags().GetString("port")
		if err != nil {
			fmt.Printf("error getting the port: %v", err)
		}

		quiet, err := cmd.Flags().GetBool("quiet")
		if err != nil {
			fmt.Printf("error getting quite value: %v", err)
		}

		ch := make(chan string)
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			for msg := range ch {
				fmt.Println("\n" + msg)
			}
		}()

		err = internal.StartServer(internal.StartServerParams{
			Ch:      ch,
			Wg:      &wg,
			Address: address,
			Port:    port,
			Quiet:   quiet,
		})

		if err != nil {
			wg.Wait()
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("address", "", "address to run the server on")
	serveCmd.Flags().String("port", "", "port to run the server on")
	serveCmd.Flags().BoolP("quiet", "q", false, "detach the server from the terminal")
}

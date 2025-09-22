package cmd

import (
	"fmt"
	"os"

	"github.com/madeinly/core/internal/flows"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "core",
	Short: "simplest cms commands",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("testing the core")

	},
}

/* este es el verdadero lugar para el router*/
/*
Recordar: el verdadero entry point a la aplicacion son los gateways
asi que aqui es donde debe estar la funcion que la inicializa
y por lo tanto aqui debe colocarse la logica que redirija al cmd correcto
*/

// Initialize the app
func Execute() {

	//3 elementos porque como minimo se necesita
	// [app] [feature] [action]
	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments")
		return
	}

	if os.Args[1] != "core" {
		flows.StartCmd()
		return
	}

	// Remove element at index 1 (second element)
	newArgs := make([]string, 0, len(os.Args)-1)
	newArgs = append(newArgs, os.Args[0])     // Keep first element
	newArgs = append(newArgs, os.Args[2:]...) // Skip second element, add rest
	os.Args = newArgs

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

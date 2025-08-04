package internal

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/madeinly/core/internal/files"
	"github.com/madeinly/core/internal/server"
	"github.com/madeinly/core/internal/settings"
)

func StartServer(
	ch chan<- string,
	wg *sync.WaitGroup,
	address string,
	port string,
	quiet bool) error {

	RunChecks(ch)

	defer func() {
		close(ch)
		wg.Wait()
	}()

	// Integrity Validations ==================//
	troubleFiles, err := files.FilesIntegrity()

	if err != nil {
		return fmt.Errorf("could not validate files:%w", err)
	}

	if len(troubleFiles) > 0 {

		ch <- "Some files seems to be lack of permissions or ownership"

		for _, file := range troubleFiles {
			ch <- file
		}

		return fmt.Errorf("")
	}

	// Settings handles ====================== //
	go settings.WatchSettings()

	currentSettings := settings.GetSettings()

	// debuggers  ====================== //

	if currentSettings.Debug {

		currentSettingsJson, err := json.MarshalIndent(currentSettings, "", " ")

		if err != nil {
			close(ch)
			return err
		}
		ch <- "Current Settings:"

		ch <- string(currentSettingsJson)
	}

	// start server pre-flight  ====================== //

	if port == "" {
		port = currentSettings.Port
	}

	if address == "" {
		address = currentSettings.Address
	}

	ch <- fmt.Sprintf("The server is running on %s:%s", address, port)

	// Server Launch TTY took over by the server listener  ====================== //
	//[!TODO]: run the server without attaching to the tty
	if quiet {
		return server.Start(address, port)
	}

	return server.Start(address, port)
}

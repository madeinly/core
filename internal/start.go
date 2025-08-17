package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/madeinly/core/internal/files"
	"github.com/madeinly/core/internal/server"
	"github.com/madeinly/core/internal/settings"
)

var ErrFileIntegrity = errors.New("file integrity check failed, see logs for details")

type StartServerParams struct {
	Ch      chan<- string
	Wg      *sync.WaitGroup
	Address string
	Port    string
	Quiet   bool
}

func StartServer(ctx context.Context, params StartServerParams) error {

	RunChecks(params.Ch)

	defer func() {
		close(params.Ch)
		params.Wg.Wait()
	}()

	// Integrity Validations ==================//
	troubleFiles, err := files.FilesIntegrity()

	if err != nil {
		return fmt.Errorf("could not validate files:%w", err)
	}

	if len(troubleFiles) > 0 {

		params.Ch <- "Some files seems to be lack of permissions or ownership"

		for _, file := range troubleFiles {
			params.Ch <- file
		}

		return ErrFileIntegrity
	}

	// Settings handles ====================== //
	go settings.WatchSettings(ctx)

	currentSettings := settings.GetSettings()

	// debuggers  ====================== //

	if currentSettings.Debug {

		currentSettingsJson, err := json.MarshalIndent(currentSettings, "", " ")

		if err != nil {
			return err
		}
		params.Ch <- "Current Settings:"

		params.Ch <- string(currentSettingsJson)
	}

	// start server pre-flight  ====================== //

	if params.Port == "" {
		params.Port = currentSettings.Port
	}

	if params.Address == "" {
		params.Address = currentSettings.Address
	}

	params.Ch <- fmt.Sprintf("The server is running on %s:%s", params.Address, params.Port)

	// Server Launch TTY took over by the server listener  ====================== //
	//[!TODO]: run the server without attaching to the tty
	return server.Start(params.Address, params.Port)
}

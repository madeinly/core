package core

import (
	"database/sql"
	"os"
	"path"

	"github.com/madeinly/core/internal/cmd"
	"github.com/madeinly/core/internal/db"
	"github.com/madeinly/core/internal/features"
	"github.com/madeinly/core/models"

	"github.com/madeinly/core/fatal"
	"github.com/madeinly/core/logger"
	"github.com/madeinly/core/settings"
	"github.com/madeinly/core/validation"
)

// Exposing features from internal packages

// Fatal will terminate the program immediately if an error occurs.
var Fatal = fatal.OnErr

// Log provides a way to log errors and warnings.
var Log = logger.Log

// Settings provides access to the application's configuration in a structured format.
var Settings = settings.GetSettings

// RawSettings provides access to the application's configuration as a raw map.
var RawSettings = settings.GetRawSettings

// Validate provides a validation system for accumulating input.
var Validate = validation.New

func Start(featuresAvailable models.Features) {

	features.RegisterFeatures(featuresAvailable)

	cmd.CmdRouter()

}

func DB() *sql.DB {

	db := db.GetDB()

	return db
}

func RootPath() string {

	binPath, err := os.Executable()

	if err != nil {
		fatal.OnErr(err, "error getting the root path: %v", err)
	}

	return path.Dir(binPath)
}

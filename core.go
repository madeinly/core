package core

import (
	"github.com/madeinly/core/internal/cmd"
	"github.com/madeinly/core/internal/db"
	"github.com/madeinly/core/internal/email"
	"github.com/madeinly/core/internal/fatal"
	"github.com/madeinly/core/internal/features"
	"github.com/madeinly/core/internal/files"
	"github.com/madeinly/core/internal/logger"
	"github.com/madeinly/core/internal/settings"
	"github.com/madeinly/core/models"
	"github.com/madeinly/core/validation"
)

// --- Exposed Core Features ---

// Fatal logs a critical error, attempts to send a notification, and then
// immediately terminates the application. This should be used for unrecoverable
// errors where the application cannot continue safely.
//
// Example:
//
//	if err != nil {
//	    core.Fatal(err, "Failed to connect to the database")
//	}
var Fatal = fatal.OnErr

// Log records a non-critical error or warning to the error log file.
// This is for situations that are unexpected but do not require the
// application to stop.
//
// Example:
//
//	core.Log(err, "Could not process optional analytics event", "event-id-123")
var Log = logger.Log

// Settings returns the application's configuration, parsed from the settings.toml
// file into a structured format. It is loaded once and cached for performance.
//
// Example:
//
//	appPort := core.Settings().Port
var Settings = settings.GetSettings

// RawSettings returns the application's configuration as a raw map[string]interface{}.
// This is useful for packages that need to parse their own specific settings
// from the settings.toml file.
//
// Example:
//
//	raw, err := core.RawSettings()
//	if err == nil {
//	    mySetting := raw["my_package_setting"].(string)
//	}
var RawSettings = settings.GetRawSettings

// Validate returns a new validation bag. This is used to accumulate validation
// errors for a set of inputs, allowing for a full pre-flight response.
//
// Example:
//
//	bag := core.Validate()
//	bag.Validate(value, rule)
var Validate = validation.New

// SendEmail sends an email using the SMTP settings configured in settings.toml.
//
// Example:
//
//	err := core.SendEmail("user@example.com", "Welcome!", "Your account is ready.")
var SendEmail = email.Send

// RootPath returns the absolute path to the directory where the application
// executable is located.
var RootPath = files.RootPath

// DB returns a pointer to the established database connection.
var DB = db.GetDB

// --- Application Lifecycle ---

// Start initializes the application, runs pre-flight checks, registers features,
// and starts the command router.
func Start(featuresAvailable models.Features) {

	features.RegisterFeatures(featuresAvailable)

	cmd.CmdRouter()
}

// --- Utility Functions ---

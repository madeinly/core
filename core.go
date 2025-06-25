package core

import (
	"database/sql"

	"github.com/madeinly/core/internal/cmd"
	"github.com/madeinly/core/internal/db"
	"github.com/madeinly/core/internal/features"
	"github.com/madeinly/core/internal/settings"
	"github.com/madeinly/core/models"
)

func Start(featuresAvailable models.Features) {

	features.RegisterFeatures(featuresAvailable)

	cmd.CmdRouter()

}

var Settings = settings.Settings

func DB() *sql.DB {

	db := db.GetDB()

	return db
}

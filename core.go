package core

import (
	"database/sql"

	"github.com/madeinly/core/internal/cmd"
	"github.com/madeinly/core/internal/db"
	"github.com/madeinly/core/internal/features"
	"github.com/madeinly/core/models"
)

type api struct{} // unexported, keeps docs

var Madeinly api // exported value

func (api) Start(featuresAvailable models.Features) {

	features.RegisterFeatures(featuresAvailable)

	cmd.CmdRouter()

}

func (api) DB() *sql.DB {

	db := db.GetDB()

	return db
}

func Start(featuresAvailable models.Features) {

	features.RegisterFeatures(featuresAvailable)

	cmd.CmdRouter()

}

func DB() *sql.DB {

	db := db.GetDB()

	return db
}

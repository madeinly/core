package core

import (
	"database/sql"

	"github.com/MadeSimplest/core/internal/cmd"
	"github.com/MadeSimplest/core/internal/db"
	"github.com/MadeSimplest/core/internal/features"
	"github.com/MadeSimplest/core/models"
)

func Start(featuresAvailable models.Features) {

	features.RegisterFeatures(featuresAvailable)

	cmd.CmdRouter()

}

func DB() *sql.DB {

	db := db.GetDB()

	return db
}

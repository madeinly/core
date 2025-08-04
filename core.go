package core

import (
	"database/sql"
	"os"
	"path"

	"github.com/madeinly/core/fatal"
	"github.com/madeinly/core/internal/cmd"
	"github.com/madeinly/core/internal/db"
	"github.com/madeinly/core/internal/features"
	"github.com/madeinly/core/models"
)

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
		fatal.OnErr(err, "error getting the root path", err.Error())
	}

	return path.Dir(binPath)
}

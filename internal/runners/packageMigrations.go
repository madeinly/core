package runners

import (
	"fmt"
	"log"

	"github.com/madeinly/core/internal/db"
	"github.com/madeinly/core/models"
)

// RunAll executes all registered migrations
func RunMigrations(features models.Features) error {

	dbConn := db.GetDB()

	_, err := dbConn.Exec(db.InitialSchema)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Execute each migration in order
	for _, feature := range features {
		// Check if already executed
		var exists bool
		err := dbConn.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM _migrations WHERE name = ?)`,
			feature.Migration.Name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if exists {
			log.Printf("Migration %s already applied - skipping", feature.Migration.Name)
			continue
		}

		// Execute migration
		_, err = dbConn.Exec(feature.Migration.Schema)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", feature.Migration.Name, err)
		}

		// Record migration
		_, err = dbConn.Exec(`
			INSERT INTO _migrations (name) VALUES (?)`,
			feature.Migration.Name)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", feature.Migration.Name, err)
		}

		log.Printf("Successfully applied migration: %s", feature.Migration.Name)
	}

	return nil
}

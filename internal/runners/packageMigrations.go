package runners

import (
	"fmt"
	"log"

	"github.com/madeinly/core/internal/db"
	"github.com/madeinly/core/models"
)

// RunAll executes all registered migrations
func RunMigrations(features models.Features) error {

	db := db.GetDB()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS _migrations (
			name TEXT PRIMARY KEY,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Execute each migration in order
	for _, feature := range features {
		// Check if already executed
		var exists bool
		err := db.QueryRow(`
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
		_, err = db.Exec(feature.Migration.Schema)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", feature.Migration.Name, err)
		}

		// Record migration
		_, err = db.Exec(`
			INSERT INTO _migrations (name) VALUES (?)`,
			feature.Migration.Name)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", feature.Migration.Name, err)
		}

		log.Printf("Successfully applied migration: %s", feature.Migration.Name)
	}

	return nil
}

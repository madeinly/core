package db

import (
	_ "embed" // Required for //go:embed

	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	DB     *sql.DB
	DBOnce sync.Once
)

// Migration represents a single database migration
type Migration struct {
	Name   string
	Schema string
}

//go:embed queries/initial_schema.sql
var InitialSchema string

/*
go run wont work as go will create the binary at /tmp/gobuild**** and dispose the folder so the db
- will work on compiles
- will work with air
- will work when pointing directly to the ./tmp/main

if tmp is the folder it goes to ../ for convenience of air.
*/
func getDatabasePath() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("failed to get executable path: %w", err))
	}

	exeDir := filepath.Dir(exePath)
	dbName := "database.sqlite" // You can make this configurable if needed

	// If running in a tmp dir (like with `air`), place DB one level up
	if filepath.Base(exeDir) == "tmp" {
		return filepath.Join(exeDir, "..", dbName)
	}

	// Otherwise, place DB next to the binary
	return filepath.Join(exeDir, dbName)
}

func GetDB() *sql.DB {
	DBOnce.Do(func() {
		dbPath := getDatabasePath()

		// Ensure the directory exists (in case custom nested paths are used later)
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			panic(fmt.Errorf("failed to create database directory: %w", err))
		}

		// Open SQLite with WAL mode and foreign keys
		dataSource := fmt.Sprintf("file:%s?_foreign_keys=on&_journal_mode=WAL", dbPath)
		var err error
		DB, err = sql.Open("sqlite", dataSource)
		if err != nil {
			panic(fmt.Errorf("failed to open database: %w", err))
		}

		// Recommended SQLite settings
		DB.SetMaxOpenConns(1)
		DB.SetMaxIdleConns(1)
		DB.SetConnMaxLifetime(0)

		if err = DB.Ping(); err != nil {
			panic(fmt.Errorf("failed to ping database: %w", err))
		}
	})
	return DB
}

package db

import (
	"database/sql"
	_ "embed"
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

var dbName string = "db.sqlite"

func GetDB() *sql.DB {
	DBOnce.Do(func() {
		binaryPath, err := os.Executable()

		dbPath := filepath.Join(filepath.Dir(binaryPath), dbName)

		if err != nil {
			fmt.Println(err.Error())

		}

		// Ensure the directory exists (in case custom nested paths are used later)
		if err = os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			panic(fmt.Errorf("failed to create database directory: %w", err))
		}

		dataSource := fmt.Sprintf("file:%s?_foreign_keys=on&_journal_mode=WAL", dbPath)

		DB, err = sql.Open("sqlite", dataSource)
		if err != nil {
			panic(fmt.Errorf("failed to open database: %w", err))
		}

		DB.SetMaxOpenConns(1)
		DB.SetMaxIdleConns(1)
		DB.SetConnMaxLifetime(0)

		if err = DB.Ping(); err != nil {
			panic(fmt.Errorf("failed to ping database: %w", err))
		}
	})
	return DB
}

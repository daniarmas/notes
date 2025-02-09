package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/config"
)

// Open opens a database specified by its database driver name and a driver-specific data source name, usually consisting of at least a database name and connection information.
func Open(ctx context.Context, cfg *config.Configuration, showLog bool) *sql.DB {
	var db *sql.DB
	var err error

	// Retry settings
	const maxRetries = 5
	const initialBackoff = 1 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = sql.Open("pgx", cfg.DatabaseUrl)
		if err != nil {
			msg := fmt.Sprintf("Attempt %d: Failed to open the database: %v", attempt, err)
			clog.Warn(ctx, msg, err)
			time.Sleep(initialBackoff * time.Duration(attempt))
			continue
		}

		// Set connection pool settings
		db.SetMaxOpenConns(12)
		db.SetMaxIdleConns(12)
		db.SetConnMaxLifetime(30 * time.Minute)
		db.SetConnMaxIdleTime(30 * time.Minute)

		err = db.Ping()
		if err != nil {
			msg := fmt.Sprintf("Attempt %d: Failed to open the database: %v", attempt, err)
			clog.Warn(ctx, msg, err)
			time.Sleep(initialBackoff * time.Duration(attempt))
			continue
		}
		if showLog {
			clog.Info(ctx, "Connected to the database", nil)
		}
		return db
	}

	msg := fmt.Sprintf("Exceeded maximum retries. Failed to connect to the database: %v", err)
	clog.Error(ctx, msg, err)
	return nil
}

func Close(ctx context.Context, db *sql.DB, showLog bool) {
	db.Close()
	if showLog {
		clog.Info(ctx, "Database connection closed", nil)
	}
}

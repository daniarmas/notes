package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/daniarmas/clogg"
)

// Open opens a database specified by its database driver name and a driver-specific data source name, usually consisting of at least a database name and connection information.
func Open(ctx context.Context, databaseUrl string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	// Retry settings
	const maxRetries = 5
	const initialBackoff = 1 * time.Second
	const maxBackoff = 10 * time.Second
	backoff := initialBackoff

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = sql.Open("pgx", databaseUrl)
		if err != nil {
			clogg.Warn(ctx, "Failed to open the database", clogg.Int("attempt", attempt), clogg.String("error", err.Error()))
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * 1.5)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		// Set connection pool settings
		db.SetMaxOpenConns(12)
		db.SetMaxIdleConns(12)
		db.SetConnMaxLifetime(30 * time.Minute)
		db.SetConnMaxIdleTime(30 * time.Minute)

		err = db.PingContext(ctx)
		if err != nil {
			clogg.Warn(ctx, "Failed to open the database", clogg.Int("attempt", attempt), clogg.String("error", err.Error()))
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * 1.5)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Close(ctx context.Context, db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}

package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/daniarmas/notes/internal/config"
)

// Open opens a database specified by its database driver name and a driver-specific data source name, usually consisting of at least a database name and connection information.
func Open(cfg *config.Configuration) *sql.DB {
	var db *sql.DB
	var err error

	// Retry settings
	const maxRetries = 5
	const initialBackoff = 1 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = sql.Open("pgx", cfg.DatabaseUrl)
		if err != nil {
			log.Printf("Attempt %d: Failed to open the database: %v", attempt, err)
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
			log.Printf("Attempt %d: Failed to connect to the database: %v", attempt, err)
			time.Sleep(initialBackoff * time.Duration(attempt))
			continue
		}

		log.Println("Connected to the database")
		return db
	}

	log.Fatalf("Exceeded maximum retries. Failed to connect to the database: %v", err)
	return nil
}

func Close(db *sql.DB) {
	db.Close()
	log.Println("Database connection closed")
}

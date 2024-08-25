package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/daniarmas/notes/internal/config"
)

// Open opens a database specified by its database driver name and a driver-specific data source name, usually consisting of at least a database name and connection information.
func Open(cfg *config.Configuration) *sql.DB {
	// Database connection
	db, err := sql.Open("pgx", cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(12)
	db.SetMaxIdleConns(12)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Connected to the database")
	return db
}

func Close(db *sql.DB) {
	db.Close()
	log.Println("Database connection closed")
}

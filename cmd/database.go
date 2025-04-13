/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/daniarmas/clogg"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/database"
	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Set up clogg
		handler := slog.NewJSONHandler(os.Stdout, nil)
		logger := clogg.GetLogger(clogg.LoggerConfig{
			BufferSize: 100,
			Handler:    handler,
		})
		defer logger.Shutdown()

		// Config
		cfg := config.LoadServerConfig()

		// Database connection
		db, err := database.Open(ctx, cfg.DatabaseUrl)
		if err != nil {
			clogg.Error(ctx, "error opening database", clogg.String("error", err.Error()))
			os.Exit(1)
		}
		defer database.Close(ctx, db)

		// Create notes_database if not exists
		stmt, err := db.Prepare(`
			DO $$
			BEGIN
			   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'postgres') THEN
			      PERFORM pg_sleep(0.1); -- Workaround for the DO check
			      CREATE DATABASE postgres;
			   END IF;
			END
			$$;`)
		if err != nil {
			clogg.Error(ctx, "error preparing sql to create notes_database", clogg.String("error", err.Error()))
		}
		_, err = stmt.Exec()
		if err != nil {
			clogg.Error(ctx, "error creating notes_database", clogg.String("error", err.Error()))
		}

		// Create users table if not exists
		stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
    		id UUID DEFAULT gen_random_uuid(),
    		name VARCHAR NOT NULL,
    		email VARCHAR NOT NULL UNIQUE,
			password VARCHAR NOT NULL,
    		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    		update_time TIMESTAMP,
			CONSTRAINT users_pk PRIMARY KEY (id)
		);`)
		if err != nil {
			clogg.Error(ctx, "error preparing sql to create users table", clogg.String("error", err.Error()))
		}
		_, err = stmt.Exec()
		if err != nil {
			clogg.Error(ctx, "error creating users table", clogg.String("error", err.Error()))
		}

		// Create refresh tokens table if not exists
		stmt, err = db.Prepare(`
			CREATE TABLE IF NOT EXISTS refresh_tokens (
				id UUID DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
    			create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    			update_time TIMESTAMP,
				CONSTRAINT refresh_tokens_pk PRIMARY KEY (id),
				CONSTRAINT fk_user
        			FOREIGN KEY (user_id) 
        			REFERENCES users(id)
        			ON DELETE CASCADE
			)
		`)
		if err != nil {
			clogg.Error(ctx, "error preparing sql to create refresh_tokens table", clogg.String("error", err.Error()))
		}
		_, err = stmt.Exec()
		if err != nil {
			clogg.Error(ctx, "error creating refresh_tokens table", clogg.String("error", err.Error()))
		}

		// Create access tokens table if not exists
		stmt, err = db.Prepare(`
			CREATE TABLE IF NOT EXISTS access_tokens (
				id UUID DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
				refresh_token_id UUID NOT NULL,
    			create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    			update_time TIMESTAMP,
				CONSTRAINT access_tokens_pk PRIMARY KEY (id),
				CONSTRAINT fk_user
        			FOREIGN KEY (user_id) 
        			REFERENCES users(id)
        			ON DELETE CASCADE,
				CONSTRAINT fk_refresh_token
        			FOREIGN KEY (refresh_token_id) 
        			REFERENCES refresh_tokens(id)
        			ON DELETE CASCADE
			)
		`)
		if err != nil {
			clogg.Error(ctx, "error preparing sql to create access_tokens table", clogg.String("error", err.Error()))
		}
		_, err = stmt.Exec()
		if err != nil {
			clogg.Error(ctx, "error creating access_tokens table", clogg.String("error", err.Error()))
		}

		// Create notes table if not exists
		stmt, err = db.Prepare(`
			CREATE TABLE IF NOT EXISTS notes (
				id UUID DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
				title VARCHAR,
				content VARCHAR,
    			create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    			update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    			delete_time TIMESTAMP,
				CONSTRAINT notes_pk PRIMARY KEY (id),
				CONSTRAINT fk_user
        			FOREIGN KEY (user_id) 
        			REFERENCES users(id)
        			ON DELETE CASCADE
			)
		`)
		if err != nil {
			clogg.Error(ctx, "error preparing sql to create notes table", clogg.String("error", err.Error()))
		}
		_, err = stmt.Exec()
		if err != nil {
			clogg.Error(ctx, "error creating notes table", clogg.String("error", err.Error()))
		}

		// Create files table if not exists
		stmt, err = db.Prepare(`
			CREATE TABLE IF NOT EXISTS files (
				id UUID DEFAULT gen_random_uuid(),
				processed_file VARCHAR,
				original_file VARCHAR NOT NULL,
				note_id UUID NOT NULL,
				create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
				update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
				delete_time TIMESTAMP,
				CONSTRAINT pk PRIMARY KEY (id),
				CONSTRAINT fk_note
					FOREIGN KEY (note_id) 
					REFERENCES notes(id)
					ON DELETE CASCADE
			)
		`)
		if err != nil {
			clogg.Error(ctx, "error preparing sql to create files table", clogg.String("error", err.Error()))
		}
		_, err = stmt.Exec()
		if err != nil {
			clogg.Error(ctx, "error creating files table", clogg.String("error", err.Error()))
		}

		clogg.Info(ctx, "Database tables created successfully")
	},
}

func init() {
	createCmd.AddCommand(databaseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// databaseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// databaseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

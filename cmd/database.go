/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

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
		// Config
		cfg := config.LoadConfig()

		// Database connection
		db := database.Open(cfg, false)
		defer database.Close(db, false)

		// Create notes_database if not exists
		stmt, err := db.Prepare("CREATE DATABASE IF NOT EXISTS postgres;")
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Fatal(err)
		}

		// Create users table if not exists
		stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
    		id UUID DEFAULT gen_random_uuid(),
    		name STRING NOT NULL,
    		email STRING NOT NULL UNIQUE,
			password STRING NOT NULL,
    		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    		update_time TIMESTAMP,
			CONSTRAINT pk PRIMARY KEY (id)
		);`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Fatal(err)
		}

		// Create refresh tokens table if not exists
		stmt, err = db.Prepare(`
			CREATE TABLE IF NOT EXISTS refresh_tokens (
				id UUID DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
    			create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    			update_time TIMESTAMP,
				CONSTRAINT pk PRIMARY KEY (id),
				CONSTRAINT fk_user
        			FOREIGN KEY (user_id) 
        			REFERENCES users(id)
        			ON DELETE CASCADE
			)
		`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Fatal(err)
		}

		// Create access tokens table if not exists
		stmt, err = db.Prepare(`
			CREATE TABLE IF NOT EXISTS access_tokens (
				id UUID DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
				refresh_token_id UUID NOT NULL,
    			create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    			update_time TIMESTAMP,
				CONSTRAINT pk PRIMARY KEY (id),
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
			log.Fatal(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Fatal(err)
		}

		// Create notes table if not exists
		stmt, err = db.Prepare(`
			CREATE TABLE IF NOT EXISTS notes (
				id UUID DEFAULT gen_random_uuid(),
				user_id UUID NOT NULL,
				title STRING,
				content STRING,
    			create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    			update_time TIMESTAMP,
    			delete_time TIMESTAMP,
				CONSTRAINT pk PRIMARY KEY (id),
				CONSTRAINT fk_user
        			FOREIGN KEY (user_id) 
        			REFERENCES users(id)
        			ON DELETE CASCADE
			)
		`)
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec()
		if err != nil {
			log.Fatal(err)
		}
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

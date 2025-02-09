/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/data"
	"github.com/daniarmas/notes/internal/database"
	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with mock data",
	Long:  `Seed the database with mock data to simplify testing the app`,
	Run: func(cmd *cobra.Command, args []string) {
		// Context
		ctx := context.Background()

		// Config
		cfg := config.LoadServerConfig()

		// Database connection
		db := database.Open(ctx, cfg, false)
		defer database.Close(ctx, db, false)

		// Get sqlc queries
		queries := database.New(db)

		// Datasources
		hashDs := data.NewBcryptHashDatasource()

		// Create users
		user1Pass, _ := hashDs.Hash("user1")
		user2Pass, _ := hashDs.Hash("user2")
		_, err := queries.CreateUser(ctx, database.CreateUserParams{Name: "user1", Email: "user1@email.com", Password: user1Pass})
		if err != nil {
			clog.Error(ctx, "error creating user 1", err)
		}
		_, err = queries.CreateUser(ctx, database.CreateUserParams{Name: "user2", Email: "user2@email.com", Password: user2Pass})
		if err != nil {
			clog.Error(ctx, "error creating user 2", err)
		}
		clog.Info(ctx, "Database seeding completed successfully", nil)
	},
}

func init() {
	createCmd.AddCommand(seedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

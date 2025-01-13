/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/daniarmas/notes/internal/clog"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/data"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/oss"
	"github.com/spf13/cobra"
)

// Declare a variable to hold the flag values
var files []string

// processFileCmd represents the processFile command
var processFileCmd = &cobra.Command{
	Use:   "process-files",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Custom logger
		clog.NewClog()

		// Config
		cfg := config.LoadServerConfig()

		// Database connection
		db := database.Open(cfg, true)
		defer database.Close(db, true)

		// Database queries
		dbQueries := database.New(db)

		// Object storage service
		oss := oss.NewDigitalOceanWithMinio(cfg)
		// Healthcheck
		if err := oss.HealthCheck(); err != nil {
			clog.Error(ctx, "error checking object storage service health", err)
		}

		// Datasources
		fileDatabaseDs := data.NewFileDatabaseDs(dbQueries)

		// // Repositories
		fileRepository := domain.NewFileRepository(fileDatabaseDs, oss, cfg)

		// Process files
		if _, err := fileRepository.Process(ctx, "original/6f36c25d-2e1a-47cb-a308-7fa223fc1df2.jpg"); err != nil {
			clog.Error(ctx, "error processing file", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(processFileCmd)
	processFileCmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "Files to process")
}

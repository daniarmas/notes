/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/daniarmas/clogg"
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
		db, err := database.Open(ctx, cfg)
		if err != nil {
			clogg.Error(ctx, "error opening database", clogg.String("error", err.Error()))
			os.Exit(1)
		}
		defer database.Close(ctx, db)

		// Database queries
		dbQueries := database.New(db)

		// Object storage service
		oss := oss.NewDigitalOceanWithMinio(cfg)

		// Healthcheck
		if err := oss.HealthCheck(); err != nil {
			clogg.Error(ctx, "error checking object storage service health", clogg.String("error", err.Error()))
		}

		// Datasources
		fileDatabaseDs := data.NewFileDatabaseDs(dbQueries)

		// Repositories
		fileRepository := domain.NewFileRepository(fileDatabaseDs, oss, cfg)

		// Access files
		files, err := cmd.Flags().GetStringSlice("files")
		if err != nil {
			clogg.Error(ctx, "error getting files flag", clogg.String("error", err.Error()))
			os.Exit(1)
		}

		// Process files concurrently
		for _, file := range files {
			// Start the sql transaction
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				clogg.Error(ctx, "error starting transaction", clogg.String("error", err.Error()))
				os.Exit(0)
			}

			// Defer the transaction rollback or commit
			defer func() {
				if p := recover(); p != nil {
					tx.Rollback()
					panic(p)
				} else if err != nil {
					tx.Rollback()
				} else {
					err = tx.Commit()
				}
			}()

			if err := fileRepository.Process(ctx, tx, file); err != nil {
				clogg.Error(ctx, "error processing file", clogg.String("error", err.Error()))
			}
		}

		// Command completed successfully
		clogg.Info(ctx, "files processed successfully")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(processFileCmd)
	processFileCmd.Flags().StringSliceVarP(&files, "files", "f", []string{}, "Files to process")
	// Mark the files flag as required
	processFileCmd.MarkFlagRequired("files")
}

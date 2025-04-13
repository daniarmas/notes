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

// tablesCmd represents the tables command
var tablesCmd = &cobra.Command{
	Use:   "tables",
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
			BufferSize: 5,
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
		stmt, err := db.Prepare("DELETE FROM users;")
		if err != nil {
			clogg.Error(ctx, "error preparing sql to delete users", clogg.String("error", err.Error()))
			os.Exit(1)
		}
		_, err = stmt.Exec()
		if err != nil {
			clogg.Error(ctx, "error deleting users", clogg.String("error", err.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	deleteCmd.AddCommand(tablesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tablesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tablesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

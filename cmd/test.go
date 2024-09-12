/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"time"

	"github.com/daniarmas/notes/internal/cache"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/data"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Context
		ctx := context.Background()

		// Config
		cfg := config.LoadConfig()

		// Cache connection
		rdb := cache.OpenRedis(cfg)
		defer cache.CloseRedis(rdb)

		// Datasources
		userCacheDs := data.NewUserCacheDs(rdb)

		// Seed cache with mocked data...
		uuid1, _ := uuid.NewRandom()
		userCacheDs.CreateUser(ctx, &domain.User{Id: uuid1, Name: "user1", Email: "user1@email.com", Password: "user1", CreateTime: time.Now(), UpdateTime: time.Now()})
		userCacheDs.GetUserById(ctx, uuid1)
		userCacheDs.UpdateUser(ctx, &domain.User{Id: uuid1, Name: "user2", Email: "user1@email.com", Password: "user1", CreateTime: time.Now(), UpdateTime: time.Now()})
		userCacheDs.DeleteUser(ctx, uuid1)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

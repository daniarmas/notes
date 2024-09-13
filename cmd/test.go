/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/daniarmas/notes/internal/cache"
	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/data"
	"github.com/daniarmas/notes/internal/database"
	"github.com/daniarmas/notes/internal/domain"
	"github.com/daniarmas/notes/internal/service"
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

		// Database connection
		db := database.Open(cfg, true)
		defer database.Close(db, true)

		// Database queries
		dbQueries := database.New(db)

		// Cache connection
		rdb := cache.OpenRedis(cfg)
		defer cache.CloseRedis(rdb)

		// Datasources
		hashDatasource := data.NewBcryptHashDatasource()
		jwtDatasource := domain.NewJWTDatasource(cfg)
		userCacheDs := data.NewUserCacheDs(rdb)
		userDatabaseDs := data.NewUserDatabaseDs(dbQueries)
		accessTokenCacheDs := data.NewAccessTokenTokenCacheDs(rdb)
		accessTokenDatabaseDs := data.NewAccessTokenDatabaseDs(dbQueries)
		refreshTokenCacheDs := data.NewRefreshTokenCacheDs(rdb)
		refreshTokenDatabaseDs := data.NewRefreshTokenDatabaseDs(dbQueries)

		// Repositories
		userRepository := domain.NewUserRepository(&userCacheDs, &userDatabaseDs)
		accessTokenRepository := domain.NewAccessTokenRepository(&accessTokenCacheDs, &accessTokenDatabaseDs)
		refreshTokenRepository := domain.NewRefreshTokenRepository(&refreshTokenCacheDs, &refreshTokenDatabaseDs)

		// Services
		authenticationService := service.NewAuthenticationService(jwtDatasource, hashDatasource, userRepository, accessTokenRepository, refreshTokenRepository)

		// Tests
		res, err := authenticationService.SignIn(ctx, "user1@email.com", "user1")
		if err != nil {
			log.Println(err)
		} else {
			log.Println(res)
		}
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

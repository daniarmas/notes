package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/daniarmas/notes/internal/config"
	"github.com/daniarmas/notes/internal/server"
	"github.com/spf13/cobra"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.LoadConfig()
	srv := server.NewServer(cfg)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.AppRestHost, cfg.AppRestPort),
		Handler: srv,
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		if err := run(ctx, os.Stdout, os.Args); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

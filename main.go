package main

import (
	"embed"
	"fmt"

	"github.com/captain-corp/captain/cmd"
	"github.com/captain-corp/captain/config"
	"github.com/captain-corp/captain/db"
	"github.com/captain-corp/captain/server"
	"github.com/captain-corp/captain/system"

	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/cobra"
)

//go:embed embedded/admin/static/css/*
//go:embed embedded/admin/static/js/*
//go:embed embedded/admin/static/img/*
//go:embed embedded/admin/static/fonts/*
//go:embed embedded/admin/templates/includes/*
//go:embed embedded/admin/templates/*
//go:embed embedded/public/templates/*
//go:embed embedded/public/static/css/*
//go:embed embedded/public/static/js/*
//go:embed embedded/public/static/img/*
var embeddedFS embed.FS

var (
	initDevDB  bool
	configFile string
	serverHost string
	serverPort int
)

func main() {
	var rootCmd = &cobra.Command{Use: "captain"}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Captain v%s\n", system.Version)
			fmt.Printf("Commit: %s\n", system.Commit)
			fmt.Printf("Built: %s\n", system.Date)
		},
	}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the server",
		Run:   runServer,
	}

	runCmd.Flags().BoolVarP(&initDevDB, "init-dev-db", "i", false, "Initialize the development database with test data")
	runCmd.Flags().StringVarP(&serverHost, "bind", "b", "", "Address to bind to (overrides config)")
	runCmd.Flags().IntVarP(&serverPort, "port", "p", 0, "Server port (overrides config)")

	var userCmd = &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	var userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run:   cmd.CreateUser,
	}

	var userUpdatePasswordCmd = &cobra.Command{
		Use:   "update-password",
		Short: "Update user password",
		Run:   cmd.UpdateUserPassword,
	}

	userCmd.AddCommand(userCreateCmd, userUpdatePasswordCmd)
	rootCmd.AddCommand(runCmd, userCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	var srv *server.Server

	// Load config
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Override config with command line flags if provided
	if serverHost != "" {
		cfg.Server.Host = serverHost
	}
	if serverPort != 0 {
		cfg.Server.Port = serverPort
	}

	// Initialize database
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize development database if requested
	if initDevDB {
		if err := db.InsertTestData(database); err != nil {
			log.Fatalf("Failed to insert test data: %v", err)
		}
	}

	// Validate S3 configuration if S3 provider is selected
	if err := cfg.ValidateS3Config(); err != nil {
		log.Fatalf("S3 configuration error: %v", err)
	}

	// Create and start server
	if srv, err = server.New(database, cfg, embeddedFS); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Run the server
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

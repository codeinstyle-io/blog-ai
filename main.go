package main

import (
	"fmt"
	"log"
	"os"

	"codeinstyle.io/captain/cli"
	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/handlers"
	"codeinstyle.io/captain/middleware"
	"codeinstyle.io/captain/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	initDevDB  bool
	configFile string
)

func main() {
	var rootCmd = &cobra.Command{Use: "captain"}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file path")

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the server",
		Run:   runServer,
	}

	runCmd.Flags().BoolVarP(&initDevDB, "init-dev-db", "i", false, "Initialize the development database with test data")

	var userCmd = &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	var userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run:   cli.CreateUser,
	}

	var userUpdatePasswordCmd = &cobra.Command{
		Use:   "update-password",
		Short: "Update user password",
		Run:   cli.UpdateUserPassword,
	}

	userCmd.AddCommand(userCreateCmd, userUpdatePasswordCmd)
	rootCmd.AddCommand(runCmd, userCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database := db.InitDB(cfg)
	r := gin.Default()

	if initDevDB {
		err := db.InsertTestData(database)
		if err != nil {
			log.Printf("Failed to insert test data: %v\n", err)
		}
	}

	// Serve static files
	r.Static("/static", "static")

	// Custom template functions
	r.SetFuncMap(utils.GetTemplateFuncs())

	// Load templates
	r.LoadHTMLGlob("templates/**/*")

	// Register routes with config
	handlers.RegisterPublicRoutes(r, database, cfg)
	handlers.RegisterAdminRoutes(r, database, cfg)

	// Add middleware to load menu items
	r.Use(middleware.LoadMenuItems(database))

	fmt.Printf("Server running on http://%s:%d\n", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

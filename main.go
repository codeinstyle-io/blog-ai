package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/handlers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

var (
	port      int
	host      string
	initDevDB bool
)

func main() {
	var rootCmd = &cobra.Command{Use: "captain"}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Runs the server",
		Run:   runServer,
	}

	runCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	runCmd.Flags().StringVarP(&host, "bind", "b", "localhost", "Host to run the server on")
	runCmd.Flags().BoolVarP(&initDevDB, "init-dev-db", "i", false, "Initialize the development database with test data")

	var userCmd = &cobra.Command{
		Use:   "user",
		Short: "User management commands",
	}

	var userCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run:   createUser,
	}

	var userUpdatePasswordCmd = &cobra.Command{
		Use:   "update-password",
		Short: "Update user password",
		Run:   updateUserPassword,
	}

	userCmd.AddCommand(userCreateCmd, userUpdatePasswordCmd)
	rootCmd.AddCommand(runCmd, userCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	database := db.InitDB()
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
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	})

	// Load templates
	r.LoadHTMLGlob("templates/**/*")

	// Register routes
	handlers.RegisterPublicRoutes(r, database)

	handlers.RegisterAdminRoutes(r, database)

	fmt.Printf("Server running on http://%s:%d\n", host, port)
	r.Run(fmt.Sprintf("%s:%d", host, port))
}

func createUser(cmd *cobra.Command, args []string) {
	var firstName, lastName, email, password string
	fmt.Print("First Name: ")
	fmt.Scanln(&firstName)
	fmt.Print("Last Name: ")
	fmt.Scanln(&lastName)
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	fmt.Print("Password: ")
	fmt.Scanln(&password)

	if err := validateUserInput(firstName, lastName, email, password); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	database := db.InitDB()
	user := db.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashedPassword),
	}

	if err := database.Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v\n", err)
		return
	}

	fmt.Println("User created successfully")
}

func updateUserPassword(cmd *cobra.Command, args []string) {
	var email, oldPassword, newPassword, confirmPassword string

	fmt.Print("Email: ")
	fmt.Scanln(&email)

	database := db.InitDB()
	var user db.User
	if err := database.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println("User not found")
		return
	}

	fmt.Print("Old Password: ")
	fmt.Scanln(&oldPassword)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		fmt.Println("Incorrect password")
		return
	}

	fmt.Print("New Password: ")
	fmt.Scanln(&newPassword)
	fmt.Print("Confirm Password: ")
	fmt.Scanln(&confirmPassword)

	if newPassword != confirmPassword {
		fmt.Println("Passwords don't match")
		return
	}

	if err := validatePassword(newPassword); err != nil {
		fmt.Printf("Password validation error: %v\n", err)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := database.Save(&user).Error; err != nil {
		log.Printf("Failed to update password: %v\n", err)
		return
	}

	fmt.Println("Password updated successfully")
}

func validateUserInput(firstName, lastName, email, password string) error {
	if len(firstName) < 1 || len(firstName) > 255 {
		return fmt.Errorf("first name must be between 1 and 255 characters")
	}
	if len(lastName) < 1 || len(lastName) > 255 {
		return fmt.Errorf("last name must be between 1 and 255 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(strings.ToLower(email)) {
		return fmt.Errorf("invalid email format")
	}

	return validatePassword(password)
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character (!@#$%%^&*(),.?\":{}|<>)")
	}
	return nil
}

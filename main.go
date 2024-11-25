package main

import (
	"fmt"
	"log"
	"text/template"

	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.InitDB()

	if err := db.InsertTestData(database); err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	r := gin.Default()
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

	// Register all routes
	handlers.RegisterPublicRoutes(r, database)
	handlers.RegisterAdminRoutes(r, database)

	fmt.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}

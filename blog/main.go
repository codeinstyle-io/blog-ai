package main

import (
	"fmt"
	"log"

	"codeinstyle.io/blog/db"
	"codeinstyle.io/blog/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	database := db.InitDB()

	if err := db.InsertTestData(database); err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static", "static")

	// Register all routes
	handlers.RegisterPublicRoutes(r, database)
	handlers.RegisterAdminRoutes(r, database)

	fmt.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}

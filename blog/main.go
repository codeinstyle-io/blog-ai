package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeinstyle.io/blog/db"
	"codeinstyle.io/blog/handlers"
	"codeinstyle.io/blog/middleware"
	"codeinstyle.io/blog/types"
	"github.com/gin-gonic/gin"
)

func getSkills() ([]types.SkillSection, error) {
	data, err := os.ReadFile("data/skills.json")
	if err != nil {
		return nil, err
	}

	var skills []types.SkillSection
	err = json.Unmarshal(data, &skills)
	if err != nil {
		return nil, err
	}

	return skills, nil
}

func main() {
	database := db.InitDB()

	// Insert test data if the database is empty
	if err := db.InsertTestData(database); err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	postHandlers := handlers.NewPostHandlers(database)
	authHandlers := handlers.NewAuthHandlers(database)

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static", "static")

	// Auth routes
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.tmpl", nil) })
	r.POST("/login", authHandlers.Login)
	r.GET("/logout", authHandlers.Logout)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired(database))
	{
		authorized.GET("/posts/create", func(c *gin.Context) { c.HTML(http.StatusOK, "create_post.tmpl", nil) })
		authorized.POST("/posts", postHandlers.CreatePost)
		authorized.GET("/posts/edit/:id", func(c *gin.Context) {
			id := c.Param("id")
			var post db.Post
			if err := database.First(&post, id).Error; err != nil {
				c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
					"title": "Post not found",
				})
				return
			}
			c.HTML(http.StatusOK, "edit_post.tmpl", gin.H{
				"post": post,
			})
		})
		authorized.PUT("/posts/:id", func(c *gin.Context) {
			var post db.Post
			id := c.Param("id")
			if err := database.First(&post, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
				return
			}
			if err := c.ShouldBindJSON(&post); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := database.Save(&post).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
				return
			}
			c.JSON(http.StatusOK, post)
		})
	}

	// Public routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	r.GET("/posts", postHandlers.ListPosts)

	r.GET("/posts/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		post, err := db.GetPostBySlug(database, slug)

		if err != nil {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
				"title": "Post not found",
			})
			return
		}
		c.HTML(http.StatusOK, "post.tmpl", gin.H{
			"title": post.Title,
			"post":  post,
		})
	})

	r.GET("/skills", func(c *gin.Context) {
		skills, err := getSkills()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read skills data"})
			return
		}
		c.JSON(http.StatusOK, skills)
	})

	r.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact.tmpl", gin.H{})
	})

	r.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl", gin.H{})
	})

	fmt.Println("Server running on http://localhost:8080")

	r.Run(":8080")
}

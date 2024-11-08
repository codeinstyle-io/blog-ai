package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeinstyle.io/blog/db"
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

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static", "static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	r.GET("/posts", func(c *gin.Context) {
		posts, err := db.GetPosts(database, 5)

		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
				"title": "Error",
			})
			return
		}
		c.HTML(http.StatusOK, "posts.tmpl", gin.H{
			"title": "Latest Posts",
			"posts": posts,
		})
	})

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

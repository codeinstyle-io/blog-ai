package db

import (
	"encoding/json"
	"fmt"
	mathrand "math/rand/v2"
	"os"
	"time"

	"captain-corp/captain/models"

	"gorm.io/gorm"
)

// Add helper function for random tag selection
func getRandomTags(tags []models.Tag, min, max int) []models.Tag {
	if len(tags) == 0 {
		return []models.Tag{}
	}

	// Get random count between min and max
	count := min + mathrand.IntN(max-min+1)
	if count > len(tags) {
		count = len(tags)
	}

	// Shuffle tags
	shuffled := make([]models.Tag, len(tags))
	copy(shuffled, tags)
	mathrand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:count]
}

type testData struct {
	Tags  []string   `json:"tags"`
	Posts []testPost `json:"posts"`
}

type testPost struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"content"`
	PublishedAt string `json:"publishedAt"`
	Visible     bool   `json:"visible"`
	Excerpt     string `json:"excerpt"`
}

func InsertTestData(db *gorm.DB) error {
	var count int64
	err := db.Model(&models.Post{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		// Create default author for test posts
		author := models.User{
			FirstName: "Test",
			LastName:  "Author",
			Email:     "test@example.com",
			Password:  "hashed_password", // In real app, this should be properly hashed
		}
		if err := db.FirstOrCreate(&author, models.User{Email: "test@example.com"}).Error; err != nil {
			return err
		}

		// Read test data
		data, err := os.ReadFile("data/test_posts.json")
		if err != nil {
			return err
		}

		var testData testData
		if err := json.Unmarshal(data, &testData); err != nil {
			return err
		}

		// Create tags
		tags := make([]models.Tag, len(testData.Tags))
		for i, name := range testData.Tags {
			tag := models.Tag{Name: name}
			if err := db.FirstOrCreate(&tag, models.Tag{Name: name}).Error; err != nil {
				return err
			}
			tags[i] = tag
		}

		// Create posts
		for _, p := range testData.Posts {
			// Parse relative date
			days := 0
			if n, err := fmt.Sscanf(p.PublishedAt, "-%dd", &days); err != nil || n != 1 {
				return fmt.Errorf("invalid publishedAt format: %s", p.PublishedAt)
			}

			post := models.Post{
				Title:       p.Title,
				Slug:        p.Slug,
				Content:     p.Content,
				PublishedAt: time.Now().AddDate(0, 0, -days),
				Visible:     p.Visible,
				Excerpt:     &p.Excerpt,
				Tags:        getRandomTags(tags, 2, 4),
				AuthorID:    author.ID, // Set the author ID for test posts
			}

			if err := db.Create(&post).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

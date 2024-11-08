package db

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	err = db.AutoMigrate(&Post{}, &Tag{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetPosts(db *gorm.DB, limit int) ([]Post, error) {
	var posts []Post
	result := db.Preload("Tags").
		Where("published_at <= ?", time.Now()).
		Where("visible = ?", true).
		Order("id desc").
		Limit(limit).
		Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func GetPostBySlug(db *gorm.DB, slug string) (Post, error) {
	var post Post
	result := db.Preload("Tags").
		Where("slug = ?", slug).
		Where("published_at <= ?", time.Now()).
		Where("visible = ?", true).
		First(&post)
	if result.Error != nil {
		return post, result.Error
	}
	return post, nil
}

func InsertTestData(db *gorm.DB) error {
	var count int64
	err := db.Model(&Post{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		// Create tags
		tags := []Tag{
			{Name: "Go"},
			{Name: "Programming"},
			{Name: "Tutorial"},
		}
		for _, tag := range tags {
			err := db.FirstOrCreate(&tag, Tag{Name: tag.Name}).Error
			if err != nil {
				return err
			}
		}

		// Create test posts with associated tags and excerpts
		excerpt1 := "This is an excerpt for the first post."
		testPosts := []Post{
			{
				Title:       "First Post",
				Slug:        "first-post",
				Content:     "This is the first post.",
				PublishedAt: time.Now(),
				Visible:     true,
				Excerpt:     &excerpt1,
				Tags:        []Tag{tags[0], tags[1]},
			},
			{
				Title:       "Second Post",
				Slug:        "second-post",
				Content:     "This is the second post.",
				PublishedAt: time.Now(),
				Visible:     true,
				// Excerpt is nil; an extract of Content will be used
				Tags: []Tag{tags[1], tags[2]},
			},
			{
				Title:       "Third Post",
				Slug:        "third-post",
				Content:     "This is the third post.",
				PublishedAt: time.Now(),
				Visible:     true,
				Excerpt:     nil, // Excerpt is nil
				Tags:        []Tag{tags[0], tags[2]},
			},
		}

		for _, post := range testPosts {
			err := db.Create(&post).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

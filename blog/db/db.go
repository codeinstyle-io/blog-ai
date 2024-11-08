package db

import (
	"log"

	"codeinstyle.io/blog/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	err = db.AutoMigrate(&types.Post{}, &types.Tag{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetPosts(db *gorm.DB, limit int) ([]types.Post, error) {
	var posts []types.Post
	result := db.Order("id desc").Limit(limit).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func GetPostBySlug(db *gorm.DB, slug string) (types.Post, error) {
	var post types.Post
	result := db.Where("slug = ?", slug).First(&post)
	if result.Error != nil {
		return post, result.Error
	}
	return post, nil
}

func InsertTestData(db *gorm.DB) error {
	var count int64
	err := db.Model(&types.Post{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		testPosts := []types.Post{
			{Title: "First Post", Slug: "first-post", Content: "This is the first post."},
			{Title: "Second Post", Slug: "second-post", Content: "This is the second post."},
			{Title: "Third Post", Slug: "third-post", Content: "This is the third post."},
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

package db

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	err = db.AutoMigrate(&Post{}, &Tag{}, &User{})
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

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

// New function to generate random session token
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GetUserByToken(db *gorm.DB, token string) (*User, error) {
	var user User
	result := db.Where("session_token = ?", token).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUserSessionToken(db *gorm.DB, user *User) error {
	token, err := generateSessionToken()
	if err != nil {
		return err
	}
	user.SessionToken = token
	return db.Save(user).Error
}

func InsertTestData(db *gorm.DB) error {
	var count int64
	err := db.Model(&Post{}).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		// Create test user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		token, err := generateSessionToken()
		if err != nil {
			return err
		}

		testUser := &User{
			Email:        "test@example.com",
			Password:     string(hashedPassword),
			SessionToken: token,
		}
		if err := db.Create(testUser).Error; err != nil {
			return err
		}

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

package repository

import (
	"testing"
	"time"

	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate the schema
	err = db.ExecuteMigrations(database)
	require.NoError(t, err)

	return database
}

func TestPostRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)

	post := &models.Post{
		Title:   "Test Post",
		Slug:    "test-post",
		Content: "Test Content",
	}

	err := repo.Create(post)
	assert.NoError(t, err)
	assert.NotZero(t, post.ID)
}

func TestPostRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)

	post := &models.Post{
		Title:   "Test Post",
		Slug:    "test-post",
		Content: "Test Content",
	}

	err := repo.Create(post)
	require.NoError(t, err)

	found, err := repo.FindByID(post.ID)
	assert.NoError(t, err)
	assert.Equal(t, post.Title, found.Title)
	assert.Equal(t, post.Slug, found.Slug)
	assert.Equal(t, post.Content, found.Content)
}

func TestPostRepository_FindBySlug(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)

	post := &models.Post{
		Title:   "Test Post",
		Slug:    "test-post",
		Content: "Test Content",
	}

	err := repo.Create(post)
	require.NoError(t, err)

	found, err := repo.FindBySlug(post.Slug)
	assert.NoError(t, err)
	assert.Equal(t, post.Title, found.Title)
	assert.Equal(t, post.Slug, found.Slug)
	assert.Equal(t, post.Content, found.Content)
}

func TestPostRepository_FindVisible(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPostRepository(db)

	posts := []*models.Post{
		{
			Title:       "Visible Post 1",
			Slug:        "visible-post-1",
			Content:     "Content 1",
			PublishedAt: time.Now(),
			Visible:     true,
		},
		{
			Title:       "Hidden Post",
			Slug:        "hidden-post",
			Content:     "Content 2",
			PublishedAt: time.Now(),
			Visible:     false,
		},
		{
			Title:       "Visible Post 2",
			Slug:        "visible-post-2",
			Content:     "Content 3",
			PublishedAt: time.Now(),
			Visible:     true,
		},
	}

	for _, p := range posts {
		err := repo.Create(p)
		require.NoError(t, err)
	}

	found, total, err := repo.FindVisible(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total) // Only visible posts
	assert.Len(t, found, 2)
}

package repository

import (
	"testing"
	"time"

	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/models"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// Use a fixed time for testing
	timezone := "Europe/Paris"
	loc, err := time.LoadLocation(timezone)
	require.NoError(t, err)
	now := time.Date(2024, 12, 18, 22, 17, 0, 0, loc) // Fixed time from the context

	posts := []*models.Post{
		{
			Title:       "Visible Post 1",
			Slug:        "visible-post-1",
			Content:     "Content 1",
			PublishedAt: now.Add(-24 * time.Hour), // Yesterday
			Visible:     true,
		},
		{
			Title:       "Hidden Post",
			Slug:        "hidden-post",
			Content:     "Content 2",
			PublishedAt: now,
			Visible:     false,
		},
		{
			Title:       "Visible Post 2",
			Slug:        "visible-post-2",
			Content:     "Content 3",
			PublishedAt: now.Add(-1 * time.Hour), // 1 hour ago
			Visible:     true,
		},
		{
			Title:       "Scheduled Post",
			Slug:        "scheduled-post",
			Content:     "Content 4",
			PublishedAt: now.Add(24 * time.Hour), // Tomorrow
			Visible:     true,
		},
	}

	for _, p := range posts {
		err := repo.Create(p)
		require.NoError(t, err)
	}

	found, total, err := repo.FindVisiblePaginated(1, 10, timezone)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total) // Only visible and published posts
	assert.Len(t, found, 2)

	// Verify that only visible and published posts are returned
	titles := []string{found[0].Title, found[1].Title}
	assert.Contains(t, titles, "Visible Post 1")
	assert.Contains(t, titles, "Visible Post 2")
}

package repository

import (
	"strings"
	"time"

	"codeinstyle.io/captain/models"
	"gorm.io/gorm"

	"codeinstyle.io/captain/utils"
)

// PostRepository handles database operations for posts
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new post repository
func NewPostRepository(db *gorm.DB) models.PostRepository {
	return &PostRepository{
		db: db,
	}
}

// Create creates a new post
func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// Update updates an existing post
func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// Delete deletes a post
func (r *PostRepository) Delete(post *models.Post) error {
	return r.db.Delete(post).Error
}

// FindByID finds a post by ID
func (r *PostRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Tags").Joins("Author").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// FindBySlug finds a post by slug
func (r *PostRepository) FindBySlug(slug string) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Tags").Joins("Author").Where("slug = ?", slug).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) CountByAuthor(user *models.User) (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Where("author_id = ?", user.ID).Count(&count).Error
	return count, err
}

func (r *PostRepository) CountByTag(tagID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ?", tagID).
		Count(&count).Error
	return count, err
}

// FindByTag finds all posts with a specific tag slug
func (r *PostRepository) FindByTag(tag string) ([]*models.Post, error) {
	var posts []*models.Post
	err := r.db.Preload("Tags").Joins("Author").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.slug = ?", tag).
		Find(&posts).Error
	return posts, err
}

// FindVisiblePaginated finds all visible posts with pagination
func (r *PostRepository) FindVisiblePaginated(page, perPage int, timezone string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * perPage

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)

	query := r.db.Model(&models.Post{}).
		Where("visible = ? AND published_at <= ?", true, now)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err = query.Preload("Tags").Joins("Author").
		Order("published_at desc").
		Offset(offset).
		Limit(perPage).
		Find(&posts).Error

	return posts, total, err
}

// FindVisibleByTag finds all visible posts with a specific tag
func (r *PostRepository) FindVisibleByTag(tagID uint, page, perPage int, timezone string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * perPage

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)

	query := r.db.Model(&models.Post{}).
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ? AND visible = ? AND published_at <= ?", tagID, true, now)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err = query.Preload("Tags").Joins("Author").
		Order("posts.published_at desc").
		Offset(offset).
		Limit(perPage).
		Find(&posts).Error

	return posts, total, err
}

// FindAll finds all posts
func (r *PostRepository) FindAll() ([]*models.Post, error) {
	var posts []*models.Post
	err := r.db.Preload("Tags").Joins("Author").
		Order("posts.created_at desc").
		Find(&posts).Error
	return posts, err
}

// FindRecent finds the most recent posts
func (r *PostRepository) FindRecent(limit int) ([]*models.Post, error) {
	var posts []*models.Post
	err := r.db.Preload("Tags").Joins("Author").
		Order("posts.created_at desc").
		Limit(limit).
		Find(&posts).Error
	return posts, err
}

// FindAllPaginated finds all posts with pagination
func (r *PostRepository) FindAllPaginated(page, perPage int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * perPage

	query := r.db.Model(&models.Post{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Tags").Joins("Author").
		Order("published_at desc").
		Offset(offset).
		Limit(perPage).
		Find(&posts).Error

	return posts, total, err
}

// FindAllByTag finds all posts with a specific tag (including non-visible)
func (r *PostRepository) FindAllByTag(tagID uint, page, perPage int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * perPage

	query := r.db.Model(&models.Post{}).
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ?", tagID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Tags").Joins("Author").
		Order("published_at desc").
		Offset(offset).
		Limit(perPage).
		Find(&posts).Error

	return posts, total, err
}

// AssociateTags associates tags with a post
func (r *PostRepository) AssociateTags(post *models.Post, tags []string) error {

	var tagsToSave []*models.Tag
	for _, tag := range tags {
		if strings.TrimSpace(tag) == "" {
			continue
		}

		var existingTag models.Tag
		slug := utils.Slugify(tag)
		err := r.db.Where("slug = ?", slug).First(&existingTag).Error
		if err != nil {
			existingTag = models.Tag{
				Name: tag,
				Slug: slug,
			}
			if err := r.db.Create(&existingTag).Error; err != nil {
				return err
			}
		}
		tagsToSave = append(tagsToSave, &existingTag)
	}

	assoc := r.db.Model(post).Association("Tags")

	if len(tagsToSave) == 0 {
		if err := assoc.Clear(); err != nil {
			return err
		}
		return nil
	}

	return assoc.Replace(tagsToSave)
}

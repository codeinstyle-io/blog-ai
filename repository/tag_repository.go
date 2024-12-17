package repository

import (
	"codeinstyle.io/captain/models"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *gorm.DB) models.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(tag *models.Tag) error {
	return r.db.Delete(&models.Tag{}, tag).Error
}

func (r *tagRepository) FindPostsAndCount() ([]struct {
	models.Tag
	PostCount int64
}, error) {
	var tags []struct {
		models.Tag
		PostCount int64
	}

	result := r.db.Model(&models.Tag{}).
		Select("tags.*, count(post_tags.post_id) as post_count").
		Joins("left join post_tags on post_tags.tag_id = tags.id").
		Group("tags.id").
		Find(&tags)

	if result.Error != nil {
		return nil, result.Error
	}

	return tags, nil
}

func (r *tagRepository) FindBySlug(slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindAll() ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.Find(&tags).Error
	return tags, err
}

func (r *tagRepository) FindByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

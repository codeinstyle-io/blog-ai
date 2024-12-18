package models

// PostRepository defines the interface for post operations
type PostRepository interface {
	Create(post *Post) error
	Update(post *Post) error
	Delete(post *Post) error
	FindByID(id uint) (*Post, error)
	FindBySlug(slug string) (*Post, error)
	FindByTag(tag string) ([]*Post, error)
	FindAllPaginated(page, perPage int) ([]Post, int64, error)
	FindVisiblePaginated(page, perPage int, timezone string) ([]Post, int64, error)
	FindAllByTag(tagID uint, page, perPage int) ([]Post, int64, error)
	FindVisibleByTag(tagID uint, page, perPage int, timezone string) ([]Post, int64, error)
	FindAll() ([]*Post, error)
	FindRecent(limit int) ([]*Post, error)
	AssociateTags(post *Post, tags []string) error
	CountByAuthor(user *User) (int64, error)
	CountByTag(tagID uint) (int64, error)
}

// TagRepository defines the interface for tag operations
type TagRepository interface {
	Create(tag *Tag) error
	Update(tag *Tag) error
	Delete(tag *Tag) error
	FindByID(id uint) (*Tag, error)
	FindBySlug(slug string) (*Tag, error)
	FindByName(name string) (*Tag, error)
	FindAll() ([]*Tag, error)
	FindPostsAndCount() ([]struct {
		Tag
		PostCount int64
	}, error)
}

// UserRepository defines the interface for user operations
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	CountByEmail(email string) (int64, error)
	CountAll() (int64, error)
}

// PageRepository defines the interface for page operations
type PageRepository interface {
	Create(page *Page) error
	Update(page *Page) error
	Delete(page *Page) error
	FindByID(id uint) (*Page, error)
	FindBySlug(slug string) (*Page, error)
	FindAll() ([]*Page, error)
	CountRelatedMenuItems(id uint, count *int64) error
}

// MenuItemRepository defines the interface for menu item operations
type MenuItemRepository interface {
	Create(item *MenuItem) error
	CreateAll(items []MenuItem) error
	Update(item *MenuItem) error
	Delete(item *MenuItem) error
	DeleteAll() error
	FindByID(id uint) (*MenuItem, error)
	FindAll() ([]*MenuItem, error)
	UpdatePositions(startPosition int) error
	GetNextPosition() int
	Move(id uint, direction string) error
}

// SettingsRepository defines the interface for settings operations
type SettingsRepository interface {
	Get() (*Settings, error)
	Update(settings *Settings) error
	Create(settings Settings) error
}

// MediaRepository defines the interface for media operations
type MediaRepository interface {
	Create(media *Media) error
	Update(media *Media) error
	Delete(media *Media) error
	FindByPath(path string) (*Media, error)
	FindByID(id uint) (*Media, error)
	FindAll() ([]*Media, error)
}

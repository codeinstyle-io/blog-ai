package models

// PostRepository defines the interface for post operations
type PostRepository interface {
	Create(post *Post) error
	Update(post *Post) error
	Delete(post *Post) error
	FindByID(id uint) (*Post, error)
	FindBySlug(slug string) (*Post, error)
	FindByTag(tag string) ([]*Post, error)
	FindVisible(page, perPage int) ([]Post, int64, error)
	FindVisibleByTag(tagID uint, page, perPage int) ([]Post, int64, error)
	FindAll() ([]*Post, error)
	FindRecent(limit int) ([]*Post, error)
}

// TagRepository defines the interface for tag operations
type TagRepository interface {
	Create(tag *Tag) error
	Update(tag *Tag) error
	Delete(tag *Tag) error
	FindByID(id uint) (*Tag, error)
	FindBySlug(slug string) (*Tag, error)
	FindAll() ([]*Tag, error)
}

// UserRepository defines the interface for user operations
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	Delete(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	FindBySessionToken(token string) (*User, error)
	FindAll() ([]*User, error)
}

// PageRepository defines the interface for page operations
type PageRepository interface {
	Create(page *Page) error
	Update(page *Page) error
	Delete(page *Page) error
	FindByID(id uint) (*Page, error)
	FindBySlug(slug string) (*Page, error)
	FindAll() ([]*Page, error)
}

// MenuItemRepository defines the interface for menu item operations
type MenuItemRepository interface {
	Create(item *MenuItem) error
	Update(item *MenuItem) error
	Delete(item *MenuItem) error
	DeleteAll() error
	FindByID(id uint) (*MenuItem, error)
	FindAll() ([]*MenuItem, error)
	UpdatePositions(startPosition int) error
	GetNextPosition() int
}

// SettingsRepository defines the interface for settings operations
type SettingsRepository interface {
	Get() (*Settings, error)
	Update(settings *Settings) error
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

// SessionRepository defines the interface for session operations
type SessionRepository interface {
	Create(session *Session) error
	FindByToken(token string) (*Session, error)
	DeleteByToken(token string) error
	DeleteExpired() error
}

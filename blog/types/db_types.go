package types

type Post struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `gorm:"not null"`
	Slug    string `gorm:"uniqueIndex;not null"`
	Content string `gorm:"not null"`
	Tags    []Tag  `gorm:"many2many:post_tags;"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}

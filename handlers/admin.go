package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandlers struct {
	db     *gorm.DB
	config *config.Config
}

func NewAdminHandlers(database *gorm.DB, config *config.Config) *AdminHandlers {
	return &AdminHandlers{
		db:     database,
		config: config,
	}
}

// ListTags shows all tags and their post counts
func (h *AdminHandlers) ListTags(c *gin.Context) {
	var tags []struct {
		db.Tag
		PostCount int64
	}

	result := h.db.Model(&db.Tag{}).
		Select("tags.*, count(post_tags.post_id) as post_count").
		Joins("left join post_tags on post_tags.tag_id = tags.id").
		Group("tags.id").
		Find(&tags)

	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_tags.tmpl", gin.H{
		"title": "Tags",
		"tags":  tags,
	})
}

// DeleteTag removes a tag without affecting posts
func (h *AdminHandlers) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&db.Tag{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}
	c.Redirect(http.StatusFound, "/admin/tags")
}

// ListUsers shows all users (except sensitive data)
func (h *AdminHandlers) ListUsers(c *gin.Context) {
	var users []db.User
	if err := h.db.Select("id, first_name, last_name, email, created_at, updated_at").Find(&users).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_users.tmpl", gin.H{
		"title": "Users",
		"users": users,
	})
}

// ShowCreatePost displays the post creation form
func (h *AdminHandlers) ShowCreatePost(c *gin.Context) {
	var tags []db.Tag
	if err := h.db.Find(&tags).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	c.HTML(http.StatusOK, "admin_create_post.tmpl", gin.H{
		"title": "Create Post",
		"tags":  tags,
	})
}

func (h *AdminHandlers) CreatePost(c *gin.Context) {
	// Get the logged in user
	userInterface, exists := c.Get("user")
	if !exists {
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", gin.H{
			"error": "User session not found",
		})
		return
	}
	user := userInterface.(*db.User)

	var post db.Post

	// Parse form data
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	content := c.PostForm("content")
	publishedAt := c.PostForm("publishedAt")
	var parsedTime time.Time

	if publishedAt == "" {
		// If no date provided, use current time
		parsedTime = time.Now().In(h.config.GetLocation())
	} else {
		var err error
		parsedTime, err = time.ParseInLocation("2006-01-02T15:04", publishedAt, h.config.GetLocation())
		if err != nil {
			c.HTML(http.StatusBadRequest, "admin_create_post.tmpl", gin.H{
				"error": "Invalid date format",
			})
			return
		}
	}
	visible := c.PostForm("visible") == "on"

	// Basic validation
	if title == "" || slug == "" || content == "" {
		c.HTML(http.StatusBadRequest, "admin_create_post.tmpl", gin.H{
			"error": "All fields are required",
		})
		return
	}

	// Create post object
	post = db.Post{
		Title:       title,
		Slug:        slug,
		Content:     content,
		PublishedAt: parsedTime.UTC(),
		Visible:     visible,
		AuthorID:    user.ID, // Set the author ID
	}

	// Handle tags
	var tagNames []string
	tagsJSON := c.PostForm("tags")
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &tagNames); err != nil {
			c.HTML(http.StatusBadRequest, "admin_create_post.tmpl", gin.H{
				"error": "Invalid tags format",
				"post":  post,
			})
			return
		}
	}

	// Create/get tags and associate
	var tags []db.Tag
	for _, name := range tagNames {
		var tag db.Tag
		result := h.db.Where(db.Tag{Name: name}).FirstOrCreate(&tag)
		if result.Error != nil {
			c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", gin.H{
				"error": "Failed to create tag",
				"post":  post,
			})
			return
		}
		tags = append(tags, tag)
	}
	post.Tags = tags

	// Create post with transaction to ensure atomic operation
	tx := h.db.Begin()
	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", gin.H{
			"error": "Failed to create post",
			"post":  post,
		})
		return
	}

	if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin_create_post.tmpl", gin.H{
			"error": "Failed to associate tags",
			"post":  post,
		})
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, "/admin/posts")
}

// ListPosts shows all posts for admin
func (h *AdminHandlers) ListPosts(c *gin.Context) {
	var posts []db.Post
	if err := h.db.Preload("Tags").Preload("Author").Find(&posts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	// Convert UTC times to configured timezone for display
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(h.config.GetLocation())
	}

	c.HTML(http.StatusOK, "admin_posts.tmpl", gin.H{
		"title": "Posts",
		"posts": posts,
	})
}

// ListPostsByTag shows all posts for a specific tag
func (h *AdminHandlers) ListPostsByTag(c *gin.Context) {
	tagID := c.Param("id")

	var tag db.Tag
	if err := h.db.First(&tag, tagID).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	var posts []db.Post
	if err := h.db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagID).
		Preload("Tags").
		Preload("Author").
		Find(&posts).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	// Convert UTC times to configured timezone for display
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(h.config.GetLocation())
	}

	data := gin.H{
		"title": "Posts tagged with " + tag.Name,
		"tag":   tag,
		"posts": posts,
	}
	data = h.addCommonData(c, data)

	c.HTML(http.StatusOK, "admin_tag_posts.tmpl", data)
}

// ConfirmDeletePost shows deletion confirmation page
func (h *AdminHandlers) ConfirmDeletePost(c *gin.Context) {
	id := c.Param("id")
	var post db.Post
	if err := h.db.First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_confirm_delete.tmpl", gin.H{
		"title": "Confirm Delete",
		"post":  post,
	})
}

// DeletePost removes a post and its tag associations
func (h *AdminHandlers) DeletePost(c *gin.Context) {
	id := c.Param("id")

	// Start transaction
	tx := h.db.Begin()

	// Delete post_tags associations
	if err := tx.Exec("DELETE FROM post_tags WHERE post_id = ?", id).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
			"error": "Failed to delete post tags",
		})
		return
	}

	// Delete post
	if err := tx.Delete(&db.Post{}, id).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{
			"error": "Failed to delete post",
		})
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, "/admin/posts")
}

func (h *AdminHandlers) EditPost(c *gin.Context) {
	id := c.Param("id")
	var post db.Post

	if err := h.db.Preload("Tags").First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	var allTags []db.Tag
	if err := h.db.Find(&allTags).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	// Convert UTC time to configured timezone for display
	post.PublishedAt = post.PublishedAt.In(h.config.GetLocation())

	c.HTML(http.StatusOK, "admin_edit_post.tmpl", gin.H{
		"title":   "Edit Post",
		"post":    post,
		"allTags": allTags,
	})
}

func (h *AdminHandlers) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post db.Post

	if err := h.db.First(&post, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{})
		return
	}

	// Update fields
	post.Title = c.PostForm("title")
	post.Slug = c.PostForm("slug")
	post.Content = c.PostForm("content")
	post.Visible = c.PostForm("visible") == "on"

	// Parse the published date in configured timezone
	publishedAt, err := time.ParseInLocation("2006-01-02T15:04", c.PostForm("publishedAt"), h.config.GetLocation())
	post.PublishedAt = publishedAt
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin_edit_post.tmpl", gin.H{
			"error": "Invalid date format",
			"post":  post,
		})
		return
	}

	// Handle tags
	var tagNames []string
	tagsJSON := c.PostForm("tags")
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &tagNames); err != nil {
			c.HTML(http.StatusBadRequest, "admin_edit_post.tmpl", gin.H{
				"error": "Invalid tags format",
				"post":  post,
			})
			return
		}
	}

	// Update tags
	var tags []db.Tag
	for _, name := range tagNames {
		var tag db.Tag
		h.db.FirstOrCreate(&tag, db.Tag{Name: name})
		tags = append(tags, tag)
	}
	post.Tags = tags
	post.PublishedAt = publishedAt.UTC()

	// Update with transaction
	tx := h.db.Begin()
	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin_edit_post.tmpl", gin.H{
			"error": "Failed to update post",
			"post":  post,
		})
		return
	}

	if err := tx.Model(&post).Association("Tags").Replace(tags); err != nil {
		tx.Rollback()
		c.HTML(http.StatusInternalServerError, "admin_edit_post.tmpl", gin.H{
			"error": "Failed to update tags",
			"post":  post,
		})
		return
	}

	tx.Commit()
	c.Redirect(http.StatusFound, "/admin/posts")
}

func (h *AdminHandlers) Index(c *gin.Context) {
	var postCount, tagCount, userCount int64
	var recentPosts []db.Post

	// Get counts
	h.db.Model(&db.Post{}).Count(&postCount)
	h.db.Model(&db.Tag{}).Count(&tagCount)
	h.db.Model(&db.User{}).Count(&userCount)

	// Get 5 most recent posts
	h.db.Order("published_at desc").Limit(5).Find(&recentPosts)

	data := gin.H{
		"title":       "Dashboard",
		"postCount":   postCount,
		"tagCount":    tagCount,
		"userCount":   userCount,
		"recentPosts": recentPosts,
	}

	data = h.addCommonData(c, data)

	c.HTML(http.StatusOK, "admin_index.tmpl", data)
}

// Add response struct
type tagResponse struct {
	Id   uint   `json:"id"`   // lowercase for JS
	Name string `json:"name"` // lowercase for JS
}

func (h *AdminHandlers) GetTags(c *gin.Context) {
	var dbTags []db.Tag

	if err := h.db.Find(&dbTags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	// Transform to response format
	tags := make([]tagResponse, len(dbTags))
	for i, tag := range dbTags {
		tags[i] = tagResponse{
			Id:   tag.ID,
			Name: tag.Name,
		}
	}

	c.JSON(http.StatusOK, tags)
}

// handlers/admin.go
func (h *AdminHandlers) SavePreferences(c *gin.Context) {
	var prefs struct {
		Theme string `json:"theme"`
	}

	if err := c.BindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preferences"})
		return
	}

	// Save theme preference in cookie
	c.SetCookie("admin_theme", prefs.Theme, 3600*24*365, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Add to all admin handlers:
func (h *AdminHandlers) addCommonData(c *gin.Context, data gin.H) gin.H {
	if data == nil {
		data = gin.H{}
	}

	theme, _ := c.Cookie("admin_theme")
	if theme == "" {
		theme = "light"
	}

	data["theme"] = theme
	return data
}

// ShowCreateTag displays the tag creation form
func (h *AdminHandlers) ShowCreateTag(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create_tag.tmpl", gin.H{
		"title": "Create Tag",
	})
}

// CreateTag handles tag creation
func (h *AdminHandlers) CreateTag(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		c.HTML(http.StatusBadRequest, "admin_create_tag.tmpl", gin.H{
			"error": "Tag name is required",
		})
		return
	}

	tag := db.Tag{Name: name}
	if err := h.db.Create(&tag).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_tag.tmpl", gin.H{
			"error": "Failed to create tag",
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/tags")
}

// Page handlers
func (h *AdminHandlers) ListPages(c *gin.Context) {
	var pages []db.Page
	if err := h.db.Find(&pages).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_pages.tmpl", h.addCommonData(c, gin.H{
		"title": "Pages",
		"pages": pages,
	}))
}

func (h *AdminHandlers) CreatePage(c *gin.Context) {
	page := db.Page{
		Title:       c.PostForm("title"),
		Slug:        c.PostForm("slug"),
		Content:     c.PostForm("content"),
		ContentType: c.PostForm("content_type"),
		Visible:     c.PostForm("visible") == "on",
	}

	if err := h.db.Create(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_page.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to create page",
			"page":  page,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/pages")
}

func (h *AdminHandlers) ShowCreatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create_page.tmpl", h.addCommonData(c, gin.H{
		"title": "Create Page",
	}))
}

func (h *AdminHandlers) EditPage(c *gin.Context) {
	id := c.Param("id")
	var page db.Page
	if err := h.db.First(&page, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", nil)
		return
	}

	c.HTML(http.StatusOK, "admin_edit_page.tmpl", h.addCommonData(c, gin.H{
		"title": "Edit Page",
		"page":  page,
	}))
}

func (h *AdminHandlers) UpdatePage(c *gin.Context) {
	id := c.Param("id")
	var page db.Page
	if err := h.db.First(&page, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", nil)
		return
	}

	page.Title = c.PostForm("title")
	page.Slug = c.PostForm("slug")
	page.Content = c.PostForm("content")
	page.ContentType = c.PostForm("content_type")
	page.Visible = c.PostForm("visible") == "on"

	if err := h.db.Save(&page).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_edit_page.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update page",
			"page":  page,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/pages")
}

func (h *AdminHandlers) DeletePage(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&db.Page{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Menu handlers
func (h *AdminHandlers) ListMenuItems(c *gin.Context) {
	var items []db.MenuItem
	if err := h.db.Preload("Page").Order("position").Find(&items).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "admin_menus.tmpl", h.addCommonData(c, gin.H{
		"title":     "Menu Items",
		"menuItems": items,
	}))
}

func (h *AdminHandlers) CreateMenuItem(c *gin.Context) {
	item := db.MenuItem{
		Label:    c.PostForm("label"),
		Position: h.getNextMenuPosition(),
	}

	// Handle either URL or Page reference
	if pageID := c.PostForm("page_id"); pageID != "" {
		pid := parseUint(pageID)
		item.PageID = &pid
	} else if url := c.PostForm("url"); url != "" {
		item.URL = &url
	}

	if err := h.db.Create(&item).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin_create_menu_item.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to create menu item",
			"item":  item,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/menus")
}

func (h *AdminHandlers) ShowCreateMenuItem(c *gin.Context) {
	var pages []db.Page
	if err := h.db.Find(&pages).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", nil)
		return
	}

	c.HTML(http.StatusOK, "admin_create_menu_item.tmpl", h.addCommonData(c, gin.H{
		"title": "Create Menu Item",
		"pages": pages,
	}))
}

func (h *AdminHandlers) MoveMenuItem(c *gin.Context) {
	id := c.Param("id")
	direction := c.Param("direction")

	var item db.MenuItem
	if err := h.db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
		return
	}

	var targetItem db.MenuItem
	query := h.db
	if direction == "up" {
		query = query.Where("position < ?", item.Position).Order("position DESC")
	} else {
		query = query.Where("position > ?", item.Position).Order("position ASC")
	}

	if err := query.First(&targetItem).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot move item further"})
		return
	}

	// Swap positions
	tempPosition := item.Position
	item.Position = targetItem.Position
	targetItem.Position = tempPosition

	tx := h.db.Begin()
	if err := tx.Save(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item position"})
		return
	}

	if err := tx.Save(&targetItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target position"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *AdminHandlers) ConfirmDeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	var item db.MenuItem
	if err := h.db.First(&item, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", nil)
		return
	}

	c.HTML(http.StatusOK, "admin_delete_menu_item.tmpl", h.addCommonData(c, gin.H{
		"title": "Delete Menu Item",
		"item":  item,
	}))
}

func (h *AdminHandlers) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	var item db.MenuItem
	if err := h.db.First(&item, id).Error; err != nil {
		if c.Request.Header.Get("Accept") == "application/json" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
		} else {
			c.HTML(http.StatusNotFound, "404.tmpl", nil)
		}
		return
	}

	if err := h.db.Delete(&item).Error; err != nil {
		if c.Request.Header.Get("Accept") == "application/json" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu item"})
		} else {
			c.HTML(http.StatusInternalServerError, "500.tmpl", nil)
		}
		return
	}

	if c.Request.Header.Get("Accept") == "application/json" {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.Redirect(http.StatusFound, "/admin/menus")
	}
}

func (h *AdminHandlers) EditMenuItem(c *gin.Context) {
	id := c.Param("id")
	var item db.MenuItem
	if err := h.db.Preload("Page").First(&item, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", nil)
		return
	}

	var pages []db.Page
	if err := h.db.Find(&pages).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", nil)
		return
	}

	c.HTML(http.StatusOK, "admin_edit_menu_item.tmpl", h.addCommonData(c, gin.H{
		"title": "Edit Menu Item",
		"item":  item,
		"pages": pages,
	}))
}

func (h *AdminHandlers) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")
	var item db.MenuItem
	if err := h.db.First(&item, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", nil)
		return
	}

	item.Label = c.PostForm("label")

	// Reset both URL and PageID
	item.URL = nil
	item.PageID = nil

	// Handle either URL or Page reference
	if pageID := c.PostForm("page_id"); pageID != "" {
		pid := parseUint(pageID)
		item.PageID = &pid
	} else if url := c.PostForm("url"); url != "" {
		item.URL = &url
	}

	if err := h.db.Save(&item).Error; err != nil {
		var pages []db.Page
		h.db.Find(&pages)
		c.HTML(http.StatusInternalServerError, "admin_edit_menu_item.tmpl", h.addCommonData(c, gin.H{
			"error": "Failed to update menu item",
			"item":  item,
			"pages": pages,
		}))
		return
	}

	c.Redirect(http.StatusFound, "/admin/menus")
}

func parseUint(pageID string) uint {
	pid, err := strconv.ParseUint(pageID, 10, 32)
	if err != nil {
		return 0
	}
	return uint(pid)

}

func (h *AdminHandlers) getNextMenuPosition() int {
	var maxPos struct{ Max int }
	h.db.Model(&db.MenuItem{}).Select("COALESCE(MAX(position), -1) as max").Scan(&maxPos)
	return maxPos.Max + 1
}

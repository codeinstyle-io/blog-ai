package handlers

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/db"
	"codeinstyle.io/captain/system"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	chromaCSS string
)

type PublicHandlers struct {
	db       *gorm.DB
	config   *config.Config
	settings *db.Settings
}

func NewPublicHandlers(database *gorm.DB, cfg *config.Config) *PublicHandlers {
	var settings db.Settings
	database.First(&settings)

	return &PublicHandlers{
		db:       database,
		config:   cfg,
		settings: &settings,
	}
}

func (h *PublicHandlers) GetChromaCSS(c *gin.Context) {
	// Generate ETag based on the chroma style name
	etag := fmt.Sprintf("\"%x\"", md5.Sum([]byte(h.settings.ChromaStyle)))

	// Check If-None-Match header first
	if match := c.GetHeader("If-None-Match"); match != "" {
		if match == etag {
			c.Status(http.StatusNotModified)
			return
		}
	}

	// Generate CSS
	style := styles.Get(h.settings.ChromaStyle)
	if style == nil {
		style = styles.Fallback
	}
	formatter := html.New(html.WithClasses(true))
	buf := new(bytes.Buffer)
	if err := formatter.WriteCSS(buf, style); err != nil {
		// If there's an error, use empty CSS
		chromaCSS = ""
		return
	}
	chromaCSS = buf.String()

	// Set headers
	c.Header("ETag", etag)
	c.Header("Content-Type", "text/css")
	c.String(http.StatusOK, chromaCSS)
}

func (h *PublicHandlers) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	now := time.Now()

	var post db.Post
	if err := h.db.Preload("Tags").Preload("Author").
		Where("slug = ? AND visible = ? AND published_at <= ?", slug, true, now).
		First(&post).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(gin.H{
				"title": "Post not found",
			}))
			return
		}
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
			"title": "Error",
		}))
		return
	}

	// Render markdown content
	post.Content = renderMarkdown(post.Content)

	// Convert UTC time to configured timezone for display
	loc, err := time.LoadLocation(h.settings.Timezone)
	if err != nil {
		loc = time.UTC
	}
	post.PublishedAt = post.PublishedAt.In(loc)

	c.HTML(http.StatusOK, "post.tmpl", h.addCommonData(gin.H{
		"title": post.Title,
		"post":  post,
	}))
}

func (h *PublicHandlers) ListPosts(c *gin.Context) {
	// Get page number from query params
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get posts per page from settings
	perPage := h.settings.PostsPerPage

	now := time.Now()

	var total int64
	h.db.Model(&db.Post{}).Where("visible = ? AND published_at <= ?", true, now).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	offset := (page - 1) * perPage

	var posts []db.Post
	result := h.db.Preload("Tags").Preload("Author").
		Where("visible = ? AND published_at <= ?", true, now).
		Order("published_at desc").
		Offset(offset).
		Limit(perPage).
		Find(&posts)

	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	processPostsContent(posts)
	processPostsPublishedAt(posts)

	c.HTML(http.StatusOK, "posts.tmpl", h.addCommonData(gin.H{
		"title":       "Latest Articles",
		"posts":       posts,
		"currentPage": page,
		"totalPages":  totalPages,
	}))
}

func (h *PublicHandlers) ListPostsByTag(c *gin.Context) {
	tagName := c.Param("tag")
	// Get page number from query params
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get posts per page from settings
	perPage := h.settings.PostsPerPage

	now := time.Now()

	var total int64
	h.db.Model(&db.Post{}).
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.name = ? AND posts.visible = ? AND posts.published_at <= ?", tagName, true, now).
		Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	offset := (page - 1) * perPage

	var posts []db.Post
	result := h.db.Preload("Tags").Preload("Author").
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("tags.name = ? AND posts.visible = ? AND posts.published_at <= ?", tagName, true, now).
		Order("published_at desc").
		Offset(offset).
		Limit(perPage).
		Find(&posts)

	if result.Error != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", gin.H{})
		return
	}

	processPostsContent(posts)
	processPostsPublishedAt(posts)

	c.HTML(http.StatusOK, "tag_posts.tmpl", h.addCommonData(gin.H{
		"title":       fmt.Sprintf("Posts tagged with #%s", tagName),
		"tag":         tagName,
		"posts":       posts,
		"currentPage": page,
		"totalPages":  totalPages,
	}))
}

func (h *PublicHandlers) GetPageBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var page db.Page
	if err := h.db.Where("slug = ? AND visible = true", slug).First(&page).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", nil)
		return
	}

	// Render content based on type
	if page.ContentType == "markdown" {
		page.Content = renderMarkdown(page.Content)
	}

	c.HTML(http.StatusOK, "page.tmpl", h.addCommonData(gin.H{
		"title": page.Title,
		"page":  page,
	}))
}

func (h *PublicHandlers) addCommonData(data gin.H) gin.H {
	// Get menu items
	var menuItems []db.MenuItem
	h.db.Preload("Page").Order("position").Find(&menuItems)

	var settings db.Settings
	h.db.First(&settings)
	h.settings = &settings

	// Add menu items to the data
	data["menuItems"] = menuItems

	// Add site config from settings
	data["config"] = gin.H{
		"SiteTitle":    h.settings.Title,
		"SiteSubtitle": h.settings.Subtitle,
		"Theme":        h.settings.Theme,
	}

	// Add version information
	data["version"] = system.Version

	return data
}

func processPostsPublishedAt(posts []db.Post) {
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(time.UTC)
	}
}

func processPostsContent(posts []db.Post) {
	for i := range posts {
		if posts[i].Excerpt != nil && *posts[i].Excerpt != "" {
			continue
		}
		// Truncate content to ~200 chars if no excerpt
		content := posts[i].Content
		if len(content) > 200 {
			content = content[:200]
			if idx := strings.LastIndex(content, " "); idx > 0 {
				content = content[:idx]
			}
			content += "..."
		}
		posts[i].Excerpt = &content
		// Render markdown for excerpt
		if posts[i].Excerpt != nil {
			rendered := renderMarkdown(*posts[i].Excerpt)
			posts[i].Excerpt = &rendered
		}
	}
}

// New functions for markdown processing
func renderMarkdown(content string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	opts := mdhtml.RendererOptions{
		Flags:          mdhtml.CommonFlags,
		RenderNodeHook: renderHook,
	}
	htmlRenderer := mdhtml.NewRenderer(opts)

	md := []byte(content)
	parsed := p.Parse(md)
	return string(markdown.Render(parsed, htmlRenderer))
}

func renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok && entering {
		language := string(code.Info)
		content := string(code.Literal)

		// Default to plain text if language not specified
		if language == "" {
			language = "text"
		}

		lexer := lexers.Get(language)
		if lexer == nil {
			lexer = lexers.Get("text")
		}

		// Use a light theme by default
		style := styles.Get("github")
		if style == nil {
			style = styles.Fallback
		}

		formatter := html.New(html.WithClasses(true))
		iterator, err := lexer.Tokenise(nil, content)
		if err != nil {
			// Fallback to plain text on error
			if _, err := io.WriteString(w, "<pre><code>"+content+"</code></pre>"); err != nil {
				return ast.GoToNext, true
			}
			return ast.GoToNext, true
		}

		buf := new(bytes.Buffer)
		err = formatter.Format(buf, style, iterator)
		if err != nil {
			if _, err := io.WriteString(w, "<pre><code>"+content+"</code></pre>"); err != nil {
				return ast.GoToNext, true
			}
			return ast.GoToNext, true
		}

		if _, err := io.WriteString(w, buf.String()); err != nil {
			return ast.GoToNext, true
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

// ...other post handlers like GetPostBySlug, EditPost, etc...

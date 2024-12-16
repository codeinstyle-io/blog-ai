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

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/models"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/system"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// PublicHandlers handles all public routes
type PublicHandlers struct {
	*BaseHandler
}

// NewPublicHandlers creates a new public handlers instance
func NewPublicHandlers(repos *repository.Repositories, cfg *config.Config) *PublicHandlers {
	return &PublicHandlers{
		BaseHandler: NewBaseHandler(repos, cfg),
	}
}

func (h *PublicHandlers) GetChromaCSS(c *gin.Context) {
	// Generate ETag based on the chroma style name
	var chromaCSS string
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
	post, err := h.postRepo.FindBySlug(slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(gin.H{
			"title": "Post not found",
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

	posts, total, err := h.postRepo.FindVisible(page, h.settings.PostsPerPage)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{}))
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(h.settings.PostsPerPage)))

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
	tagSlug := c.Param("slug")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	tag, err := h.tagRepo.FindBySlug(tagSlug)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
			"title": "Error",
		}))
		return
	}

	posts, total, err := h.postRepo.FindVisibleByTag(tag.ID, page, h.settings.PostsPerPage)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.tmpl", h.addCommonData(gin.H{
			"title": "Error",
		}))
		return
	}

	// Process posts
	processPostsPublishedAt(posts)
	processPostsContent(posts)

	totalPages := int(math.Ceil(float64(total) / float64(h.settings.PostsPerPage)))

	c.HTML(http.StatusOK, "tag_posts.tmpl", h.addCommonData(gin.H{
		"title":      fmt.Sprintf("Posts tagged with %s", tag.Name),
		"posts":      posts,
		"tag":        tag,
		"page":       page,
		"totalPages": totalPages,
	}))
}

func (h *PublicHandlers) GetPageBySlug(c *gin.Context) {
	slug := c.Param("slug")
	page, err := h.pageRepo.FindBySlug(slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", h.addCommonData(gin.H{
			"title": "Page not found",
		}))
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
	menuItems, _ := h.menuRepo.FindAll()

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

func processPostsPublishedAt(posts []models.Post) {
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(time.UTC)
	}
}

func processPostsContent(posts []models.Post) {
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

// renderMarkdown converts markdown content to HTML
func renderMarkdown(content string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(content))

	htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank
	opts := mdhtml.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: renderHook,
	}
	renderer := mdhtml.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
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

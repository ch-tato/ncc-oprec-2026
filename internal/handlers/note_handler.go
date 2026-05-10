package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"notes-app/internal/database"
	"notes-app/internal/models"
	"notes-app/internal/services"

	"github.com/gin-gonic/gin"
)

// --- Validation constants ---

const (
	maxTitleLen   = 200
	maxContentLen = 50000
)

// --- Helpers ---

// parseID extracts and validates the :id path parameter.
func parseID(c *gin.Context) (uint, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Invalid ID",
			"Message": "The note ID provided is not valid.",
		})
		return 0, false
	}
	return uint(id), true
}

// findNote fetches a note by ID, rendering a 404 page on failure.
func findNote(c *gin.Context, id uint) (*models.Note, bool) {
	var note models.Note
	if err := database.DB.First(&note, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"Title":   "Note Not Found",
			"Message": "The note you are looking for does not exist.",
		})
		return nil, false
	}
	return &note, true
}

// --- Handlers ---

// IndexHandler lists all notes ordered by most recent.
func IndexHandler(c *gin.Context) {
	var notes []models.Note
	if err := database.DB.Order("updated_at DESC").Find(&notes).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Error",
			"Message": "Failed to load notes.",
		})
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "My Notes",
		"Notes": notes,
	})
}

// NewNoteHandler renders the create-note form.
func NewNoteHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", gin.H{
		"Title": "New Note",
	})
}

// CreateNoteHandler processes the create-note form submission.
func CreateNoteHandler(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	content := strings.TrimSpace(c.PostForm("content"))

	// Validate inputs.
	if title == "" || content == "" {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"Title":   "New Note",
			"Error":   "Title and content are required.",
			"NoteTitle":   title,
			"NoteContent": content,
		})
		return
	}
	if len(title) > maxTitleLen {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"Title":   "New Note",
			"Error":   "Title must be 200 characters or fewer.",
			"NoteTitle":   title,
			"NoteContent": content,
		})
		return
	}
	if len(content) > maxContentLen {
		c.HTML(http.StatusBadRequest, "create.html", gin.H{
			"Title":   "New Note",
			"Error":   "Content must be 50,000 characters or fewer.",
			"NoteTitle":   title,
			"NoteContent": content,
		})
		return
	}

	note := models.Note{Title: title, Content: content}
	if err := database.DB.Create(&note).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "create.html", gin.H{
			"Title":   "New Note",
			"Error":   "Failed to save note. Please try again.",
			"NoteTitle":   title,
			"NoteContent": content,
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// ViewNoteHandler displays a single note with rendered Markdown.
func ViewNoteHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	note, ok := findNote(c, id)
	if !ok {
		return
	}

	rendered, err := services.RenderMarkdown(note.Content)
	if err != nil {
		rendered = "<p>Error rendering markdown.</p>"
	}

	c.HTML(http.StatusOK, "view.html", gin.H{
		"Title":           note.Title,
		"Note":            note,
		"RenderedContent": template.HTML(rendered),
	})
}

// EditNoteHandler renders the edit-note form.
func EditNoteHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	note, ok := findNote(c, id)
	if !ok {
		return
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"Title": "Edit Note",
		"Note":  note,
	})
}

// UpdateNoteHandler processes the edit-note form submission.
func UpdateNoteHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	note, ok := findNote(c, id)
	if !ok {
		return
	}

	title := strings.TrimSpace(c.PostForm("title"))
	content := strings.TrimSpace(c.PostForm("content"))

	if title == "" || content == "" {
		c.HTML(http.StatusBadRequest, "edit.html", gin.H{
			"Title": "Edit Note",
			"Error": "Title and content are required.",
			"Note":  note,
		})
		return
	}
	if len(title) > maxTitleLen {
		c.HTML(http.StatusBadRequest, "edit.html", gin.H{
			"Title": "Edit Note",
			"Error": "Title must be 200 characters or fewer.",
			"Note":  note,
		})
		return
	}
	if len(content) > maxContentLen {
		c.HTML(http.StatusBadRequest, "edit.html", gin.H{
			"Title": "Edit Note",
			"Error": "Content must be 50,000 characters or fewer.",
			"Note":  note,
		})
		return
	}

	note.Title = title
	note.Content = content
	if err := database.DB.Save(note).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "edit.html", gin.H{
			"Title": "Edit Note",
			"Error": "Failed to update note. Please try again.",
			"Note":  note,
		})
		return
	}

	c.Redirect(http.StatusFound, "/notes/"+strconv.FormatUint(uint64(note.ID), 10))
}

// DeleteNoteHandler removes a note and redirects to the index.
func DeleteNoteHandler(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if _, ok := findNote(c, id); !ok {
		return
	}

	if err := database.DB.Delete(&models.Note{}, id).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Title":   "Error",
			"Message": "Failed to delete note.",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

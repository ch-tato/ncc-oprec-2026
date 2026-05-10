package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"notes-app/internal/database"
	"notes-app/internal/models"
	"notes-app/internal/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.LoadHTMLGlob("../../templates/*")
	routes.Setup(r)

	return r
}

func setupTestDB() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Gagal menghubungkan ke in-memory database: " + err.Error())
	}
	db.AutoMigrate(&models.Note{})
	database.DB = db
}

func TestIndexHandler(t *testing.T) {
	setupTestDB()
	r := setupTestRouter()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Diharapkan status 200 OK, mendapat %d", w.Code)
	}
}

func TestCreateNoteHandler_Success(t *testing.T) {
	setupTestDB()
	r := setupTestRouter()

	formData := url.Values{}
	formData.Set("title", "Catatan Test")
	formData.Set("content", "Ini adalah isi catatan test")

	req, _ := http.NewRequest("POST", "/notes", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Diharapkan status 302 Found, mendapat %d", w.Code)
	}
	var count int64
	database.DB.Model(&models.Note{}).Count(&count)
	if count != 1 {
		t.Errorf("Diharapkan ada 1 catatan di database, mendapat %d", count)
	}
}

func TestCreateNoteHandler_ValidationFail(t *testing.T) {
	setupTestDB()
	r := setupTestRouter()

	formData := url.Values{}
	formData.Set("title", "")
	formData.Set("content", "")

	req, _ := http.NewRequest("POST", "/notes", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Diharapkan status 400 Bad Request, mendapat %d", w.Code)
	}
}

func TestViewNoteHandler_NotFound(t *testing.T) {
	setupTestDB()
	r := setupTestRouter()
	req, _ := http.NewRequest("GET", "/notes/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Diharapkan status 404 Not Found, mendapat %d", w.Code)
	}
}

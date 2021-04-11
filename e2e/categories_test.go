package e2e

import (
	"food_ordering_backend/category"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var cat = &category.Category{}
var testCategories = []category.Category{{Model: gorm.Model{ID: 1}, Title: "Salads"}, {Model: gorm.Model{ID: 2}, Title: "Burgers"}, {Model: gorm.Model{ID: 3}, Title: "Pizza"}, {Model: gorm.Model{ID: 4}, Title: "Drinks"}}

func TestCategories(t *testing.T) {
	db := database.MustGetTest()
	r := router.Setup(db)

	t.Cleanup(cleanup(db))

	t.Run("GET /categories", func(t *testing.T) {
		t.Run("should return JSON array of categories", func(t *testing.T) {
			t.Cleanup(cleanup(db))
			it := assert.New(t)
			var categories []category.Category

			db.Find(&categories)
			require.Len(t, categories, 0)

			db.Create(&testCategories)
			db.Find(&categories)
			require.Len(t, categories, len(testCategories))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/categories", nil)
			r.ServeHTTP(w, req)

			it.Equal(http.StatusOK, w.Code)
			it.Contains(w.Header().Get("Content-Type"), "application/json")
			it.Equal(`[{"id":1,"title":"Salads","removable":false},{"id":2,"title":"Burgers","removable":false},{"id":3,"title":"Pizza","removable":false},{"id":4,"title":"Drinks","removable":false}]`, w.Body.String())
		})
	})

	t.Run("POST /categories", func(t *testing.T) {
		t.Run("should add category to db and return it", func(t *testing.T) {
			t.Cleanup(cleanup(db))
			it := assert.New(t)

			var categories []category.Category
			db.Find(&categories)
			initialLen := len(categories)
			json := `{"id":69,"title":"Pizza","removable":false}`

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/categories", strings.NewReader(json))
			r.ServeHTTP(w, req)

			it.Equal(http.StatusOK, w.Code)
			it.Contains(w.Header().Get("Content-Type"), "application/json")
			it.Equal(json, w.Body.String())

			var last category.Category

			db.Last(&last)
			it.False(last.Removable)
			it.Equal(last.Title, "Pizza")
			it.Equal(last.ID, uint(69))

			db.Find(&categories)
			it.Len(categories, initialLen+1)
		})
	})
}

func cleanup(db *gorm.DB) func() {
	return func() {
		db.Exec("TRUNCATE categories;")
	}
}

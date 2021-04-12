package e2e

import (
	"fmt"
	"food_ordering_backend/category"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var testCategories = []category.Category{{Model: gorm.Model{ID: 1}, Title: "Salads"}, {Model: gorm.Model{ID: 2}, Title: "Burgers"}, {Model: gorm.Model{ID: 3}, Title: "Pizza"}, {Model: gorm.Model{ID: 4}, Title: "Drinks"}}
var db = database.MustGetTest()
var r = router.Setup(db)

func TestCategories(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)

	t.Run("GET /categories", func(t *testing.T) {
		send := sendReq(http.MethodGet, "/categories")
		t.Run("should return JSON array of categories", func(t *testing.T) {
			t.Cleanup(cleanup)
			it := assert.New(t)
			var categories []category.Category

			db.Find(&categories)
			require.Len(t, categories, 0)

			db.Create(&testCategories)
			db.Find(&categories)
			require.Len(t, categories, len(testCategories))

			resp := send("")

			it.Equal(http.StatusOK, resp.Code)
			it.Contains(resp.Header().Get("Content-Type"), "application/json")
			it.Equal(
				`[{"id":1,"title":"Salads","removable":false},{"id":2,"title":"Burgers","removable":false},{"id":3,"title":"Pizza","removable":false},{"id":4,"title":"Drinks","removable":false}]`,
				resp.Body.String(),
			)
		})
	})

	t.Run("GET /categories/:id", func(t *testing.T) {
		t.Cleanup(cleanup)
		sendWithParam := func(id uint) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodGet, "/categories/"+param)("")
		}

		t.Run("should return category with provided id", func(t *testing.T) {
			t.Cleanup(cleanup)
			cleanup()
			it := assert.New(t)
			cat := testCategories[0]

			db.Create(&cat)
			var c category.Category
			db.First(&c, cat.ID)
			require.Equal(t, c.ID, cat.ID)

			resp := sendWithParam(cat.ID)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(
				fmt.Sprintf(`{"id":%d,"title":%q,"removable":%t}`, cat.ID, cat.Title, cat.Removable),
				resp.Body.String(),
			)
		})
	})

	t.Run("POST /categories", func(t *testing.T) {
		t.Cleanup(cleanup)
		send := sendReq(http.MethodPost, "/categories")

		t.Run("should add category to db and return it", func(t *testing.T) {
			t.Cleanup(cleanup)
			it := assert.New(t)

			var categories []category.Category
			db.Find(&categories)
			initialLen := len(categories)
			json := `{"id":69,"title":"Pizza","removable":false}`

			resp := send(json)

			it.Equal(http.StatusCreated, resp.Code)
			it.Contains(resp.Header().Get("Content-Type"), "application/json")
			it.Equal(json, resp.Body.String())

			var last category.Category

			db.Last(&last)
			it.False(last.Removable)
			it.Equal(last.Title, "Pizza")
			it.Equal(last.ID, uint(69))

			db.Find(&categories)
			it.Len(categories, initialLen+1)
		})

		t.Run("should return 400 if provided json isn't correct", func(t *testing.T) {
			it := assert.New(t)
			json := `{"title": 123}`

			resp := send(json)
			it.Equal(http.StatusBadRequest, resp.Code)
		})
	})
}

func sendReq(method, target string) func(body string) *httptest.ResponseRecorder {
	return func(body string) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(method, target, strings.NewReader(body))
		r.ServeHTTP(w, req)
		return w
	}
}

func cleanup() {
	db.Exec("TRUNCATE categories;")
}

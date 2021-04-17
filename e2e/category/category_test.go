package category_test

import (
	"fmt"
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var testCategories = []category.Category{{Model: gorm.Model{ID: 1}, Title: "Salads", Removable: true}, {Model: gorm.Model{ID: 2}, Title: "Burgers", Removable: true}, {Model: gorm.Model{ID: 3}, Title: "Pizza", Removable: true}, {Model: gorm.Model{ID: 4}, Title: "Drinks", Removable: true}}
var db = database.MustGetTest()
var r = router.Setup(db)

func TestCategories(t *testing.T) {
	db.Logger.LogMode(logger.Silent)

	t.Run("GET /categories", func(t *testing.T) {
		send := sendReq(http.MethodGet, "/categories")

		t.Run("should return JSON array of categories", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			resp := send("")

			it.Equal(http.StatusOK, resp.Code)
			it.Contains(resp.Header().Get("Content-Type"), "application/json")
			it.Equal(
				`[{"id":1,"title":"Salads","removable":true},{"id":2,"title":"Burgers","removable":true},{"id":3,"title":"Pizza","removable":true},{"id":4,"title":"Drinks","removable":true}]`,
				resp.Body.String(),
			)
		})
	})

	t.Run("GET /categories/:id", func(t *testing.T) {
		sendWithParam := func(id uint) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodGet, "/categories/"+param)("")
		}

		t.Run("should return category with provided id", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)

			for _, cat := range testCategories {
				resp := sendWithParam(cat.ID)
				it.Equal(http.StatusOK, resp.Code)
				it.Equal(
					fmt.Sprintf(`{"id":%d,"title":%q,"removable":%t}`, cat.ID, cat.Title, cat.Removable),
					resp.Body.String(),
				)
			}

		})

		runFindByIDTests(t)
	})

	t.Run("POST /categories", func(t *testing.T) {
		send := sendReq(http.MethodPost, "/categories")

		t.Run("should add category to db and return it", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
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

			var categories []category.Category
			if it.NoError(db.Find(&categories).Error) {
				it.Len(categories, len(testCategories)+1)
			}
		})

		t.Run("should return 400 if provided json isn't correct", func(t *testing.T) {
			it := assert.New(t)
			json := `{"title": 123}`

			resp := send(json)
			it.Equal(http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 409 if category already exists", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			json := `{"id":13,"title":"Salads","removable":false}`

			resp := send(json)
			it.Equal(http.StatusCreated, resp.Code)

			resp = send(json)
			it.Equal(http.StatusConflict, resp.Code)
		})
	})

	t.Run("PUT /categories/:id", func(t *testing.T) {
		sendWithParam := func(id uint, body string) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodPut, "/categories/"+param)(body)
		}

		t.Run("should update category in db based on provided json", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			testCategory := testCategories[0]
			updateJSON := `{"title":"Sushi","removable":true}`

			resp := sendWithParam(testCategory.ID, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(fmt.Sprintf(`{"id":%d,"title":"Sushi","removable":true}`, testCategory.ID), resp.Body.String())
		})

		t.Run("should ignore id in provided json", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			testCategory := testCategories[0]
			updateJSON := `{"id":420,"title":"Sushi","removable":true}`
			require.NotEqual(t, testCategory.ID, 420)

			resp := sendWithParam(testCategory.ID, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(fmt.Sprintf(`{"id":%d,"title":"Sushi","removable":true}`, testCategory.ID), resp.Body.String())
		})

		runFindByIDTests(t)
	})

	t.Run("DELETE /categories/:id", func(t *testing.T) {
		sendWithParam := func(id uint) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodDelete, "/categories/"+param)("")
		}

		t.Run("should removed a category with provided ID from db", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			testCat := testCategories[len(testCategories)/2]
			resp := sendWithParam(testCat.ID)

			it.Equal(http.StatusOK, resp.Code)

			var categories []category.Category
			db.Find(&categories)
			it.Len(categories, len(testCategories)-1)

			for _, c := range categories {
				it.NotEqual(c.ID, testCat.ID)
			}
		})

		t.Run("should return 403 if category isn't removable", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			c := category.Category{Title: "Fish", Removable: false}

			db.Create(&c)

			resp := sendWithParam(c.ID)
			it.Equal(http.StatusForbidden, resp.Code)
		})
	})
}

func runFindByIDTests(t *testing.T) {
	t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
		setupDB(t)
		resp := sendReq(http.MethodGet, "/categories/some-random-id")("")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 404 if category with provided id doesn't exist", func(t *testing.T) {
		setupDB(t)
		resp := sendReq(http.MethodGet, "/categories/69")("")
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}

func setupDB(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)
	db.Create(&testCategories)
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
	db.Exec("TRUNCATE categories CASCADE;")
}

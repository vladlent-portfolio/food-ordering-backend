package category_test

import (
	"fmt"
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/database"
	"food_ordering_backend/e2e/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm/logger"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var testCategories = []category.Category{{ID: 1, Title: "Salads", Removable: true}, {ID: 2, Title: "Burgers", Removable: true}, {ID: 3, Title: "Pizza", Removable: true}, {ID: 4, Title: "Drinks", Removable: true}}
var db = database.MustGetTest()

func TestCategories(t *testing.T) {
	db.Logger.LogMode(logger.Silent)

	t.Run("GET /categories", func(t *testing.T) {
		send := testutils.SendReq(http.MethodGet, "/categories")

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
			return testutils.SendReq(http.MethodGet, "/categories/"+param)("")
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
		send := testutils.ReqWithCookie(http.MethodPost, "/categories")

		t.Run("should add category to db and return it", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			setupDB(t)
			it := assert.New(t)
			json := `{"id":69,"title":"Pizza","removable":false}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := send(c, json)

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
			testutils.SetupUsersDB(t)
			it := assert.New(t)
			json := `{"title": 123}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := send(c, json)
			it.Equal(http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 409 if category already exists", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			setupDB(t)
			it := assert.New(t)
			json := `{"id":13,"title":"Salads","removable":false}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := send(c, json)
			it.Equal(http.StatusCreated, resp.Code)

			resp = send(c, json)
			it.Equal(http.StatusConflict, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodPost, "/categories", true)
	})

	t.Run("PUT /categories/:id", func(t *testing.T) {
		sendWithParam := func(id uint, body string, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodPut, "/categories/"+param)(c, body)
		}

		t.Run("should update category in db based on provided json", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			setupDB(t)
			it := assert.New(t)
			testCategory := testCategories[0]
			updateJSON := `{"title":"Sushi","removable":true}`

			_, c := testutils.LoginAsRandomAdmin(t)

			resp := sendWithParam(testCategory.ID, updateJSON, c)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(fmt.Sprintf(`{"id":%d,"title":"Sushi","removable":true}`, testCategory.ID), resp.Body.String())
		})

		t.Run("should ignore id in provided json", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			setupDB(t)
			it := assert.New(t)
			testCategory := testCategories[0]
			updateJSON := `{"id":420,"title":"Sushi","removable":true}`
			require.NotEqual(t, testCategory.ID, 420)

			_, c := testutils.LoginAsRandomAdmin(t)

			resp := sendWithParam(testCategory.ID, updateJSON, c)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(fmt.Sprintf(`{"id":%d,"title":"Sushi","removable":true}`, testCategory.ID), resp.Body.String())
		})

		testutils.RunAuthTests(t, http.MethodPut, "/categories/69", true)
		runFindByIDTests(t)
	})

	t.Run("DELETE /categories/:id", func(t *testing.T) {
		sendWithParam := func(id uint, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodDelete, "/categories/"+param)(c, "")
		}

		t.Run("should removed a category with provided ID from db", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			setupDB(t)
			it := assert.New(t)
			testCat := testCategories[len(testCategories)/2]
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := sendWithParam(testCat.ID, c)

			it.Equal(http.StatusOK, resp.Code)

			var categories []category.Category
			db.Find(&categories)
			it.Len(categories, len(testCategories)-1)

			for _, c := range categories {
				it.NotEqual(c.ID, testCat.ID)
			}
		})

		t.Run("should return 403 if category isn't removable", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			setupDB(t)
			it := assert.New(t)
			c := category.Category{ID: 69, Title: "Seafood", Removable: false}
			_, cookie := testutils.LoginAsRandomAdmin(t)

			if it.NoError(db.Create(&c).Error) {
				resp := sendWithParam(c.ID, cookie)
				it.Equal(http.StatusForbidden, resp.Code)
			}
		})

		testutils.RunAuthTests(t, http.MethodDelete, "/categories/69", true)
	})
}

func runFindByIDTests(t *testing.T) {
	t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
		setupDB(t)
		resp := testutils.SendReq(http.MethodGet, "/categories/some-random-id")("")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 404 if category with provided id doesn't exist", func(t *testing.T) {
		setupDB(t)
		resp := testutils.SendReq(http.MethodGet, "/categories/69")("")
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}

func setupDB(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)
	require.NoError(t, db.Create(&testCategories).Error)
}

func cleanup() {
	db.Exec("TRUNCATE categories CASCADE;")
}

package dish_test

import (
	"fmt"
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
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

var db = database.MustGetTest()
var r = router.Setup(db)
var testCategories = []category.Category{
	{Model: gorm.Model{ID: 1}, Title: "Salads", Removable: true},
	{Model: gorm.Model{ID: 2}, Title: "Burgers", Removable: true},
	{Model: gorm.Model{ID: 3}, Title: "Pizza", Removable: true},
	{Model: gorm.Model{ID: 4}, Title: "Drinks", Removable: true},
}
var testDishes = []dish.Dish{
	{Model: gorm.Model{ID: 1}, Title: "Fresh and Healthy Salad", Price: 2.65, CategoryID: 1},
	{Model: gorm.Model{ID: 2}, Title: "Crunchy Cashew Salad", Price: 3.22, CategoryID: 1},
	{Model: gorm.Model{ID: 3}, Title: "Hamburger", Price: 1.99, CategoryID: 2},
	{Model: gorm.Model{ID: 4}, Title: "Cheeseburger", Price: 2.28, CategoryID: 2},
	{Model: gorm.Model{ID: 5}, Title: "Margherita", Price: 4.20, CategoryID: 3},
	{Model: gorm.Model{ID: 6}, Title: "4 Cheese", Price: 4.69, CategoryID: 3},
	{Model: gorm.Model{ID: 7}, Title: "Pepsi 2L", Price: 1.50, CategoryID: 4},
	{Model: gorm.Model{ID: 8}, Title: "Orange Juice 2L", Price: 2, CategoryID: 4},
}

var testCategoriesJSON = `[{"id":1,"title":"Salads","removable":true},{"id":2,"title":"Burgers","removable":true},{"id":3,"title":"Pizza","removable":true},{"id":4,"title":"Drinks","removable":true}]`
var testDishesJSON = `[{"id":1,"title":"Fresh and Healthy Salad","price":2.65,"category_id":1,"category":{"id":1,"title":"Salads","removable":true}},{"id":2,"title":"Crunchy Cashew Salad","price":3.22,"category_id":1,"category":{"id":1,"title":"Salads","removable":true}},{"id":3,"title":"Hamburger","price":1.99,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}},{"id":4,"title":"Cheeseburger","price":2.28,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}},{"id":5,"title":"Margherita","price":4.2,"category_id":3,"category":{"id":3,"title":"Pizza","removable":true}},{"id":6,"title":"4 Cheese","price":4.69,"category_id":3,"category":{"id":3,"title":"Pizza","removable":true}},{"id":7,"title":"Pepsi 2L","price":1.5,"category_id":4,"category":{"id":4,"title":"Drinks","removable":true}},{"id":8,"title":"Orange Juice 2L","price":2,"category_id":4,"category":{"id":4,"title":"Drinks","removable":true}}]`

func TestDishes(t *testing.T) {
	cleanup()
	t.Cleanup(cleanup)
	db.Logger.LogMode(logger.Silent)

	t.Run("GET /dishes", func(t *testing.T) {
		send := sendReq(http.MethodGet, "/dishes")
		t.Run("should return JSON array of dishes with their respective categories", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			resp := send("")

			it.Equal(http.StatusOK, resp.Code)
			it.Contains(resp.Header().Get("Content-Type"), "application/json")
			it.Equal(
				testDishesJSON,
				resp.Body.String(),
			)
		})
	})

	t.Run("GET /dishes/:id", func(t *testing.T) {
		sendWithParam := func(id uint) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodGet, "/dishes/"+param)("")
		}

		t.Run("should return dish with provided id", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)

			for i, testDish := range testDishes {
				resp := sendWithParam(testDish.ID)
				testCat := testCategories[i/2]
				it.Equal(http.StatusOK, resp.Code)
				it.Equal(
					fmt.Sprintf(
						`{"id":%d,"title":%q,"price":%v,"category_id":%d,"category":{"id":%d,"title":%q,"removable":%t}}`,
						testDish.ID, testDish.Title, testDish.Price, testDish.CategoryID, testCat.ID, testCat.Title, testCat.Removable,
					),
					resp.Body.String(),
				)
			}

		})

		runFindByIDTests(t)
	})

	t.Run("POST /dishes", func(t *testing.T) {
		send := sendReq(http.MethodPost, "/dishes")

		t.Run("should add dish to db and return it", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			initialLen := len(testDishes)
			reqJSON := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2}`
			respJSON := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}}`

			resp := send(reqJSON)

			it.Equal(http.StatusCreated, resp.Code)
			it.Contains(resp.Header().Get("Content-Type"), "application/json")
			it.Equal(respJSON, resp.Body.String())

			var last dish.Dish
			var dishes []dish.Dish

			err := db.Preload("Category").Last(&last).Error
			require.NoError(t, err)

			it.Equal(4.56, last.Price)
			it.Equal("Double Cheeseburger", last.Title)
			it.Equal(uint(69), last.ID)
			it.Equal(uint(2), last.CategoryID)
			it.Equal(uint(2), last.Category.ID)
			it.Equal("Burgers", last.Category.Title)
			it.True(last.Category.Removable)

			db.Find(&dishes)
			it.Len(dishes, initialLen+1)
		})

		t.Run("should return 400 if provided json isn't correct", func(t *testing.T) {
			it := assert.New(t)
			json := `{"title": 123}`

			resp := send(json)
			it.Equal(http.StatusBadRequest, resp.Code)
		})

		// TODO: Add test for negative price
		//t.Run("should ", func(t *testing.T) {
		//    it := assert.New(t)
		//
		//})

		t.Run("should return 409 if dish already exists", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			json := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2}`

			resp := send(json)
			it.Equal(http.StatusCreated, resp.Code)

			resp = send(json)
			it.Equal(http.StatusConflict, resp.Code)
		})
	})

	t.Run("PUT /categories/:id", func(t *testing.T) {
		t.Cleanup(cleanup)
		sendWithParam := func(id uint, body string) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodPut, "/categories/"+param)(body)
		}

		setup := func(t *testing.T, initial category.Category) category.Category {
			var c category.Category
			db.Create(&initial)
			db.First(&c, initial.ID)
			require.Equal(t, c.ID, initial.ID)
			return c
		}

		t.Run("should update category in db based on provided json", func(t *testing.T) {
			cleanup()
			t.Cleanup(cleanup)
			it := assert.New(t)
			testCategory := testCategories[0]
			updateJSON := `{"title":"Sushi","removable":true}`

			c := setup(t, testCategory)

			resp := sendWithParam(c.ID, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(fmt.Sprintf(`{"id":%d,"title":"Sushi","removable":true}`, testCategory.ID), resp.Body.String())
		})

		t.Run("should ignore id in provided json", func(t *testing.T) {
			cleanup()
			t.Cleanup(cleanup)
			it := assert.New(t)
			testCategory := testCategories[0]
			updateJSON := `{"id":420,"title":"Sushi","removable":true}`
			require.NotEqual(t, testCategory.ID, 420)

			c := setup(t, testCategory)

			resp := sendWithParam(c.ID, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(fmt.Sprintf(`{"id":%d,"title":"Sushi","removable":true}`, testCategory.ID), resp.Body.String())
		})

		runFindByIDTests(t)
	})

	t.Run("DELETE /categories/:id", func(t *testing.T) {
		t.Cleanup(cleanup)
		sendWithParam := func(id uint) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodDelete, "/categories/"+param)("")
		}

		t.Run("should removed a category with provided ID from db", func(t *testing.T) {
			t.Cleanup(cleanup)
			cleanup()
			it := assert.New(t)
			var categories []category.Category

			db.Find(&categories)
			require.Len(t, categories, 0)

			db.Create(&testCategories)
			db.Find(&categories)
			require.Len(t, categories, len(testCategories))

			testCat := testCategories[len(testCategories)/2]
			resp := sendWithParam(testCat.ID)

			it.Equal(http.StatusOK, resp.Code)

			db.Find(&categories)
			it.Len(categories, len(testCategories)-1)

			for _, c := range categories {
				it.NotEqual(c.ID, testCat.ID)
			}
		})

		t.Run("should return 403 if category isn't removable", func(t *testing.T) {
			t.Cleanup(cleanup)
			cleanup()
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
		resp := sendReq(http.MethodGet, "/dishes/some-random-id")("")
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("should return 404 if category with provided id doesn't exist", func(t *testing.T) {
		setupDB(t)
		resp := sendReq(http.MethodGet, "/dishes/69")("")
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}

func setupDB(t *testing.T) {
	t.Cleanup(cleanup)
	cleanup()
	req := require.New(t)

	req.NoError(db.Create(&testCategories).Error)
	req.NoError(db.Create(&testDishes).Error)
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

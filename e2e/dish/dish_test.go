package dish_test

import (
	"encoding/json"
	"fmt"
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/database"
	"food_ordering_backend/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
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
	{Model: gorm.Model{ID: 1}, Title: "Fresh and Healthy Salad", Price: 2.65, CategoryID: 1, Category: testCategories[0]},
	{Model: gorm.Model{ID: 2}, Title: "Crunchy Cashew Salad", Price: 3.22, CategoryID: 1, Category: testCategories[0]},
	{Model: gorm.Model{ID: 3}, Title: "Hamburger", Price: 1.99, CategoryID: 2, Category: testCategories[1]},
	{Model: gorm.Model{ID: 4}, Title: "Cheeseburger", Price: 2.28, CategoryID: 2, Category: testCategories[1]},
	{Model: gorm.Model{ID: 5}, Title: "Margherita", Price: 4.20, CategoryID: 3, Category: testCategories[2]},
	{Model: gorm.Model{ID: 6}, Title: "4 Cheese", Price: 4.69, CategoryID: 3, Category: testCategories[2]},
	{Model: gorm.Model{ID: 7}, Title: "Pepsi 2L", Price: 1.50, CategoryID: 4, Category: testCategories[3]},
	{Model: gorm.Model{ID: 8}, Title: "Orange Juice 2L", Price: 2, CategoryID: 4, Category: testCategories[3]},
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

		t.Run("should return dishes filtered by provided category id", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)

			for i, c := range testCategories {
				dishes := testDishes[i*2 : i*2+2]
				resp := sendReq(http.MethodGet, fmt.Sprintf("/dishes?cid=%d", c.ID))("")

				it.Equal(http.StatusOK, resp.Code)
				it.Contains(resp.Header().Get("Content-Type"), "application/json")

				var dtos []dish.DTO
				err := json.NewDecoder(resp.Body).Decode(&dtos)
				require.NoError(t, err)
				require.Len(t, dtos, len(dishes))

				for i, d := range dishes {
					dto := dtos[i]
					it.Equal(d.ID, dto.ID)
					it.Equal(d.Title, dto.Title)
					it.Equal(d.Price, dto.Price)
					it.Equal(d.CategoryID, dto.CategoryID)
					it.Equal(d.Category.ID, dto.Category.ID)
					it.Equal(d.Category.Title, dto.Category.Title)
					it.Equal(d.Category.Removable, dto.Category.Removable)
				}
			}
		})

		t.Run("should return 400 if provided category id isn't a number", func(t *testing.T) {
			resp := sendReq(http.MethodGet, "/dishes?cid=hello")("")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return an empty array if category doesn't exist", func(t *testing.T) {
			it := assert.New(t)
			resp := sendReq(http.MethodGet, "/dishes?cid=228")("")
			it.Equal(http.StatusOK, resp.Code)
			it.Equal("[]", resp.Body.String())
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

		t.Run("should return 409 if dish already exists", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			json := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2}`

			resp := send(json)
			it.Equal(http.StatusCreated, resp.Code)

			resp = send(json)
			it.Equal(http.StatusConflict, resp.Code)
		})

		negativePriceTest(t, http.MethodPost)
	})

	t.Run("PUT /dishes/:id", func(t *testing.T) {
		sendWithParam := func(id uint, body string) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodPut, "/dishes/"+param)(body)
		}

		t.Run("should update category in db based on provided json", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			updateJSON := `{"title":"Double Cheeseburger","price":4.56,"category_id":2}`
			respJSON := `{"id":4,"title":"Double Cheeseburger","price":4.56,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}}`

			resp := sendWithParam(4, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(respJSON, resp.Body.String())
		})

		t.Run("should ignore id in provided json", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			updateJSON := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2}`
			respJSON := `{"id":4,"title":"Double Cheeseburger","price":4.56,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}}`

			resp := sendWithParam(4, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(respJSON, resp.Body.String())
		})

		t.Run("should correctly handle category change", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			updateJSON := `{"id":69,"title":"Meat Supreme","price":3.22,"category_id":3}`
			respJSON := `{"id":4,"title":"Meat Supreme","price":3.22,"category_id":3,"category":{"id":3,"title":"Pizza","removable":true}}`

			resp := sendWithParam(4, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(respJSON, resp.Body.String())
		})

		runFindByIDTests(t)
		negativePriceTest(t, http.MethodPut)
	})

	t.Run("DELETE /dishes/:id", func(t *testing.T) {
		sendWithParam := func(id uint) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return sendReq(http.MethodDelete, "/dishes/"+param)("")
		}

		t.Run("should remove a dish with provided ID from db", func(t *testing.T) {
			setupDB(t)
			it := assert.New(t)
			randomIndex := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(testDishes))
			testDish := testDishes[randomIndex]
			resp := sendWithParam(testDish.ID)

			it.Equal(http.StatusOK, resp.Code)
			var dishes []dish.Dish

			db.Find(&dishes)
			it.Len(dishes, len(testDishes)-1)

			for _, d := range dishes {
				it.NotEqual(d.ID, testDish.ID)
			}
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

func negativePriceTest(t *testing.T, method string) {
	t.Run("should return 400 if price is < 0", func(t *testing.T) {
		setupDB(t)
		it := assert.New(t)
		json := `{"id":1,"title":"Meat Supreme","price":-3.22,"category_id":3}`
		var resp *httptest.ResponseRecorder

		if method == http.MethodPost {
			resp = sendReq(method, "/dishes")(json)
		} else {
			resp = sendReq(method, "/dishes/1")(json)
		}

		it.Equal(http.StatusBadRequest, resp.Code)
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

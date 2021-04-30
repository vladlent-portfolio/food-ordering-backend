package dish_test

import (
	"encoding/json"
	"fmt"
	"food_ordering_backend/common"
	"food_ordering_backend/config"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/database"
	"food_ordering_backend/e2e/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var db = database.MustGetTest()

func TestDishes(t *testing.T) {
	t.Run("GET /dishes", func(t *testing.T) {
		send := testutils.SendReq(http.MethodGet, "/dishes")
		t.Run("should return JSON array of dishes with their respective categories", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			resp := send("")

			it.Equal(http.StatusOK, resp.Code)
			it.Contains(resp.Header().Get("Content-Type"), "application/json")
			it.Equal(
				testutils.TestDishesJSON,
				resp.Body.String(),
			)
		})

		t.Run("should return dishes filtered by provided category id", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)

			for i, c := range testutils.TestCategories {
				dishes := testutils.TestDishes[i*2 : i*2+2]
				resp := testutils.SendReq(http.MethodGet, fmt.Sprintf("/dishes?cid=%d", c.ID))("")

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
			resp := testutils.SendReq(http.MethodGet, "/dishes?cid=hello")("")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return an empty array if category doesn't exist", func(t *testing.T) {
			it := assert.New(t)
			resp := testutils.SendReq(http.MethodGet, "/dishes?cid=228")("")
			it.Equal(http.StatusOK, resp.Code)
			it.Equal("[]", resp.Body.String())
		})
	})

	t.Run("GET /dishes/:id", func(t *testing.T) {
		sendWithParam := func(id uint) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.SendReq(http.MethodGet, "/dishes/"+param)("")
		}

		t.Run("should return dish with provided id", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)

			for i, testDish := range testutils.TestDishes {
				resp := sendWithParam(testDish.ID)
				testCat := testutils.TestCategories[i/2]
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

		t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
			resp := testutils.SendReq(http.MethodGet, "/dishes/some-random-id")("")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if dish with provided id doesn't exist", func(t *testing.T) {
			resp := testutils.SendReq(http.MethodGet, "/dishes/69")("")
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})
	})

	t.Run("POST /dishes", func(t *testing.T) {
		send := testutils.ReqWithCookie(http.MethodPost, "/dishes")

		t.Run("should add dish to db and return it", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			initialLen := len(testutils.TestDishes)
			reqJSON := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2}`
			respJSON := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}}`

			_, c := testutils.LoginAsRandomAdmin(t)
			resp := send(c, reqJSON)

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
			testutils.SetupUsersDB(t)
			it := assert.New(t)
			json := `{"title": 123}`

			_, c := testutils.LoginAsRandomAdmin(t)
			resp := send(c, json)

			it.Equal(http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 409 if dish already exists", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			json := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := send(c, json)
			it.Equal(http.StatusCreated, resp.Code)

			resp = send(c, json)
			it.Equal(http.StatusConflict, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodPost, "/dishes", true)
		negativePriceTest(t, http.MethodPost)
	})

	t.Run("PATCH /dishes/:id/upload", func(t *testing.T) {
		upload := func(id uint, c *http.Cookie, fileName string, file io.Reader) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.UploadReqWithCookie(http.MethodPatch, "/dishes/"+param+"/upload", "image")(c, fileName, file)
		}

		t.Run("should upload an image, update dish in db and return a link to image", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			testutils.SetupUsersDB(t)
			t.Cleanup(testutils.CleanupStaticFolder)
			it := assert.New(t)
			d := testutils.TestDishes[4]

			img, err := os.Open(testutils.PathToFile("./img/pizza.png"))
			require.NoError(t, err)
			defer img.Close()

			fileName := filepath.Base(img.Name())
			stat, err := img.Stat()
			require.NoError(t, err)
			expectedName := fmt.Sprintf("%d.png", d.ID)

			_, cookie := testutils.LoginAsRandomAdmin(t)

			resp := upload(d.ID, cookie, fileName, img)

			if it.Equal(http.StatusOK, resp.Code) {
				link, err := url.Parse(resp.Body.String())

				if it.NoError(err, "expected valid link to image in response") {
					it.Equal(expectedName, path.Base(link.String()), "expected filename to be 'dish_id'+'file_extension'")
					resp := testutils.SendReq(http.MethodGet, link.String())("")

					if it.Equal(http.StatusOK, resp.Code) {
						it.Contains(resp.Header().Get("Content-Type"), "image/png", "expected served image to have correct Content-Type")
						it.Equal(stat.Size(), resp.Result().ContentLength, "expected served image to be the same size as uploaded")
					}
				}
			}

			if it.NoError(db.First(&d).Error) {
				if it.NotNil(d.Image) {
					it.Equal(expectedName, *d.Image, "expected filename to be 'dish_id'+'file_extension'")
				}
			}

			if it.DirExists(config.CategoriesImgDirAbs) {
				it.FileExists(filepath.Join(config.CategoriesImgDirAbs, expectedName))
			}
		})

		t.Run("should return 415 if file type is not supported", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := upload(3, c, "img.json", strings.NewReader("{}"))
			assert.Equal(t, http.StatusUnsupportedMediaType, resp.Code)
		})

		t.Run("should return 413 if file size is too big", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			t.Cleanup(testutils.CleanupStaticFolder)
			img, err := os.Open(testutils.PathToFile("./img/big-image.jpg"))

			if assert.NoError(t, err) {
				_, c := testutils.LoginAsRandomAdmin(t)

				resp := upload(3, c, "photo.png", img)
				assert.Equal(t, http.StatusRequestEntityTooLarge, resp.Code)
			}
		})

		t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPatch, "/dishes/some-random-id/upload")(c, "")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if dishes with provided id doesn't exist", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPatch, "/dishes/69/upload")(c, "")
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodPatch, "/dishes/1337/upload", true)
	})

	t.Run("PUT /dishes/:id", func(t *testing.T) {
		sendWithParam := func(id uint, c *http.Cookie, body string) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodPut, "/dishes/"+param)(c, body)
		}

		t.Run("should update dish in db based on provided json", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			updateJSON := `{"title":"Double Cheeseburger","price":4.56,"category_id":2}`
			respJSON := `{"id":4,"title":"Double Cheeseburger","price":4.56,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := sendWithParam(4, c, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(respJSON, resp.Body.String())
		})

		t.Run("should ignore id in provided json", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			updateJSON := `{"id":69,"title":"Double Cheeseburger","price":4.56,"category_id":2}`
			respJSON := `{"id":4,"title":"Double Cheeseburger","price":4.56,"category_id":2,"category":{"id":2,"title":"Burgers","removable":true}}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := sendWithParam(4, c, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(respJSON, resp.Body.String())
		})

		t.Run("should correctly handle category change", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			updateJSON := `{"id":69,"title":"Meat Supreme","price":3.22,"category_id":3}`
			respJSON := `{"id":4,"title":"Meat Supreme","price":3.22,"category_id":3,"category":{"id":3,"title":"Pizza","removable":true}}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := sendWithParam(4, c, updateJSON)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(respJSON, resp.Body.String())
		})

		t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPut, "/dishes/some-random-id/upload")(c, "")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if category with provided id doesn't exist", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPut, "/dishes/69/upload")(c, "")
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodPut, "/dishes/1", true)
		negativePriceTest(t, http.MethodPut)
	})

	t.Run("DELETE /dishes/:id", func(t *testing.T) {
		sendWithParam := func(id uint, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodDelete, "/dishes/"+param)(c, "")
		}

		t.Run("should remove a dish with provided ID from db", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			testDishes := testutils.TestDishes
			randomIndex := common.RandomInt(len(testDishes))
			testDish := testDishes[randomIndex]
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := sendWithParam(testDish.ID, c)

			it.Equal(http.StatusOK, resp.Code)
			var dishes []dish.Dish

			db.Find(&dishes)
			it.Len(dishes, len(testDishes)-1)

			for _, d := range dishes {
				it.NotEqual(d.ID, testDish.ID)
			}
		})

		t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodDelete, "/dishes/some-random-id/upload")(c, "")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if dish with provided id doesn't exist", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodDelete, "/dishes/69/upload")(c, "")
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodDelete, "/dishes/1", true)
	})
}

func negativePriceTest(t *testing.T, method string) {
	t.Run("should return 400 if price is < 0", func(t *testing.T) {
		testutils.SetupUsersDB(t)
		testutils.SetupDishesAndCategories(t)

		it := assert.New(t)
		json := `{"id":1,"title":"Meat Supreme","price":-3.22,"category_id":3}`
		_, c := testutils.LoginAsRandomAdmin(t)
		var resp *httptest.ResponseRecorder

		if method == http.MethodPost {
			resp = testutils.ReqWithCookie(method, "/dishes")(c, json)
		} else {
			resp = testutils.ReqWithCookie(method, "/dishes/1")(c, json)
		}

		it.Equal(http.StatusBadRequest, resp.Code)
	})
}

package category_test

import (
	"encoding/json"
	"fmt"
	"food_ordering_backend/config"
	"food_ordering_backend/controllers/category"
	"food_ordering_backend/database"
	"food_ordering_backend/tests/testutils"
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

func TestCategories(t *testing.T) {
	t.Run("GET /categories", func(t *testing.T) {
		send := testutils.SendReq(http.MethodGet, "/categories")

		t.Run("should return sorted by id array of categories", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			resp := send("")

			if it.Equal(http.StatusOK, resp.Code) {
				it.Contains(resp.Header().Get("Content-Type"), "application/json")
				var dtos []category.DTO

				if it.NoError(json.NewDecoder(resp.Body).Decode(&dtos)) {
					it.Len(dtos, len(testutils.TestCategories))

					ids := make([]uint, len(dtos))
					for i, dto := range dtos {
						ids[i] = dto.ID
					}
					it.IsIncreasing(ids, "expected array to be sorted by id")

					for _, dto := range dtos {
						c := testutils.FindTestCategoryByID(dto.ID)
						it.Equal(c.ID, dto.ID)
						it.Equal(c.Title, dto.Title)
						it.Equal(c.Removable, dto.Removable)
						it.Equal(imgURL(*c.Image), *dto.Image)
					}
				}
			}
		})

	})

	t.Run("GET /categories/:id", func(t *testing.T) {
		sendWithParam := func(id uint) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.SendReq(http.MethodGet, "/categories/"+param)("")
		}

		t.Run("should return category with provided id", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)

			for _, cat := range testutils.TestCategories {
				resp := sendWithParam(cat.ID)
				it.Equal(http.StatusOK, resp.Code)
				it.Equal(
					fmt.Sprintf(`{"id":%d,"title":%q,"removable":%t,"image":%q}`, cat.ID, cat.Title, cat.Removable, imgURL(*cat.Image)),
					resp.Body.String(),
				)
			}

		})

		t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
			resp := testutils.SendReq(http.MethodGet, "/categories/some-random-id")("")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if category with provided id doesn't exist", func(t *testing.T) {
			resp := testutils.SendReq(http.MethodGet, "/categories/69")("")
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})

	})

	t.Run("POST /categories", func(t *testing.T) {
		send := testutils.ReqWithCookie(http.MethodPost, "/categories")

		t.Run("should add category to db and return it", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			json := `{"id":69,"title":"Seafood","removable":true}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := send(c, json)

			it.Equal(http.StatusCreated, resp.Code)
			it.Contains(resp.Header().Get("Content-Type"), "application/json")
			it.Equal(json, resp.Body.String())

			var last category.Category

			db.Last(&last)
			it.True(last.Removable)
			it.Equal(last.Title, "Seafood")
			it.Equal(last.ID, uint(69))

			var categories []category.Category
			if it.NoError(db.Find(&categories).Error) {
				it.Len(categories, len(testutils.TestCategories)+1)
			}
		})

		t.Run("should trim title", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			_, c := testutils.LoginAsRandomAdmin(t)

			tests := []string{"  Pizza", "Pizza   ", "   Pizza   "}
			for _, tc := range tests {
				json := fmt.Sprintf(`{"title":%q}`, tc)
				resp := send(c, json)
				it.Equal(http.StatusConflict, resp.Code)
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
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			json := `{"title":"Seafood","removable":false}`
			_, c := testutils.LoginAsRandomAdmin(t)

			resp := send(c, json)
			it.Equal(http.StatusCreated, resp.Code)

			resp = send(c, json)
			it.Equal(http.StatusConflict, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodPost, "/categories", true)
	})

	t.Run("PATCH /categories/:id/upload", func(t *testing.T) {

		t.Run("should upload an image, update category in db and return a link to image", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			testutils.SetupUsersDB(t)
			t.Cleanup(testutils.CleanupStaticFolder)
			it := assert.New(t)
			c := testutils.TestCategories[2]
			c.Image = nil
			require.NoError(t, db.Save(&c).Error)

			img, err := os.Open(testutils.PathToFile("./img/pizza.png"))
			require.NoError(t, err)
			defer img.Close()

			fileName := filepath.Base(img.Name())
			stat, err := img.Stat()
			require.NoError(t, err)
			expectedName := fmt.Sprintf("%d.png", c.ID)

			_, cookie := testutils.LoginAsRandomAdmin(t)

			resp := upload(c.ID, cookie, fileName, img)

			if it.Equal(http.StatusOK, resp.Code) {
				link, err := url.Parse(resp.Body.String())

				if it.NoError(err, "expected valid link to image in response") {
					it.Equal(expectedName, filepath.Base(link.String()), "expected filename to be 'category_id'+'file_extension'")
					resp := testutils.SendReq(http.MethodGet, link.String())("")

					if it.Equal(http.StatusOK, resp.Code) {
						it.Contains(resp.Header().Get("Content-Type"), "image/png", "expected served image to have correct Content-Type")
						it.Equal(stat.Size(), resp.Result().ContentLength, "expected served image to be the same size as uploaded")
					}
				}
			}

			if it.NoError(db.First(&c).Error) {
				if it.NotNil(c.Image) {
					it.Equal(expectedName, *c.Image, "expected filename to be 'category_id'+'file_extension'")
				}
			}

			if it.DirExists(config.CategoriesImgDirAbs) {
				it.FileExists(filepath.Join(config.CategoriesImgDirAbs, expectedName))
			}
		})

		t.Run("should replace previous image with a new one", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			testutils.SetupUsersDB(t)
			t.Cleanup(testutils.CleanupStaticFolder)
			it := assert.New(t)
			_, cookie := testutils.LoginAsRandomAdmin(t)

			c := testutils.TestCategories[2]
			c.Image = nil
			require.NoError(t, db.Save(&c).Error)

			img, err := os.Open(testutils.PathToFile("./img/pizza.png"))
			require.NoError(t, err)
			defer img.Close()

			oldName := fmt.Sprintf("%d.png", c.ID)
			resp := upload(c.ID, cookie, filepath.Base(img.Name()), img)

			if it.Equal(http.StatusOK, resp.Code) && it.FileExists(filepath.Join(config.CategoriesImgDirAbs, oldName)) {
				newImg, err := os.Open(testutils.PathToFile("./img/hawaiian.webp"))
				require.NoError(t, err)
				defer newImg.Close()
				newName := fmt.Sprintf("%d.webp", c.ID)

				resp = upload(c.ID, cookie, filepath.Base(newImg.Name()), newImg)

				if it.Equal(http.StatusOK, resp.Code) {
					it.NoFileExists(filepath.Join(config.CategoriesImgDirAbs, oldName))
					it.FileExists(filepath.Join(config.CategoriesImgDirAbs, newName))
				}

				if it.NoError(db.First(&c).Error) {
					if it.NotNil(c.Image) {
						it.Equal(newName, *c.Image, "expected category to have a new image")
					}
				}
			}
		})

		t.Run("should return 415 if file type is not supported", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			c := testutils.TestCategories[2]
			c.Image = nil
			require.NoError(t, db.Save(&c).Error)

			_, cookie := testutils.LoginAsRandomAdmin(t)
			resp := upload(c.ID, cookie, "img.json", strings.NewReader("{}"))
			assert.Equal(t, http.StatusUnsupportedMediaType, resp.Code)
		})

		t.Run("should return 413 if file size is too big", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			t.Cleanup(testutils.CleanupStaticFolder)
			c := testutils.TestCategories[2]
			c.Image = nil
			require.NoError(t, db.Save(&c).Error)

			img, err := os.Open(testutils.PathToFile("./img/big-image.jpg"))

			if assert.NoError(t, err) {
				_, cookie := testutils.LoginAsRandomAdmin(t)

				resp := upload(c.ID, cookie, "photo.png", img)
				assert.Equal(t, http.StatusRequestEntityTooLarge, resp.Code)
			}
		})

		t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPatch, "/categories/some-random-id/upload")(c, "")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if category with provided id doesn't exist", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPatch, "/categories/69/upload")(c, "")
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodPatch, "/categories/1337/upload", true)
	})

	t.Run("PUT /categories/:id", func(t *testing.T) {
		sendWithParam := func(id uint, body string, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodPut, "/categories/"+param)(c, body)
		}

		t.Run("should update category in db based on provided json", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			testCategory := testutils.TestCategories[0]
			updateJSON := `{"title":"Sushi","removable":true}`

			_, c := testutils.LoginAsRandomAdmin(t)

			resp := sendWithParam(testCategory.ID, updateJSON, c)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(
				fmt.Sprintf(
					`{"id":%d,"title":"Sushi","removable":true,"image":%q}`,
					testCategory.ID,
					imgURL(*testCategory.Image),
				),
				resp.Body.String(),
			)
		})

		t.Run("should ignore id in provided json", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			testCategory := testutils.TestCategories[0]
			updateJSON := `{"id":420,"title":"Sushi","removable":true}`
			require.NotEqual(t, testCategory.ID, 420)

			_, c := testutils.LoginAsRandomAdmin(t)

			resp := sendWithParam(testCategory.ID, updateJSON, c)
			it.Equal(http.StatusOK, resp.Code)
			it.Equal(fmt.Sprintf(`{"id":%d,"title":"Sushi","removable":true,"image":%q}`, testCategory.ID, imgURL(*testCategory.Image)), resp.Body.String())
		})

		t.Run("should trim title", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			_, c := testutils.LoginAsRandomAdmin(t)

			tests := []string{"  Pizza", "Pizza   ", "   Pizza   "}
			for _, tc := range tests {
				json := fmt.Sprintf(`{"title":%q}`, tc)
				resp := sendWithParam(2, json, c)
				it.Equal(http.StatusConflict, resp.Code)
			}
		})

		t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPut, "/categories/some-random-id")(c, "")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if category with provided id doesn't exist", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPatch, "/categories/69")(c, "")
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodPut, "/categories/69", true)
	})

	t.Run("DELETE /categories/:id", func(t *testing.T) {
		sendWithParam := func(id uint, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodDelete, "/categories/"+param)(c, "")
		}

		t.Run("should removed a category with provided ID from db", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			testCategories := testutils.TestCategories
			testCat := testCategories[len(testCategories)/2]
			testCat.Image = nil
			require.NoError(t, db.Save(&testCat).Error)
			require.NoError(t, db.Exec("UPDATE dishes SET image = NULL").Error)

			_, c := testutils.LoginAsRandomAdmin(t)
			resp := sendWithParam(testCat.ID, c)

			if it.Equal(http.StatusOK, resp.Code) {
				var categories []category.Category
				db.Find(&categories)
				it.Len(categories, len(testCategories)-1)

				for _, c := range categories {
					it.NotEqual(c.ID, testCat.ID)
				}
			}

		})

		t.Run("should delete image if category had one", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			testutils.SetupUsersDB(t)
			t.Cleanup(testutils.CleanupStaticFolder)
			it := assert.New(t)
			cat := testutils.FindTestCategoryByID(3)
			cat.Image = nil
			require.NoError(t, db.Save(&cat).Error)
			require.NoError(t, db.Exec("UPDATE dishes SET image = NULL").Error)

			img, err := os.Open(testutils.PathToFile("./img/pizza.png"))
			require.NoError(t, err)
			defer img.Close()

			_, cookie := testutils.LoginAsRandomAdmin(t)

			resp := upload(cat.ID, cookie, filepath.Base(img.Name()), img)

			if it.Equal(http.StatusOK, resp.Code) {
				name := path.Base(resp.Body.String())

				resp = sendWithParam(cat.ID, cookie)
				if it.Equal(http.StatusOK, resp.Code) {
					it.NoFileExists(filepath.Join(config.CategoriesImgDirAbs, name))
				}
			}
		})

		t.Run("should delete images for all associated dishes", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			testutils.SetupUsersDB(t)
			t.Cleanup(testutils.CleanupStaticFolder)
			it := assert.New(t)
			cat := testutils.FindTestCategoryByID(3)
			cat.Image = nil
			require.NoError(t, db.Save(&cat).Error)

			img, err := os.Open(testutils.PathToFile("./img/pizza.png"))
			require.NoError(t, err)
			defer img.Close()
			imgName := filepath.Base(img.Name())

			_, cookie := testutils.LoginAsRandomAdmin(t)
			resp := upload(cat.ID, cookie, imgName, img)

			if it.Equal(http.StatusOK, resp.Code) {
				img.Seek(0, 0)
				catImageName := path.Base(resp.Body.String())

				var dishesImages []string
				dishes := testutils.FindTestDishesByCategoryID(cat.ID)

				for _, d := range dishes {
					d.Image = nil
					if it.NoError(db.Save(&d).Error) {
						resp = uploadDishImage(d.ID, cookie, imgName, img)
						img.Seek(0, 0)
						if it.Equal(http.StatusOK, resp.Code) {
							dishesImages = append(dishesImages, path.Base(resp.Body.String()))
						}
					}
				}

				resp = sendWithParam(cat.ID, cookie)
				if it.Equal(http.StatusOK, resp.Code) {
					it.NoFileExists(filepath.Join(config.CategoriesImgDirAbs, catImageName))

					for _, dishImage := range dishesImages {
						it.NoFileExists(filepath.Join(config.DishesImgDirAbs, dishImage))
					}
				}
			}
		})

		t.Run("should return 403 if category isn't removable", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			c := category.Category{ID: 69, Title: "Seafood", Removable: false}
			_, cookie := testutils.LoginAsRandomAdmin(t)

			if it.NoError(db.Create(&c).Error) {
				resp := sendWithParam(c.ID, cookie)
				it.Equal(http.StatusForbidden, resp.Code)
			}
		})

		t.Run("should return 403 if the dish from corresponding category has already been used in some order", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			it := assert.New(t)
			require.NoError(t, db.Exec("UPDATE categories SET image = NULL").Error)
			require.NoError(t, db.Exec("UPDATE dishes SET image = NULL").Error)

			cat := testutils.FindTestCategoryByID(1)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := sendWithParam(cat.ID, c)

			if it.Equal(http.StatusForbidden, resp.Code) {
				it.NotEmpty(resp.Body.String())
			}
		})

		t.Run("should return 400 if provided id isn't valid", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodDelete, "/categories/some-random-id")(c, "")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if category with provided id doesn't exist", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodDelete, "/categories/69")(c, "")
			assert.Equal(t, http.StatusNotFound, resp.Code)
		})

		testutils.RunAuthTests(t, http.MethodDelete, "/categories/69", true)
	})
}

func upload(id uint, c *http.Cookie, fileName string, file io.Reader) *httptest.ResponseRecorder {
	param := strconv.Itoa(int(id))
	return testutils.UploadReqWithCookie(http.MethodPatch, "/categories/"+param+"/upload", "image")(c, fileName, file)
}

func uploadDishImage(id uint, c *http.Cookie, fileName string, file io.Reader) *httptest.ResponseRecorder {
	param := strconv.Itoa(int(id))
	return testutils.UploadReqWithCookie(http.MethodPatch, "/dishes/"+param+"/upload", "image")(c, fileName, file)
}

func imgURL(name string) string {
	return category.PathToImg(name)
}

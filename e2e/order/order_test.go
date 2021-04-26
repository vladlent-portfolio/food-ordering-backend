package order_test

import (
	"encoding/json"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/order"
	"food_ordering_backend/controllers/user"
	"food_ordering_backend/database"
	"food_ordering_backend/e2e/testutils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var db = database.MustGetTest()
var orderRepo = order.ProvideRepository(db)

func TestOrders(t *testing.T) {
	t.Run("GET /orders", func(t *testing.T) {
		send := testutils.ReqWithCookie(http.MethodGet, "/orders")

		t.Run("should return a list of orders for specific user", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			c := testutils.LoginAs(t, testutils.TestUsersDTOs[0])
			verifyResponse(t, 3, send(c, ""))
		})

		t.Run("should return all orders if requester is admin", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			verifyResponse(t, len(testutils.TestOrders), send(c, ""))
		})

		testutils.RunAuthTests(t, http.MethodGet, "/orders", false)
	})

	t.Run("POST /orders", func(t *testing.T) {
		send := testutils.ReqWithCookie(http.MethodPost, "/orders")

		t.Run("should create an order and return it", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			testutils.SetupDishesAndCategories(t)
			it := assert.New(t)
			reqJSON := `{"items": [{"id":  1, "quantity": 2}, {"id":  3, "quantity": 1}]}`
			userDTO, c := testutils.LoginAsRandomUser(t)

			resp := send(c, reqJSON)

			if it.Equal(http.StatusCreated, resp.Code) {
				var dto order.ResponseDTO
				if it.NoError(json.NewDecoder(resp.Body).Decode(&dto)) {
					it.NotZero(dto.ID)
					it.NotZero(dto.CreatedAt)
					it.Equal(dto.CreatedAt, dto.UpdatedAt)
					it.Equal(order.StatusCreated, dto.Status)
					it.Equal(userDTO.Email, dto.User.Email)
					it.Equal(7.29, dto.Total)

					if it.Len(dto.Items, 2) {
						item1, item2 := dto.Items[0], dto.Items[1]
						it.NotZero(item1.ID)
						it.NotZero(item2.ID)

						it.NotZero(item1.OrderID)
						it.Equal(item1.OrderID, item2.OrderID)

						it.Equal(uint(1), item1.DishID)
						it.Equal(uint(3), item2.DishID)

						it.Equal(dish.ToDTO(testutils.FindTestDishByID(1)), item1.Dish)
						it.Equal(dish.ToDTO(testutils.FindTestDishByID(3)), item2.Dish)
					}
				}
			}
		})

		t.Run("should return 422 if json is incorrect or contains validation errors", func(t *testing.T) {
			testutils.SetupUsersDB(t)
			it := assert.New(t)
			tests := []string{
				`{"items":[{"id":  1, "quantity": 2}, {"id":  3, "quantity": `,
				`{"items":[{"id":  1, "quantity": 2}, {"id":  3, "quantity": -1}]}`,
				`{"items":[{"id":  1, "quantity": 2}, {"id":  3, "quantity": 0}]}`,
				`{"items":[]}`,
			}
			_, c := testutils.LoginAsRandomUser(t)

			for _, req := range tests {
				resp := send(c, req)
				it.Equal(http.StatusUnprocessableEntity, resp.Code)
			}
		})
		t.Run("should return 400 if dish with provided id doesn't exist", func(t *testing.T) {
			testutils.SetupDishesAndCategories(t)
			testutils.SetupUsersDB(t)
			it := assert.New(t)
			req := `{"items":[{"id":  1, "quantity": 2}, {"id": 233, "quantity": 1}]}`
			_, c := testutils.LoginAsRandomUser(t)

			resp := send(c, req)
			it.Equal(http.StatusBadRequest, resp.Code)
			it.Contains("Dish with id 233 doesn't exist", resp.Body.String())
		})

		testutils.RunAuthTests(t, http.MethodPost, "/orders", false)
	})

	t.Run("PATCH /orders/:id", func(t *testing.T) {
		sendWithParam := func(id uint, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodPatch, "/orders/"+param+"/cancel")(c, "")
		}

		t.Run("/orders/:id/cancel", func(t *testing.T) {
			t.Run("should change orders status to 'canceled' and return it", func(t *testing.T) {
				testutils.SetupOrdersDB(t)
				it := assert.New(t)
				tests := []struct {
					userDTO user.AuthDTO
					orderID uint
				}{
					{testutils.TestUsersDTOs[0], 2},
					{testutils.TestUsersDTOs[2], 4},
					{testutils.TestUsersDTOs[0], 5},
				}

				for _, tc := range tests {
					c := testutils.LoginAs(t, tc.userDTO)
					resp := sendWithParam(tc.orderID, c)

					if it.Equal(http.StatusOK, resp.Code) {
						o, err := orderRepo.FindByID(tc.orderID)

						if it.NoError(err) {
							it.Equal(order.StatusCanceled, o.Status)

							// Checking that only order status was changed
							unmodified := testutils.FindTestOrderByID(tc.orderID)
							it.NotEqual(unmodified.UpdatedAt, o.UpdatedAt)
							it.Equal(unmodified.CreatedAt, o.CreatedAt)
							it.Equal(unmodified.User, o.User)
							it.Equal(unmodified.UserID, o.UserID)
							it.Equal(unmodified.Total, o.Total)
							it.Equal(unmodified.Items, o.Items)
						}
					}
				}

			})

			t.Run("should return 304 if order is already canceled or done", func(t *testing.T) {
				it := assert.New(t)
				tests := []struct {
					userDTO user.AuthDTO
					orderID uint
				}{
					{testutils.TestUsersDTOs[0], 1},
					{testutils.TestUsersDTOs[3], 2},
				}

				for _, tc := range tests {
					c := testutils.LoginAs(t, tc.userDTO)
					resp := sendWithParam(tc.orderID, c)

					if it.Equal(http.StatusNotModified, resp.Code) {
						o, err := orderRepo.FindByID(tc.orderID)
						if it.NoError(err) {
							unmodified := testutils.FindTestOrderByID(tc.orderID)
							it.Equal(unmodified.UpdatedAt, o.UpdatedAt)
							it.Equal(unmodified.Status, o.Status)
							it.Equal(unmodified.CreatedAt, o.CreatedAt)
							it.Equal(unmodified.User, o.User)
							it.Equal(unmodified.UserID, o.UserID)
							it.Equal(unmodified.Total, o.Total)
							it.Equal(unmodified.Items, o.Items)
						}
					}
				}
			})

			testutils.RunAuthTests(t, http.MethodPut, "/orders/69/cancel", false)
		})
	})

	t.Run("PUT /orders/:id", func(t *testing.T) {
		//sendWithParam := func(id uint, body string, c *http.Cookie) *httptest.ResponseRecorder {
		//	param := strconv.Itoa(int(id))
		//	return testutils.ReqWithCookie(http.MethodPut, "/orders/"+param)(c, body)
		//}

		t.Run("should change order and return modified version", func(t *testing.T) {
			//it := assert.New(t)

		})

		testutils.RunAuthTests(t, http.MethodPut, "/orders/69", true)
	})
}

func verifyResponse(t *testing.T, expectedLen int, resp *httptest.ResponseRecorder) {
	it := assert.New(t)
	if it.Equal(http.StatusOK, resp.Code) {

		var orders []order.ResponseDTO
		if it.NoError(json.NewDecoder(resp.Body).Decode(&orders)) {
			it.Len(orders, expectedLen)

			for _, o := range orders {
				it.NotZero(o.CreatedAt)
				it.NotZero(o.UpdatedAt)

				if it.NotZero(o.ID) {
					testOrder := testutils.FindTestOrderByID(o.ID)
					it.Equal(testOrder.UserID, o.UserID)
					it.Equal(order.ToItemsResponseDTO(testOrder.Items), o.Items)
					it.Equal(testOrder.Status, o.Status)
					it.Equal(testOrder.Total, o.Total)
					it.Equal(user.ToResponseDTO(testOrder.User), o.User)
				}
			}
		}
	}

}

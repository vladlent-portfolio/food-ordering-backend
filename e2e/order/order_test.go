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
	"time"
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
			it.Contains(resp.Body.String(), "Dish with id 233 doesn't exist")
		})

		testutils.RunAuthTests(t, http.MethodPost, "/orders", false)
	})

	t.Run("PATCH /orders/:id", func(t *testing.T) {
		sendWithParam := func(id uint, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodPatch, "/orders/"+param+"/cancel")(c, "")
		}

		t.Run("/orders/:id/cancel", func(t *testing.T) {
			t.Run("should change orders status to 'canceled'", func(t *testing.T) {
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
						verifyStatusChange(t, tc.orderID, order.StatusCanceled)
					}
				}
			})

			t.Run("should return 200 if admin is changing status", func(t *testing.T) {
				testutils.SetupOrdersDB(t)
				it := assert.New(t)
				_, c := testutils.LoginAsRandomAdmin(t)
				resp := sendWithParam(4, c)

				if it.Equal(http.StatusOK, resp.Code) {
					verifyStatusChange(t, 4, order.StatusCanceled)
				}
			})

			t.Run("should return 304 if order is already canceled or done", func(t *testing.T) {
				testutils.SetupOrdersDB(t)
				it := assert.New(t)
				tests := []struct {
					userDTO user.AuthDTO
					orderID uint
					status  order.Status
				}{
					{testutils.TestUsersDTOs[0], 1, order.StatusDone},
					{testutils.TestUsersDTOs[1], 3, order.StatusCanceled},
				}

				for _, tc := range tests {
					c := testutils.LoginAs(t, tc.userDTO)
					resp := sendWithParam(tc.orderID, c)

					if it.Equal(http.StatusNotModified, resp.Code) {
						verifyStatusNotChange(t, tc.orderID, tc.status)
					}
				}
			})

			t.Run("should return 400 if order id is invalid", func(t *testing.T) {
				testutils.SetupOrdersDB(t)
				_, c := testutils.LoginAsRandomUser(t)
				resp := testutils.ReqWithCookie(http.MethodPatch, "/orders/randomtext/cancel")(c, "")
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			})

			t.Run("should return 403 if user tries to change another users order", func(t *testing.T) {
				testutils.SetupOrdersDB(t)
				c := testutils.LoginAs(t, testutils.TestUsersDTOs[2])
				resp := sendWithParam(5, c)
				assert.Equal(t, http.StatusForbidden, resp.Code)
				verifyStatusNotChange(t, 5, order.StatusCreated)
			})

			t.Run("should return 404 if order with provided id doesn't exist", func(t *testing.T) {
				testutils.SetupOrdersDB(t)
				it := assert.New(t)
				_, c := testutils.LoginAsRandomUser(t)
				resp := sendWithParam(1337, c)

				if it.Equal(http.StatusNotFound, resp.Code) {
					it.Contains(resp.Body.String(), "Order with id 1337 doesn't exist")
				}
			})

			testutils.RunAuthTests(t, http.MethodPatch, "/orders/69/cancel", false)
		})
	})

	t.Run("PUT /orders/:id", func(t *testing.T) {
		sendWithParam := func(id uint, body string, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodPut, "/orders/"+param)(c, body)
		}

		t.Run("should change order and return modified version", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			it := assert.New(t)
			initialOrder := testutils.TestOrders[3]
			_, c := testutils.LoginAsRandomAdmin(t)
			dto := order.UpdateDTO{
				Status: order.StatusInProgress,
				UserID: 1,
				Total:  1337,
				Items:  []order.ItemCreateDTO{{ID: 5, Quantity: 4}, {ID: 1, Quantity: 20}},
			}
			body, _ := json.Marshal(&dto)

			resp := sendWithParam(initialOrder.ID, string(body), c)

			if it.Equal(http.StatusOK, resp.Code) {
				var respDTO order.ResponseDTO
				if it.NoError(json.NewDecoder(resp.Body).Decode(&respDTO)) {
					it.Equal(initialOrder.ID, respDTO.ID)
					it.Equal(initialOrder.CreatedAt, respDTO.CreatedAt)
					it.NotEqual(respDTO.UpdatedAt, initialOrder.CreatedAt)
					it.Equal(dto.Status, respDTO.Status)
					it.Equal(dto.UserID, respDTO.UserID)
					it.Equal(dto.Total, respDTO.Total)

					if it.Len(respDTO.Items, 2) {
						it.Equal(dto.Items[0].ID, respDTO.Items[0].ID)
						it.Equal(dto.Items[0].Quantity, respDTO.Items[0].Quantity)
						it.Equal(dto.Items[1].ID, respDTO.Items[1].ID)
						it.Equal(dto.Items[1].Quantity, respDTO.Items[1].Quantity)
					}

				}
				orderInDB, err := orderRepo.FindByID(initialOrder.ID)
				if it.NoError(err) {
					it.Equal(order.ToResponseDTO(orderInDB), respDTO)
				}

				var initialOrderItems []order.Item
				itemsIDs := []uint{initialOrder.Items[0].ID, initialOrder.Items[1].ID}
				if it.NoError(db.Find(&initialOrderItems, itemsIDs).Error) {
					it.Len(initialOrderItems, 0, "expected previous order items to be deleted")
				}
			}
		})

		t.Run("should return 422 if json in request is malformed or invalid", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			orderID := uint(2)
			tests := []struct{ reqJSON, reason string }{
				{`{"status": 2, "user_id": 1, "total": 23, "items": [{"id"}]}`, "malformed"},
				{`{"status": -1, "user_id": 1, "total": 12.88, "items": [{"id": 2, "quantity": 4}]}`, "status < 0"},
				{`{"status": 4, "user_id": 1, "total": 23, "items": [{"id": 2, "quantity": 4}]}`, "status > 4"},
				{`{"user_id": 1, "total": 23, "items": [{"id": 2, "quantity": 4}]}`, "no status field"},
				{`{"status": 4, "total": 23, "items": [{"id": 2, "quantity": 4}]}`, "no user id"},
				{`{"status": 4, "user_id": 1, "total": 12.88, "items": [{"id": 2, "quantity": 4}]}`, "total < 0"},
				{`{"status": 4, "user_id": 1, "total": 12.88, "items": []}`, "empty items array"},
				{`{"status": 4, "user_id": 1, "total": 12.88}`, "no items field"},
				{`{"status": 4, "user_id": 1, "total": 12.88, "items": [{"quantity": 4}]}`, "no id field in items"},
				{`{"status": 4, "user_id": 1, "total": 12.88, "items": [{"id": 2, "quantity": 0}]}`, "quantity is 0"},
				{`{"status": 4, "user_id": 1, "total": 12.88, "items": [{"id": 2, "quantity": -1}]}`, "quantity < 0"},
			}
			_, c := testutils.LoginAsRandomAdmin(t)

			for _, tc := range tests {
				resp := sendWithParam(orderID, tc.reqJSON, c)
				assert.Equalf(t, http.StatusUnprocessableEntity, resp.Code, "expected request with %q to return 422", tc.reason)
			}
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
					it.Equal(testOrder.User.ID, o.User.ID)
					it.Equal(testOrder.User.Email, o.User.Email)
					it.Equal(testOrder.User.IsAdmin, o.User.IsAdmin)
					// There can be slight difference between cached user and user from db
					// so we compare string representation instead
					it.Equal(
						testOrder.User.CreatedAt.Format(time.RFC3339),
						o.User.CreatedAt.Format(time.RFC3339),
					)
				}
			}
		}
	}
}

func verifyStatusChange(t *testing.T, orderID uint, newStatus order.Status) {
	it := assert.New(t)
	o, err := orderRepo.FindByID(orderID)

	if it.NoError(err) {
		it.Equal(newStatus, o.Status)

		unmodified := testutils.FindTestOrderByID(orderID)
		it.NotEqual(unmodified.UpdatedAt, o.UpdatedAt)
		it.NotEqual(o.CreatedAt, o.UpdatedAt)
		testutils.UsersEqual(t, unmodified.User, o.User)
		it.Equal(unmodified.UserID, o.UserID)
		it.Equal(unmodified.Total, o.Total)
		it.Equal(unmodified.Items, o.Items)
	}
}
func verifyStatusNotChange(t *testing.T, orderID uint, status order.Status) {
	it := assert.New(t)
	o, err := orderRepo.FindByID(orderID)

	if it.NoError(err) {
		it.Equal(status, o.Status)

		unmodified := testutils.FindTestOrderByID(orderID)
		// There can be slight difference between cached user and user from db
		// so we compare string representation instead
		it.Equal(unmodified.CreatedAt.Format(time.RFC3339), o.CreatedAt.Format(time.RFC3339))
		it.Equal(unmodified.UpdatedAt.Format(time.RFC3339), o.UpdatedAt.Format(time.RFC3339))
		testutils.UsersEqual(t, unmodified.User, o.User)
		it.Equal(unmodified.UserID, o.UserID)
		it.Equal(unmodified.Total, o.Total)
		it.Equal(unmodified.Items, o.Items)
	}
}

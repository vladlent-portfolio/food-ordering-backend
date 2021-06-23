package order_test

import (
	"encoding/json"
	"fmt"
	"food_ordering_backend/controllers/dish"
	"food_ordering_backend/controllers/order"
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

		t.Run("should work with pagination", func(t *testing.T) {
			it := assert.New(t)
			testutils.SetupOrdersDB(t)

			orders := make([]order.Order, len(testutils.TestOrders))
			copy(orders, testutils.TestOrders)
			testutils.SortOrdersByID(orders)

			tests := []struct {
				query          string
				expectedOrders []order.Order
				expectedPage   int
				expectedLimit  int
			}{
				{"limit=2", orders[:2], 0, 2},
				{"limit=3&page=2", orders[3:], 2, 3},
				{"page=2", nil, 2, 10},
				{"", orders, 0, 10},
			}
			_, c := testutils.LoginAsRandomAdmin(t)

			for _, tc := range tests {
				resp := testutils.ReqWithCookie(http.MethodGet, "/orders?"+tc.query)(c, "")

				if it.Equal(http.StatusOK, resp.Code) {
					var dto order.DTOsWithPagination

					if it.NoError(json.NewDecoder(resp.Body).Decode(&dto)) {
						it.Equal(tc.expectedPage, dto.Pagination.Page)
						it.Equal(tc.expectedLimit, dto.Pagination.Limit)
						it.Equal(len(orders), dto.Pagination.Total)

						for i, o := range dto.Orders {
							isEqualOrder(t, tc.expectedOrders[i], o)
						}
					}
				}
			}
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
		sendWithParam := func(c *http.Cookie, id uint, status order.Status) *httptest.ResponseRecorder {
			uri := fmt.Sprintf("/orders/%d?status=%d", id, status)
			return testutils.ReqWithCookie(http.MethodPatch, uri)(c, "")
		}

		t.Run("should change orders status", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			it := assert.New(t)
			id := testutils.TestOrders[1].ID
			_, c := testutils.LoginAsRandomAdmin(t)

			for _, status := range order.Statuses {
				resp := sendWithParam(c, id, status)

				if it.Equal(http.StatusNoContent, resp.Code) {
					verifyStatusChange(t, id, status)
				}
			}
		})

		t.Run("should return 304 if order hasn't been changed", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			it := assert.New(t)
			_, c := testutils.LoginAsRandomAdmin(t)

			for _, o := range testutils.TestOrders {
				resp := sendWithParam(c, o.ID, o.Status)

				if it.Equal(http.StatusNotModified, resp.Code) {
					verifyStatusNotChange(t, o.ID, o.Status)
				}
			}
		})

		t.Run("should return 400 if order id is invalid", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := testutils.ReqWithCookie(http.MethodPatch, "/orders/randomtext?status=1")(c, "")
			assert.Equal(t, http.StatusBadRequest, resp.Code)
		})

		t.Run("should return 404 if order with provided id doesn't exist", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			it := assert.New(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			resp := sendWithParam(c, 1337, order.StatusInProgress)

			if it.Equal(http.StatusNotFound, resp.Code) {
				it.Contains(resp.Body.String(), "Order with id 1337 doesn't exist")
			}

		})

		t.Run("should return 422 is status is invalid", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			it := assert.New(t)
			tests := []string{"-23", "-1", "4", "10", "45", "done", "create"}
			_, c := testutils.LoginAsRandomAdmin(t)

			for _, test := range tests {
				resp := testutils.ReqWithCookie(http.MethodPatch, "/orders/randomtext?status="+test)(c, "")
				it.Equal(http.StatusUnprocessableEntity, resp.Code)
			}
		})

		testutils.RunAuthTests(t, http.MethodPatch, "/orders/69?status=1", true)
	})

	t.Run("PUT /orders/:id", func(t *testing.T) {
		sendWithParam := func(id uint, body string, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodPut, "/orders/"+param)(c, body)
		}

		t.Run("should change order and return modified version", func(t *testing.T) {
			testutils.SetupOrdersDB(t)
			it := assert.New(t)
			initialOrder := testutils.FindTestOrderByID(4)
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
					it.True(testutils.EqualTimestamps(initialOrder.CreatedAt, respDTO.CreatedAt))
					it.NotEqual(initialOrder.UpdatedAt, respDTO.UpdatedAt)
					it.Equal(dto.Status, respDTO.Status)
					it.Equal(dto.UserID, respDTO.UserID)
					it.Equal(dto.Total, respDTO.Total)

					if it.Len(respDTO.Items, 2) {
						it.NotEqual(dto.Items[0].ID, respDTO.Items[0].ID)
						it.Equal(dto.Items[0].Quantity, respDTO.Items[0].Quantity)
						it.NotEqual(dto.Items[1].ID, respDTO.Items[1].ID)
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

		t.Run("should return 404 if order with provided id doesn't exist", func(t *testing.T) {
			it := assert.New(t)
			testutils.SetupOrdersDB(t)
			_, c := testutils.LoginAsRandomAdmin(t)
			ids := []uint{1337, 420}

			for _, id := range ids {
				resp := sendWithParam(id, "", c)
				if it.Equal(http.StatusNotFound, resp.Code) {
					it.Contains(resp.Body.String(), fmt.Sprintf("Order with id %d doesn't exist", id))
				}
			}
		})

		testutils.RunAuthTests(t, http.MethodPut, "/orders/69", true)
	})
}

func verifyResponse(t *testing.T, expectedLen int, resp *httptest.ResponseRecorder) {
	it := assert.New(t)
	if it.Equal(http.StatusOK, resp.Code) {

		var dto order.DTOsWithPagination
		if it.NoError(json.NewDecoder(resp.Body).Decode(&dto)) {
			it.Len(dto.Orders, expectedLen)

			ids := make([]uint, len(dto.Orders))

			for i, o := range dto.Orders {
				ids[i] = o.ID
			}

			it.IsIncreasing(ids, "expected orders to be sorted by id")

			for _, o := range dto.Orders {
				it.NotZero(o.CreatedAt)
				it.NotZero(o.UpdatedAt)

				if it.NotZero(o.ID) {
					testOrder := testutils.FindTestOrderByID(o.ID)
					isEqualOrder(t, testOrder, o)
				}
			}
		}
	}
}

func isEqualOrder(t *testing.T, o order.Order, dto order.ResponseDTO) {
	it := assert.New(t)
	it.Equal(o.UserID, dto.UserID)
	it.Equal(order.ToItemsResponseDTO(o.Items), dto.Items)
	it.Equal(o.Status, dto.Status)
	it.Equal(o.Total, dto.Total)
	it.Equal(o.User.ID, dto.User.ID)
	it.Equal(o.User.Email, dto.User.Email)
	it.Equal(o.User.IsAdmin, dto.User.IsAdmin)
	it.True(testutils.EqualTimestamps(o.User.CreatedAt, dto.User.CreatedAt))
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
		it.True(testutils.EqualTimestamps(unmodified.CreatedAt, o.CreatedAt))
		it.True(testutils.EqualTimestamps(unmodified.UpdatedAt, o.UpdatedAt))
		testutils.UsersEqual(t, unmodified.User, o.User)
		it.Equal(unmodified.UserID, o.UserID)
		it.Equal(unmodified.Total, o.Total)
		it.Equal(unmodified.Items, o.Items)
	}
}

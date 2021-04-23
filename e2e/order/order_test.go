package order

import (
	"encoding/json"
	"food_ordering_backend/controllers/order"
	"food_ordering_backend/e2e/testutils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestOrders(t *testing.T) {
	t.Run("GET /orders", func(t *testing.T) {
		//send := testutils.ReqWithCookie(http.MethodGet, "/orders")

		t.Run("should return a list of orders for specific user", func(t *testing.T) {
			//it := assert.New(t)

		})

		t.Run("should return all orders if requester is admin", func(t *testing.T) {
			//it := assert.New(t)

		})

		testutils.RunAuthTests(t, http.MethodGet, "/orders", false)
	})

	t.Run("POST /orders", func(t *testing.T) {
		send := testutils.ReqWithCookie(http.MethodPost, "/orders")

		t.Run("should create an order and return it", func(t *testing.T) {
			it := assert.New(t)
			// TODO: Use TestDishes instead of hardcoded

			reqJSON := `{"items": [{"id":  1, "quantity": 2}, {"id":  3, "quantity": 1}]}`
			userDTO, c := testutils.LoginAsRandomUser(t)

			resp := send(c, reqJSON)

			if it.Equal(http.StatusCreated, resp.Code) {
				var dto order.ResponseDTO
				if it.NoError(json.NewDecoder(resp.Body).Decode(&dto)) {
					it.NotZero(dto.ID)
					it.NotZero(dto.CreatedAt)
					it.Equal(dto.CreatedAt, dto.UpdatedAt)
					// TODO: Add order status check

					it.Equal(userDTO.Email, dto.User.Email)
					it.Equal(7.29, dto.Total)
					// TODO: Check returned dishes

					it.Len(dto.Items, 2)
				}
			}
		})

		t.Run("should return 422 if json is incorrect or contains validation errors", func(t *testing.T) {
			it := assert.New(t)
			malformed := `{"items":[{"id":  1, "quantity": 2}, {"id":  3, "quantity": `
			invalid := `{"items":[{"id":  1, "quantity": 2}, {"id":  3, "quantity": -1}]}`
			zeroQuantity := `{"items":[{"id":  1, "quantity": 2}, {"id":  3, "quantity": 0}]}`

			_, c := testutils.LoginAsRandomUser(t)

			for _, req := range []string{malformed, invalid, zeroQuantity} {
				resp := send(c, req)
				it.Equal(http.StatusUnprocessableEntity, resp.Code)
			}
		})
		testutils.RunAuthTests(t, http.MethodPost, "/orders", false)
	})

	t.Run("PATCH /orders/:id", func(t *testing.T) {

		t.Run("/orders/:id/cancel", func(t *testing.T) {
			//sendWithParam := func(id uint, body string, c *http.Cookie) *httptest.ResponseRecorder {
			//	param := strconv.Itoa(int(id))
			//	return testutils.ReqWithCookie(http.MethodPatch, "/orders/"+param+"/cancel")(c, body)
			//}

			t.Run("should change orders status to 'canceled'", func(t *testing.T) {
				//it := assert.New(t)

			})

			t.Run("should return 304 if order is already canceled", func(t *testing.T) {
				//it := assert.New(t)

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

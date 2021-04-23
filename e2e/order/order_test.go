package order

import (
	"food_ordering_backend/e2e/testutils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestOrders(t *testing.T) {
	t.Run("GET /orders", func(t *testing.T) {
		send := testutils.ReqWithCookie(http.MethodGet, "/orders")

		t.Run("should return a list of orders for specific user", func(t *testing.T) {
			it := assert.New(t)

		})

		t.Run("should return all orders if requester is admin", func(t *testing.T) {
			it := assert.New(t)

		})

		testutils.RunAuthTests(t, http.MethodGet, "/orders", false)
	})

	t.Run("POST /orders", func(t *testing.T) {
		send := testutils.ReqWithCookie(http.MethodPost, "/orders")

		t.Run("should create an order and return it", func(t *testing.T) {
			it := assert.New(t)

		})
		testutils.RunAuthTests(t, http.MethodPost, "/orders", false)
	})

	t.Run("PATCH /orders/:id", func(t *testing.T) {

		t.Run("/orders/:id/cancel", func(t *testing.T) {
			sendWithParam := func(id uint, body string, c *http.Cookie) *httptest.ResponseRecorder {
				param := strconv.Itoa(int(id))
				return testutils.ReqWithCookie(http.MethodPatch, "/orders/"+param+"/cancel")(c, body)
			}

			t.Run("should change orders status to 'canceled'", func(t *testing.T) {
				it := assert.New(t)

			})

			t.Run("should return 304 if order is already canceled", func(t *testing.T) {
				it := assert.New(t)

			})

			testutils.RunAuthTests(t, http.MethodPut, "/orders/69/cancel", false)
		})
	})

	t.Run("PUT /orders/:id", func(t *testing.T) {
		sendWithParam := func(id uint, body string, c *http.Cookie) *httptest.ResponseRecorder {
			param := strconv.Itoa(int(id))
			return testutils.ReqWithCookie(http.MethodPut, "/orders/"+param)(c, body)
		}

		t.Run("should change order and return modified version", func(t *testing.T) {
			it := assert.New(t)

		})

		testutils.RunAuthTests(t, http.MethodPut, "/orders/69", true)
	})
}

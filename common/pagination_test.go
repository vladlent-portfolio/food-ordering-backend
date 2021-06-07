package common

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestExtractPagination(t *testing.T) {
	t.Run("should parse query and return Pagination instance", func(t *testing.T) {
		it := assert.New(t)
		tests := []struct{ limit, page int }{
			{1, 3},
			{34, 657},
			{23, 90},
		}

		for _, tc := range tests {
			q := newQuery(tc.page, tc.limit)
			p := ExtractPagination(q, 0)

			it.Equal(tc.page, p.Page())
			it.Equal(tc.limit, p.Limit())
		}

	})

	t.Run("should use default value for limit", func(t *testing.T) {
		it := assert.New(t)
		tests := []struct{ defaultLimit, offset, page int }{
			{3, 23, 45},
			{0, 1, 24},
			{99, 45, 32},
		}

		for _, tc := range tests {
			q := newQuery(tc.page, 0)
			p := ExtractPagination(q, tc.defaultLimit)

			it.Equal(tc.page, p.Page())
			it.Equal(tc.defaultLimit, p.Limit())
		}
	})

	t.Run("should use 0 as a default value for page", func(t *testing.T) {
		it := assert.New(t)
		tests := []struct{ offset, limit int }{
			{3, 23},
			{4576, 1},
			{99, 45},
		}

		for _, tc := range tests {
			q := newQuery(0, tc.limit)
			p := ExtractPagination(q, 123)

			it.Equal(0, p.Page())
			it.Equal(tc.limit, p.Limit())
		}

	})
}

type query map[string]string

func (q query) setValue(name string, value int) {
	if value != 0 {
		q[name] = strconv.Itoa(value)
	}
}

func (q query) limit() int {
	limit, _ := strconv.Atoi(q["limit"])
	return limit
}

func (q query) page() int {
	page, _ := strconv.Atoi(q["page"])
	return page
}

func (q query) Query(key string) string {
	return q[key]
}

func newQuery(page, limit int) query {
	q := query{}
	q.setValue("page", page)
	q.setValue("limit", limit)
	return q
}

package common

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestExtractPagination(t *testing.T) {
	t.Run("should parse query and return Pagination instance", func(t *testing.T) {
		it := assert.New(t)
		tests := []struct{ limit, offset, page int }{
			{1, 2, 3},
			{34, 45, 657},
			{23, 67, 90},
		}

		for _, tc := range tests {
			q := newQuery(tc.page, tc.limit, tc.offset)
			p := ExtractPagination(q, 0, 0)

			it.Equal(tc.page, p.Page())
			it.Equal(tc.limit, p.Limit())
			it.Equal(tc.offset, p.Offset())
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
			q := newQuery(tc.page, 0, tc.offset)
			p := ExtractPagination(q, tc.defaultLimit, 0)

			it.Equal(tc.page, p.Page())
			it.Equal(tc.defaultLimit, p.Limit())
			it.Equal(tc.offset, p.Offset())
		}
	})

	t.Run("should use default value for offset", func(t *testing.T) {
		it := assert.New(t)
		tests := []struct{ defaultOffset, limit, page int }{
			{3, 23, 45},
			{0, 1, 24},
			{99, 45, 32},
		}

		for _, tc := range tests {
			q := newQuery(tc.page, tc.limit, 0)
			p := ExtractPagination(q, 0, tc.defaultOffset)

			it.Equal(tc.page, p.Page())
			it.Equal(tc.limit, p.Limit())
			it.Equal(tc.defaultOffset, p.Offset())
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
			q := newQuery(0, tc.limit, tc.offset)
			p := ExtractPagination(q, 123, 346)

			it.Equal(0, p.Page())
			it.Equal(tc.limit, p.Limit())
			it.Equal(tc.offset, p.Offset())
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

func (q query) offset() int {
	offset, _ := strconv.Atoi(q["offset"])
	return offset
}

func (q query) page() int {
	page, _ := strconv.Atoi(q["page"])
	return page
}

func (q query) Query(key string) string {
	return q[key]
}

func newQuery(page, limit, offset int) query {
	q := query{}
	q.setValue("page", page)
	q.setValue("limit", limit)
	q.setValue("offset", offset)
	return q
}

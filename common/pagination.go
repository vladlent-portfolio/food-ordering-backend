package common

import "strconv"

// QueryParser extracts query parameter from request by it's respective name
// using Query method.
type QueryParser interface {
	Query(key string) string
}

type Pagination struct {
	page   int
	limit  int
	offset int
}

func (p *Pagination) Page() int {
	return p.page
}

func (p *Pagination) Limit() int {
	return p.limit
}

func (p *Pagination) Offset() int {
	return p.offset
}

func ExtractPagination(parser QueryParser, defaultLimit, defaultOffset int) *Pagination {
	return &Pagination{
		page:   parseDefault(parser, "page", 0),
		limit:  parseDefault(parser, "limit", defaultLimit),
		offset: parseDefault(parser, "offset", defaultOffset),
	}
}

func parseDefault(parser QueryParser, key string, defaultVal int) int {
	query := parser.Query(key)

	if query == "" {
		return defaultVal
	}

	if val, err := strconv.Atoi(query); err == nil {
		return val
	}

	return defaultVal
}

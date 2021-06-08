package common

import (
	"gorm.io/gorm"
	"strconv"
)

type Paginator interface {
	Page() int
	Limit() int
}

// QueryParser extracts query parameter from request by it's respective name
// using Query method.
type QueryParser interface {
	Query(key string) string
}

type PaginationDTO struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type Pagination struct {
	page  int
	limit int
}

func (p *Pagination) Page() int {
	return p.page
}

func (p *Pagination) Limit() int {
	return p.limit
}

func ExtractPagination(parser QueryParser, defaultLimit int) *Pagination {
	return &Pagination{
		page:  parseDefault(parser, "page", 0),
		limit: parseDefault(parser, "limit", defaultLimit),
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

func WithPagination(p Paginator) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Page() * p.Limit()).Limit(p.Limit())
	}
}

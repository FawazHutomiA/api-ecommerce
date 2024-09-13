package sqlx

import (
	"example/pkg/helper"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PaginationMetadata struct {
	DB *sqlx.DB
}

func NewPaginationMetadata(db *sqlx.DB) *PaginationMetadata {
	return &PaginationMetadata{DB: db}
}

func (p *PaginationMetadata) GetPagination(query string, param helper.PaginationParams, dest interface{}, args ...interface{}) (helper.Pagination, error) {
	page := param.Page
	if param.Page < 1 {
		page = 1
	}
	limit := param.Limit
	if param.Limit < 1 {
		limit = 10
	}
	// get total data
	totalQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as total", query)
	var total int
	err := p.DB.Get(&total, totalQuery, args...)
	if err != nil {
		return helper.Pagination{}, err
	}

	//search data
	// if param.Search != "" {
	// 	query = fmt.Sprintf("%s WHERE %s", query, param.Search)
	// }

	// order by and sort by
	// if param.OrderBy != "" {
	// 	query = fmt.Sprintf("%s ORDER BY %s %s", query, param.OrderBy, param.SortBy)
	// }

	// get data
	offset := (page - 1) * limit
	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
	err = p.DB.Select(dest, query, args...)
	if err != nil {
		return helper.Pagination{}, err
	}

	return helper.Pagination{
		CurrentPage:  page,
		PageSize:     limit,
		FirstPage:    1,
		LastPage:     total/limit + 1,
		TotalRecords: total,
		Records:      dest,
	}, nil
}

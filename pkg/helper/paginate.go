package helper

import (
	"net/http"
	"strconv"
)

type PaginationParams struct {
	Page    int
	Limit   int
	Search  string
	SortBy  string
	OrderBy string
}

func (p *PaginationParams) GetPaginateParam(r *http.Request) PaginationParams {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		p.Page = 1
	} else {
		p.Page = page
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		p.Limit = 10
	} else {
		p.Limit = limit
	}
	search := r.URL.Query().Get("search")
	if search != "" {
		p.Search = search
	}
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy != "" {
		p.SortBy = sortBy
	} else {
		sortBy = "created_at"
		p.SortBy = sortBy
	}
	orderBy := r.URL.Query().Get("order_by")
	if orderBy != "" {
		p.OrderBy = orderBy
	} else {
		orderBy = "desc"
		p.OrderBy = orderBy
	}

	return PaginationParams{
		Page:    page,
		Limit:   limit,
		Search:  search,
		SortBy:  sortBy,
		OrderBy: orderBy,
	}
}

type Pagination struct {
	CurrentPage  int         `json:"current_page"`
	PageSize     int         `json:"page_size"`
	FirstPage    int         `json:"first_page"`
	LastPage     int         `json:"last_page"`
	TotalRecords int         `json:"total_records"`
	Records      interface{} `json:"records"`
}

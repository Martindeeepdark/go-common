package sql

import (
	"fmt"
	"math"
)

// Pager represents a pagination result
type Pager struct {
	Page      int   `json:"page"`       // Current page number (1-based)
	PageSize  int   `json:"page_size"`  // Number of items per page
	Total     int64 `json:"total"`      // Total number of items
	TotalPage int   `json:"total_page"` // Total number of pages
}

// NewPager creates a new pager with default values
func NewPager(page, pageSize int) *Pager {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	return &Pager{
		Page:     page,
		PageSize: pageSize,
	}
}

// SetTotal sets the total number of items and calculates total pages
func (p *Pager) SetTotal(total int64) {
	p.Total = total
	if p.PageSize > 0 {
		p.TotalPage = int(math.Ceil(float64(total) / float64(p.PageSize)))
	} else {
		p.TotalPage = 0
	}
}

// Offset returns the OFFSET value for SQL query
func (p *Pager) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the LIMIT value for SQL query
func (p *Pager) Limit() int {
	return p.PageSize
}

// HasNext returns true if there is a next page
func (p *Pager) HasNext() bool {
	return p.Page < p.TotalPage
}

// HasPrev returns true if there is a previous page
func (p *Pager) HasPrev() bool {
	return p.Page > 1
}

// PageInfo returns pagination information for API response
func (p *Pager) PageInfo() map[string]interface{} {
	return map[string]interface{}{
		"page":       p.Page,
		"page_size":  p.PageSize,
		"total":      p.Total,
		"total_page": p.TotalPage,
	}
}

// String returns the string representation
func (p *Pager) String() string {
	return fmt.Sprintf("Page: %d/%d, PageSize: %d, Total: %d",
		p.Page, p.TotalPage, p.PageSize, p.Total)
}

// PaginateResult wraps a slice with pagination info
type PaginateResult struct {
	Items interface{} `json:"items"`
	Pager *Pager      `json:"pager"`
}

// NewPaginateResult creates a new paginated result
func NewPaginateResult(items interface{}, pager *Pager) *PaginateResult {
	return &PaginateResult{
		Items: items,
		Pager: pager,
	}
}

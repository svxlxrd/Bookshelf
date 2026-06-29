package domain

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewPagination(page, limit, total int) Pagination {
	pagination := Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	if limit > 0 {
		pagination.TotalPages = (total + limit - 1) / limit
	}

	return pagination
}
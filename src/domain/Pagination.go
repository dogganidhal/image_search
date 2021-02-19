package domain

type Pagination struct {
	From int64 `json:"from"`
	Size int64 `json:"size"`
}

type PaginatedResource struct {
	Total      int64         `json:"total"`
	Pagination Pagination    `json:"pagination"`
	Items      []interface{} `json:"items"`
}

func DefaultPaginationFrom() int64 {
	return 0
}

func DefaultPaginationSize() int64 {
	return 20
}

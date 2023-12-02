package model

type Pagination[T any] struct {
	Item       []T   `json:"item,omitempty"`
	TotalItems int64 `json:"total_items,omitempty"`
}

type TodoPagination Pagination[TodoItem]

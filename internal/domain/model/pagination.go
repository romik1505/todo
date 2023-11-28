package model

type Pagination[T any] struct {
	Item       []T
	TotalItems int64
}

package dto

import (
	"time"
)

const TodoTableName = "todos"

type TodoItem struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Date        *time.Time `db:"date"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	TotalItems  int64      `db:"total_items"`
}

type TodoFilter struct {
	Date   *time.Time `json:"date,omitempty" form:"date"`
	Status string     `json:"status,omitempty" form:"status"`
	Page   int64      `json:"page,omitempty" form:"page"`
	Limit  int64      `json:"limit,omitempty" form:"limit"`
}

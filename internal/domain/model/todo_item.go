package model

import "time"

type TodoItem struct {
	ID          int64      `json:"id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Date        time.Time  `json:"date,omitempty"`
	Status      TodoStatus `json:"status,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type TodoStatus string

var (
	TodoStatusCompleted = "completed"
	TodoStatusPending   = "pending"
)

func (t TodoItem) Bind() error {
	return nil
}

package model

import (
	"fmt"
	"time"
)

type TodoItem struct {
	ID          int64      `json:"id,omitempty"`
	Title       string     `json:"title,omitempty" form:"title"`
	Description string     `json:"description,omitempty" form:"description"`
	Date        *time.Time `json:"date,omitempty" form:"date" time_format:"2006-01-02"`
	Status      TodoStatus `json:"status,omitempty" form:"status"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type TodoStatus string

var (
	TodoStatusCompleted = "completed"
	TodoStatusPending   = "pending"
)

const (
	TodoTitleField       = "title"
	TodoDescriptionField = "description"
	TodoDateField        = "date"
	TodoStatusField      = "status"
)

var TodoFields = []string{
	TodoTitleField,
	TodoDescriptionField,
	TodoDateField,
	TodoStatusField,
}

func (t *TodoItem) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("title must be set")
	}
	if t.Date == nil {
		return fmt.Errorf("date must be set")
	}
	if t.Status == "" {
		return fmt.Errorf("status must be set")
	}
	return nil
}

func (t *TodoItem) EditableFields() []string {
	res := make([]string, 0)
	if t.Title != "" {
		res = append(res, TodoTitleField)
	}

	if t.Description != "" {
		res = append(res, TodoDescriptionField)
	}

	if t.Date != nil {
		res = append(res, TodoDateField)
	}

	if t.Status != "" {
		res = append(res, TodoStatusField)
	}

	return res
}

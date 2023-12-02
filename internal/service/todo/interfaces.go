package todo

import (
	"context"
	"errors"
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
)

type (
	Service interface {
		CreateTodo(ctx context.Context, item *model.TodoItem) error
		GetTodoByID(ctx context.Context, id int64) (model.TodoItem, error)
		UpdateTodo(ctx context.Context, item *model.TodoItem) error
		DeleteTodo(ctx context.Context, id int64) error
		ListTodos(ctx context.Context, filter dto.TodoFilter) (model.TodoPagination, error)
	}

	Repository interface {
		CreateTodo(ctx context.Context, item *dto.TodoItem) error
		GetTodoByID(ctx context.Context, id int64) (dto.TodoItem, error)
		UpdateTodo(ctx context.Context, item *dto.TodoItem, updatedFields []string) error
		DeleteTodo(ctx context.Context, id int64) error
		ListTodos(ctx context.Context, filter dto.TodoFilter) ([]dto.TodoItem, int64, error)
	}
)

var (
	ErrValidation   = errors.New("validation error")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal error")
	ErrEmptyContent = errors.New("empty content")
)

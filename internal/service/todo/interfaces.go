package todo

import (
	"context"
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
)

type (
	Service interface {
		CreateTodo(ctx context.Context, item *model.TodoItem) error
		GetTodoByID(ctx context.Context, id int64) (model.TodoItem, error)
		UpdateTodo(ctx context.Context, item *model.TodoItem) error
		DeleteTodo(ctx context.Context, id int64) error
		ListTodos(ctx context.Context, filter dto.TodoFilter) (model.Pagination[model.TodoItem], error)
	}

	Repository interface {
		CreateTodo(ctx context.Context, item *dto.TodoItem) error
		GetTodoByID(ctx context.Context, id int64) (dto.TodoItem, error)
		UpdateTodo(ctx context.Context, item *dto.TodoItem, updatedFields []string) error
		DeleteTodo(ctx context.Context, id int64) error
		ListTodos(ctx context.Context, filter dto.TodoFilter) ([]dto.TodoItem, int64, error)
	}
)

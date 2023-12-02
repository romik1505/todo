package todo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
	"todo-list/internal/util/converter"
)

type TodoService struct {
	TodoRepo Repository
}

func NewTodoService(tr Repository) *TodoService {
	return &TodoService{
		TodoRepo: tr,
	}
}

func (t *TodoService) CreateTodo(ctx context.Context, item *model.TodoItem) error {
	err := item.Validate()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	todoDto := converter.ConvertTodoToDTO(*item)
	err = t.TodoRepo.CreateTodo(ctx, &todoDto)
	if err != nil {
		return err
	}

	*item = converter.ConvertTodoToModel(todoDto)
	return nil
}

func (t *TodoService) GetTodoByID(ctx context.Context, id int64) (model.TodoItem, error) {
	if id <= 0 {
		return model.TodoItem{}, ErrValidation
	}

	td, err := t.TodoRepo.GetTodoByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.TodoItem{}, ErrNotFound
		}
		return model.TodoItem{}, err
	}

	return converter.ConvertTodoToModel(td), nil
}

func (t *TodoService) UpdateTodo(ctx context.Context, item *model.TodoItem) error {
	fields := item.EditableFields()

	todoDto := converter.ConvertTodoToDTO(*item)
	if err := t.TodoRepo.UpdateTodo(ctx, &todoDto, fields); err != nil {
		return err
	}

	*item = converter.ConvertTodoToModel(todoDto)
	return nil
}

func (t *TodoService) DeleteTodo(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrValidation
	}

	err := t.TodoRepo.DeleteTodo(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (t *TodoService) ListTodos(ctx context.Context, filter dto.TodoFilter) (model.TodoPagination, error) {
	items, totalItems, err := t.TodoRepo.ListTodos(ctx, filter)
	if err != nil {
		return model.TodoPagination{}, err
	}

	if totalItems == 0 {
		return model.TodoPagination{}, ErrNotFound
	}

	return model.TodoPagination{
		Item:       converter.ConvertTodoToModels(items),
		TotalItems: totalItems,
	}, nil
}

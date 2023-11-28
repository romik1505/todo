package todo

import (
	"context"
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
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
	//TODO implement me
	panic("implement me")
}

func (t *TodoService) GetTodoByID(ctx context.Context, id int64) (model.TodoItem, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TodoService) UpdateTodo(ctx context.Context, item *model.TodoItem) error {
	//TODO implement me
	panic("implement me")
}

func (t *TodoService) DeleteTodo(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (t *TodoService) ListTodos(ctx context.Context, filter dto.TodoFilter) (model.Pagination[model.TodoItem], error) {
	//TODO implement me
	panic("implement me")
}

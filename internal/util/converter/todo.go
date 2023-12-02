package converter

import (
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
)

func ConvertTodoToDTO(inp model.TodoItem) dto.TodoItem {
	return dto.TodoItem{
		ID:          inp.ID,
		Title:       inp.Title,
		Description: inp.Description,
		Date:        inp.Date,
		Status:      string(inp.Status),
	}
}

func ConvertTodoToModel(inp dto.TodoItem) model.TodoItem {
	return model.TodoItem{
		ID:          inp.ID,
		Title:       inp.Title,
		Description: inp.Description,
		Date:        inp.Date,
		Status:      model.TodoStatus(inp.Status),
		CreatedAt:   inp.CreatedAt,
		UpdatedAt:   inp.UpdatedAt,
	}
}

func ConvertTodoToModels(inp []dto.TodoItem) []model.TodoItem {
	res := make([]model.TodoItem, len(inp))

	for i, v := range inp {
		res[i] = ConvertTodoToModel(v)
	}

	return res
}

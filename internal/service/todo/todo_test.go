package todo

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
	"todo-list/internal/util/pointer"
	mock_todo "todo-list/pkg/mocks/service/todo"
)

func TestTodoService_CreateTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_todo.NewMockRepository(ctrl)
	s := NewTodoService(repo)
	now := time.Now()

	t.Run("casual creation todo", func(t *testing.T) {
		input := &model.TodoItem{
			Title:       "Полить цветы",
			Description: "Взять лейку. Наполнить водой. Полить цветы.",
			Date:        pointer.Pointer(time.Date(2023, time.February, 11, 0, 0, 0, 0, time.UTC)),
			Status:      "pending",
		}
		expectDto := &dto.TodoItem{
			Title:       "Полить цветы",
			Description: "Взять лейку. Наполнить водой. Полить цветы.",
			Date:        pointer.Pointer(time.Date(2023, time.February, 11, 0, 0, 0, 0, time.UTC)),
			Status:      "pending",
		}
		repo.EXPECT().CreateTodo(gomock.Any(), expectDto).DoAndReturn(func(ctx context.Context, inp *dto.TodoItem) {
			inp.CreatedAt = now
			inp.ID = 12

		}).Return(nil)

		err := s.CreateTodo(context.Background(), input)
		require.NoError(t, err)
		require.Equal(t, &model.TodoItem{
			ID:          12,
			Title:       "Полить цветы",
			Description: "Взять лейку. Наполнить водой. Полить цветы.",
			Date:        pointer.Pointer(time.Date(2023, time.February, 11, 0, 0, 0, 0, time.UTC)),
			CreatedAt:   now,
			Status:      "pending",
		}, input)
	})

	t.Run("validation error", func(t *testing.T) {
		input := &model.TodoItem{
			Title:       "",
			Description: "Взять лейку. Наполнить водой. Полить цветы.",
			Date:        pointer.Pointer(time.Date(2023, time.February, 11, 0, 0, 0, 0, time.UTC)),
			Status:      "pending",
		}
		err := s.CreateTodo(context.Background(), input)
		require.ErrorIs(t, err, ErrValidation)
	})

	t.Run("database error", func(t *testing.T) {
		repo.EXPECT().CreateTodo(gomock.Any(), gomock.Any()).Return(sql.ErrConnDone)
		err := s.CreateTodo(context.Background(), &model.TodoItem{
			Title:       "title",
			Description: "desc",
			Date:        pointer.Pointer(time.Date(2023, time.February, 11, 0, 0, 0, 0, time.UTC)),
			Status:      "pending",
		})
		require.ErrorIs(t, sql.ErrConnDone, err)
	})
}

func TestTodoService_DeleteTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_todo.NewMockRepository(ctrl)
	s := NewTodoService(repo)

	t.Run("invalid id", func(t *testing.T) {
		err := s.DeleteTodo(context.Background(), 0)
		require.ErrorIs(t, err, ErrValidation)
	})

	t.Run("success deletion todo item", func(t *testing.T) {
		repo.EXPECT().DeleteTodo(gomock.Any(), int64(123)).Return(nil)
		err := s.DeleteTodo(context.Background(), int64(123))
		require.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		repo.EXPECT().DeleteTodo(gomock.Any(), gomock.Any()).Return(sql.ErrConnDone)
		err := s.DeleteTodo(context.Background(), int64(22))
		require.Error(t, sql.ErrConnDone, err)
	})
}

func TestTodoService_GetTodoByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_todo.NewMockRepository(ctrl)
	s := NewTodoService(repo)
	now := time.Now()

	t.Run("invalid id", func(t *testing.T) {
		_, err := s.GetTodoByID(context.Background(), int64(0))
		require.Error(t, err, ErrValidation)
	})

	t.Run("success get todo item", func(t *testing.T) {
		repo.EXPECT().GetTodoByID(gomock.Any(), int64(22)).Return(
			dto.TodoItem{
				ID:          22,
				Title:       "title 22",
				Description: "description 22",
				Date:        pointer.Pointer(time.Date(2012, 01, 02, 0, 0, 0, 0, time.UTC)),
				Status:      "pending",
				CreatedAt:   now,
			},
			nil,
		)
		res, err := s.GetTodoByID(context.Background(), int64(22))
		require.NoError(t, err)
		require.Equal(t, model.TodoItem{
			ID:          22,
			Title:       "title 22",
			Description: "description 22",
			Date:        pointer.Pointer(time.Date(2012, 01, 02, 0, 0, 0, 0, time.UTC)),
			Status:      "pending",
			CreatedAt:   now,
		}, res)
	})

	t.Run("todo not found", func(t *testing.T) {
		repo.EXPECT().GetTodoByID(gomock.Any(), int64(404)).Return(dto.TodoItem{}, sql.ErrNoRows)
		_, err := s.GetTodoByID(context.Background(), int64(404))
		require.Error(t, err, ErrNotFound)
	})

	t.Run("database error", func(t *testing.T) {
		repo.EXPECT().GetTodoByID(gomock.Any(), int64(500)).Return(dto.TodoItem{}, sql.ErrConnDone)
		_, err := s.GetTodoByID(context.Background(), int64(500))
		require.Error(t, sql.ErrConnDone, err)
	})
}

func TestTodoService_UpdateTodo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_todo.NewMockRepository(ctrl)
	s := NewTodoService(repo)

	t.Run("update all fields", func(t *testing.T) {
		inp := &model.TodoItem{
			ID:          33,
			Title:       "title 33",
			Description: "description 33",
			Date:        pointer.Pointer(time.Date(2010, 12, 01, 0, 0, 0, 0, time.UTC)),
			Status:      "complete",
		}
		inpDto := &dto.TodoItem{
			ID:          33,
			Title:       "title 33",
			Description: "description 33",
			Date:        pointer.Pointer(time.Date(2010, 12, 01, 0, 0, 0, 0, time.UTC)),
			Status:      "complete",
		}
		repo.EXPECT().
			UpdateTodo(
				gomock.Any(),
				inpDto,
				[]string{model.TodoTitleField, model.TodoDescriptionField, model.TodoDateField, model.TodoStatusField},
			).DoAndReturn(func(ctx context.Context, item *dto.TodoItem, updatedFields []string) {
			item.UpdatedAt = pointer.Pointer(time.Now())
		}).Return(nil)

		err := s.UpdateTodo(context.Background(), inp)
		require.NoError(t, err)
		require.NotNil(t, inp.UpdatedAt)
	})

	t.Run("database error", func(t *testing.T) {
		repo.EXPECT().UpdateTodo(gomock.Any(), &dto.TodoItem{}, []string{}).Return(sql.ErrConnDone)
		err := s.UpdateTodo(context.Background(), &model.TodoItem{})
		require.Error(t, sql.ErrConnDone, err)
	})
}

func TestTodoService_ListTodos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock_todo.NewMockRepository(ctrl)
	s := NewTodoService(repo)

	t.Run("ok case", func(t *testing.T) {
		filter := dto.TodoFilter{
			Date:   pointer.Pointer(time.Date(2023, 02, 01, 0, 0, 0, 0, time.UTC)),
			Status: "completed",
			Page:   2,
			Limit:  2,
		}

		repo.EXPECT().ListTodos(gomock.Any(), filter).Return([]dto.TodoItem{
			{
				ID:          23,
				Title:       "title 23",
				Description: "desc 23",
				Date:        pointer.Pointer(time.Date(2023, 02, 01, 0, 0, 0, 0, time.UTC)),
				Status:      "completed",
				CreatedAt:   time.Date(2023, 02, 01, 10, 20, 30, 123, time.UTC),
			},
			{
				ID:          33,
				Title:       "title 33",
				Description: "desc 33",
				Date:        pointer.Pointer(time.Date(2023, 02, 01, 0, 0, 0, 0, time.UTC)),
				Status:      "completed",
				CreatedAt:   time.Date(2023, 02, 01, 10, 30, 30, 123, time.UTC),
			},
		}, int64(13), nil)
		res, err := s.ListTodos(context.Background(), filter)
		require.NoError(t, err)
		require.Equal(t, model.TodoPagination{
			Item: []model.TodoItem{
				{
					ID:          23,
					Title:       "title 23",
					Description: "desc 23",
					Date:        pointer.Pointer(time.Date(2023, 02, 01, 0, 0, 0, 0, time.UTC)),
					Status:      "completed",
					CreatedAt:   time.Date(2023, 02, 01, 10, 20, 30, 123, time.UTC),
				},
				{
					ID:          33,
					Title:       "title 33",
					Description: "desc 33",
					Date:        pointer.Pointer(time.Date(2023, 02, 01, 0, 0, 0, 0, time.UTC)),
					Status:      "completed",
					CreatedAt:   time.Date(2023, 02, 01, 10, 30, 30, 123, time.UTC),
				},
			},
			TotalItems: 13,
		}, res)
	})

	t.Run("database error", func(t *testing.T) {
		repo.EXPECT().ListTodos(gomock.Any(), gomock.Any()).Return(nil, int64(0), sql.ErrConnDone)
		_, err := s.ListTodos(context.Background(), dto.TodoFilter{})
		require.Error(t, err, sql.ErrConnDone)
	})
}

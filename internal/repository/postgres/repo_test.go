package postgres

import (
	"context"
	"database/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
	"todo-list/internal/util/pointer"
)

func mustTruncate(t *testing.T) {
	_, err := repo.DB.Exec("DELETE FROM todos;")
	require.NoError(t, err)
}

func mustCreateTodo(t *testing.T, item *dto.TodoItem) {
	err := repo.CreateTodo(context.Background(), item)
	require.NoError(t, err)
}

func mustCreateTodos(t *testing.T, items []dto.TodoItem) {
	for i := range items {
		err := repo.CreateTodo(context.Background(), &items[i])
		require.NoError(t, err)
	}
}

func TestTodoRepository_CreateTodo(t *testing.T) {
	t.Run("ok case", func(t *testing.T) {
		mustTruncate(t)
		input := &dto.TodoItem{
			Title:       "title 111",
			Description: "desc 111",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "completed",
		}

		mustCreateTodo(t, input)
		require.NotNil(t, input.CreatedAt)
		require.NotNil(t, input.ID)

		out := &dto.TodoItem{}
		err := repo.DB.QueryRowx("SELECT * FROM todos where id=$1", input.ID).StructScan(out)
		require.NoError(t, err)
		require.Equal(t, input, out)
		mustTruncate(t)
	})
}

func TestTodoRepository_DeleteTodo(t *testing.T) {
	t.Run("delete existed todo", func(t *testing.T) {
		mustTruncate(t)
		inp1 := &dto.TodoItem{
			Title:       "title 1",
			Description: "desc 1",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "completed",
		}
		inp2 := &dto.TodoItem{
			Title:       "title 2",
			Description: "desc 2",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "pending",
		}

		mustCreateTodo(t, inp1)
		mustCreateTodo(t, inp2)

		err := repo.DeleteTodo(context.Background(), inp1.ID)
		require.NoError(t, err)

		res := &dto.TodoItem{}
		err = repo.DB.QueryRowx("SELECT * FROM todos WHERE id = $1", inp1.ID).StructScan(res)
		require.ErrorIs(t, sql.ErrNoRows, err)
		mustTruncate(t)
	})

	t.Run("delete not existed item", func(t *testing.T) {
		mustTruncate(t)

		inp1 := &dto.TodoItem{
			Title:       "title 1",
			Description: "desc 1",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "completed",
		}
		inp2 := &dto.TodoItem{
			Title:       "title 2",
			Description: "desc 2",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "pending",
		}

		mustCreateTodo(t, inp1)
		mustCreateTodo(t, inp2)

		err := repo.DeleteTodo(context.Background(), 404)
		require.ErrorIs(t, err, sql.ErrNoRows)

		mustTruncate(t)
	})
}

func TestTodoRepository_GetTodoByID(t *testing.T) {
	t.Run("ok case", func(t *testing.T) {
		mustTruncate(t)

		inp1 := &dto.TodoItem{
			Title:       "title 1",
			Description: "desc 1",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "completed",
		}
		inp2 := &dto.TodoItem{
			Title:       "title 2",
			Description: "desc 2",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "pending",
		}

		mustCreateTodo(t, inp1)
		mustCreateTodo(t, inp2)

		res, err := repo.GetTodoByID(context.Background(), inp1.ID)
		require.NoError(t, err)
		require.Equal(t, res, *inp1)

		mustTruncate(t)
	})

	t.Run("get not existed item", func(t *testing.T) {
		_, err := repo.GetTodoByID(context.Background(), 404)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func TestTodoRepository_ListTodos(t *testing.T) {
	mustTruncate(t)

	input := []dto.TodoItem{
		{
			Title:       "title 1",
			Description: "desc 1",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "completed",
		},
		{
			Title:       "title 2",
			Description: "desc 2",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "pending",
		},
		{
			Title:       "title 3",
			Description: "desc 3",
			Date:        pointer.Pointer(time.Date(2023, 12, 2, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "completed",
		},
		{
			Title:       "title 4",
			Description: "desc 4",
			Date:        pointer.Pointer(time.Date(2023, 12, 2, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "pending",
		},
	}

	mustCreateTodos(t, input)

	tests := []struct {
		name      string
		filter    dto.TodoFilter
		want      []dto.TodoItem
		wantTotal int64
		wantErr   bool
	}{
		{
			name:      "case without filter",
			filter:    dto.TodoFilter{},
			want:      input,
			wantTotal: 4,
			wantErr:   false,
		},
		{
			name: "date filter",
			filter: dto.TodoFilter{
				Date: pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			},
			want:      input[:2],
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "status filter",
			filter: dto.TodoFilter{
				Status: "completed",
			},
			want:      []dto.TodoItem{input[0], input[2]},
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "date & status filter",
			filter: dto.TodoFilter{
				Date:   pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
				Status: "completed",
			},
			want:      []dto.TodoItem{input[0]},
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "limit filter",
			filter: dto.TodoFilter{
				Limit: 2,
			},
			want:      input[:2],
			wantTotal: 4,
			wantErr:   false,
		},
		{
			name: "page filter",
			filter: dto.TodoFilter{
				Page:  2,
				Limit: 2,
			},
			want:      input[2:],
			wantTotal: 4,
			wantErr:   false,
		},
		{
			name: "not found",
			filter: dto.TodoFilter{
				Date: pointer.Pointer(time.Date(2023, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			},
			want:      nil,
			wantTotal: 0,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := repo.ListTodos(context.Background(), tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTodos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, len(tt.want), len(got))
			for i := range tt.want {
				require.Empty(t, cmp.Diff(tt.want[i], got[i], cmpopts.IgnoreFields(dto.TodoItem{}, "TotalItems")))
			}
			require.Equal(t, tt.wantTotal, got1)
		})
	}

	mustTruncate(t)
}

func TestTodoRepository_UpdateTodo(t *testing.T) {
	mustTruncate(t)

	input := []dto.TodoItem{
		{
			Title:       "title 1",
			Description: "desc 1",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "completed",
		},
		{
			Title:       "title 2",
			Description: "desc 2",
			Date:        pointer.Pointer(time.Date(2023, 12, 1, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "pending",
		},
		{
			Title:       "title 3",
			Description: "desc 3",
			Date:        pointer.Pointer(time.Date(2023, 12, 2, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "completed",
		},
		{
			Title:       "title 4",
			Description: "desc 4",
			Date:        pointer.Pointer(time.Date(2023, 12, 2, 0, 0, 0, 0, time.FixedZone("", 0))),
			Status:      "pending",
		},
	}

	mustCreateTodos(t, input)

	t.Run("update title", func(t *testing.T) {
		err := repo.UpdateTodo(context.Background(), &dto.TodoItem{
			ID:    input[0].ID,
			Title: "updated title 1",
		}, []string{model.TodoTitleField})
		require.NoError(t, err)

		item, err := repo.GetTodoByID(context.Background(), input[0].ID)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(input[0], item, cmpopts.IgnoreFields(dto.TodoItem{}, "Title", "UpdatedAt")))
		require.Equal(t, item.Title, "updated title 1")
		require.NotNil(t, item.UpdatedAt)
	})

	t.Run("update description", func(t *testing.T) {
		err := repo.UpdateTodo(context.Background(), &dto.TodoItem{
			ID:          input[1].ID,
			Description: "updated description 2",
		}, []string{model.TodoDescriptionField})
		require.NoError(t, err)

		item, err := repo.GetTodoByID(context.Background(), input[1].ID)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(input[1], item, cmpopts.IgnoreFields(dto.TodoItem{}, "Description", "UpdatedAt")))
		require.Equal(t, item.Description, "updated description 2")
		require.NotNil(t, item.UpdatedAt)
	})

	t.Run("update date", func(t *testing.T) {
		date := pointer.Pointer(time.Date(2023, 1, 1, 0, 0, 0, 0, time.FixedZone("", 0)))
		err := repo.UpdateTodo(context.Background(), &dto.TodoItem{
			ID:   input[2].ID,
			Date: date,
		}, []string{model.TodoDateField})
		require.NoError(t, err)

		item, err := repo.GetTodoByID(context.Background(), input[2].ID)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(input[2], item, cmpopts.IgnoreFields(dto.TodoItem{}, "Date", "UpdatedAt")))
		require.Equal(t, item.Date, date)
		require.NotNil(t, item.UpdatedAt)
	})

	t.Run("status date", func(t *testing.T) {
		err := repo.UpdateTodo(context.Background(), &dto.TodoItem{
			ID:     input[3].ID,
			Status: "pending",
		}, []string{model.TodoStatusField})
		require.NoError(t, err)

		item, err := repo.GetTodoByID(context.Background(), input[3].ID)
		require.NoError(t, err)
		require.Empty(t, cmp.Diff(input[3], item, cmpopts.IgnoreFields(dto.TodoItem{}, "Status", "UpdatedAt")))
		require.Equal(t, item.Status, "pending")
		require.NotNil(t, item.UpdatedAt)
	})

	mustTruncate(t)
}

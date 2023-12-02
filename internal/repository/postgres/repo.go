package postgres

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"strings"
	"time"
	"todo-list/internal/config"
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
)

type TodoRepository struct {
	DB *sqlx.DB
}

const (
	DefaultLimit = 100
)

func NewPostgresTodoRepository() *TodoRepository {
	connect, err := sqlx.Connect(config.Config.DBConfig.Driver, config.Config.DBConfig.ConnectionString())
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	if config.Config.AppLevel != "test" {
		if err := goose.SetDialect(config.Config.DBConfig.Driver); err != nil {
			log.Fatalf("goose set dialect error: %v", err)
		}

		if err := goose.Up(connect.DB, "migrations"); err != nil {
			log.Fatalf("goose up :%s", err)
		}
	}

	return &TodoRepository{
		DB: connect,
	}
}

func (s *TodoRepository) Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(s.DB)
}

func (s *TodoRepository) CreateTodo(ctx context.Context, item *dto.TodoItem) error {
	q := s.Builder().Insert("todos").SetMap(map[string]interface{}{
		model.TodoTitleField:       item.Title,
		model.TodoDescriptionField: item.Description,
		model.TodoDateField:        item.Date,
		model.TodoStatusField:      item.Status,
	}).Suffix("RETURNING id, created_at")

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	if err = s.DB.QueryRowxContext(ctx, query, args...).StructScan(item); err != nil {
		return err
	}

	return nil
}

func (s *TodoRepository) GetTodoByID(ctx context.Context, id int64) (dto.TodoItem, error) {
	q := s.Builder().Select("*").From("todos").Where(sq.Eq{"id": id})
	query, args, err := q.ToSql()
	if err != nil {
		return dto.TodoItem{}, err
	}

	var res dto.TodoItem
	if err = s.DB.QueryRowxContext(ctx, query, args...).StructScan(&res); err != nil {
		return res, err
	}

	return res, nil
}

var col = map[string]func(item *dto.TodoItem) interface{}{
	model.TodoTitleField:       func(item *dto.TodoItem) interface{} { return item.Title },
	model.TodoDescriptionField: func(item *dto.TodoItem) interface{} { return item.Description },
	model.TodoDateField:        func(item *dto.TodoItem) interface{} { return item.Date },
	model.TodoStatusField:      func(item *dto.TodoItem) interface{} { return item.Status },
}

func (s *TodoRepository) UpdateTodo(ctx context.Context, item *dto.TodoItem, updatedFields []string) error {
	query := s.Builder().Update("todos").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": item.ID}).Suffix("RETURNING id, title, description, date, status, created_at, updated_at")

	for _, fieldToUpdate := range updatedFields {
		getter, ok := col[fieldToUpdate]
		if !ok {
			return fmt.Errorf("field not found")
		}
		query = query.Set(fieldToUpdate, getter(item))
	}

	q, args, err := query.ToSql()
	if err != nil {
		return err
	}

	err = s.DB.QueryRowxContext(ctx, q, args...).StructScan(item)
	if err != nil {
		return err
	}

	return nil
}

func (s *TodoRepository) DeleteTodo(ctx context.Context, id int64) error {
	q := s.Builder().Delete("todos").Where(sq.Eq{"id": id})
	res, err := q.ExecContext(ctx)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func applyTodoFilter(s sq.SelectBuilder, f dto.TodoFilter) sq.SelectBuilder {
	if f.Date != nil {
		s = s.Where(sq.Eq{model.TodoDateField: f.Date})
	}

	if f.Status != "" {
		s = s.Where(sq.Eq{model.TodoStatusField: f.Status})
	}

	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 || f.Limit > 10_000 {
		f.Limit = DefaultLimit
	}

	s = s.Limit(uint64(f.Limit)).Offset(uint64((f.Page - 1) * f.Limit))
	return s
}

func (s *TodoRepository) ListTodos(ctx context.Context, filter dto.TodoFilter) ([]dto.TodoItem, int64, error) {
	q := s.Builder().Select(
		"id", strings.Join(model.TodoFields, ", "), "created_at", "updated_at",
		"COUNT(*) OVER() as total_items").
		From("todos").
		OrderBy("id")

	q = applyTodoFilter(q, filter)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := s.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	todos := make([]dto.TodoItem, 0)

	for rows.Next() {
		var buff dto.TodoItem
		err = rows.StructScan(&buff)
		if err != nil {
			return nil, 0, err
		}

		todos = append(todos, buff)
	}
	var totalItems int64
	if len(todos) != 0 {
		totalItems = todos[0].TotalItems
	}

	return todos, totalItems, nil
}

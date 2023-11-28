package postgres

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"todo-list/internal/config"
	"todo-list/internal/domain/dto"
)

type TodoRepository struct {
	DB *sqlx.DB
}

func NewPostgresTodoRepository() *TodoRepository {
	connect, err := sqlx.Connect(config.Config.DBConfig.Driver, config.Config.DBConfig.ConnectionString())
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	//
	//if err := goose.SetDialect(config.Config.DBConfig.Driver); err != nil {
	//	log.Fatalf("goose set dialect error: %v", err)
	//}
	//
	//if err := goose.Up(connect.DB, "migrations"); err != nil {
	//	log.Fatalf("goose up :%s", err)
	//}

	return &TodoRepository{
		DB: connect,
	}
}

func (s *TodoRepository) Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(s.DB)
}

func (s *TodoRepository) CreateTodo(ctx context.Context, item *dto.TodoItem) error {
	//TODO implement me
	panic("implement me")
}

func (s *TodoRepository) GetTodoByID(ctx context.Context, id int64) (dto.TodoItem, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TodoRepository) UpdateTodo(ctx context.Context, item *dto.TodoItem, updatedFields []string) error {
	//TODO implement me
	panic("implement me")
}

func (s *TodoRepository) DeleteTodo(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *TodoRepository) ListTodos(ctx context.Context, filter dto.TodoFilter) ([]dto.TodoItem, int64, error) {
	//TODO implement me
	panic("implement me")
}

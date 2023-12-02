package main

import (
	_ "todo-list/docs"
	"todo-list/internal/repository/postgres"
	"todo-list/internal/server"
	"todo-list/internal/service/todo"
)

// @title TodoList API
// @version         1.0
// @host localhost:8080
// @BasePath /api/v1
func main() {
	repo := postgres.NewPostgresTodoRepository()
	s := todo.NewTodoService(repo)
	srv := server.NewServer(s)
	_ = srv.Run()
}

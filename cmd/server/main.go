package main

import (
	"todo-list/internal/repository/postgres"
	"todo-list/internal/server"
	"todo-list/internal/service/todo"
)

func main() {
	repo := postgres.NewPostgresTodoRepository()
	s := todo.NewTodoService(repo)
	srv := server.NewServer(s)
	srv.Run()
}

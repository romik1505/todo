package postgres

import (
	"os"
	"testing"
)

var (
	repo *TodoRepository
)

func TestMain(m *testing.M) {
	repo = NewPostgresTodoRepository()
	code := m.Run()
	os.Exit(code)
}

package v1

import (
	"github.com/gin-gonic/gin"
	"todo-list/internal/service/todo"
)

type Handler struct {
	TodoService todo.Service
}

func NewHandler(ts todo.Service) *Handler {
	return &Handler{
		TodoService: ts,
	}
}
func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		td := v1.Group("/todo")
		{
			td.GET(":id", h.GetTodo)
			td.POST("", h.CreateTodo)
			td.PATCH("", h.UpdateTodo)
			td.DELETE(":id", h.DeleteTodo)
			td.GET("", h.ListTodos)
		}
	}
}

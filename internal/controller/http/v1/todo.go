package v1

import (
	"github.com/gin-gonic/gin"
	"todo-list/internal/domain/dto"
)

func (h *Handler) GetTodo(c *gin.Context) {
	id := c.Param("id")
	c.String(200, "%v", id)

	//TODO implement me
	//panic("implement me")
}

func (h *Handler) CreateTodo(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) UpdateTodo(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) DeleteTodo(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) ListTodos(c *gin.Context) {
	var filter dto.TodoFilter
	err := c.ShouldBind(&filter)
	if err != nil {

	}
	c.JSON(200, filter)
}

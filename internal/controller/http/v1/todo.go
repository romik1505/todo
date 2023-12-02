package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo-list/internal/domain/dto"
	"todo-list/internal/domain/model"
	"todo-list/internal/service/todo"
)

// GetTodo	godoc
//
// @Summary Get todo by id
// @Tags todo
// @Accept json
// @Produce json
// @Param id path int64 true "todo id"
// @Success 200 {object} model.TodoItem
// @Failure 400,404,500 {string} string
// @Router /todo/{id} [get]
func (h *Handler) GetTodo(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = c.Error(fmt.Errorf("%w: %v", todo.ErrValidation, err))
		return
	}

	res, err := h.TodoService.GetTodoByID(c, intID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// CreateTodo	godoc
//
// @Summary Create new todo
// @Tags todo
// @Accept json
// @Produce json
// @Param input body model.TodoItem true "todo info"
// @Success 200
// @Failure 400,404,500 {string} string
// @Router /todo [post]
func (h *Handler) CreateTodo(c *gin.Context) {
	var t model.TodoItem
	if err := c.ShouldBind(&t); err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.TodoService.CreateTodo(c, &t); err != nil {
		_ = c.Error(err)
		return
	}
}

// UpdateTodo	godoc
//
// @Summary Update todo item by id
// @Tags todo
// @Accept json
// @Produce json
// @Param input body model.TodoItem true "updated todo item"
// @Success 200
// @Failure 400,404,500 {string} string
// @Router /todo [patch]
func (h *Handler) UpdateTodo(c *gin.Context) {
	var t model.TodoItem

	if err := c.ShouldBind(&t); err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.TodoService.UpdateTodo(c, &t); err != nil {
		_ = c.Error(err)
		return
	}
}

// DeleteTodo	godoc
//
// @Summary delete todo by id
// @Tags todo
// @Accept json
// @Produce json
// @Param id path int64 true "id todo for delete"
// @Success 200
// @Failure 400,404,500 {string} string
// @Router /todo/{id} [delete]
func (h *Handler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.TodoService.DeleteTodo(c, intID); err != nil {
		_ = c.Error(err)
		return
	}
}

// ListTodos	godoc
//
// @Summary Get list todos with pagination
// @Tags todo
// @Accept json
// @Produce json
// @Param input query dto.TodoFilter true "filter for list todos"
// @Success 200,204 {object} model.TodoPagination
// @Failure 400,404,500 {string} string
// @Router /todo [get]
func (h *Handler) ListTodos(c *gin.Context) {
	var filter dto.TodoFilter
	err := c.ShouldBind(&filter)
	if err != nil {
		_ = c.Error(err)
		return
	}

	pagination, err := h.TodoService.ListTodos(c, filter)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pagination)
}

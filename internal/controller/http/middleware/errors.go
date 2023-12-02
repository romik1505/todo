package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-list/internal/service/todo"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		switch {
		case errors.Is(err.Err, todo.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err.Err, todo.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err.Err, todo.ErrEmptyContent):
			c.Status(http.StatusNoContent)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}
}

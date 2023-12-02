package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"todo-list/internal/controller/http/middleware"
	v1 "todo-list/internal/controller/http/v1"
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

func (h *Handler) NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ErrorHandler)

	handlerV1 := v1.NewHandler(h.TodoService)
	api := r.Group("/api")
	{
		handlerV1.Init(api)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

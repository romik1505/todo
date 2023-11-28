package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	http2 "todo-list/internal/controller/http"
	"todo-list/internal/service/todo"
)

type Server struct {
	httpServer  *http.Server
	todoService todo.Service
}

func NewServer(s todo.Service) Server {
	r := http2.NewHandler(s).NewRouter()

	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r,
	}

	return Server{
		httpServer:  srv,
		todoService: s,
	}
}

func (s *Server) Run() error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)
	go func(ch chan os.Signal) {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Println(err.Error())
			done <- os.Interrupt
			return
		}
	}(done)

	log.Printf("Server started on %s port", ":8080")

	<-done
	defer close(done)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log.Println("Server gracefully closed")

	return s.httpServer.Shutdown(ctx)
}

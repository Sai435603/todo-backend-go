package main

import (
	"log"

	"github.com/Sai435603/todo-backend-go/internal/app"
	"github.com/Sai435603/todo-backend-go/internal/handler"
	"github.com/Sai435603/todo-backend-go/internal/repository"
	"github.com/Sai435603/todo-backend-go/internal/server"
	"github.com/Sai435603/todo-backend-go/internal/service"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer application.DB.Close()
	todoRepo := repository.New(application.DB)
	todoSvc := service.New(todoRepo)
	todoHnd := handler.New(application.Logger, todoSvc)
	srv := server.New(application, todoHnd)
	if err := srv.Start(); err != nil {
		application.Logger.Error(
			"server stopped unexpectedly",
			"error", err,
		)
	}
}

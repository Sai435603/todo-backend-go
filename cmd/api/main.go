package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	//graceful shutdown block
	{
		go func() {
			application.Logger.Info("Server is starting...")
			if err := srv.Start(); err != nil {
				application.Logger.Error(
					"server stopped unexpectedly",
					"error", err,
				)
			}
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		sig := <-quit
		application.Logger.Warn(
			"Received shutdown signal, initiating graceful shutdown...",
			"signal", sig.String(),
		)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			application.Logger.Error(
				"Server forced to shutdown",
				"error", err,
			)
		}

		application.Logger.Info("Server exited cleanly")
	}
}

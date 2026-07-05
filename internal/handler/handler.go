package handler

import (
	"log/slog"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/service"
)

type Handler struct {
	logger      *slog.Logger
	service     *service.TodoService
	userService *service.UserService
	jwt         *auth.JWTService
	AuthHandler *AuthHandler
}

func New(
	logger *slog.Logger,
	svc *service.TodoService,
	userSvc *service.UserService,
	jwtSvc *auth.JWTService,
	authHandler *AuthHandler,
) *Handler {
	return &Handler{
		logger:      logger,
		service:     svc,
		userService: userSvc,
		jwt:         jwtSvc,
		AuthHandler: authHandler,
	}
}

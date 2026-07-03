package handler

import (
	"net/http"

	"github.com/Sai435603/todo-backend-go/internal/response"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	err := response.JSON(w, http.StatusOK,
		HealthResponse{
			Status: "ok",
		},
	)

	if err != nil {
		// Use h.logger instead of h.app.Logger
		h.logger.Error(
			"failed to write response",
			"error", err,
		)
	}
}

package response

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

func Error(w http.ResponseWriter, status int, message string) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(
		ErrorResponse{
			Success: false,
			Error: ErrorBody{
				Message: message,
			},
		},
	)
}

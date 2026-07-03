package response

import (
	"encoding/json"
	"net/http"

	"github.com/Sai435603/todo-backend-go/internal/validator"
)

// ValidationErrorResponse represents a structured validation error response.
type ValidationErrorResponse struct {
	Success bool                       `json:"success"`
	Error   ValidationErrorBody        `json:"error"`
}

// ValidationErrorBody contains the validation error details.
type ValidationErrorBody struct {
	Message string                     `json:"message"`
	Fields  validator.ValidationErrors `json:"fields"`
}

// ValidationError writes a 422 Unprocessable Entity response with field-level validation errors.
func ValidationError(w http.ResponseWriter, errors *validator.ValidationErrors) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	return json.NewEncoder(w).Encode(
		ValidationErrorResponse{
			Success: false,
			Error: ValidationErrorBody{
				Message: "validation failed",
				Fields:  *errors,
			},
		},
	)
}

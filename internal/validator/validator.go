package validator

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ValidationError represents a single field-level validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors that implements the error interface.
type ValidationErrors []ValidationError

// Error returns a human-readable string of all validation errors.
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}

	msgs := make([]string, 0, len(ve))
	for _, e := range ve {
		msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(msgs, "; ")
}

// HasErrors returns true if there are any validation errors.
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// Validator accumulates validation errors for a request.
type Validator struct {
	errors ValidationErrors
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{}
}

// AddError adds a validation error for the given field.
func (v *Validator) AddError(field, message string) {
	v.errors = append(v.errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// HasErrors returns true if the validator has accumulated any errors.
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Errors returns the accumulated validation errors.
func (v *Validator) Errors() ValidationErrors {
	return v.errors
}

// --- Validation Rule Methods ---
// Each method checks a condition and records an error if the check fails.

// Required checks that a string field is not empty after trimming whitespace.
func (v *Validator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.AddError(field, "is required")
	}
}

// MaxLength checks that a string does not exceed the given rune count.
func (v *Validator) MaxLength(field, value string, max int) {
	if utf8.RuneCountInString(value) > max {
		v.AddError(field, fmt.Sprintf("must be at most %d characters", max))
	}
}

// MinLength checks that a string has at least the given rune count.
func (v *Validator) MinLength(field, value string, min int) {
	if utf8.RuneCountInString(value) < min {
		v.AddError(field, fmt.Sprintf("must be at least %d characters", min))
	}
}

// PositiveID checks that an int64 is a valid positive identifier.
func (v *Validator) PositiveID(field string, id int64) {
	if id <= 0 {
		v.AddError(field, "must be a positive integer")
	}
}

// NotEmpty checks that a string is not the zero value (without trimming).
func (v *Validator) NotEmpty(field, value string) {
	if value == "" {
		v.AddError(field, "must not be empty")
	}
}

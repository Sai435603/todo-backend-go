package validator

const (
	// titleMaxLen is the maximum allowed length for a todo title.
	titleMaxLen = 255
	// titleMinLen is the minimum allowed length for a todo title.
	titleMinLen = 1
	// descriptionMaxLen is the maximum allowed length for a todo description.
	descriptionMaxLen = 5000
)

// CreateTodoInput holds the validated input for creating a todo.
type CreateTodoInput struct {
	Title       string
	Description string
}

// UpdateTodoInput holds the validated input for updating a todo.
type UpdateTodoInput struct {
	ID          int64
	Title       string
	Description string
	Completed   bool
}

// ValidateCreateTodo validates the fields required to create a new todo.
// Returns nil if all validations pass.
func ValidateCreateTodo(title, description string) *ValidationErrors {
	v := New()

	v.Required("title", title)
	v.MinLength("title", title, titleMinLen)
	v.MaxLength("title", title, titleMaxLen)
	v.MaxLength("description", description, descriptionMaxLen)

	if v.HasErrors() {
		errs := v.Errors()
		return &errs
	}
	return nil
}

// ValidateUpdateTodo validates the fields required to update a todo.
// Returns nil if all validations pass.
func ValidateUpdateTodo(id int64, title, description string) *ValidationErrors {
	v := New()

	v.PositiveID("id", id)
	v.Required("title", title)
	v.MinLength("title", title, titleMinLen)
	v.MaxLength("title", title, titleMaxLen)
	v.MaxLength("description", description, descriptionMaxLen)

	if v.HasErrors() {
		errs := v.Errors()
		return &errs
	}
	return nil
}

// ValidateID validates a single resource identifier.
// Returns nil if the ID is valid.
func ValidateID(id int64) *ValidationErrors {
	v := New()
	v.PositiveID("id", id)

	if v.HasErrors() {
		errs := v.Errors()
		return &errs
	}
	return nil
}

// ValidateSearchQuery validates a search query string.
// Returns nil if the query is valid.
func ValidateSearchQuery(query string) *ValidationErrors {
	v := New()
	v.Required("q", query)
	v.MaxLength("q", query, 500)

	if v.HasErrors() {
		errs := v.Errors()
		return &errs
	}
	return nil
}

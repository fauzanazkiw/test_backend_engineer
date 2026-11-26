package domain

import "errors"

var (
	ErrNotFound           = errors.New("data not found")
	ErrConflict           = errors.New("data conflict")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type ValidationErrors map[string][]string

func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{}
}

func (v ValidationErrors) Add(field, message string) {
	v[field] = append(v[field], message)
}

func (v ValidationErrors) Empty() bool {
	return len(v) == 0
}

// Pointer receiver: hanya *ValidationErrors yang menjadi error.
func (v *ValidationErrors) Error() string {
	return "validation failed"
}

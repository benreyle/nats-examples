package helpers

import "github.com/google/uuid"

type Position struct {
	ID      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Deleted bool      `json:"deleted"`
}

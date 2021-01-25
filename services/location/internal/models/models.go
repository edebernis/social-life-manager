package models

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

// ID abstracts id implementation
type ID uuid.UUID

// String serializes ID
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// MarshalJSON implements Marshaler interface
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// UnmarshalJSON implements Unmarshaler interface
func (id *ID) UnmarshalJSON(b []byte) error {
	_id, err := ParseID(string(b))
	*id = ID(_id)
	return err
}

// NewID returns a new ID using specific ID implementation
func NewID() ID {
	return ID(uuid.New())
}

// ParseID decodes input string into an ID using specific ID implementation
func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

// User of the app
type User struct {
	// Unique ID of the user.
	ID ID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// User email.
	Email string `json:"email" example:"test@no-reply.com"`
}

// NewUser creates a new user
func NewUser(id ID, email string) *User {
	return &User{
		id,
		email,
	}
}

// NewUserFromContext creates a new user from context data
func NewUserFromContext(c context.Context) *User {
	return c.Value("user").(*User)
}

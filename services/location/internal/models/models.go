package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// ID abstracts id implementation
type ID uuid.UUID

var (
	// NilID is empty ID, all zeros
	NilID = ID(uuid.Nil)
)

// String serializes ID
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// MarshalJSON implements the Marshaler interface
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// UnmarshalJSON implements the Unmarshaler interface
func (id *ID) UnmarshalJSON(b []byte) error {
	_id, err := ParseID(string(b))
	*id = ID(_id)
	return err
}

// Value implements the database/sql valuer interface
func (id ID) Value() (driver.Value, error) {
	return driver.Value(id.String()), nil
}

// Scan implements the database/sql scanner interface
func (id *ID) Scan(value interface{}) error {
	if value == nil {
		return errors.New("ID: cannot scan empty value")
	}

	dv, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("Failed to convert value: %s. %w", value, err)
	}
	v, ok := dv.(string)
	if !ok {
		return fmt.Errorf("Failed to cast driver value to string, %s", dv)
	}
	*id, err = ParseID(v)
	if err != nil {
		return fmt.Errorf("Failed to parse ID: %s. %w", v, err)
	}

	return nil
}

// NewID returns a new ID using specific ID implementation
func NewID() ID {
	return ID(uuid.New())
}

// ParseID decodes input string into an ID using specific ID implementation
func ParseID(s string) (ID, error) {
	if s == "" {
		return NilID, nil
	}

	id, err := uuid.Parse(s)
	return ID(id), err
}

type contextKey string

var (
	userContextKey = contextKey("user")
)

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
func NewUserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}

// NewContextWithUser adds user information to provided context
func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

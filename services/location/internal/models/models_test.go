package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDEmptyStringParsing(t *testing.T) {
	parsedID, err := ParseID("")

	assert.NoError(t, err)
	assert.Equal(t, NilID, parsedID)
}

func TestIDStringParsing(t *testing.T) {
	id := NewID()
	parsedID, err := ParseID(id.String())

	assert.NoError(t, err)
	assert.Equal(t, id, parsedID)
}

func TestIDJSONMarshalling(t *testing.T) {
	id := NewID()
	marshalledID, err := id.MarshalJSON()

	assert.NoError(t, err)

	unmarshalledID := new(ID)
	err = unmarshalledID.UnmarshalJSON(marshalledID)

	assert.NoError(t, err)
	assert.Equal(t, id, *unmarshalledID)
}

func TestIDSQLValuingAndScanning(t *testing.T) {
	id := NewID()
	value, err := id.Value()

	assert.NoError(t, err)

	scannedID := new(ID)
	err = scannedID.Scan(value)

	assert.NoError(t, err)
	assert.Equal(t, id, *scannedID)
}

func TestUserFromContext(t *testing.T) {
	user := NewUser(NewID(), "user@test.fr")
	ctx := NewContextWithUser(context.Background(), user)
	userFromContext, ok := NewUserFromContext(ctx)

	assert.True(t, ok)
	assert.Equal(t, user, userFromContext)
}

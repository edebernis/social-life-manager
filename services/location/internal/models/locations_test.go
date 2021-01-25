package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

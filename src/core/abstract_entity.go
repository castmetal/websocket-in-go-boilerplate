package core

import (
	"github.com/google/uuid"
)

func NewAbstractEntity(id uuid.UUID) *EntityBase {
	var _id uuid.UUID

	idString := string(id[:])

	if idString != "" {
		_id = id
	} else {
		_id = uuid.New()
	}

	var entity = &EntityBase{
		Id: _id,
	}

	return entity
}

package core

import (
	"github.com/google/uuid"
)

func NewAbstractEntity(id string) *EntityBase {
	var _id string

	if id != "" {
		_id = id
	} else {
		_id = uuid.New().String()
	}

	var entity = &EntityBase{
		Id: _id,
	}

	return entity
}

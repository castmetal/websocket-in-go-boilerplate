package in_memory_tests

import (
	"bytes"
	"encoding/json"
	"io"

	_dtos "websocket-in-go-boilerplate/src/application/dtos"

	"github.com/go-playground/validator/v10"
)

type (
	InMemoryCreateUserDTO struct {
		_dtos.ICreateUserDTO
		_dtos.UserDTO
	}
)

func (dto *InMemoryCreateUserDTO) Validate(message io.Reader) (_dtos.ICreateUserDTO, error) {
	var validate *validator.Validate
	var IDTO _dtos.ICreateUserDTO = &_dtos.CreateUserDTO{}

	messageBuffer := &bytes.Buffer{}
	messageBuffer.ReadFrom(message)

	err := json.Unmarshal(messageBuffer.Bytes(), &IDTO)
	if err != nil {
		return IDTO, err
	}

	validate = validator.New()
	err = validate.Struct(IDTO)
	if err != nil {
		return IDTO, err
	}

	return IDTO, nil
}

func (dto *InMemoryCreateUserDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}

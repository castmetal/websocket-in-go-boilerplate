package dtos

import (
	"bytes"
	"encoding/json"
	"io"

	_core "websocket-in-go-boilerplate/src/core"

	"github.com/go-playground/validator/v10"
)

type (
	ICreateUserDTO interface {
		_core.IDTO
		Validate(message io.Reader) (ICreateUserDTO, error)
		ToBytes() ([]byte, error)
	}

	CreateUserDTO struct {
		ICreateUserDTO
		UserDTO
	}
)

func (dto *CreateUserDTO) Validate(message io.Reader) (ICreateUserDTO, error) {
	var validate *validator.Validate
	var IDTO ICreateUserDTO = &CreateUserDTO{}

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

func (dto *CreateUserDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}

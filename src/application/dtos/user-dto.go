package dtos

import (
	"bytes"
	"encoding/json"
	"io"

	_core "websocket-in-go-boilerplate/src/core"

	"github.com/go-playground/validator/v10"
)

type (
	IUserDTO interface {
		_core.IDTO
		Validate(message io.Reader) (IUserDTO, error)
		ToBytes() ([]byte, error)
	}

	UserDTO struct {
		Id        string `json:"id"`
		FirstName string `json:"first_name" validate:"required,min=2"`
		LastName  string `json:"last_name" validate:"required,min=2"`
		UserName  string `json:"user_name" validate:"required,min=2"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,min=8"`
	}
)

func (dto *UserDTO) Validate(message io.Reader) (IUserDTO, error) {
	var validate *validator.Validate
	var IDTO IUserDTO = &UserDTO{}

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

func (dto *UserDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}

package dtos

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type CreateUserDTO struct {
	UserDTO
}

func (dto *CreateUserDTO) Validate(message io.Reader) (*CreateUserDTO, error) {
	var validate *validator.Validate
	DTO := CreateUserDTO{}

	messageBuffer := &bytes.Buffer{}
	messageBuffer.ReadFrom(message)

	err := json.Unmarshal(messageBuffer.Bytes(), &DTO)
	if err != nil {
		return &DTO, err
	}

	validate = validator.New()
	err = validate.Struct(DTO)
	if err != nil {
		return &DTO, err
	}

	return &DTO, nil
}

func (dto *CreateUserDTO) ToBytes() ([]byte, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	return b, nil
}

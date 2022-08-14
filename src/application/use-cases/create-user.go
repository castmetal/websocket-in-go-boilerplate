package use_cases

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	_dtos "websocket-in-go-boilerplate/src/application/dtos"
	_core "websocket-in-go-boilerplate/src/core"
	_user "websocket-in-go-boilerplate/src/domains/user"
)

type (
	CreateUser interface {
		_core.IUseCase
	}
	CreateUserRequest struct {
		CreateUser
		Repository _user.IUserRepository
		DTO        _dtos.ICreateUserDTO
	}
)

func NewCreateUser(repository _user.IUserRepository, dto _dtos.ICreateUserDTO) (_core.IUseCase, error) {
	var uc _core.IUseCase = &CreateUserRequest{
		Repository: repository,
		DTO:        dto,
	}

	return uc, nil
}

// Put here your validation message and return your struct mapper to service
func (uc *CreateUserRequest) Execute(ctx context.Context, message io.Reader) (bool, error) {
	dto, err := uc.DTO.Validate(message)
	if err != nil {
		return false, err
	}

	dtoBytes, err := dto.ToBytes()
	if err != nil {
		return false, err
	}

	var dtoReader io.Reader = bytes.NewReader(dtoBytes)

	var userDto = getUserDTO(dtoReader)

	user, err := _user.NewUserEntity(userDto)
	if err != nil {
		return false, err
	}

	_, err = uc.Repository.Create(ctx, user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getUserDTO(message io.Reader) _dtos.UserDTO {
	var userDTO _dtos.UserDTO
	messageBuffer := &bytes.Buffer{}
	messageBuffer.ReadFrom(message)

	if err := json.Unmarshal(messageBuffer.Bytes(), &userDTO); err != nil {
		panic(err)
	}

	return userDTO
}

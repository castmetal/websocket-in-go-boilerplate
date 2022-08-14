package use_cases

import (
	"context"
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
		DTO        *_dtos.CreateUserDTO
	}
)

func NewCreateUser(repository _user.IUserRepository, dto *_dtos.CreateUserDTO) (_core.IUseCase, error) {
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

	var userDto = getUserDTO(dto)

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

func getUserDTO(createUserDTO *_dtos.CreateUserDTO) _dtos.UserDTO {
	return _dtos.UserDTO{
		Id:        createUserDTO.Id,
		FirstName: createUserDTO.FirstName,
		LastName:  createUserDTO.LastName,
		UserName:  createUserDTO.UserName,
		Email:     createUserDTO.Email,
		Password:  createUserDTO.Password,
	}
}

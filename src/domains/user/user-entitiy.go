package user

import (
	"time"
	_dtos "websocket-in-go-boilerplate/src/application/dtos"
	_core "websocket-in-go-boilerplate/src/core"
)

type User struct {
	_core.EntityBase
	FirstName         string             `json:"first_name"`
	LastName          string             `json:"last_name"`
	UserName          string             `json:"user_name"`
	EncryptedPassword *EncryptedPassword `json:"encrypted_password"`
	CreatedAt         string             `json:"roles"`
	UpdatedAt         string             `json:"domain_events"`
}

func NewUserEntity(dto _dtos.UserDTO) (*User, error) {
	var user *User
	abstractEntity := _core.NewAbstractEntity(dto.Id)

	if _core.IsNullOrEmpty(dto.UserName) {
		return nil, _core.IsNullOrEmptyError("user_name")
	}

	actualDate := time.Now().Format("2006-01-02 15:04:05")

	user = &User{
		FirstName:         dto.FirstName,
		LastName:          dto.LastName,
		UserName:          dto.UserName,
		EncryptedPassword: NewEncryptedPassword(dto.Password),
		CreatedAt:         actualDate,
		UpdatedAt:         actualDate,
	}

	user.Id = abstractEntity.Id

	/*user.AddEvent(&events.UserCreated{
		Id:        user.Id,
		FirstName: firstName,
		LastName:  lastName,
		UserName:  username,
	})*/

	return user, nil
}

package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	_dtos "websocket-in-go-boilerplate/src/application/dtos"
	_core "websocket-in-go-boilerplate/src/core"
)

type User struct {
	gorm.Model
	_core.EntityBase
	FirstName         string             `json:"first_name" gorm:"type:varchar(60);column:first_name"`
	LastName          string             `json:"last_name" gorm:"type:varchar(140);column:last_name"`
	UserName          string             `json:"user_name" gorm:"type:varchar(90);unique;uniqueIndex;collumn:user_name"`
	Email             string             `json:"email" gorm:"type:varchar(150);unique;uniqueIndex;column:email"`
	EncryptedPassword *EncryptedPassword `json:"encrypted_password" gorm:"embedded"`
	CreatedAt         time.Time          `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         time.Time          `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt         gorm.DeletedAt     `json:"deleted_at" gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func NewUserEntity(dto _dtos.UserDTO) (*User, error) {
	var user *User

	uuidFromId := uuid.Must(uuid.FromBytes([]byte(dto.Id)))
	abstractEntity := _core.NewAbstractEntity(uuidFromId)

	if _core.IsNullOrEmpty(dto.UserName) {
		return nil, _core.IsNullOrEmptyError("user_name")
	}

	actualDate := time.Now()

	user = &User{
		FirstName:         dto.FirstName,
		LastName:          dto.LastName,
		UserName:          dto.UserName,
		Email:             dto.Email,
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

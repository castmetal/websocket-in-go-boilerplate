package repositories

import (
	"context"
	"fmt"
	"log"

	_core "websocket-in-go-boilerplate/src/core"
	_user "websocket-in-go-boilerplate/src/domains/user"
	_infra_db "websocket-in-go-boilerplate/src/infra/db"

	"gorm.io/gorm"
)

const _collectionName = "Users"

type userRepository struct {
	db *gorm.DB
}

func NewUserRepositoryFromConfig() _user.IUserRepository {
	db, err := _infra_db.NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Error on Database Connection: %v", err)
	}

	return newUserRepository(db)
}

func newUserRepository(db *gorm.DB) _user.IUserRepository {
	return &userRepository{db: db}
}

func (repository userRepository) FindOneById(ctx context.Context, id string) (*_user.User, error) {
	var user *_user.User

	repository.db.First(&user, "id = ?", id)
	userId := string(user.Id[:])

	if userId == "" {
		return nil, _core.NotFoundError("User")
	}

	return user, nil
}

func (repository userRepository) Create(ctx context.Context, user *_user.User) (*_user.User, error) {
	var u *_user.User

	repository.db.First(&u, "email = ?", user.Email)
	fmt.Println(u.Email)
	if u.Email != "" {
		return nil, _core.AlreadyExistsError("User")
	}

	result := repository.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repository userRepository) FindOneByEmail(ctx context.Context, email string) (*_user.User, error) {
	var user *_user.User

	repository.db.First(&user, "email = ?", email)
	userId := string(user.Id[:])
	if userId == "" {
		return nil, _core.NotFoundError("User")
	}

	return user, nil
}

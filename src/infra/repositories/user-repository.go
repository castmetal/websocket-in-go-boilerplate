package repositories

import (
	"context"
	"log"

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

	// TO DO - Get database connection from config

	return newUserRepository(db)
}

func newUserRepository(db *gorm.DB) _user.IUserRepository {
	return &userRepository{db: db}
}

func (repository userRepository) FindOneById(ctx context.Context, id string) (*_user.User, error) {
	var user *_user.User

	// your implementation here

	return user, nil
}

func (repository userRepository) Create(ctx context.Context, user *_user.User) (*_user.User, error) {
	var u *_user.User = user

	// your implementation here

	// TO DO - Domain events

	return u, nil
}

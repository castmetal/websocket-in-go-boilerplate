package core

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IWebsocket interface {
	EstabilishConnection(ctx context.Context, userId string, conn net.Conn) (bool, error)
	ExecuteUseCase(ctx context.Context, useCase IUseCase, userId string, conn net.Conn) (bool, error)
	WriteToAllClients(ctx context.Context, userId string, conn net.Conn) (bool, error)
	WriteToAnUser(ctx context.Context, userId string, conn net.Conn) (bool, error)
}

type IDTO interface {
	Validate(message io.Reader) (IDTO, error)
	ToBytes() ([]byte, error)
}

type IUseCase interface {
	Execute(message io.Reader) (bool, error)
}

type IError interface {
	error
}

type EntityBase struct {
	gorm.Model
	Id uuid.UUID `json:"id" bson:"_id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4();collumn:id"`
}

type IEntity interface {
	SetId(id uuid.UUID) *EntityBase
	GetId(entity *EntityBase) uuid.UUID
	GetEntity() *EntityBase
}

type Controller struct {
	URL    string
	Method string
	Handle func(w http.ResponseWriter, r *http.Request)
}

// TODO
type IRepository interface {
}

type IDatabase interface {
}

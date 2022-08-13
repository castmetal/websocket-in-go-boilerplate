package core

import (
	"context"
	"io"
	"net"
	"net/http"
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
	Id string `json:"id" ,bson:"_id"`
}

type IEntity interface {
	SetId(id string) *EntityBase
	GetId(entity *EntityBase) string
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

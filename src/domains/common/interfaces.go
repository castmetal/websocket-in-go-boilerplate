package common

import (
	"context"
	"net"
	"net/http"
)

type IWebsocket interface {
	SimpleSocket(ctx context.Context, userId string, conn net.Conn) (bool, error)
	WriteToAllClients(ctx context.Context, userId string, conn net.Conn) (bool, error)
	WriteToAnUser(ctx context.Context, userId string, conn net.Conn) (bool, error)
}

type IUseCase interface {
	Execute() ([]byte, error)
}

type IError interface {
	error
}

type Controller struct {
	URL    string
	Method string
	Handle func(w http.ResponseWriter, r *http.Request)
}

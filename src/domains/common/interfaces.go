package common

import (
	"context"
	"net/http"
)

type IWebsocket interface {
	SimpleSocket(ctx context.Context) (bool, error)
	WriteToAllClients(ctx context.Context) (bool, error)
	WriteToAnUser(ctx context.Context) (bool, error)
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

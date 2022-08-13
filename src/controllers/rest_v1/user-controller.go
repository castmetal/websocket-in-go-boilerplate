package controllers_rest_v1

import (
	"context"
	"net"
	"net/http"
	_config "websocket-in-go-boilerplate/src/config"
	_interfaces "websocket-in-go-boilerplate/src/domains/common"
	_user "websocket-in-go-boilerplate/src/domains/user"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"

	"github.com/gobwas/ws"
)

const (
	VERSION = "/v1"
)

type UserControllers struct {
	Controllers map[int]_interfaces.Controller
}

func CreateSimpleClientSocket(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _interfaces.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _interfaces.IWebsocket = _user.NewUserSocketService(epoll)
		userId := r.Header.Get(_config.SystemParams.AUTH_HEADER)

		conn, err := WsHandshake(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Connection"))

			return
		}

		ws.SimpleSocket(context.Background(), userId, conn)
	}
}

func WriteMessageToAnUser(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _interfaces.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _interfaces.IWebsocket = _user.NewUserSocketService(epoll)
		userId := r.Header.Get(_config.SystemParams.AUTH_HEADER)

		conn, err := WsHandshake(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Connection"))

			return
		}

		ws.WriteToAnUser(context.Background(), userId, conn)
	}
}

func WriteMessageToAllUsers(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _interfaces.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _interfaces.IWebsocket = _user.NewUserSocketService(epoll)
		userId := r.Header.Get(_config.SystemParams.AUTH_HEADER)

		conn, err := WsHandshake(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Connection"))

			return
		}

		ws.WriteToAllClients(context.Background(), userId, conn)
	}
}

func GetUserControllers(epoll *_epoll.Epoll) map[int]_interfaces.Controller {
	userController := newUserController()

	userController.add(_interfaces.Controller{
		URL:    VERSION + "/ws",
		Method: "GET",
		Handle: CreateSimpleClientSocket(epoll),
	})

	userController.add(_interfaces.Controller{
		URL:    VERSION + "/ws/writeToAnUser",
		Method: "GET",
		Handle: WriteMessageToAnUser(epoll),
	})

	userController.add(_interfaces.Controller{
		URL:    VERSION + "/ws/writeToAll",
		Method: "GET",
		Handle: WriteMessageToAllUsers(epoll),
	})

	return userController.Controllers
}

func WsHandshake(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return nil, _interfaces.InvalidConnectionError()
	}

	return conn, nil
}

func (ctl *UserControllers) add(controller _interfaces.Controller) {
	ctl.Controllers[len(ctl.Controllers)] = controller
}

func newUserController() *UserControllers {
	controllers := UserControllers{
		Controllers: make(map[int]_interfaces.Controller),
	}

	return &controllers
}

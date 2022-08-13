package controllers_rest_v1

import (
	"context"
	"net/http"
	_interfaces "websocket-in-go-boilerplate/src/domains/common"
	_user "websocket-in-go-boilerplate/src/domains/user"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"
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

		var ws _interfaces.IWebsocket = _user.NewUserSocketService(w, r, epoll)
		ws.SimpleSocket(context.Background())
	}
}

func WriteMessageToAnUser(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _interfaces.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _interfaces.IWebsocket = _user.NewUserSocketService(w, r, epoll)

		ws.WriteToAnUser(context.Background())
	}
}

func WriteMessageToAllUsers(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _interfaces.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _interfaces.IWebsocket = _user.NewUserSocketService(w, r, epoll)

		ws.WriteToAllClients(context.Background())
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

func (ctl *UserControllers) add(controller _interfaces.Controller) {
	ctl.Controllers[len(ctl.Controllers)] = controller
}

func newUserController() *UserControllers {
	controllers := UserControllers{
		Controllers: make(map[int]_interfaces.Controller),
	}

	return &controllers
}

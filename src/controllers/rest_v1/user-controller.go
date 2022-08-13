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

func NewUserController() *UserControllers {
	controllers := UserControllers{
		Controllers: make(map[int]_interfaces.Controller),
	}

	return &controllers
}

func (ctl *UserControllers) Add(controller _interfaces.Controller) {
	ctl.Controllers[len(ctl.Controllers)] = controller
}

func GetUserControllers(epoll *_epoll.Epoll) map[int]_interfaces.Controller {
	userController := NewUserController()

	userController.Add(_interfaces.Controller{
		URL:    VERSION + "/ws",
		Method: "GET",
		Handle: func(w http.ResponseWriter, r *http.Request) {
			err := _interfaces.ValidateSameMethod(w, r, "GET")
			if err != nil {
				return
			}

			var ws _interfaces.IWebsocket = _user.NewUserSocketService(w, r, epoll)
			ws.SimpleSocket(context.Background())
		},
	})

	userController.Add(_interfaces.Controller{
		URL:    VERSION + "/ws/writeToAnUser",
		Method: "GET",
		Handle: func(w http.ResponseWriter, r *http.Request) {
			err := _interfaces.ValidateSameMethod(w, r, "GET")
			if err != nil {
				return
			}

			var ws _interfaces.IWebsocket = _user.NewUserSocketService(w, r, epoll)

			ws.WriteToAnUser(context.Background())
		},
	})

	userController.Add(_interfaces.Controller{
		URL:    VERSION + "/ws/writeToAll",
		Method: "GET",
		Handle: func(w http.ResponseWriter, r *http.Request) {
			err := _interfaces.ValidateSameMethod(w, r, "GET")
			if err != nil {
				return
			}

			var ws _interfaces.IWebsocket = _user.NewUserSocketService(w, r, epoll)

			ws.WriteToAllClients(context.Background())
		},
	})

	return userController.Controllers
}

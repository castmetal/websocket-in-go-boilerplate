package controllers_rest_v1

import (
	_core "websocket-in-go-boilerplate/src/core"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"
)

func GetControllers(epoll *_epoll.Epoll) []_core.Controller {
	var controllers []_core.Controller

	userControllers := GetUserControllers(epoll)
	for _, v := range userControllers {
		controllers = append(controllers, v)
	}

	return controllers
}

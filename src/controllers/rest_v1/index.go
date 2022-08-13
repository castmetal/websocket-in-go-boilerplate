package controllers_rest_v1

import (
	_interfaces "websocket-in-go-boilerplate/src/domains/common"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"
)

func GetControllers(epoll *_epoll.Epoll) []_interfaces.Controller {
	var controllers []_interfaces.Controller

	userControllers := GetUserControllers(epoll)
	for _, v := range userControllers {
		controllers = append(controllers, v)
	}

	return controllers
}

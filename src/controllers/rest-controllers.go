package controllers

import (
	_controllers_rest_v1 "websocket-in-go-boilerplate/src/controllers/rest_v1"
	_interfaces "websocket-in-go-boilerplate/src/domains/common"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"
)

func SetRestControllers(mux *_interfaces.MyMux, epoll *_epoll.Epoll) {
	controllers := _controllers_rest_v1.GetControllers(epoll)

	for _, ctl := range controllers {
		mux.HandleFunc(ctl.URL, ctl.Handle)
	}
}

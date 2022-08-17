package controllers_rest_v1

import (
	"context"
	"net"
	"net/http"
	_dtos "websocket-in-go-boilerplate/src/application/dtos"
	_services "websocket-in-go-boilerplate/src/application/services"
	_use_cases "websocket-in-go-boilerplate/src/application/use-cases"
	_config "websocket-in-go-boilerplate/src/config"
	_core "websocket-in-go-boilerplate/src/core"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"
	_repositories "websocket-in-go-boilerplate/src/infra/repositories"

	"github.com/gobwas/ws"
)

const (
	VERSION = "/v1"
)

type UserControllers struct {
	Controllers map[int]_core.Controller
}

func EstabilishConnection(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _core.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _core.IWebsocket = _services.NewWebSocketService(epoll)
		userId := r.Header.Get(_config.SystemParams.AUTH_HEADER)

		conn, err := WsHandshake(w, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Connection"))

			return
		}

		_, err = ws.EstabilishConnection(context.Background(), userId, conn)
		if err != nil {
			conn.Close()
			return
		}
	}
}

func CreateUserWithSocket(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _core.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _core.IWebsocket = _services.NewWebSocketService(epoll)
		userId := r.Header.Get(_config.SystemParams.AUTH_HEADER)

		conn, err := WsHandshake(w, r)
		if err != nil {
			conn.Close()
			return
		}

		userRepository := _repositories.NewUserRepositoryFromConfig()

		var dto = &_dtos.CreateUserDTO{}
		useCase, err := _use_cases.NewCreateUser(userRepository, dto)
		if err != nil {
			conn.Close()
			return
		}

		_, err = ws.ExecuteUseCase(context.Background(), useCase, userId, conn)
		if err != nil {
			conn.Close()
			return
		}
	}
}

func WriteMessageToAnUser(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _core.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _core.IWebsocket = _services.NewWebSocketService(epoll)
		userId := r.Header.Get(_config.SystemParams.AUTH_HEADER)

		conn, err := WsHandshake(w, r)
		if err != nil {
			conn.Close()
			return
		}

		_, err = ws.WriteToAnUser(context.Background(), userId, conn)
		if err != nil {
			conn.Close()
			return
		}
	}
}

func WriteMessageToAllUsers(epoll *_epoll.Epoll) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := _core.ValidateSameMethod(w, r, "GET")
		if err != nil {
			return
		}

		var ws _core.IWebsocket = _services.NewWebSocketService(epoll)
		userId := r.Header.Get(_config.SystemParams.AUTH_HEADER)

		conn, err := WsHandshake(w, r)
		if err != nil {
			conn.Close()
			return
		}

		_, err = ws.WriteToAllClients(context.Background(), userId, conn)
		if err != nil {
			conn.Close()
			return
		}
	}
}

func GetUserControllers(epoll *_epoll.Epoll) map[int]_core.Controller {
	userController := newUserController()

	userController.add(_core.Controller{
		URL:    VERSION + "/ws",
		Method: "GET",
		Handle: EstabilishConnection(epoll),
	})

	userController.add(_core.Controller{
		URL:    VERSION + "/ws/createUser",
		Method: "GET",
		Handle: CreateUserWithSocket(epoll),
	})

	userController.add(_core.Controller{
		URL:    VERSION + "/ws/writeToAnUser",
		Method: "GET",
		Handle: WriteMessageToAnUser(epoll),
	})

	userController.add(_core.Controller{
		URL:    VERSION + "/ws/writeToAll",
		Method: "GET",
		Handle: WriteMessageToAllUsers(epoll),
	})

	return userController.Controllers
}

func WsHandshake(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return nil, _core.InvalidConnectionError()
	}

	return conn, nil
}

func (ctl *UserControllers) add(controller _core.Controller) {
	ctl.Controllers[len(ctl.Controllers)] = controller
}

func newUserController() *UserControllers {
	controllers := UserControllers{
		Controllers: make(map[int]_core.Controller),
	}

	return &controllers
}

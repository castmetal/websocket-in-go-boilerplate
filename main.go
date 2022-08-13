package main

import (
	_config "websocket-in-go-boilerplate/src/config"
	_ports "websocket-in-go-boilerplate/src/ports"
)

func main() {
	switch _config.SystemParams.SERVER_TYPE {
	case "ws":
		_ports.InitWebsocketServer()
	}

}

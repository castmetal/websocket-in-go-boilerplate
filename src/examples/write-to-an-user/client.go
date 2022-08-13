package main

import (
	"log"
	"net/http"
	"os"

	_config "websocket-in-go-boilerplate/src/config"

	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

func main() {
	socketUrl := "ws://localhost:" + _config.SystemParams.PORT + "/v1/ws/writeToAnUser"

	userArgId := os.Args[1]
	if userArgId == "" {
		userArgId = "1"
	}

	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, http.Header{_config.SystemParams.AUTH_HEADER: []string{userArgId}})
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}

	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from GolangDocs To An User "+userArgId+"!"))
	if err != nil {
		log.Println("Error during writing to websocket:", err)
		return
	}
}

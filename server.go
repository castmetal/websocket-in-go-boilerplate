package main

/*
* This code was adapted from https://github.com/eranyanay/1m-go-websockets
 */

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"syscall"

	_userSocket "websocket-in-go-boilerplate/src/domains/user/controllers"
	_epoll "websocket-in-go-boilerplate/src/epoll"
	_utils "websocket-in-go-boilerplate/src/utils"
)

var epoller *_epoll.Epoll

func increaseResourcesLimitations() {
	// Increase resources limitations
	var rLimit syscall.Rlimit

	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
}

func enablePprofHooks() {
	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
		log.Fatalf("pprof failed: %v", err)
	}
}

func makeEpoll() {
	var err error

	epoller, err = _epoll.MkEpoll()
	if err != nil {
		panic(err)
	}
}

func main() {
	increaseResourcesLimitations()

	// Enable pprof hooks
	go enablePprofHooks()

	// Start epoll
	makeEpoll()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userSocketCfg := _userSocket.UserSocketController{
			Response: w,
			Request:  r,
			Epoll:    epoller,
		}

		userSocketCfg.SimpleSocket()
	})

	http.HandleFunc("/ws/writeToAll", func(w http.ResponseWriter, r *http.Request) {
		userSocketCfg := _userSocket.UserSocketController{
			Response: w,
			Request:  r,
			Epoll:    epoller,
		}

		userSocketCfg.WriteToAllClients()
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		io.WriteString(w, `{"alive": true}`)
	})

	log.Printf("Websocket Server listen at port :%s\n\n", _utils.SystemParams.PORT)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", _utils.SystemParams.PORT), nil); err != nil {
		log.Fatal(err)
	}

}

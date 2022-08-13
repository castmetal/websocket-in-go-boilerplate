package ports

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
	"time"

	_config "websocket-in-go-boilerplate/src/config"
	_rest_controllers "websocket-in-go-boilerplate/src/controllers"
	_interfaces "websocket-in-go-boilerplate/src/domains/common"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"
)

var epoller *_epoll.Epoll

func InitWebsocketServer() {
	increaseResourcesLimitations()

	// Enable pprof hooks
	go enablePprofHooks()

	// Start epoll
	makeEpoll()

	middlewares := []_interfaces.Middleware{
		_interfaces.PanicRecover,
		_interfaces.WithLogger,
	}

	mux := _interfaces.NewMyMux()

	mux.Use(middlewares...)

	_rest_controllers.SetRestControllers(mux, epoller)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		io.WriteString(w, `{"alive": true}`)
	})

	log.Printf("Websocket Server listen at port :%s\n\n", _config.SystemParams.PORT)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", _config.SystemParams.PORT),
		Handler:      mux,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

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

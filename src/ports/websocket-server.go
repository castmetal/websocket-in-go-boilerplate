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
	"runtime/debug"
	"syscall"
	"time"

	_config "websocket-in-go-boilerplate/src/config"
	_rest_controllers "websocket-in-go-boilerplate/src/controllers"
	_interfaces "websocket-in-go-boilerplate/src/domains/common"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"
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

func WithLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("path:%s process start...\n", r.URL.Path)
		defer func() {
			log.Printf("path:%s process end...\n", r.URL.Path)
		}()
		handler.ServeHTTP(w, r)
	})
}

func PanicRecover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(string(debug.Stack()))
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func InitWebsocketServer() {
	increaseResourcesLimitations()

	// Enable pprof hooks
	go enablePprofHooks()

	// Start epoll
	makeEpoll()

	middlewares := []_interfaces.Middleware{
		PanicRecover,
		WithLogger,
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

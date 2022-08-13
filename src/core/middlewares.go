package core

import (
	"log"
	"net/http"
	"runtime/debug"
)

type Middleware func(http.Handler) http.Handler

type MyMux struct {
	*http.ServeMux
	middlewares []Middleware
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

func NewMyMux() *MyMux {
	return &MyMux{
		ServeMux: http.NewServeMux(),
	}
}

func (m *MyMux) Use(middlewares ...Middleware) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *MyMux) Handle(pattern string, handler http.Handler) {
	handler = applyMiddlewares(handler, m.middlewares...)
	m.ServeMux.Handle(pattern, handler)
}

func (m *MyMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	newHandler := applyMiddlewares(handler, m.middlewares...)
	m.ServeMux.Handle(pattern, newHandler)
}

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

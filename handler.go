package gorest

import (
	"fmt"
	"net/http"
)

// HandlerFunc is the type of function we pass to the router to handle a route
type HandlerFunc func(http.ResponseWriter, *http.Request) (int, interface{})

// Middleware is the type of function we use for middlewares
type Middleware func(HandlerFunc) HandlerFunc

// Handler is the http.Handler we pass to the gowww.Router handler methods
// It takes care of executing middlewares
type Handler struct {
	Handle      HandlerFunc
	Middlewares []Middleware
	Write       WriteFunc
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handle := h.Handle

	for idx := len(h.Middlewares); idx > 0; idx-- {
		fmt.Println(len(h.Middlewares), idx-1, h.Middlewares)
		handle = h.Middlewares[idx-1](handle)
	}

	code, data := handle(w, r)

	h.Write(w, code, data)
}

// NewHandler creates a new handler
func NewHandler(write WriteFunc, handle HandlerFunc, middlewares []Middleware) Handler {
	return Handler{
		Write:       write,
		Handle:      handle,
		Middlewares: middlewares,
	}
}

package gorest

import (
	"net/http"

	"github.com/gowww/fatal"
	"github.com/gowww/log"
	"github.com/gowww/router"
)

// Router is the struct type for the router
type Router struct {
	rt *router.Router
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{router.New()}
}

type routerFunc func(...interface{}) (int, interface{})

// Get is used to declare a GET route
func (r *Router) Get(w WriteFunc, path string, h HandlerFunc, m ...Middleware) {
	r.rt.Get(path, NewHandler(w, h, m))
}

// Post is used to declare a POST route
func (r *Router) Post(w WriteFunc, path string, h HandlerFunc, m ...Middleware) {
	r.rt.Post(path, NewHandler(w, h, m))
}

// Put is used to declare a PUT route
func (r *Router) Put(w WriteFunc, path string, h HandlerFunc, m ...Middleware) {
	r.rt.Put(path, NewHandler(w, h, m))
}

// Delete is used to declare a DELETE route
func (r *Router) Delete(w WriteFunc, path string, h HandlerFunc, m ...Middleware) {
	r.rt.Delete(path, NewHandler(w, h, m))
}

// GetJSON is a helper for handling GET JSON
func (r *Router) GetJSON(path string, h HandlerFunc, m ...Middleware) {
	r.Get(JSON, path, h, m...)
}

// PostJSON is a helper for handling POST JSON
func (r *Router) PostJSON(path string, h HandlerFunc, m ...Middleware) {
	r.Post(JSON, path, h, m...)
}

// PutJSON is a helper for handling PUT JSON
func (r *Router) PutJSON(path string, h HandlerFunc, m ...Middleware) {
	r.Put(JSON, path, h, m...)
}

// DeleteJSON is a helper for handling DELETE JSON
func (r *Router) DeleteJSON(path string, h HandlerFunc, m ...Middleware) {
	r.Delete(JSON, path, h, m...)
}

// GetHandler is what you call to get the internal handler to give the server
func (r *Router) GetHandler() http.Handler {
	return log.Handle(fatal.Handle(r.rt, nil), nil)
}

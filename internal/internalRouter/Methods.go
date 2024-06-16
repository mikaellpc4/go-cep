package internalRouter

import (
	"net/http"
)

func (r *Router) GET(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodGet, pattern, handler)
}

func (r *Router) POST(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodPost, pattern, handler)
}

func (r *Router) PUT(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodPut, pattern, handler)
}

func (r *Router) DELETE(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodDelete, pattern, handler)
}

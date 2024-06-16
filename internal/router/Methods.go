package router

import "net/http"

func (r *router) GET(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodGet, pattern, handler)
}

func (r *router) POST(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodPost, pattern, handler)
}

func (r *router) PUT(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodPut, pattern, handler)
}

func (r *router) DELETE(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodDelete, pattern, handler)
}

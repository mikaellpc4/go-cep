package internalRouter

import (
	"context"
	"net/http"
	"strings"
)

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var allow []string

	responseWriter := NewResponseWriter(w)

	for _, route := range r.routes {
		matches := route.pattern.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 {
			if req.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			route.handler(
				*responseWriter,
				buildContext(req, route.paramKeys, matches[1:]),
			)
			return
		}
	}

	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	http.NotFound(w, req)
}

type ContextKey string

func buildContext(req *http.Request, paramKeys, paramValues []string) *http.Request {
	ctx := req.Context()
	for i := 0; i < len(paramKeys); i++ {
		ctx = context.WithValue(ctx, ContextKey(paramKeys[i]), paramValues[i])
	}

	return req.WithContext(ctx)
}

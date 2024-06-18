package internalRouter

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

type HandlerFunc = func(ResponseWriter, *http.Request)

type route struct {
	method       string
	pattern      *regexp.Regexp
	innerHandler HandlerFunc
	paramKeys    []string
}

type Router struct {
	routes []route
	prefix string
}

func NewRouter() *Router {
	return &Router{routes: []route{}}
}

func (r *route) handler(w ResponseWriter, req *http.Request) {
	requestString := fmt.Sprint(req.Method, " ", req.URL)

	fmt.Println("recieved", requestString)

	start := time.Now()
	r.innerHandler(w, req)
	w.Time = time.Since(start).Milliseconds()
}

func (r *Router) addRoute(method, endpoint string, handler HandlerFunc) {
	paramKeys := []string{}

	if endpoint == "/" {
		endpoint = ""
	}

	pathParamsPattern := regexp.MustCompile(":([a-z]+)")
	matches := pathParamsPattern.FindAllStringSubmatch(endpoint, -1)

	if r.prefix != "" {
		endpoint = r.prefix + endpoint
	}

	fmt.Printf("created endpoint: %s\n", endpoint)

	if len(matches) > 0 {
		endpoint = pathParamsPattern.ReplaceAllLiteralString(endpoint, "([^/]+)")

		for i := 0; i < len(matches); i++ {
			paramKeys = append(paramKeys, matches[i][1])
		}
	}

	route := route{method, regexp.MustCompile("^" + endpoint + "$"), handler, paramKeys}
	r.routes = append(r.routes, route)
}

func (router *Router) GROUP(prefix string, fn func(*Router)) {
	router.prefix = router.prefix + prefix

	fn(router)

	router.prefix = ""
}

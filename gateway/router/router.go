package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Route struct {
	Prefix  string
	Targets []string
	current int
}

type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) AddRoute(prefix string, targets []string) {
	r.routes = append(r.routes, &Route{
		Prefix:  prefix,
		Targets: targets,
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if strings.HasPrefix(req.URL.Path, route.Prefix) {
			// Round-robin load balancing
			target := route.Targets[route.current % len(route.Targets)]
			route.current++

			targetURL, err := url.Parse(target)
			if err != nil {
				http.Error(w, "Bad gateway", http.StatusBadGateway)
				return
			}

			proxy := httputil.NewSingleHostReverseProxy(targetURL)
			proxy.ServeHTTP(w, req)
			return
		}
	}
	http.Error(w, "Route not found", http.StatusNotFound);
}
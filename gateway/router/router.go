package router

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/rijav1001/API-Gateway/gateway/loadbalancer"
)

type Route struct {
	Prefix  string
	LB		*loadbalancer.LoadBalancer
}

type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) AddRoute(prefix string, targets []string) error {
	lb, err := loadbalancer.New(targets)
	if err != nil {
		return err
	}
	r.routes = append(r.routes, &Route{
		Prefix:  prefix,
		LB: lb,
	})
	return nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if strings.HasPrefix(req.URL.Path, route.Prefix) {
			// Round-robin load balancing
			target := route.LB.Next()
			proxy := httputil.NewSingleHostReverseProxy(target)
			proxy.ServeHTTP(w, req)
			return
		}
	}
	http.Error(w, "Route not found", http.StatusNotFound);
}
package nertivia

import "strings"

type Router struct {
	Routes map[string]func()
	Prefix string
}

func (r *Router) Add(route string, fun func()) {
	r.Routes[route] = fun
}

func NewRouter(prefix string) *Router {
	r := new(Router)
	r.Routes = make(map[string]func())
	r.Prefix = prefix
	return r
}

func (r *Router) Route(content string) {
	if !strings.HasPrefix(content, r.Prefix) {
		return
	}
	command := strings.TrimPrefix(content, r.Prefix)
	for route, handler := range r.Routes {
		if strings.HasPrefix(command, route) {
			handler()
		}
	}
}

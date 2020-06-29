package nertiviago

import "strings"

type Router struct {
	Routes []route
	Prefix string
}

type route struct {
	Name       string
	Handler    func()
	IgnoreCase bool
}

func (r *Router) Add(routeName string, ignoreCase bool, fun func()) {
	newRoute := new(route)
	newRoute.Name = routeName
	newRoute.Handler = fun
	newRoute.IgnoreCase = ignoreCase
	r.Routes = append(r.Routes, *newRoute)
}

func NewRouter(prefix string) *Router {
	r := new(Router)
	r.Prefix = prefix
	return r
}

func (r *Router) RemovePrefixAndCommand(command string) string {
	content := strings.Split(command, " ")
	if len(content) == 0 {
		return command
	}
	return strings.Join(content[1:len(content)], " ")
}

func (r *Router) GetRoutes() []string {
	var routes []string
	for _, route := range r.Routes {
		routes = append(routes, route.Name)
	}
	return routes
}

func (r *Router) Route(content string) {
	if !strings.HasPrefix(content, r.Prefix) {
		return
	}
	command := strings.TrimPrefix(content, r.Prefix)
	for _, route := range r.Routes {
		newCommand := command
		listenTo := route.Name
		if route.IgnoreCase {
			newCommand = strings.ToLower(newCommand)
			listenTo = strings.ToLower(listenTo)
		}
		if strings.HasPrefix(newCommand, listenTo) {
			route.Handler()
		}
	}
}

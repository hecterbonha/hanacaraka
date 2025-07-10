package routes

import (
	"net/http"

	"com.hanacaraka/application/services"
	"com.hanacaraka/interfaces/http/handlers"
	"github.com/gorilla/mux"
)

type RouteRegistry struct {
	router *mux.Router
}

func NewRouteRegistry(router *mux.Router) *RouteRegistry {
	return &RouteRegistry{
		router: router,
	}
}

func (rr *RouteRegistry) RegisterAllRoutes(userService services.UserServiceInterface) {
	userHandler := handlers.NewUserHandler(userService)

	routes := []struct {
		registerFunc func()
	}{
		{registerFunc: rr.registerBasicRoutes},
		{registerFunc: func() { rr.registerUserRoutes(userHandler) }},
		{registerFunc: rr.registerAPIRoutes},
		{registerFunc: rr.registerStaticRoutes},
	}

	for _, route := range routes {
		route.registerFunc()
	}
}

func (rr *RouteRegistry) registerBasicRoutes() {
	rr.router.HandleFunc("/", homeHandler).Methods("GET")
	rr.router.HandleFunc("/ping", pingHandler).Methods("GET")
}

func (rr *RouteRegistry) registerUserRoutes(userHandler *handlers.UserHandler) {
	SetupUserRoutes(rr.router, userHandler)
}

func (rr *RouteRegistry) registerAPIRoutes() {
	SetupAPIRoutes(rr.router)
}

func (rr *RouteRegistry) registerStaticRoutes() {
	staticDir := "./static/"
	rr.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
}

type RouteGroup struct {
	Name       string
	Prefix     string
	Middleware []mux.MiddlewareFunc
	Routes     []Route
}

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
	Name    string
}

func (rr *RouteRegistry) GetRegisteredRouteGroups() []RouteGroup {
	return []RouteGroup{
		{
			Name:   "Basic Routes",
			Prefix: "",
			Routes: []Route{
				{Path: "/", Method: "GET", Name: "Home"},
				{Path: "/ping", Method: "GET", Name: "Ping"},
			},
		},
		{
			Name:   "User Routes",
			Prefix: "/users",
			Routes: []Route{
				{Path: "/users", Method: "GET", Name: "GetUsers"},
				{Path: "/users", Method: "POST", Name: "CreateUser"},
				{Path: "/users/{id}", Method: "GET", Name: "GetUser"},
				{Path: "/users/{id}", Method: "PUT", Name: "UpdateUser"},
				{Path: "/users/{id}", Method: "DELETE", Name: "DeleteUser"},
			},
		},
		{
			Name:   "API v1 Routes",
			Prefix: "/api/v1",
			Routes: []Route{
				{Path: "/health", Method: "GET", Name: "Health"},
				{Path: "/status", Method: "GET", Name: "Status"},
				{Path: "/version", Method: "GET", Name: "Version"},
			},
		},
	}
}

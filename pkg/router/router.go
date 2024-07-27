package router

import (
	"context"
	"net/http"
)

type RouteType int

const (
	Get     = 1
	Post    = 2
	Put     = 3
	Delete  = 4
	Options = 5
)

type Route struct {
	Type    RouteType
	Pattern string
	Handler func(w http.ResponseWriter, r *http.Request)
	Roles   []string
}

type RouterConfig struct {
	Port              string `yaml:"port"`
	ReadTimeout       int    `yaml:"readTimeout"`
	ReadHeaderTimeout int    `yaml:"readHeaderTimeout"`
	WriteTimeout      int    `yaml:"writeTimeout"`
	IdleTimeout       int    `yaml:"idleTimeout"`
	BaseApi           string `yaml:"baseApi"`
	Https             bool   `yaml:"https"`
	HttpsPort         string `yaml:"httpsPort"`
	Domain            string `yaml:"domain"`
	CertFile          string `yaml:"certFile"`
	KeyFile           string `yaml:"keyFile"`
}

type Router interface {
	//Start the server
	ListenAndServe(context.Context) error
	//Add a route
	AddRoutes(context.Context, ...Route)
	//Set the router's config
	SetConfig(context.Context, RouterConfig)
	//Get a param in the url
	GetUrlParam(context.Context, *http.Request, string) string
	//Get a param in the url and parse it as uint, return 0 if error or not present
	GetUintParam(ctx context.Context, req *http.Request, paramName string) uint
}

// //Add a route that requires authentication
// AddAuthenticatedRoute(context.Context, ...Route)
// //Add a route that requires authentication as admin
// AddAdminRoute(context.Context, ...Route)

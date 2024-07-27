package router

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func NewMockRouter() Router {
	r := &MockRouter{}

	return r
}

type MockRouter struct {
	config         RouterConfig
	routes         []Route
	CurrentPattern string
}

func (r *MockRouter) AddRoutes(ctx context.Context, routes ...Route) {
	r.routes = append(r.routes, routes...)
}

var GetUrlParamFunc = func(pattern string, req *http.Request, key string) string {
	var reg *regexp.Regexp
	elem := fmt.Sprintf("\\{%s\\}", key)
	// split url in two parts, first part prefix elem and second part suffix elem in url
	reg, _ = regexp.Compile(elem)
	index := reg.FindStringIndex(pattern)
	if index == nil {
		return ""
	}

	prefix := pattern[:index[0]]
	suffix := pattern[index[1]:]

	varReg := regexp.MustCompile("{[a-zA-Z]+}")
	prefix = varReg.ReplaceAllString(prefix, "[0-9a-zA-Z\\-]+")
	suffix = varReg.ReplaceAllString(suffix, "[0-9a-zA-Z\\-]+")

	prefixReg := regexp.MustCompile(prefix)
	suffixReg := regexp.MustCompile(suffix)

	value := req.RequestURI
	value = prefixReg.ReplaceAllString(value, "")
	value = suffixReg.ReplaceAllString(value, "")

	return value
}

func (r *MockRouter) GetUrlParam(ctx context.Context, req *http.Request, key string) string {
	return GetUrlParamFunc(r.CurrentPattern, req, key)
}

func (r *MockRouter) GetUintParam(ctx context.Context, req *http.Request, paramName string) uint {
	//Get id in url
	idString := r.GetUrlParam(ctx, req, paramName)
	if idString == "" {
		return 0
	}

	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		return 0
	}

	return uint(id)
}

func (r *MockRouter) SetConfig(ctx context.Context, config RouterConfig) {
	r.config = config
}

func (r *MockRouter) ListenAndServe(ctx context.Context) error {
	return nil
}

// func (r *MockRouter) AddAuthenticatedRoute(ctx context.Context, routes ...Route) {
// 	r.routes = append(r.authenticatedRoutes, routes...)
// }

// func (r *MockRouter) AddAdminRoute(ctx context.Context, routes ...Route) {
// 	r.routes = append(r.adminRoutes, routes...)
// }
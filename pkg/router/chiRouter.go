package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cleoGitHub/golem/pkg/merror"

	"github.com/go-chi/chi"
)

func NewChiRouter(tokenValidator TokenValidator) Router {
	r := &chiRouter{
		m:              chi.NewRouter(),
		tokenValidator: tokenValidator,
	}
	r.m.Use(WithCors)

	return r
}

type chiRouter struct {
	tokenValidator TokenValidator
	config         RouterConfig
	m              *chi.Mux
	routes         []Route
}

func (r *chiRouter) AddRoutes(ctx context.Context, routes ...Route) {
	r.routes = append(r.routes, routes...)
}

func (r *chiRouter) GetUrlParam(ctx context.Context, req *http.Request, key string) string {
	return chi.URLParam(req, key)
}

func (r *chiRouter) GetUintParam(ctx context.Context, req *http.Request, paramName string) uint {
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

func (r *chiRouter) SetConfig(ctx context.Context, config RouterConfig) {
	r.config = config
}

func (r *chiRouter) ListenAndServe(ctx context.Context) error {
	fs := http.FileServer(http.Dir("assets"))
	r.m.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	//Add routes to router
	for _, route := range r.routes {
		r.m.Group(func(protectedRouter chi.Router) {
			protectedRouter.Use(r.AuthorizationMiddleware(ctx, route.Roles))
			switch route.Type {
			case Get:
				protectedRouter.Get(route.Pattern, route.Handler)
			case Post:
				protectedRouter.Post(route.Pattern, route.Handler)
			case Put:
				protectedRouter.Put(route.Pattern, route.Handler)
			case Delete:
				protectedRouter.Delete(route.Pattern, route.Handler)
			case Options:
				protectedRouter.Options(route.Pattern, route.Handler)
			}
		})
	}

	if r.config.Https {
		// redirect every http request to https
		go http.ListenAndServe(":"+r.config.Port, http.HandlerFunc(r.redirect))

		srv := &http.Server{
			Handler:           r.m,
			Addr:              ":" + r.config.HttpsPort,
			ReadTimeout:       time.Duration(r.config.ReadTimeout) * time.Second,
			ReadHeaderTimeout: time.Duration(r.config.ReadHeaderTimeout) * time.Second,
			WriteTimeout:      time.Duration(r.config.WriteTimeout) * time.Second,
			IdleTimeout:       time.Duration(r.config.IdleTimeout) * time.Second,
		}

		// mlogger.Logger().Infof("Starting server, port: %s", r.config.Port)
		err := srv.ListenAndServeTLS(r.config.CertFile, r.config.KeyFile)
		if err != nil {
			return merror.Stack(err)
		}
	} else {
		srv := &http.Server{
			Handler:           r.m,
			Addr:              ":" + r.config.Port,
			ReadTimeout:       time.Duration(r.config.ReadTimeout) * time.Second,
			ReadHeaderTimeout: time.Duration(r.config.ReadHeaderTimeout) * time.Second,
			WriteTimeout:      time.Duration(r.config.WriteTimeout) * time.Second,
			IdleTimeout:       time.Duration(r.config.IdleTimeout) * time.Second,
		}

		// mlogger.Logger().Infof("Starting server, port: %s", r.config.Port)
		err := srv.ListenAndServe()
		if err != nil {
			return merror.Stack(err)
		}
	}

	return nil
}

func (r *chiRouter) redirect(w http.ResponseWriter, req *http.Request) {
	redirection := fmt.Sprintf("https://%s:%s%s", r.config.Domain, r.config.HttpsPort, req.URL.String())
	// mlogger.Logger().Debugf("Redirecting http request: %s to %s", req.Host, redirection)
	http.Redirect(w, req, redirection, http.StatusMovedPermanently)
}

func (rtr *chiRouter) AuthorizationMiddleware(ctx context.Context, roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(roles) > 0 {
				bearer := r.Header.Get("Authorization")
				if bearer == "" {
					w.WriteHeader(http.StatusForbidden)
					return
				}

				token := extractTokenFromBearer(bearer)
				if token == "" {
					w.WriteHeader(http.StatusForbidden)
					return
				}

				ok, err := rtr.tokenValidator.ValidateAtLeastOneRole(ctx, token, roles)
				if err != nil {
					// mlogger.Logger().Errorf("unable to check authorization, err: %+v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractTokenFromBearer(s string) string {
	str := strings.Replace(s, "Bearer", "", 1)
	return strings.TrimSpace(str)
}

package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

func NewServerWithMiddleware() http.Handler {

	server := NewServer()

	r := http.NewServeMux()

	ctx := context.Background()

	// Middleware section
	// openapi validation middleware
	// loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	loader := openapi3.NewLoader()
	// TODO: get the path from configuration
	spec, err := loader.LoadFromFile("../../openapi/api.yaml")
	if err != nil {
		slog.Error("error", "msg", err.Error())
		return nil
	}
	// Validate document
	_ = spec.Validate(ctx)
	valmw := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			DoNotValidateServers: true,
		})

	// logging middleware
	logmw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				ip     = r.RemoteAddr
				method = r.Method
				url    = r.URL.String()
				proto  = r.Proto
			)

			userAttrs := slog.Group("user", "ip", ip)
			requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)

			slog.Info("request received", userAttrs, requestAttrs)
			next.ServeHTTP(w, r)
		})
	}
	// get an `http.Handler` that we can use
	h := HandlerWithOptions(server,
		StdHTTPServerOptions{
			BaseRouter:  r,
			Middlewares: []MiddlewareFunc{valmw, logmw},
		})

	return h
}

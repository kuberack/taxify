package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"kuberack.com/taxify/api"
)

func main() {
	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	server := api.NewServer()

	r := http.NewServeMux()

	ctx := context.Background()

	// Middlware section
	// openapi validation middleware
	// loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	loader := openapi3.NewLoader()
	spec, _ := loader.LoadFromFile("openapi/api.yaml")
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
	h := api.HandlerWithOptions(server,
		api.StdHTTPServerOptions{
			BaseRouter:  r,
			Middlewares: []api.MiddlewareFunc{valmw, logmw},
		})

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:8080",
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())

	fmt.Printf("help")
}

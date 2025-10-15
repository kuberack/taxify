package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"kuberack.com/taxify/internal/twilio_client"
)

var rootDir string

func init() {
	var err error
	if rootDir, err = os.Getwd(); err != nil {
		log.Fatal(err.Error())
	}
}

func NewServerWithMiddleware() (http.Handler, error) {

	// Get a twilio client
	t, err := twilio_client.GetTwilioClient()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return nil, errors.New("error: unable to get a Twilio Client")
	}
	// Inject twilio client into Server
	server := NewServer(t)

	r := http.NewServeMux()

	ctx := context.Background()

	// Middleware section
	// openapi validation middleware
	// loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	loader := openapi3.NewLoader()
	// TODO: get the path from configuration
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("error: runtime.Caller")
		return nil, errors.New("error: runtime Caller")
	}
	fmt.Println("Caller:", filename)
	var spec *openapi3.T
	if strings.HasPrefix(filename, "/app/") {
		// Docker environment
		fmt.Println("Environment: docker")
		spec, err = loader.LoadFromFile("openapi/api.yaml")
	} else {
		// baremetal environment
		fmt.Println("Environment: baremetal")
		spec, err = loader.LoadFromFile("../../openapi/api.yaml")
	}
	if err != nil {
		slog.Error("error", "msg", err.Error())
		return nil, err
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

	return h, nil
}

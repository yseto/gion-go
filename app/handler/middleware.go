package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/jmoiron/sqlx"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

func CreateMiddlewareEmptyContext(dbConn *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := emptyUserContext(r.Context())
			NewDBContext(ctx, dbConn)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CreateMiddleware(jwtSignedKeyBin []byte) (func(next http.Handler) http.Handler, error) {
	spec, err := GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(jwtSignedKeyBin),
			},
			ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
				if statusCode == 401 && strings.HasPrefix(message, "security requirements failed") {
					w.WriteHeader(statusCode)
					return
				}
				w.WriteHeader(statusCode)
				w.Write([]byte(message))
			},
		})

	return validator, nil
}

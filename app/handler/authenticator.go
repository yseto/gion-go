package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt"
)

var (
	ErrNoAuthHeader          = errors.New("authorization header is missing")
	ErrInvalidAuthHeader     = errors.New("authorization header is malformed")
	ErrSecuritySchemeInvalid = errors.New("security scheme is invalid")
)

const SessionContextKey = "user"
const JwtSignedKey = "secret"

type jwtCustomClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

// https://github.com/deepmap/oapi-codegen/blob/master/examples/authenticated-api/README.md あたりを参考にした

func GetJWTFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func NewAuthenticator() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(ctx, input)
	}
}

func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "BearerAuth" {
		return ErrSecuritySchemeInvalid
	}

	bearerToken, err := GetJWTFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jwt: %w", err)
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(JwtSignedKey), nil
	}

	t := reflect.ValueOf(&jwtCustomClaims{}).Type().Elem()
	claims := reflect.New(t).Interface().(*jwtCustomClaims)
	token, err := jwt.ParseWithClaims(bearerToken, claims, keyFunc)

	if err != nil {
		return err
	}
	if err == nil && !token.Valid {
		return ErrInvalidAuthHeader
	}

	tokenClaims := token.Claims.(*jwtCustomClaims)

	eCtx := middleware.GetEchoContext(ctx)
	eCtx.Set(SessionContextKey, tokenClaims.UserID)

	return nil
}

func GenerateToken(userid uint64) (string, error) {
	claims := &jwtCustomClaims{
		userid,
		jwt.StandardClaims{},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JwtSignedKey))
}

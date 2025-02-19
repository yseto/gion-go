package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNoAuthHeader          = errors.New("authorization header is missing")
	ErrInvalidAuthHeader     = errors.New("authorization header is malformed")
	ErrSecuritySchemeInvalid = errors.New("security scheme is invalid")
)

const SessionContextKey = "user"

type claim = jwt.RegisteredClaims

// https://github.com/oapi-codegen/oapi-codegen/blob/8bbe226927c98d11457cb125d3eaf82589022e7f/examples/authenticated-api/README.md あたりを参考にした

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

func NewAuthenticator(jwtSignedKey []byte) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(ctx, input, jwtSignedKey)
	}
}

func Authenticate(ctx context.Context, input *openapi3filter.AuthenticationInput, sign []byte) error {
	if input.SecuritySchemeName != "BearerAuth" {
		return ErrSecuritySchemeInvalid
	}

	bearerToken, err := GetJWTFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jwt: %w", err)
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return sign, nil
	}

	claims := &claim{}
	token, err := jwt.ParseWithClaims(bearerToken, claims, keyFunc, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return err
	}
	if !token.Valid {
		return ErrInvalidAuthHeader
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return err
	}

	userid, err := strconv.Atoi(sub)
	if err != nil {
		return err
	}

	newUserContext(ctx, uint64(userid))

	return nil
}

func GenerateToken(userid uint64, sign []byte) (string, error) {
	claims := &claim{
		Subject:   fmt.Sprint(userid),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(sign)
}

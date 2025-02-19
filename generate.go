 // +build generate

package main

//go:generate go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.0
//go:generate oapi-codegen -package handler -generate echo-server,strict-server -o app/handler/server.gen.go app/openapi.yaml
//go:generate oapi-codegen -package handler -generate spec -o app/handler/spec.gen.go app/openapi.yaml
//go:generate oapi-codegen -package handler -generate types -o app/handler/types.gen.go app/openapi.yaml


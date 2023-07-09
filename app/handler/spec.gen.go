// Package handler provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"BAAA//8cjAUsKTUAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

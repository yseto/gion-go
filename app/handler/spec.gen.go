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

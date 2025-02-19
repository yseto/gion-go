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

// Copyright 2021 Telefonica Cybersecurity & Cloud Tech SL
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/Telefonica/golium"
	"github.com/google/uuid"
	"github.com/tidwall/sjson"
	"github.com/xeipuuv/gojsonschema"
)

// Request information of the Session.
type Request struct {
	// Endpoint of the HTTP server. It might include a base path.
	Endpoint string
	// Path of the API endpoint. This path is considered with the endpoint to invoke the HTTP server.
	Path string
	// Query parameters
	QueryParams map[string][]string
	// Request headers
	Headers map[string][]string
	// HTTP method
	Method string
	// Request body as slice of bytes
	RequestBody []byte
}

// Response information of the session.
type Response struct {
	// HTTP response
	Response *http.Response
	// Response body as slice of bytes
	ResponseBody []byte
}

// Session contains the information of a HTTP session (request and response).
type Session struct {
	Request  Request
	Response Response
}

// URL composes the endpoint, the resource, and query parameters to build a URL.
func (s *Session) URL() (*url.URL, error) {
	u, err := url.Parse(s.Request.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("Invalid endpoint URL: %s. %s", s.Request.Endpoint, err)
	}
	u.Path = path.Join(u.Path, s.Request.Path)
	params := url.Values(s.Request.QueryParams)
	u.RawQuery = params.Encode()
	return u, nil
}

// ConfigureEndpoint configures the HTTP endpoint.
func (s *Session) ConfigureEndpoint(ctx context.Context, endpoint string) error {
	s.Request.Endpoint = endpoint
	return nil
}

// ConfigurePath configures the path of the HTTP endpoint.
// It configures a resource path in the application context.
// The API endpoint and the resource path are composed when invoking the HTTP server.
func (s *Session) ConfigurePath(ctx context.Context, path string) error {
	s.Request.Path = path
	return nil
}

// ConfigureQueryParams stores a table of query parameters in the application context.
func (s *Session) ConfigureQueryParams(ctx context.Context, params map[string][]string) error {
	s.Request.QueryParams = params
	return nil
}

// ConfigureHeaders stores a table of HTTP headers in the application context.
func (s *Session) ConfigureHeaders(ctx context.Context, headers map[string][]string) error {
	s.Request.Headers = headers
	return nil
}

// ConfigureRequestBodyJSONProperties writes the body in the HTTP request as a JSON with properties.
func (s *Session) ConfigureRequestBodyJSONProperties(ctx context.Context, props map[string]interface{}) error {
	var json string
	var err error
	for key, value := range props {
		if json, err = sjson.Set(json, key, value); err != nil {
			return fmt.Errorf("Error setting property '%s' with value '%s' in the request body. %s", key, value, err)
		}
	}
	return s.ConfigureRequestBodyJSONText(ctx, json)
}

// ConfigureRequestBodyJSONText writes the body in the HTTP request as a JSON from text.
func (s *Session) ConfigureRequestBodyJSONText(ctx context.Context, message string) error {
	s.Request.RequestBody = []byte(message)
	if s.Request.Headers == nil {
		s.Request.Headers = make(map[string][]string)
	}
	s.Request.Headers["Content-Type"] = []string{"application/json"}
	return nil
}

// SendHTTPRequest sends a HTTP request using the configuration in the application context.
func (s *Session) SendHTTPRequest(ctx context.Context, method string) error {
	logger := GetLogger()
	s.Request.Method = method
	corr := uuid.New().String()
	url, err := s.URL()
	if err != nil {
		return err
	}
	reqBodyReader := bytes.NewReader(s.Request.RequestBody)
	req, err := http.NewRequest(method, url.String(), reqBodyReader)
	if err != nil {
		return fmt.Errorf("Error creating the HTTP request with method: '%s' and url: '%s'. %s", method, url, err)
	}
	req.Header = s.Request.Headers
	logger.LogRequest(req, s.Request.RequestBody, corr)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error with the HTTP request. %s", err)
	}
	defer resp.Body.Close()
	// This is dangerous for big response bodies, but is read now to make sure that the body reader is closed.
	// TODO: limit the max size of the response body.
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading the response body. %s", err)
	}
	s.Response.Response = resp
	s.Response.ResponseBody = respBodyBytes
	logger.LogResponse(resp, respBodyBytes, corr)
	return nil
}

// ValidateStatusCode validates the status code from the HTTP response.
func (s *Session) ValidateStatusCode(ctx context.Context, expectedCode int) error {
	if expectedCode != s.Response.Response.StatusCode {
		return fmt.Errorf("Status code mismatch. Expected: %d, actual: %d", expectedCode, s.Response.Response.StatusCode)
	}
	return nil
}

// ValidateResponseBodyJSONSchema validates the response body against the JSON schema.
func (s *Session) ValidateResponseBodyJSONSchema(ctx context.Context, schema string) error {
	schemasDir := golium.GetConfig().Dir.Schemas
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/%s.json", schemasDir, schema))
	documentLoader := gojsonschema.NewStringLoader(string(s.Response.ResponseBody))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("Error validating response body against schema: %s. %s", schema, err)
	}
	if result.Valid() {
		return nil
	}
	return fmt.Errorf("Invalid response body according to schema: %s. %+v", schema, result.Errors())
}

// ValidateResponseBodyJSONProperties validates a list of properties in the JSON body of the HTTP response.
func (s *Session) ValidateResponseBodyJSONProperties(ctx context.Context, props map[string]interface{}) error {
	m := golium.NewMapFromJSONBytes(s.Response.ResponseBody)
	for key, expectedValue := range props {
		value := m.Get(key)
		if value != expectedValue {
			return fmt.Errorf("Mismatch of json property '%s'. Expected: '%s', actual: '%s'", key, expectedValue, value)
		}
	}
	return nil
}

// ValidateResponseBodyEmpty validates that the response body is empty.
// It checks the Content-Length header and the response body buffer.
func (s *Session) ValidateResponseBodyEmpty(ctx context.Context) error {
	if s.Response.Response.ContentLength == 0 && len(s.Response.ResponseBody) == 0 {
		return nil
	}
	return fmt.Errorf("The response body is not empty")
}

// StoreResponseBodyJSONPropertyInContext extracts a JSON property from the HTTP response body and stores it in the context.
func (s *Session) StoreResponseBodyJSONPropertyInContext(ctx context.Context, key string, ctxtKey string) error {
	m := golium.NewMapFromJSONBytes(s.Response.ResponseBody)
	value := m.Get(key)
	golium.GetContext(ctx).Put(ctxtKey, value)
	return nil
}

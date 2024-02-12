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
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/TelefonicaTC2Tech/golium/steps/http/model"
	"github.com/TelefonicaTC2Tech/golium/steps/http/schema"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/tidwall/sjson"
	"github.com/xeipuuv/gojsonschema"

	"encoding/json"
)

const (
	parameterError = "error getting parameter from json: %w"
	withJSONError  = "error sending http request using json: %w"
	InvalidPath    = "[CONF:apiKey.invalid_apiKey]"
	NilString      = "%nil%"
	Slash          = "/"
	DefaultTestURL = "https://jsonplaceholder.typicode.com/"
)

// Neutralize HTTP parameter pollution. CWE:235
func neutralize(p string) string {
	p = strings.ReplaceAll(p, "\r", "")
	p = strings.ReplaceAll(p, "\n", "")
	return p
}

// Session contains the information of a HTTP session (request and response).
type Session struct {
	Request            model.Request
	Response           model.Response
	NoRedirect         bool
	InsecureSkipVerify bool
	Timeout            time.Duration
	Timedout           bool
}

type RequestParams struct {
	Method     string
	URL        string
	Endpoint   string
	Path       string
	APIKey     string
	BodyParams schema.Params
	Table      *godog.Table
}

// URL composes the endpoint, the resource, and query parameters to build a URL.
func (s *Session) URL() (*url.URL, error) {
	u, _ := url.Parse(s.Request.Endpoint)
	if s.Request.Path != "" {
		u.Path = path.Join(u.Path, s.Request.Path)
	}

	// /*
	//  * NOTE: path.Join removes trailing slash using Clean thus,
	//  * we need to add it if is in s.Request.Path
	//  * - Reference: https://forum.golangbridge.org/t/how-to-concatenate-paths-for-api-request/5791
	//  * - Docs: https://pkg.go.dev/path#Join
	//  */
	params := url.Values{}
	for key, values := range s.Request.QueryParams {
		for _, value := range values {
			if !params.Has(key) {
				params.Add(key, value)
			}
		}
	}
	u.RawQuery = neutralize(params.Encode())

	return u, nil
}

// ConfigureEndpoint configures the HTTP endpoint.
func (s *Session) ConfigureEndpoint(ctx context.Context, endpoint string) {
	s.Request.Endpoint = endpoint
}

// SetHTTPResponseTimeout configures a response timeout in milliseconds.
func (s *Session) SetHTTPResponseTimeout(ctx context.Context, timeout int) {
	s.Timeout = time.Duration(timeout) * time.Millisecond
}

// ConfigurePath configures the path of the HTTP endpoint.
// It configures a resource path in the application context.
// The API endpoint and the resource path are composed when invoking the HTTP server.
func (s *Session) ConfigurePath(httpPath string) {
	s.Request.AddPath(httpPath)
}

// ConfigureQueryParams stores a table of query parameters in the application context.
func (s *Session) ConfigureQueryParams(params map[string][]string) {
	s.Request.AddQueryParams(params)
}

// ConfigureHeaders stores a table of HTTP headers in the application context.
func (s *Session) ConfigureHeaders(ctx context.Context, headers map[string][]string) {
	s.Request.Headers = headers
}

func (s *Session) ConfigureCredentials(ctx context.Context, username, password string) {
	s.Request.Username = username
	s.Request.Password = password
}

// ConfigureRequestBodyJSONProperties writes the body
// in the HTTP request as a JSON with properties.
func (s *Session) ConfigureRequestBodyJSONProperties(
	ctx context.Context,
	props map[string]interface{}) error {
	var jsonRequestBody string
	var err error
	for key, value := range props {
		if jsonRequestBody, err = sjson.Set(jsonRequestBody, key, value); err != nil {
			return fmt.Errorf("failed setting property '%s' with value '%s' in the request body: %w",
				key, value, err)
		}
	}
	s.ConfigureRequestBody(ctx, jsonRequestBody)
	return nil
}

// ConfigureRequestBodyJSONText writes the body in the
// HTTP request as a JSON from text.
func (s *Session) ConfigureRequestBody(ctx context.Context, message interface{}) {
	s.Request.AddBody(message)
	s.Request.AddJSONHeaders()
}

// ConfigureRequestBodyJSONFile writes the body in the HTTP request as a JSON from file.
func (s *Session) ConfigureRequestBodyJSONFile(
	ctx context.Context,
	bodyParams schema.Params,
) error {
	message, err := schema.GetParam(bodyParams, "body")
	if err != nil {
		return fmt.Errorf(parameterError, err)
	}
	s.ConfigureRequestBody(ctx, message)
	return nil
}

// ConfigureRequestBodyJSONFileWithout writes the body in the
// HTTP request as a JSON from file without given values.
func (s *Session) ConfigureRequestBodyJSONFileWithout(
	ctx context.Context,
	bodyParams schema.Params,
	params []string) error {
	message, err := schema.GetParam(bodyParams, "body")
	if err != nil {
		return fmt.Errorf(parameterError, err)
	}
	messageMap, _ := message.(map[string]interface{})
	for _, removeParams := range params {
		delete(messageMap, removeParams)
	}
	s.ConfigureRequestBody(ctx, message)
	return nil
}

// ConfigureRequestBodyURLEncodedProperties writes the body in the
// HTTP request as x-www-form-urlencoded with properties.
func (s *Session) ConfigureRequestBodyURLEncodedProperties(
	ctx context.Context,
	props map[string][]string) error {
	data := url.Values{}
	for k, s := range props {
		for _, v := range s {
			data.Add(k, v)
		}
	}
	s.Request.RequestBody = []byte(data.Encode())
	if s.Request.Headers == nil {
		s.Request.Headers = make(map[string][]string)
	}
	s.Request.Headers["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	return nil
}

// ConfigureNoRedirection configures no redirection for the HTTP client.
func (s *Session) ConfigureNoRedirection(ctx context.Context) {
	s.NoRedirect = true
}

// ConfigureInsecureSkipVerify configures insecure skip verify for the HTTP client in HTTPS calls.
func (s *Session) ConfigureInsecureSkipVerify(ctx context.Context) {
	s.InsecureSkipVerify = true
}

// SendHTTPRequest sends a HTTP request using the configuration in the application context.
func (s *Session) SendHTTPRequest(ctx context.Context, method string) error {
	logger := GetLogger()
	s.Request.Method = method
	corr := uuid.New().String()
	u, err := s.URL()
	if err != nil {
		return err
	}
	reqBody := s.Request.GetBody()

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed creating the HTTP request with method '%s' and url '%s'. %w",
			method, u, err)
	}
	if s.Request.Headers != nil {
		hostHeaders, found := s.Request.Headers["Host"]
		if found && len(hostHeaders) > 0 {
			req.Host = neutralize(hostHeaders[0])
		}
	}
	req.Header = s.Request.Headers
	if s.Request.Username != "" || s.Request.Password != "" {
		req.SetBasicAuth(s.Request.Username, s.Request.Password)
	}
	logger.LogRequest(req, s.Request.RequestBody, corr)
	client := http.Client{Timeout: s.Timeout}
	if s.NoRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	if s.InsecureSkipVerify { // #nosec G402
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}
	resp, err := client.Do(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			logger.LogTimeout(corr)
			s.Timedout = true
			return nil
		}
		return fmt.Errorf("error with the HTTP request. %w", err)
	}
	defer resp.Body.Close()
	// This is dangerous for big response bodies,
	// but is read now to make sure that the body reader is closed.
	// TODO: limit the max size of the response body.
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed reading the response body: %w", err)
	}
	s.Response.HTTPResponse = resp
	s.Response.ResponseBody = respBodyBytes
	logger.LogResponse(resp, respBodyBytes, corr)
	return nil
}

// ValidateResponseTimedout checks if the HTTP client timed out without
// receiving a response.
func (s *Session) ValidateResponseTimedout(ctx context.Context) error {
	if !s.Timedout {
		return errors.New("no timed out")
	}
	return nil
}

// ValidateStatusCode validates the status code from the HTTP response.
func (s *Session) ValidateStatusCode(ctx context.Context, expectedCode int) error {
	if expectedCode != s.Response.HTTPResponse.StatusCode {
		return fmt.Errorf("status code mismatch: expected '%d', actual '%d'",
			expectedCode, s.Response.HTTPResponse.StatusCode)
	}
	return nil
}

// ValidateResponseHeaders checks a set of response headers.
func (s *Session) ValidateResponseHeaders(
	ctx context.Context,
	expectedHeaders map[string][]string) error {
	for expectedHeader, expectedHeaderValues := range expectedHeaders {
		for _, expectedHeaderValue := range expectedHeaderValues {
			if !golium.ContainsString(
				expectedHeaderValue, s.Response.HTTPResponse.Header.Values(expectedHeader)) {
				return fmt.Errorf("HTTP response does not have the header '%s' with value '%s'",
					expectedHeader, expectedHeaderValue)
			}
		}
	}
	return nil
}

// ValidateNotResponseHeaders checks that a set of
// response headers are not included in HTTP response.
func (s *Session) ValidateNotResponseHeaders(ctx context.Context, expectedHeaders []string) error {
	for _, expectedHeader := range expectedHeaders {
		if len(s.Response.HTTPResponse.Header.Values(expectedHeader)) > 0 {
			return fmt.Errorf("HTTP response includes the header '%s'", expectedHeader)
		}
	}
	return nil
}

// ValidateResponseBodyJSONSchema validates the response body against the JSON schema.
func (s *Session) ValidateResponseBodyJSONSchema(ctx context.Context, schemaName string) error {
	schemasPath := golium.GetConfig().Dir.Schemas
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/%s.json",
		schemasPath, schemaName))
	documentLoader := gojsonschema.NewStringLoader(string(s.Response.ResponseBody))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("failed validating response body against schema '%s': %w", schemaName, err)
	}
	if result.Valid() {
		return nil
	}
	return fmt.Errorf(
		"invalid response body according to schema '%s': %+v", schemaName, result.Errors())
}

// ValidateResponseFromJSONFile validates the response body against the response from JSON File.
func (s *Session) ValidateResponseFromJSONFile(
	response interface{},
	respDataLocation string) error {
	respBody := s.Response.ResponseBody
	if respDataLocation != "" {
		respBodyAux := golium.NewMapFromJSONBytes(respBody)
		respBodyDataLoc := respBodyAux.Get(respDataLocation)
		respBody = []byte(fmt.Sprint(respBodyDataLoc))
	}

	switch resp := response.(type) {
	case string:
		if string(respBody) != response {
			return fmt.Errorf("received body does not match expected, \n%s\n vs \n%s", response,
				respBody)
		}
	case map[string]interface{}:

		var realResponse interface{}

		if err := json.Unmarshal(respBody, &realResponse); err != nil {
			return fmt.Errorf("error unmarshalling response body: %w", err)
		}

		if !reflect.DeepEqual(response, realResponse) {
			return fmt.Errorf("expected JSON does not match real response, \n%v\n vs \n%s", response,
				realResponse)
		}
	default:
		return fmt.Errorf("body content should be string or map: %v", resp)
	}
	return nil
}

// ValidateResponseBodyJSONFile validates the response body against the JSON in File.
func (s *Session) ValidateResponseBodyJSONFile(
	ctx context.Context,
	responseParams schema.Params,
	respDataLocation string) error {
	jsonResponseBody, err := schema.GetParam(responseParams, "response")
	if err != nil {
		return fmt.Errorf(parameterError, err)
	}
	return s.ValidateResponseFromJSONFile(jsonResponseBody, respDataLocation)
}

// ValidateResponseBodyJSONFileWithout validates
// the response body against the JSON in File without params.
func (s *Session) ValidateResponseBodyJSONFileWithout(
	ctx context.Context,
	responseParams schema.Params,
	respDataLocation string, t *godog.Table) error {
	jsonResponseBody, err := schema.DeleteResponseFields(ctx, responseParams, t)
	if err != nil {
		return fmt.Errorf(
			"error deleting response %s fields from %s schema: %w",
			responseParams.Code, responseParams.File, err)
	}
	return s.ValidateResponseFromJSONFile(jsonResponseBody, respDataLocation)
}

// ValidateResponseBodyJSONFileModifying validates the response body
// against the JSON in File modifying params.
func (s *Session) ValidateResponseBodyJSONFileModifying(
	ctx context.Context, responseParams schema.Params,
	t *godog.Table,
) error {
	jsonResponseBody, err := schema.ModifyResponse(ctx, responseParams, t)
	if err != nil {
		return fmt.Errorf(
			"error modifying response %s from %s schema: %w",
			responseParams.Code, responseParams.File, err)
	}
	return s.ValidateResponseFromJSONFile(jsonResponseBody, "")
}

// ValidateResponseBodyJSONProperties validates a list
// of properties in the JSON body of the HTTP response.
func (s *Session) ValidateResponseBodyJSONProperties(
	ctx context.Context,
	props map[string]interface{}) error {
	m := golium.NewMapFromJSONBytes(s.Response.ResponseBody)
	for key, expectedValue := range props {
		value := m.Get(key)
		if value != expectedValue {
			return fmt.Errorf("mismatch of json property '%s': expected '%s', actual '%s'",
				key, expectedValue, value)
		}
	}
	return nil
}

// ValidateResponseBodyEmpty validates that the response body is empty.
// It checks the Content-Length header and the response body buffer.
func (s *Session) ValidateResponseBodyEmpty(ctx context.Context) error {
	if s.Response.HTTPResponse.ContentLength <= 0 && len(s.Response.ResponseBody) == 0 {
		return nil
	}
	return errors.New("response body is not empty")
}

// ValidateResponseBodyText validates that the response body payload is the expected text.
func (s *Session) ValidateResponseBodyText(ctx context.Context, expectedText string) error {
	if expectedText == string(s.Response.ResponseBody) {
		return nil
	}
	return fmt.Errorf("response payload: '%v' is not the expected: '%s'",
		s.Response.ResponseBody, expectedText)
}

// StoreResponseBodyJSONPropertyInContext extracts a JSON property from
// the HTTP response body and stores it in the context.
func (s *Session) StoreResponseBodyJSONPropertyInContext(
	ctx context.Context, key, ctxtKey string) error {
	m := golium.NewMapFromJSONBytes(s.Response.ResponseBody)
	value := m.Get(key)
	golium.GetContext(ctx).Put(ctxtKey, value)
	return nil
}

// StoreResponseHeaderInContext stores in context a header of the HTTP response.
// If the header does not exist, the context value is empty.
// This method does not support multiple headers with the same name. It just stores one of them.
func (s *Session) StoreResponseHeaderInContext(ctx context.Context, header, ctxtKey string) error {
	h := s.Response.HTTPResponse.Header.Get(header)
	golium.GetContext(ctx).Put(ctxtKey, h)
	return nil
}

// SendRequestWithBody send request using body from JSON file located in schemas.
func (s *Session) SendRequestWithBody(
	ctx context.Context,
	uRL, method, endpoint string, bodyParams schema.Params, apiKey string,
) error {
	// Build request
	s.Request = model.NewRequest(method, uRL, endpoint, true)
	// Configure request JSON Body
	message, err := schema.GetBody(bodyParams)
	if err != nil {
		return fmt.Errorf("error getting body: %w", err)
	}
	s.Request.AddBody(message)
	// Configure authorization headers
	s.Request.AddAuthorization(apiKey, "")
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// SendRequestWithBodyWithoutFields send request using body from JSON file located in schemas
// without fields.
func (s *Session) SendRequestWithBodyWithoutFields(
	ctx context.Context,
	uRL, method, endpoint string, bodyParams schema.Params, apiKey string, t *godog.Table,
) error {
	// Configure request JSON Body
	message, err := schema.DeleteBodyFields(ctx, bodyParams, t)
	if err != nil {
		return fmt.Errorf("error deleting body fields: %w", err)
	}
	// Build request
	s.Request = setRequestWithBodyMessage(method, uRL, endpoint, apiKey, message)
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// SendRequestWithBodyModifyingFields send request using body from JSON file located in schemas
// modifying fields.
func (s *Session) SendRequestWithBodyModifyingFields(
	ctx context.Context,
	uRL, method, endpoint string, bodyParams schema.Params, apiKey string, t *godog.Table,
) error {
	// Configure request JSON Body
	message, err := schema.ModifyBody(ctx, bodyParams, t)
	if err != nil {
		return fmt.Errorf("error modifying body fields: %w", err)
	}
	// Build request
	s.Request = setRequestWithBodyMessage(method, uRL, endpoint, apiKey, message)
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

func setRequestWithBodyMessage(
	method, uRL, endpoint, apiKey string,
	message interface{},
) model.Request {
	request := model.NewRequest(method, uRL, endpoint, true)
	request.AddBody(message)
	// Configure authorization headers
	request.AddAuthorization(apiKey, "")
	return request
}

// SendRequestWithQueryParams send request using with query params.
func (s *Session) SendRequestWithQueryParams(
	ctx context.Context,
	uRL, method, endpoint, apiKey string,
	t *godog.Table,
) error {
	// Build request
	s.Request = model.NewRequest(method, uRL, endpoint, true)
	// Configure authorization headers
	s.Request.AddAuthorization(apiKey, "")
	// Configure Query Params
	params, err := golium.ConvertTableToMultiMap(ctx, t)
	if err != nil {
		return err
	}
	s.Request.AddQueryParams(params)
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// SendRequestWithFilters send request using filters with query params.
func (s *Session) SendRequestWithFilters(
	ctx context.Context,
	uRL, method, endpoint, apiKey, filters string,
) error {
	// Build request
	s.Request = model.NewRequest(method, uRL, endpoint, true)
	// Configure authorization headers
	s.Request.AddAuthorization(apiKey, "")
	// Configure Query Params from filters
	var params url.Values
	params, _ = url.ParseQuery(filters)
	s.Request.AddQueryParams(params)
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// SendRequestWithPath send request with path.
func (s *Session) SendRequestWithPath(
	ctx context.Context,
	uRL, method, endpoint, requestPath, apiKey string,
) error {
	// Build request
	s.Request = model.NewRequest(method, uRL, endpoint, true)
	// Configure authorization headers
	s.Request.AddAuthorization(apiKey, "")
	s.Request.AddPath(requestPath)
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// SendRequestWithPathAndBody send request with path and JSON body.
func (s *Session) SendRequestWithPathAndBody(
	ctx context.Context,
	uRL, method, endpoint, requestPath string, bodyParams schema.Params, apiKey string,
) error {
	// Build request
	s.Request = model.NewRequest(method, uRL, endpoint, true)
	// Configure request JSON Body
	message, err := schema.GetParam(bodyParams, "body")
	if err != nil {
		return fmt.Errorf(parameterError, err)
	}
	s.Request.AddBody(message)
	// Configure authorization headers
	s.Request.AddAuthorization(apiKey, "")
	s.Request.AddPath(requestPath)
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// SendRequestWithoutBackslash send request without backslash.
func (s *Session) SendRequestWithoutBackslash(
	ctx context.Context,
	uRL, method, endpoint, apiKey string,
) error {
	// Build request
	s.Request = model.NewRequest(method, uRL, endpoint, false)
	// Configure authorization headers
	s.Request.AddAuthorization(apiKey, "")
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// SendRequestWithMultipartBody send request using multipart body.
func (s *Session) SendRequestWithMultipartBody(
	ctx context.Context,
	reqParams RequestParams, fileField, file string,
) error {
	// Build request
	s.Request = model.NewRequest(reqParams.Method, reqParams.URL, reqParams.Endpoint, true)
	s.Request.AddPath(reqParams.Path)
	// Add form fields
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add file field if not empty names
	if err := schema.WriteMultipartBodyFile(w, fileField, file); err != nil {
		return fmt.Errorf("error writing multipart body file: %w", err)
	}
	// Add form fields if they are included
	if err := schema.WriteMultipartBodyFields(ctx, w, reqParams.Table); err != nil {
		return fmt.Errorf("error writing multipart body fields: %w", err)
	}
	// Close multipart writer
	w.Close()
	s.Request.AddMultipartBody(b)
	// Set multipart content type
	s.Request.SetContentType(w.FormDataContentType())
	// Configure authorization headers
	s.Request.AddAuthorization(reqParams.APIKey, "")
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, reqParams.Method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// SendRequestWithoutBackslash send request without backslash.
func (s *Session) SendRequest(
	ctx context.Context,
	uRL, method, endpoint, apiKey string,
) error {
	// Build request
	s.Request = model.NewRequest(method, uRL, endpoint, true)
	// Configure authorization headers
	s.Request.AddAuthorization(apiKey, "")
	// Send HTTP Request
	if err := s.SendHTTPRequest(ctx, method); err != nil {
		return fmt.Errorf(withJSONError, err)
	}
	return nil
}

// GetURL returns URL from Configuration or Context
func (s *Session) GetURL(ctx context.Context) (string, error) {
	URL := golium.ValueAsString(ctx, "[CONF:url]")
	if URL == "<nil>" {
		URL = golium.ValueAsString(ctx, "[CTXT:url]")
	}
	if URL == NilString {
		return "", fmt.Errorf("url shall be initialized in Configuration or Context")
	}
	return URL, nil
}

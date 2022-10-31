package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	Slash = "/"
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
	// Multipart body
	MultipartBody *bytes.Buffer
	// Username for basic authentication
	Username string
	// Password for basic authentication
	Password string
}

func NewRequest(
	method, url, endpoint string,
	backslash bool,
) Request {
	var request = Request{}
	request.Headers = make(map[string][]string)
	request.Headers["Content-Type"] = []string{"application/json"}
	request.Method = method
	request.Endpoint = url + NormalizeEndpoint(endpoint, backslash)
	return request
}

func (r *Request) SetContentType(contentType string) {
	r.Headers["Content-Type"] = []string{contentType}
}

func (r *Request) AddBody(message interface{}) {
	stringMessage := reflect.TypeOf(message)
	if stringMessage.Kind() == reflect.String {
		r.RequestBody = []byte(message.(string))
		return
	}
	r.RequestBody, _ = json.Marshal(message)
}

func (r *Request) AddMultipartBody(mBody bytes.Buffer) {
	r.MultipartBody = &mBody
}

func (r *Request) GetBody() io.Reader {
	var readBody io.Reader
	if r.RequestBody != nil {
		readBody = bytes.NewReader(r.RequestBody)
	}
	if r.MultipartBody != nil {
		readBody = r.MultipartBody
	}
	return readBody
}

func (r *Request) AddAuthorization(apiKey, jwtValue string) {
	if apiKey != "" {
		r.Headers["X-API-KEY"] = []string{apiKey}
	} else {
		delete(r.Headers, "X-API-KEY")
	}
	if jwtValue != "" {
		r.Headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", jwtValue)}
	} else {
		delete(r.Headers, "Authorization")
	}
}

func (r *Request) AddQueryParams(params map[string][]string) {
	r.QueryParams = params
}

func (r *Request) AddPath(path string) {
	r.Path = path
}

// AddJSONHeaders adds json headers to Request if they are null
func (r *Request) AddJSONHeaders() {
	if r.Headers == nil {
		r.Headers = make(map[string][]string)
	}
	r.Headers["Content-Type"] = []string{"application/json"}
}

// normalizeEndpoint Normalize Endpoint considering ending backslash need.
func NormalizeEndpoint(endpoint string, backslash bool) string {
	if !backslash {
		return strings.TrimSuffix(endpoint, Slash)
	}
	if strings.HasSuffix(endpoint, Slash) {
		return endpoint
	}
	return endpoint + Slash
}

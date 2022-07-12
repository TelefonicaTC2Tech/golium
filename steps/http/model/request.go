package model

import (
	"encoding/json"
	"fmt"
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
	// Username for basic authentication
	Username string
	// Password for basic authentication
	Password string
}

func NewRequest(
	method, url, endpoint, path string,
	backslash bool,
) Request {
	var request = Request{}
	request.Headers = make(map[string][]string)
	request.Headers["Content-Type"] = []string{"application/json"}
	request.Method = method
	request.Endpoint = url + normalizeEndpoint(endpoint, endpoint, backslash)
	request.RequestBody = nil
	request.Path = path

	return request
}

func (r *Request) AddBody(message interface{}) {
	r.RequestBody, _ = json.Marshal(message)
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

// normalizeEndpoint Normalize Endpoint considering ending backslash need.
func normalizeEndpoint(endpoint, request string, backslash bool) string {
	if !backslash {
		return strings.TrimSuffix(endpoint, Slash)
	}
	if strings.HasSuffix(endpoint, Slash) {
		return endpoint
	}
	return endpoint + Slash
}

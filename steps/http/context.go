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
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// ClientContextKey defines a type to store the ClientContext in context.Context.
type ClientContextKey string

var clientContextKey ClientContextKey = "clientContextKey"

// ClientContext contains the context for the HTTP client steps (e.g. to validate the response).
type ClientContext struct {
	Request struct {
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
	Response struct {
		// HTTP response
		Response *http.Response
		// Response body as slice of bytes
		ResponseBody []byte
	}
}

// InitializeClientContext adds the ClientContext to the context.
// The new context is returned because context is immutable.
func InitializeClientContext(ctx context.Context) context.Context {
	var clientContext ClientContext
	return context.WithValue(ctx, clientContextKey, &clientContext)
}

// GetClientContext returns the ClientContext stored in context.
// Note that the context should be previously initialized with InitializeClientContext function.
func GetClientContext(ctx context.Context) *ClientContext {
	return ctx.Value(clientContextKey).(*ClientContext)
}

// URL composes the endpoint, the resource, and query parameters to build a URL.
func (cc ClientContext) URL() (*url.URL, error) {
	u, err := url.Parse(cc.Request.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("Invalid endpoint URL: %s. %s", cc.Request.Endpoint, err)
	}
	u.Path = path.Join(u.Path, cc.Request.Path)
	params := url.Values(cc.Request.QueryParams)
	u.RawQuery = params.Encode()
	return u, nil
}

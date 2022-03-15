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

import "encoding/json"

// MockRequest contains the instruction to configure the behavior of the HTTP mock server.
// The document configures which request is going to be attended (e.g. the path and method)
// and the response to be generated by the mock.
type MockRequest struct {
	// Permanent is true if the configuration is permanent.
	// If permanent is false, the mockRequest is removed after matching the first request.
	Permanent bool     `json:"permanent"`
	Request   Request  `json:"request"`
	Response  Response `json:"response"`
	// Latency is the duration in milliseconds to wait to deliver the response.
	// If 0, there is no latency to apply.
	// If negative, there will be no response (timeout simulation).
	Latency int `json:"latency"`
}

// Request configures the filter for the request of the MockRequest.
type Request struct {
	Method  string              `json:"method,omitempty"`
	Path    string              `json:"path,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

// Response configures which response if the request filter applies.
type Response struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func (m MockRequest) String() string {
	b, _ := json.Marshal(&m)
	return string(b)
}

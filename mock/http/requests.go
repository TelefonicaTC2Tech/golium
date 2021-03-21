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
	"net/http"

	"github.com/Telefonica/golium"
)

type MockRequests struct {
	mockRequests []*MockRequest
}

// PushMockRequest adds a MockRequest to the list.
func (m *MockRequests) PushMockRequest(mockRequest *MockRequest) {
	m.mockRequests = append(m.mockRequests, mockRequest)
}

func (m *MockRequests) MatchMockRequest(r *http.Request) *MockRequest {
	for _, mockRequest := range m.mockRequests {
		if matchMockRequest(r, mockRequest) {
			return mockRequest
		}
	}
	return nil
}

func matchMockRequest(r *http.Request, mockRequest *MockRequest) bool {
	mr := mockRequest.Request
	if mr.Method != "" && r.Method != mr.Method {
		return false
	}
	if mr.Path != "" && r.URL.Path != mr.Path {
		return false
	}
	if len(mr.Headers) != 0 {
		for header, values := range mr.Headers {
			rValues := r.Header[header]
			for _, value := range values {
				if !golium.ContainsString(value, rValues) {
					return false
				}
			}
		}
	}
	return true
}

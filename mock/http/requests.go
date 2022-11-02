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
	"regexp"
	"strings"
	"sync"

	"github.com/TelefonicaTC2Tech/golium"
)

type MockRequests struct {
	mockRequests []*MockRequest
	mutex        sync.Mutex
}

// PushMockRequest adds a MockRequest to the list.
func (m *MockRequests) PushMockRequest(mockRequest *MockRequest) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.mockRequests = append(m.mockRequests, mockRequest)
}

// CleanMockRequests removes all the mockRequests from the list.
func (m *MockRequests) CleanMockRequests() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.mockRequests = nil
}

// MatchMockRequest finds the first mockRequest matching the HTTP request.
func (m *MockRequests) MatchMockRequest(r *http.Request) *MockRequest {
	for _, mockRequest := range m.mockRequests {
		if matchMockRequest(r, mockRequest) {
			return mockRequest
		}
	}
	return nil
}

// Remove a mockRequest. It returns true if it was found and removed.
func (m *MockRequests) RemoveMockRequest(mockRequest *MockRequest) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for i, mr := range m.mockRequests {
		if mr != mockRequest {
			continue
		}
		m.mockRequests = append(m.mockRequests[:i], m.mockRequests[i+1:]...)
		return true
	}
	return false
}

func matchMockRequest(r *http.Request, mockRequest *MockRequest) bool {
	mr := mockRequest.Request
	if mr.Method != "" && r.Method != mr.Method {
		return false
	}
	if !matchPath(r, mockRequest) {
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

var defaultPattern = regexp.MustCompile(`<\*>`)

func matchPath(r *http.Request, mockRequest *MockRequest) bool {
	mr := mockRequest.Request
	if mr.Path == "" {
		return false
	}
	corePath := strings.TrimSuffix(mr.Path, "<*>")
	if defaultPattern.FindString(mr.Path) != "" && strings.HasPrefix(r.URL.Path, corePath) {
		return true
	}
	if mr.Path != r.URL.Path {
		return false
	}
	return true
}

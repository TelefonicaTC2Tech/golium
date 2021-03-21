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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Server type for the HTTP mock server.
type Server struct {
	Port         int
	mockRequests MockRequests
	logger       *logrus.Entry
}

// NewServer creates an instance of Server.
func NewServer(port int) *Server {
	return &Server{
		Port:         port,
		mockRequests: MockRequests{},
		logger:       logrus.WithField("mock", "http"),
	}
}

// Start the HTTP mock server.
// Note that it block the current goroutine with http.ListenAndServe function.
func (s *Server) Start() error {
	http.HandleFunc("/_mock/requests", s.handleMockRequest)
	http.HandleFunc("/", s.handle)
	addr := fmt.Sprintf(":%d", s.Port)
	s.logger.Infof("Starting server at '%s'", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleMockRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var mockRequest MockRequest
		if err := json.NewDecoder(r.Body).Decode(&mockRequest); err != nil {
			s.logger.Errorf("Failed decoding mockRequest: %s", err)
			return
		}
		s.logger.Infof("Pushing mockRequest: %s", mockRequest)
		s.mockRequests.PushMockRequest(&mockRequest)
	case http.MethodDelete:
		s.mockRequests.CleanMockRequests()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	mockRequest := s.mockRequests.MatchMockRequest(r)
	if mockRequest == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if !mockRequest.Permanent {
		s.mockRequests.RemoveMockRequest(mockRequest)
	}
	if mockRequest.Latency < 0 {
		s.logger.Info("Simulating response timeout")
		time.Sleep(time.Duration(5) * time.Minute)
		return
	}
	if mockRequest.Latency > 0 {
		time.Sleep(time.Duration(mockRequest.Latency) * time.Millisecond)
	}
	resp := mockRequest.Response
	for header, values := range resp.Headers {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(resp.Status)
	if _, err := w.Write([]byte(resp.Body)); err != nil {
		s.logger.Errorf("Failed writing the response body: %s", err)
	}
}

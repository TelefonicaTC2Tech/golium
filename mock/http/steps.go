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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
)

// Steps to initialize common steps.
type Steps struct {
}

// InitializeSteps initializes all the steps.
func (cs Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	scenCtx.Step(`^I mock the HTTP request at "([^"]*)" for path "([^"]*)" with status "(\d+)" and JSON body$`, func(server, path string, status int, message *godog.DocString) error {
		if http.StatusText(status) == "" {
			return fmt.Errorf("status code to return not valid: %d", status)
		}
		return MockRequestSimple(ctx, golium.ValueAsString(ctx, server), golium.ValueAsString(ctx, path), status, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I mock the HTTP request at "([^"]*)" with the JSON$`, func(server string, body *godog.DocString) error {
		var mockRequest MockRequest
		if err := json.Unmarshal([]byte(golium.ValueAsString(ctx, body.Content)), &mockRequest); err != nil {
			return fmt.Errorf("failed unmarshalling to mockRequest: %w", err)
		}
		return sendMockRequest(ctx, golium.ValueAsString(ctx, server), &mockRequest)
	})
	return ctx
}

// MockRequestSimple mocks a HTTP request by sending a command to the HTTP mock server.
// It only configures the request path, the response status code and the response body.
// The "application/json" content type is applied implicitly to the response.
func MockRequestSimple(ctx context.Context, server, path string, status int, body string) error {
	mockRequest := &MockRequest{
		Request: Request{
			Method: http.MethodGet,
			Path:   path,
		},
		Response: Response{
			Status: status,
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body: body,
		},
	}
	return sendMockRequest(ctx, server, mockRequest)
}

func sendMockRequest(ctx context.Context, server string, mockRequest *MockRequest) error {
	if mockRequest.Latency < 0 {
		return fmt.Errorf("invalid latency: %d", mockRequest.Latency)
	}
	u, err := url.Parse(server)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %s: %w", server, err)
	}
	u.Path = path.Join(u.Path, "/_mock/requests")
	body, err := json.Marshal(mockRequest)
	if err != nil {
		return fmt.Errorf("failed marshalling mockRequest to json: %w", err)
	}
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed sending mockRequest to mock server: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

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
	"net/url"
	"strings"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
)

const (
	configRequestBodyError = "error configuring request body: "
	sendingRequestError    = "error sending request: "
	confEndpoint           = "[CONF:endpoints.%s.api-key]"
)

// SendRequest send request using credential parameters if needed.
func (s *Session) SendRequest(
	ctx context.Context,
	method, endpoint, apiKey, jwtValue string,
) error {
	if err := s.ConfigureURLAndEndpoint(ctx, endpoint); err != nil {
		return err
	}

	if err := s.ConfigureRequestHeaders(ctx, apiKey, jwtValue); err != nil {
		return err
	}

	if err := s.GoliumInterface.SendHTTPRequest(ctx, method); err != nil {
		return err
	}
	return nil
}

// SendRequestTo send request using credential parameters if they are defined in Configuration.
func (s *Session) SendRequestTo(ctx context.Context, method, request string) error {
	apiKey := s.GoliumInterface.ValueAsString(ctx, fmt.Sprintf(confEndpoint, request))
	endpoint := s.NormalizeEndpoint(ctx, request, true)
	return s.SendRequest(ctx, method, endpoint, apiKey, "")
}

// SendRequestToWithID
// send request using credential parameters if they are defined in Configuration.
func (s *Session) SendRequestToWithPath(ctx context.Context, method, request, path string) error {
	apiKey := s.GoliumInterface.ValueAsString(ctx, fmt.Sprintf(confEndpoint, request))
	endpoint := s.NormalizeEndpoint(ctx, request, true)
	s.ConfigurePath(ctx, path)
	return s.SendRequest(ctx, method, endpoint, apiKey, "")
}

// SendRequestToWithoutBackslash send request without last backslash
// using credential parameters if they are defined in Configuration.
func (s *Session) SendRequestToWithoutBackslash(ctx context.Context, method, request string) error {
	apiKey := s.GoliumInterface.ValueAsString(ctx, fmt.Sprintf(confEndpoint, request))
	endpoint := s.NormalizeEndpoint(ctx, request, false)
	return s.SendRequest(ctx, method, endpoint, apiKey, "")
}

// SendRequestToWithAPIKEY send request using valid or invalid API key.
func (s *Session) SendRequestToWithAPIKEY(
	ctx context.Context,
	method, request, apiKeyFlag string,
) error {
	var apiKey string
	if apiKeyFlag == "valid" {
		apiKey = s.GoliumInterface.ValueAsString(ctx, fmt.Sprintf(confEndpoint, request))
	} else {
		apiKey = InvalidPath
	}
	endpoint := s.NormalizeEndpoint(ctx, request, true)
	return s.SendRequest(ctx, method, endpoint, apiKey, "")
}

// SendRequestToWithQuery send request using query parameters table.
func (s *Session) SendRequestToWithQueryParamsTable(ctx context.Context,
	method, request string,
	t *godog.Table,
) error {
	var err error
	params, err := golium.ConvertTableToMultiMap(ctx, t)
	if err != nil {
		return err
	}
	return s.sendRequestToWithQueryParams(ctx, method, request, params)
}

// SendRequestToWithFilters Send request with filters string.
// Filters string starts without "?" and filters are provided with
// format "key=value" and split by "&".
// Example: "key1=value1&key2=value2".
func (s *Session) SendRequestToWithFilters(ctx context.Context,
	method, request, filters string,
) error {
	var queryParams url.Values
	queryParams, _ = url.ParseQuery(filters)

	return s.sendRequestToWithQueryParams(ctx, method, request, queryParams)
}

// SendRequestUsingJSON send request using body from JSON file located in schemas.
func (s *Session) SendRequestUsingJSON(ctx context.Context, method, request, code string) error {
	if err := s.ConfigureRequestBodyJSONFile(ctx, code, request); err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("%s%v", configRequestBodyError, err))
	}
	if err := s.SendRequestTo(ctx, method, request); err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("%s%v", sendingRequestError, err))
	}
	return nil
}

// SendRequestWithIDUsingJSON send request using body from JSON file located in schemas.
func (s *Session) SendRequestWithPathUsingJSON(
	ctx context.Context,
	method, request, path, code string,
) error {
	if err := s.ConfigureRequestBodyJSONFile(ctx, code, request); err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("%s%v", configRequestBodyError, err))
	}
	if err := s.SendRequestToWithPath(ctx, method, request, path); err != nil {
		return fmt.Errorf("error sending request with path: %v", err)
	}
	return nil
}

// SendRequestUsingJSONWithout
// Send request using body from JSON file located in schemas without parameters.
func (s *Session) SendRequestUsingJSONWithout(ctx context.Context,
	method, request, code string,
	t *godog.Table,
) error {
	var err error
	params, err := golium.ConvertTableColumnToArray(ctx, t)
	if err != nil {
		return err
	}
	if err = s.ConfigureRequestBodyJSONFileWithout(
		ctx, code, request, params,
	); err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("%s%v", configRequestBodyError, err))
	}
	if err = s.SendRequestTo(ctx, method, request); err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("%s%v", sendingRequestError, err))
	}
	return nil
}

// SendRequestUsingJSONModifying
// Send request using body from JSON file located in schemas modifying parameters.
func (s *Session) SendRequestUsingJSONModifying(
	ctx context.Context,
	method, request, code string,
	t *godog.Table,
) error {
	message, err := GetParamFromJSON(ctx, request, code, "body")
	if err != nil {
		return fmt.Errorf("error getting parameter from json: %w", err)
	}
	messageMap, _ := message.(map[string]interface{})

	params, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return err
	}

	for key, value := range params {
		_, present := messageMap[key]
		if !present {
			return fmt.Errorf("error modifying param : param %v does not exists", key)
		}
		messageMap[key] = value
	}
	s.AddToRequestMessageFromJSONFile(message)
	if err := s.SendRequestTo(ctx, method, request); err != nil {
		return fmt.Errorf("%s", fmt.Sprintf("%s%v", sendingRequestError, err))
	}
	return nil
}

// ConfigureRequestHeaders Configure headers when apikey or jwt are provided.
func (s *Session) ConfigureRequestHeaders(ctx context.Context, apiKey, jwtValue string) error {
	var headers = map[string][]string{"Content-Type": {"application/json"}}
	if s.Request.Headers != nil {
		headers = s.Request.Headers
	}
	if apiKey != "" {
		headers["X-API-KEY"] = []string{apiKey}
	} else {
		delete(headers, "X-API-KEY")
	}
	if jwtValue != "" {
		headers["Authorization"] = []string{fmt.Sprintf("Bearer %s", jwtValue)}
	} else {
		delete(headers, "Authorization")
	}
	s.ConfigureHeaders(ctx, headers)
	return nil
}

// NormalizeEndpoint Normalize Endpoint considering ending backslash need.
func (s *Session) NormalizeEndpoint(ctx context.Context, request string, backslash bool) string {
	endpointConf := fmt.Sprintf("[CONF:endpoints.%s.api-endpoint]", request)
	endpoint := s.GoliumInterface.ValueAsString(ctx, endpointConf)
	if !backslash {
		return strings.TrimSuffix(endpoint, Slash)
	}
	if strings.HasSuffix(endpoint, Slash) {
		return endpoint
	}
	return endpoint + Slash
}

// ConfigureURLAndEndpoint Configure endpoint and path keeping existing pre configuration.
func (s *Session) ConfigureURLAndEndpoint(ctx context.Context, endpoint string) error {
	// Get context URL.
	URL, err := s.GetURL(ctx)
	if err != nil {
		return fmt.Errorf("error getting url: %v", err)
	}
	endpoint = URL + endpoint
	s.ConfigureEndpoint(ctx, endpoint)

	return nil
}

// sendRequestToWithQueryParams Sends request with query params once they
// are formatted to map map[string][]string from table or input string.
func (s *Session) sendRequestToWithQueryParams(ctx context.Context,
	method, request string,
	params map[string][]string,
) error {
	s.ConfigureQueryParams(ctx, params)
	return s.SendRequestTo(ctx, method, request)
}

// GetURL returns URL from Configuration or Context
func (s *Session) GetURL(ctx context.Context) (string, error) {
	URL := s.GoliumInterface.ValueAsString(ctx, "[CONF:url]")
	if URL == "<nil>" {
		URL = s.GoliumInterface.ValueAsString(ctx, "[CTXT:url]")
	}
	if URL == NilString {
		return "", fmt.Errorf("url shall be initialized in Configuration or Context")
	}
	return URL, nil
}

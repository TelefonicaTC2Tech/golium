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
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/require"
)

const (
	JSONpostsfile = `
[
	{
		"code": "example1",
		"body": {
			"title": "title_to_replace",
			"body": "bar1",
			"userId": 1
		},
		"response": {
			"message": "Validation error.",
			"details": {
			  "field_to_replace": {
				"message": "Field 'field_to_replace' has invalid type, expected type_to_replace",
				"code": "incorrect_type"
			  }
			},
			"status_code": 400
		}
	},
	{
		"code": "example4",
		"body": {
			"title": "title_to_replace",
			"body": "bar1",
			"userId": 1
		},
		"response": 1
	},
	{
		"code": "example5",
		"body": {
			"title": "title_to_replace",
			"body": "bar1",
			"userId": 1
		},
		"response": "field_to_replace and type_to_replace has been replaced"
	},
	{
		"code": "example6",
		"body": {
			"title": "title_to_replace",
			"body": "bar1",
			"userId": 1
		},
		"response": "field_to_replace and type_to_replace has been replaced with match error"
	}								
]
`
	requestUsingJSONFile = `
[
	{
		"code": "example1",
		"body": {
			"empty": "",
			"boolean": false,
			"list": [
			  { "attribute": "attribute0", "value": "value0"},
			  { "attribute": "attribute1", "value": "value1"},
			  { "attribute": "attribute2", "value": "value2"}
			]
		},
		"response": {
			"boolean": false, 
			"empty": "", 
			"list": [
				{ "attribute": "attribute0", "value": "value0"},
				{ "attribute": "attribute1", "value": "value1"},
				{ "attribute": "attribute2", "value": "value2"}
			]
		}
	}
]
`
	requestUsingJSONWithoutFile = `
[
	{
		"code": "example1",
		"body": {
			"empty": "",
			"boolean": false,
			"list": [
			  { "attribute": "attribute0", "value": "value0"},
			  { "attribute": "attribute1", "value": "value1"},
			  { "attribute": "attribute2", "value": "value2"}
			]
		},
		"response": {
			"boolean": false, 
			"empty": "", 
			"list": [
				{ "attribute": "attribute0", "value": "value0"},
				{ "attribute": "attribute1", "value": "value1"},
				{ "attribute": "attribute2", "value": "value2"}
			]
		}
	}
]
`
	requestUsingJSONModifying = `
[
	{
		"code": "example1",
		"body": {
			"empty": "",
			"boolean": false,
			"list": [
			  { "attribute": "attribute0", "value": "value0"},
			  { "attribute": "attribute1", "value": "value1"},
			  { "attribute": "attribute2", "value": "value2"}
			]
		},
		"response": {
			"boolean": false, 
			"empty": "", 
			"list": [
				{ "attribute": "attribute0", "value": "value0"},
				{ "attribute": "attribute1", "value": "value1"},
				{ "attribute": "attribute2", "value": "value2"}
			]
		}
	}
]
`
	schemasPath     = "./schemas"
	httpbinURLslash = "https://httpbin.org/"
	healthEndpoint  = "anything/health/"
	healthRequest   = "health"
	jsonCodeError   = "error configuring request body: error getting parameter from json: " +
		"param value: 'body' not found in '[map[body:map[boolean:false empty: " +
		"list:[map[attribute:attribute0 value:value0] map[attribute:attribute1 value:value1]" +
		" map[attribute:attribute2 value:value2]]] code:example1 response:map[boolean:false " +
		"empty: list:[map[attribute:attribute0 value:value0] map[attribute:attribute1 " +
		"value:value1] map[attribute:attribute2 value:value2]]]]]' due to error: " +
		"value for param: 'body' with code: 'not_valid_code' not found"
)

func TestSendRequest(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tcs := []struct {
		name         string
		configURL    string
		existingPath string
		fakeRequest  string
		expectedErr  error
	}{
		{
			name:         "nil url",
			configURL:    "<nil>",
			existingPath: NilString,
			expectedErr: fmt.Errorf(
				"error getting url: url shall be initialized in Configuration or Context"),
		},
		{
			name:         "request error",
			configURL:    "https://wrong.org/",
			existingPath: NilString,
			fakeRequest:  "error",
			expectedErr: fmt.Errorf(
				"error with the HTTP request. %v", "fake_error"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ValuesAsString = map[string]string{
				"[CONF:url]":                      tc.configURL,
				"[CTXT:url]":                      NilString,
				"[CONF:endpoints.health.api-key]": "valid",
			}
			FakeResponse = tc.fakeRequest

			ctx, s := setGoliumContextAndService()

			s.Request.Path = tc.existingPath

			err := s.SendRequest(ctx, http.MethodGet, healthEndpoint, "", "")

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestTo(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tcs := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "valid apikey",
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestHTTPBinContext("")

			err := s.SendRequestTo(ctx, http.MethodGet, healthRequest)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestToWithPath(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tcs := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "valid apikey",
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestHTTPBinContext("")

			err := s.SendRequestToWithPath(ctx, http.MethodGet, healthRequest, "1")

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestToWithoutBackslash(t *testing.T) {
	os.MkdirAll("./logs", os.ModePerm)
	defer os.RemoveAll("./logs/")

	tcs := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "no_backslash",
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestHTTPBinContext("")

			err := s.SendRequestToWithoutBackslash(ctx, "GET", "health")

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestToWithAPIKEY(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tcs := []struct {
		name        string
		apikey      string
		apikeyFlag  string
		expectedErr error
	}{
		{
			name:        "valid flag",
			apikey:      "valid",
			apikeyFlag:  "valid",
			expectedErr: nil,
		},
		{
			name:        "not valid flag",
			apikey:      "not_valid",
			apikeyFlag:  "else_flag",
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ValuesAsString = map[string]string{
				"[CONF:url]":                           httpbinURLslash,
				"[CTXT:url]":                           NilString,
				"[CONF:endpoints.health.api-endpoint]": healthEndpoint,
				"[CONF:endpoints.health.api-key]":      tc.apikey,
			}

			FakeResponse = ""

			ctx, s := setGoliumContextAndService()

			err := s.SendRequestToWithAPIKEY(ctx, http.MethodGet, healthRequest, tc.apikeyFlag)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestToWithQuery(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tcs := []struct {
		name        string
		table       *godog.Table
		expectedErr error
	}{
		{
			name: "valid apikey",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"field", "test"},
			}),
			expectedErr: nil,
		},
		{
			name: "input table error",
			table: golium.NewTable([][]string{
				{"field", "test"},
			}),
			expectedErr: fmt.Errorf("cannot remove header: %v",
				fmt.Errorf("table must have at least one header and one useful row")),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestHTTPBinContext("")

			err := s.SendRequestToWithQueryParamsTable(ctx, http.MethodGet, healthRequest, tc.table)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestToWithFilters(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tcs := []struct {
		name        string
		filters     string
		expectedErr error
	}{
		{
			name:        "single filter",
			filters:     "field=test",
			expectedErr: nil,
		},
		{
			name:        "multiple filter",
			filters:     "field=test&field2=test2",
			expectedErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestHTTPBinContext("")

			err := s.SendRequestToWithFilters(ctx, http.MethodGet, healthRequest, tc.filters)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestWithPathUsingJSON(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(requestUsingJSONFile), os.ModePerm)
	defer os.RemoveAll(schemasPath)

	tcs := []struct {
		name        string
		code        string
		fakeRequest string
		expectedErr error
	}{
		{
			name:        "valid json code",
			code:        "example1",
			expectedErr: nil,
		},
		{
			name:        "not valid json code",
			code:        "not_valid_code",
			expectedErr: fmt.Errorf(jsonCodeError),
		},
		{
			name:        "request error",
			code:        "example1",
			fakeRequest: "error",
			expectedErr: fmt.Errorf(
				"error sending request with path: error with the HTTP request. %v", "fake_error"),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestHTTPBinContext(tc.fakeRequest)

			err := s.SendRequestWithPathUsingJSON(ctx, http.MethodPost, healthRequest, "8", tc.code)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestUsingJSONWithout(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(requestUsingJSONWithoutFile), os.ModePerm)
	defer os.RemoveAll(schemasPath)

	tcs := []struct {
		name        string
		code        string
		fakeRequest string
		table       *godog.Table
		expectedErr error
	}{
		{
			name: "valid param",
			code: "example1",
			table: golium.NewTable([][]string{
				{"parameter"},
				{"boolean"},
			}),
			expectedErr: nil,
		},
		{
			name: "wrong code",
			table: golium.NewTable([][]string{
				{"parameter"},
				{"boolean"},
			}),
			code: "example2",
			expectedErr: fmt.Errorf(
				"error configuring request body: error getting parameter from json: param value:" +
					" 'body' not found in '[map[body:map[boolean:false empty: list:[map[attribute:" +
					"attribute0 value:value0] map[attribute:attribute1 value:value1] map[attribute:" +
					"attribute2 value:value2]]] code:example1 response:map[boolean:false empty: " +
					"list:[map[attribute:attribute0 value:value0] map[attribute:attribute1 value:" +
					"value1] map[attribute:attribute2 value:value2]]]]]' due to error: value for " +
					"param: 'body' with code: 'example2' not found"),
		},
		{
			name: "request error",
			table: golium.NewTable([][]string{
				{"parameter"},
				{"boolean"},
			}),
			code:        "example1",
			fakeRequest: "error",
			expectedErr: fmt.Errorf("error sending request: error with the HTTP request. %v", "fake_error"),
		},
		{
			name: "input table error",
			table: golium.NewTable([][]string{
				{"boolean"},
			}),
			code: "example1",
			expectedErr: fmt.Errorf("cannot remove header: %v",
				fmt.Errorf("table must have at least one header and one useful row")),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestHTTPBinContext(tc.fakeRequest)

			err := s.SendRequestUsingJSONWithout(ctx, http.MethodPost, healthRequest,
				tc.code, tc.table)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSendRequestUsingJSONModifying(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(requestUsingJSONModifying), os.ModePerm)
	defer os.RemoveAll(schemasPath)

	tcs := []struct {
		name        string
		code        string
		fakeRequest string
		table       *godog.Table
		expectedErr error
	}{
		{
			name: "valid code",
			code: "example1",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"boolean", "true"},
			}),
			expectedErr: nil,
		},
		{
			name: "wrong code",
			code: "example2",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"boolean", "true"},
			}),
			expectedErr: fmt.Errorf(
				"error getting parameter from json: %w",
				fmt.Errorf("param value: 'body' not found in '[map[body:map[boolean:false empty:"+
					" list:[map[attribute:attribute0 value:value0] map[attribute:attribute1 value:"+
					"value1] map[attribute:attribute2 value:value2]]] code:example1 response:map["+
					"boolean:false empty: list:[map[attribute:attribute0 value:value0] map["+
					"attribute:attribute1 value:value1] map[attribute:attribute2 value:value2]]]]]'"+
					" due to error: %w",
					fmt.Errorf("value for param: 'body' with code: 'example2' not found"))),
		},
		{
			name: "wrong param",
			code: "example1",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"wrong_key", "true"},
			}),
			expectedErr: fmt.Errorf("error modifying param : param wrong_key does not exists"),
		},
		{
			name: "request error",
			code: "example1",
			table: golium.NewTable([][]string{
				{"parameter", "value"},
				{"boolean", "true"},
			}),
			fakeRequest: "error",
			expectedErr: fmt.Errorf("error sending request: error with the HTTP request. %v", "fake_error"),
		},
		{
			name: "input table error",
			code: "example1",
			table: golium.NewTable([][]string{
				{"boolean", "true"},
			}),
			fakeRequest: "error",
			expectedErr: fmt.Errorf("cannot remove header: %v",
				fmt.Errorf("table must have at least one header and one useful row")),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctx, s := setRequestToTestHTTPBinContext(tc.fakeRequest)

			err := s.SendRequestUsingJSONModifying(ctx, http.MethodPost, healthRequest, tc.code, tc.table)

			require.Equal(t, tc.expectedErr, err)
		})
	}
}
func setRequestToTestHTTPBinContext(fakeResponse string) (context.Context, *Session) {
	ValuesAsString = map[string]string{
		"[CONF:url]":                           httpbinURLslash,
		"[CTXT:url]":                           NilString,
		"[CONF:endpoints.health.api-endpoint]": healthEndpoint,
		"[CONF:endpoints.health.api-key]":      "valid",
	}
	FakeResponse = fakeResponse

	ctx, s := setGoliumContextAndService()
	return ctx, s
}

func setGoliumContextAndService() (context.Context, *Session) {
	ctxGolium := InitializeContext(context.Background())
	ctx := InitializeContext(ctxGolium)
	s := GetSession(ctx)
	s.GoliumInterface = GoliumInterfaceMock{}
	return ctx, s
}

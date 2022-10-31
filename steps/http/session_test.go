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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/TelefonicaTC2Tech/golium"
	"github.com/TelefonicaTC2Tech/golium/steps/http/model"
	"github.com/TelefonicaTC2Tech/golium/steps/http/schema"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	httpbinURL         = "https://httpbin.org"
	httpbinURLslash    = "https://httpbin.org/"
	httpSelfSignedURL  = "https://self-signed.badssl.com"
	healthRequest      = "health"
	logsPath           = "./logs"
	schemasPath        = "./schemas"
	JSONhttpFileValues = `
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

	JSONhttpResponse = `{
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
	}`

	JSONFile = `{
		"boolean": false, 
		"empty": "", 
		"list": [
			{ "attribute": "attribute0", "value": "value0"},
			{ "attribute": "attribute1", "value": "value1"},
			{ "attribute": "attribute2", "value": "value2"}
		]
	}`

	JSONhttpFileBadFormat = `
	[
		{
			"code": "example1",
			"body": {
	`
	validateModifyingResponseFile = `
	[
		{
			"code": "example1",
			"body": {
				"title": "foo1",
				"body": "bar1",
				"userId": 1
			},
			"response": {
				"id": 101,
				"title": "foo2",
				"body": "bar1",
				"userId": 1
			}
		}
	]
	`
	validateModifyingResponse = `
	{
		"id": 101,
		"title": "foo1",
		"body": "bar1",
		"userId": 1
	}
	`
)

func TestURL(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		want    *url.URL
	}{
		{
			name: "url with path",
			path: "/testing",
		},
		{
			name: "url with empty path",
			path: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.Request.Endpoint = httpbinURL
			s.Request.Path = tt.path
			result, _ := s.URL()
			assert.Equal(t, result.Path, tt.path)
		})
	}
}

func TestConfigureRequestBodyJSONProperties(t *testing.T) {
	expectedResult := make(map[string]interface{})
	expectedResult["John"] = "182"

	failedResult := make(map[string]interface{})
	failedResult[""] = "182"

	tests := []struct {
		name    string
		props   map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Should setting property and return nil",
			props:   expectedResult,
			wantErr: false,
		},
		{
			name:    "Should error setting property",
			props:   failedResult,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			err := s.ConfigureRequestBodyJSONProperties(ctx, tt.props)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ConfigureRequestBodyJSONProperties() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestConfigureRequestBodyJSONFile(t *testing.T) {
	golium.GetConfig().Dir.Schemas = schemasPath
	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tests := []struct {
		name    string
		params  schema.Params
		wantErr bool
	}{
		{
			name: "Should add request message from JSON file",
			params: schema.Params{
				File: "http",
				Code: "example1",
			},
			wantErr: false,
		},
		{
			name: "Should return error bad code",
			params: schema.Params{
				File: "http",
				Code: "bad code",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			if err := s.ConfigureRequestBodyJSONFile(
				ctx, tt.params); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ConfigureRequestBodyJSONFile() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestConfigureRequestBodyJSONFileWithout(t *testing.T) {
	golium.GetConfig().Dir.Schemas = schemasPath
	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tests := []struct {
		name    string
		params  schema.Params
		fields  []string
		wantErr bool
	}{
		{
			name: "should remove the parameter from the message",
			params: schema.Params{
				File: "http",
				Code: "example1",
			},
			fields:  []string{"boolean"},
			wantErr: false,
		},
		{
			name: "should return error getting key json",
			params: schema.Params{
				File: "http",
				Code: "badcode",
			},
			fields:  []string{"boolean"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			if err := s.ConfigureRequestBodyJSONFileWithout(
				ctx, tt.params, tt.fields); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ConfigureRequestBodyJSONFileWithout() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestConfigureRequestBodyURLEncodedProperties(t *testing.T) {
	expectedResult := make(map[string][]string)
	expectedResult["list"] = []string{"testing"}

	tests := []struct {
		name    string
		props   map[string][]string
		wantErr bool
	}{
		{
			name:    "should remove the parameter from the message",
			props:   expectedResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			if err := s.ConfigureRequestBodyURLEncodedProperties(
				ctx, tt.props); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ConfigureRequestBodyURLEncodedProperties() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestSendHTTPRequest(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tests := []struct {
		name       string
		endpoint   string
		method     string
		host       string
		preHeader  bool
		wantErr    bool
		selfSigned bool
	}{
		{
			name:       "testing with correct method",
			endpoint:   httpbinURL,
			method:     "POST",
			host:       "",
			preHeader:  false,
			wantErr:    false,
			selfSigned: false,
		},
		{
			name:       "testing empty endpoint",
			endpoint:   "",
			method:     "POST",
			host:       "",
			preHeader:  false,
			wantErr:    true,
			selfSigned: false,
		},
		{
			name:       "testing invalid method",
			endpoint:   "httpbinURL",
			method:     "invalid Method",
			host:       "",
			preHeader:  false,
			wantErr:    true,
			selfSigned: false,
		},
		{
			name:       "testing headers auth",
			endpoint:   "httpbinURL",
			method:     "POST",
			host:       "httpbin.org",
			preHeader:  true,
			wantErr:    true,
			selfSigned: false,
		},
		{
			name:       "testing skip verify cert",
			endpoint:   httpSelfSignedURL,
			method:     "GET",
			host:       "",
			preHeader:  false,
			wantErr:    false,
			selfSigned: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			s.Request.Endpoint = tt.endpoint
			s.Request.Username = "QA"
			s.Request.Password = "QATesting#"
			s.NoRedirect = true
			if tt.selfSigned {
				s.ConfigureInsecureSkipVerify(ctx)
			}
			if tt.preHeader {
				s.Request.Headers = map[string][]string{
					"Content-Type":  {"application/json"},
					"Authorization": {"Bearer 1234567890ABCD"},
					"host":          {tt.host},
				}
			}
			s.Request.Headers = make(map[string][]string)
			s.Request.Headers["host"] = []string{tt.host}
			if err := s.SendHTTPRequest(ctx, tt.method); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.SendHTTPRequest() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateResponseHeaders(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	tests := []struct {
		name                string
		contentType         string
		responseContentType string
		wantErr             bool
	}{
		{
			name:                "testing correct headers",
			contentType:         "application/json",
			responseContentType: "application/json",
			wantErr:             false,
		},
		{
			name:                "testing incorrect headers",
			contentType:         "application/json",
			responseContentType: "failcontentType",
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			s.Request.Headers = map[string][]string{
				"Content-Type": {tt.contentType},
			}
			s.ConfigureHeaders(ctx, s.Request.Headers)

			header := http.Header{}
			header.Add("Content-Type", tt.responseContentType)
			s.Response.HTTPResponse = &http.Response{Header: header}

			err := s.ValidateResponseHeaders(ctx, s.Request.Headers)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ValidateResponseHeaders() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateResponseFromJSONFile(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	var response interface{}
	if err := json.Unmarshal([]byte(JSONFile), &response); err != nil {
		t.Error("error Unmarshaling expected response body: %w", err)
	}

	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	unexpectedResponse := `{
		"boolean": true, 
		"empty": "", 
		"list": [
			{ "attribute": "attribute0", "value": "value0"},
			{ "attribute": "attribute1", "value": "value1"},
			{ "attribute": "attribute2", "value": "value2"}
		]
	}`

	var incorrect = `
	{
		"boolean": false, 
		"empty": ""
	}`

	tests := []struct {
		name             string
		response         interface{}
		respDataLocation string
		responseBody     string
		wantErr          bool
	}{
		{
			name:             "testing validate response type string",
			response:         JSONFile,
			responseBody:     JSONhttpResponse,
			respDataLocation: "response",
			wantErr:          false,
		},
		{
			name:             "testing unexpected response body ",
			response:         JSONFile,
			responseBody:     unexpectedResponse,
			respDataLocation: "",
			wantErr:          true,
		},
		{
			name:             "testing validate response type interface",
			response:         response,
			responseBody:     JSONhttpResponse,
			respDataLocation: "response",
			wantErr:          false,
		},
		{
			name:             "testing incorrect response type interface",
			response:         response,
			responseBody:     incorrect,
			respDataLocation: "response",
			wantErr:          true,
		},
		{
			name:             "testing unexpected response type interface",
			response:         response,
			responseBody:     unexpectedResponse,
			respDataLocation: "",
			wantErr:          true,
		},
		{
			name:             "Response body content should be string or map",
			response:         777,
			responseBody:     JSONhttpResponse,
			respDataLocation: "",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.Response.ResponseBody = []byte(tt.responseBody)
			if err := s.ValidateResponseFromJSONFile(
				tt.response, tt.respDataLocation); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ValidateResponseFromJSONFile() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateResponseBodyJSONFile(t *testing.T) {
	golium.GetConfig().Dir.Schemas = schemasPath

	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/httpBadFormat.json", []byte(JSONhttpFileBadFormat), os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tests := []struct {
		name             string
		params           schema.Params
		responseBody     string
		respDataLocation string
		wantErr          bool
	}{
		{
			name: "Should return selected value from JSON file",
			params: schema.Params{
				File: "http",
				Code: "example1",
			},
			responseBody:     JSONhttpResponse,
			respDataLocation: "response",
			wantErr:          false,
		},
		{
			name: "Should return a error unmarsharlling JSON file",
			params: schema.Params{
				File: "httpBadFormat",
				Code: "example1",
			},
			responseBody:     JSONhttpResponse,
			respDataLocation: "response",
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			s.Response.ResponseBody = []byte(tt.responseBody)
			if err := s.ValidateResponseBodyJSONFile(
				ctx, tt.params, tt.respDataLocation); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateResponseBodyJSONFile() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateResponseBodyJSONFileWithout(t *testing.T) {
	JSONhttpResponseWithout := `{
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
			"empty": "", 
			"list": [
				{ "attribute": "attribute0", "value": "value0"},
				{ "attribute": "attribute1", "value": "value1"},
				{ "attribute": "attribute2", "value": "value2"}
			]
		}
		}`

	golium.GetConfig().Dir.Schemas = schemasPath

	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tests := []struct {
		name             string
		params           schema.Params
		responseBody     string
		respDataLocation string
		t                *godog.Table
		wantErr          bool
	}{
		{
			name: "should remove the parameter from the file",
			params: schema.Params{
				File: "http",
				Code: "example1",
			},

			responseBody:     JSONhttpResponseWithout,
			respDataLocation: "response",
			t: golium.NewTable([][]string{
				{"parameter"},
				{"boolean"},
			}),
			wantErr: false,
		},
		{
			name: "Error deleting response",
			params: schema.Params{
				File: "http",
				Code: "wrong_code",
			},
			responseBody:     JSONhttpResponseWithout,
			respDataLocation: "response",
			t: golium.NewTable([][]string{
				{"parameter"},
				{"boolean"},
			}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			s.Response.ResponseBody = []byte(tt.responseBody)
			if err := s.ValidateResponseBodyJSONFileWithout(
				ctx, tt.params, tt.respDataLocation, tt.t); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateResponseBodyJSONFileWithout() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateResponseBodyJSONProperties(t *testing.T) {
	props := map[string]interface{}{
		"boolean":          false,
		"empty":            "",
		"list.0.attribute": "attribute0",
		"list.0.value":     "value0",
		"list.1.attribute": "attribute1",
		"list.1.value":     "value1",
		"list.2.attribute": "attribute2",
		"list.2.value":     "value2",
	}

	tests := []struct {
		name         string
		responseBody string
		props        map[string]interface{}
		wantErr      bool
	}{
		{
			name:         "testing validate response body json with correct properties",
			responseBody: JSONFile,
			props:        props,
			wantErr:      false,
		},
		{
			name:         "testing validate response body json with incorrect properties",
			responseBody: JSONFile,
			props:        map[string]interface{}{"boolean": true},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			s.Response.ResponseBody = []byte(tt.responseBody)
			if err := s.ValidateResponseBodyJSONProperties(ctx, tt.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateResponseBodyJSONProperties() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateResponseBodyEmpty(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		wantErr      bool
	}{
		{
			name:         "testing response body empty",
			responseBody: "",
			wantErr:      false,
		},
		{
			name:         "Should return error body is not empty",
			responseBody: JSONFile,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			s.Response.ResponseBody = []byte(tt.responseBody)
			s.Response.HTTPResponse = &http.Response{}
			s.Response.HTTPResponse.ContentLength = 0
			if err := s.ValidateResponseBodyEmpty(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateResponseBodyEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendRequestWithBody(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		ctx      context.Context
		uRL      string
		endpoint string
		params   schema.Params
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error getting body from file",
			args: args{
				params: schema.Params{
					Code: "not_valid_code",
					File: healthRequest,
				},
				endpoint: healthRequest,
				uRL:      httpbinURLslash,
			},
			wantErr: true,
		},
		{
			name: "Error sending HTTP Request",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				endpoint: healthRequest,
				uRL:      "wrongURL",
			},
			wantErr: true,
		},
		{
			name: "Happy path",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				endpoint: healthRequest,
				uRL:      httpbinURLslash,
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &Session{}
			if err := s.SendRequestWithBody(
				tc.args.ctx, tc.args.uRL, http.MethodPost,
				tc.args.endpoint, tc.args.params, "validApiKey",
			); (err != nil) != tc.wantErr {
				t.Errorf("Session.SendRequestWithBody() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestSendRequestWithBodyWithoutFields(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		ctx      context.Context
		uRL      string
		endpoint string
		params   schema.Params
		table    *godog.Table
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error converting table",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				uRL:      "wrongURL",
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error sending HTTP Request",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				uRL:      "wrongURL",
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter"},
					{"boolean"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Happy path",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				uRL:      httpbinURLslash,
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter"},
					{"boolean"},
				}),
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &Session{}
			if err := s.SendRequestWithBodyWithoutFields(
				tc.args.ctx,
				tc.args.uRL, http.MethodPost, tc.args.endpoint, tc.args.params, "validApiKey",
				tc.args.table,
			); (err != nil) != tc.wantErr {
				t.Errorf(
					"Session.SendRequestWithBodyWithoutFields() error = %v, wantErr %v",
					err, tc.wantErr)
			}
		})
	}
}

func TestSendRequestWithBodyModifyingFields(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		ctx      context.Context
		uRL      string
		endpoint string
		params   schema.Params
		table    *godog.Table
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error modifying body",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				uRL:      "wrongURL",
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter", "value"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error sending HTTP Request",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				uRL:      "wrongURL",
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter", "value"},
					{"boolean", "true"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Happy path",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				uRL:      httpbinURLslash,
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter", "value"},
					{"boolean", "true"},
				}),
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &Session{}
			if err := s.SendRequestWithBodyModifyingFields(
				tc.args.ctx,
				tc.args.uRL, http.MethodPost, tc.args.endpoint, tc.args.params, "validApiKey",
				tc.args.table,
			); (err != nil) != tc.wantErr {
				t.Errorf(
					"Session.SendRequestWithBodyModifyingFields() error = %v, wantErr %v",
					err, tc.wantErr)
			}
		})
	}
}

func TestSendRequestWithQueryParams(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		ctx      context.Context
		uRL      string
		endpoint string
		table    *godog.Table
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error converting table",
			args: args{
				uRL:      "wrongURL",
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter", "value"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Error sending HTTP Request",
			args: args{
				uRL:      "wrongURL",
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter", "value"},
					{"field", "test"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Happy path",
			args: args{
				uRL:      httpbinURLslash,
				endpoint: healthRequest,
				table: golium.NewTable([][]string{
					{"parameter", "value"},
					{"field", "test"},
				}),
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &Session{}
			if err := s.SendRequestWithQueryParams(
				tc.args.ctx,
				tc.args.uRL, http.MethodPost, tc.args.endpoint, "validApiKey", tc.args.table,
			); (err != nil) != tc.wantErr {
				t.Errorf("Session.SendRequestWithQueryParams() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestSendRequestWithFilters(t *testing.T) {
	s := &Session{}
	testBasicRequestWithParam(t, "field=test&field2=test2", s.SendRequestWithFilters)
}

func TestSendRequestWithPath(t *testing.T) {
	s := &Session{}
	testBasicRequestWithParam(t, "1", s.SendRequestWithPath)
}

func TestSendRequestWithPathAndBody(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		ctx      context.Context
		uRL      string
		endpoint string
		params   schema.Params
		path     string
	}
	testCases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error getting body from file",
			args: args{
				params: schema.Params{
					Code: "not_valid_code",
					File: healthRequest,
				},
				uRL:      httpbinURLslash,
				endpoint: healthRequest,
				path:     "1",
			},
			wantErr: true,
		},
		{
			name: "Error sending HTTP Request",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				uRL:      "wrongURL",
				endpoint: healthRequest,
			},
			wantErr: true,
		},
		{
			name: "Happy path",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: healthRequest,
				},
				uRL:      httpbinURLslash,
				endpoint: healthRequest,
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &Session{}
			if err := s.SendRequestWithPathAndBody(
				tc.args.ctx,
				tc.args.uRL, http.MethodPost, tc.args.endpoint,
				tc.args.path, tc.args.params, "validApiKey",
			); (err != nil) != tc.wantErr {
				t.Errorf("Session.SendRequestWithPathAndBody() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestSendRequestWithoutBackslash(t *testing.T) {
	s := Session{}
	testBasicRequest(t, s.SendRequestWithoutBackslash)
}

func TestSendRequest(t *testing.T) {
	s := Session{}
	testBasicRequest(t, s.SendRequest)
}

func testBasicRequest(t *testing.T, f func(ctx context.Context,
	uRL, method, endpoint, apiKey string) error) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	testCases := getBasicRequestTestCases("")
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := f(
				tc.args.ctx,
				tc.args.uRL, http.MethodPost, tc.args.endpoint, "validApiKey",
			); (err != nil) != tc.wantErr {
				t.Errorf("Session.SendRequest() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func testBasicRequestWithParam(t *testing.T, param string, f func(ctx context.Context,
	uRL, method, endpoint, apiKey, param string) error) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/health.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	testCases := getBasicRequestTestCases(param)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := f(
				tc.args.ctx,
				tc.args.uRL, http.MethodPost, tc.args.endpoint, "validApiKey",
				tc.args.param,
			); (err != nil) != tc.wantErr {
				t.Errorf("Session.SendRequestParam() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func getBasicRequestTestCases(param string) []requestTC {
	testCases := []requestTC{
		{
			name: "Error sending HTTP Request",
			args: requestArgs{
				uRL:      "wrongURL",
				endpoint: healthRequest,
				param:    param,
			},
			wantErr: true,
		},
		{
			name: "Happy path",
			args: requestArgs{
				uRL:      httpbinURLslash,
				endpoint: healthRequest,
				param:    param,
			},
			wantErr: false,
		},
	}
	return testCases
}

type requestArgs struct {
	ctx      context.Context
	uRL      string
	endpoint string
	param    string
}
type requestTC struct {
	name    string
	args    requestArgs
	wantErr bool
}

func TestGetURL(t *testing.T) {
	tcs := []struct {
		name        string
		configURL   string
		contextURL  string
		expectedURL string
		expectedErr error
	}{
		{
			name:        "valid conf url",
			configURL:   httpbinURL,
			contextURL:  httpbinURL,
			expectedURL: httpbinURL,
			expectedErr: nil,
		},
		{
			name:        "valid contextURL",
			configURL:   "<nil>",
			contextURL:  httpbinURL,
			expectedURL: httpbinURL,
			expectedErr: nil,
		},
		{
			name:        "nil url",
			configURL:   "<nil>",
			contextURL:  NilString,
			expectedURL: "",
			expectedErr: fmt.Errorf("url shall be initialized in Configuration or Context"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Call the tested function
			s := Session{}
			monkey.Patch(golium.ValueAsString, func(ctx context.Context, s string) string {
				if s == "[CONF:url]" {
					return tc.configURL
				} else if s == "[CTXT:url]" {
					return tc.contextURL
				}
				return ""
			})
			resultURL, resulterr := s.GetURL(context.Background())

			// Check expected behavior
			require.Equal(t, tc.expectedURL, resultURL)
			require.Equal(t, tc.expectedErr, resulterr)
		})
	}
}

func TestValidateResponseBodyJSONFileModifying(t *testing.T) {
	os.MkdirAll(schemasPath, os.ModePerm)
	os.WriteFile("./schemas/posts.json", []byte(validateModifyingResponseFile), os.ModePerm)
	defer os.RemoveAll(schemasPath)
	type args struct {
		params schema.Params
		t      *godog.Table
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error modifying schema response",
			args: args{
				params: schema.Params{
					Code: "not_valid_code",
					File: "posts",
				},
				t: golium.NewTable([][]string{
					{"parameter"},
				}),
			},
			wantErr: true,
		},
		{
			name: "Happy Path",
			args: args{
				params: schema.Params{
					Code: "example1",
					File: "posts",
				},
				t: golium.NewTable([][]string{
					{"parameter", "value"},
					{"title", "foo1"},
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Response: model.Response{
					ResponseBody: []byte(validateModifyingResponse),
				},
			}
			if err := s.ValidateResponseBodyJSONFileModifying(
				context.Background(), tt.args.params, tt.args.t,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ValidateResponseBodyJSONFileModifying() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSendRequestWithMultipartBody(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	os.WriteFile("./test.txt", []byte("test file content"), os.ModePerm)
	defer os.Remove("./test.txt")
	type args struct {
		uRL       string
		method    string
		path      string
		fileField string
		file      string
		t         *godog.Table
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy path",
			args: args{
				uRL:       "https://httpbin.org/",
				method:    "POST",
				fileField: "testFieldField",
				file:      "./test.txt",
				path:      "post",
				t: golium.NewTable(
					[][]string{
						{"field", "value"},
						{"fieldA", "valueA"},
						{"fieldB", "valueB"},
					},
				),
			},
			wantErr: false,
		},
		{
			name: "Wrong Write Multipart Body File",
			args: args{
				uRL:       "https://httpbin.org/",
				fileField: "testFieldField",
				file:      "./fake.txt",
			},
			wantErr: true,
		},
		{
			name: "Error with Multipart body fields",
			args: args{
				uRL:       "https://httpbin.org/",
				method:    "POST",
				fileField: "testFieldField",
				file:      "./test.txt",
				t: golium.NewTable(
					[][]string{
						{"field", "value"},
					},
				),
			},
			wantErr: true,
		},
		{
			name: "Error sending request",
			args: args{
				uRL:       "https://httpbin.org/",
				method:    "ÑÑÑ",
				fileField: "testFieldField",
				file:      "./test.txt",
				path:      "post",
				t: golium.NewTable(
					[][]string{
						{"field", "value"},
						{"fieldA", "valueA"},
						{"fieldB", "valueB"},
					},
				),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			if err := s.SendRequestWithMultipartBody(
				context.Background(),
				RequestParams{
					URL:      tt.args.uRL,
					Method:   tt.args.method,
					Endpoint: "",
					APIKey:   "",
					Table:    tt.args.t,
				},
				tt.args.fileField,
				tt.args.file,
			); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.SendRequestWithMultipartBody() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

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
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/stretchr/testify/assert"
)

const httpbinURL = "https://httpbin.org"
const httpSelfSignedURL = "https://self-signed.badssl.com"
const logsPath = "./logs"

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
	golium.GetConfig().Dir.Schemas = schemasDir
	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tests := []struct {
		name     string
		fileName string
		code     string
		wantErr  bool
	}{
		{
			name:     "Should add request message from JSON file",
			fileName: "http",
			code:     "example1",
			wantErr:  false,
		},
		{
			name:     "Should return error bad code",
			fileName: "http",
			code:     "bad code",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			if err := s.ConfigureRequestBodyJSONFile(
				ctx, tt.code, tt.fileName); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.ConfigureRequestBodyJSONFile() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestConfigureRequestBodyJSONFileWithout(t *testing.T) {
	golium.GetConfig().Dir.Schemas = schemasDir
	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tests := []struct {
		name     string
		fileName string
		code     string
		params   []string
		wantErr  bool
	}{
		{
			name:     "should remove the parameter from the message",
			fileName: "http",
			code:     "example1",
			params:   []string{"boolean"},
			wantErr:  false,
		},
		{
			name:     "should return error getting key json",
			fileName: "http",
			code:     "badcode",
			params:   []string{"boolean"},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			if err := s.ConfigureRequestBodyJSONFileWithout(
				ctx, tt.code, tt.fileName, tt.params); (err != nil) != tt.wantErr {
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
			s.Response.Response = &http.Response{Header: header}

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
	if err := json.Unmarshal([]byte(JSON), &response); err != nil {
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
			response:         JSON,
			responseBody:     JSONhttpResponse,
			respDataLocation: "response",
			wantErr:          false,
		},
		{
			name:             "testing unexpected response body ",
			response:         JSON,
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
	golium.GetConfig().Dir.Schemas = schemasDir

	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/httpBadFormat.json", []byte(JSONhttpFileBadFormat), os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tests := []struct {
		name             string
		code             string
		file             string
		responseBody     string
		respDataLocation string
		wantErr          bool
	}{
		{
			name:             "Should return selected value from JSON file",
			file:             "http",
			code:             "example1",
			responseBody:     JSONhttpResponse,
			respDataLocation: "response",
			wantErr:          false,
		},
		{
			name:             "Should return a error unmarsharlling JSON file",
			file:             "httpBadFormat",
			code:             "example1",
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
				ctx, tt.code, tt.file, tt.respDataLocation); (err != nil) != tt.wantErr {
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

	golium.GetConfig().Dir.Schemas = schemasDir

	os.MkdirAll("./schemas", os.ModePerm)
	os.WriteFile("./schemas/http.json", []byte(JSONhttpFileValues), os.ModePerm)
	defer os.RemoveAll("./schemas/")

	tests := []struct {
		name             string
		code             string
		file             string
		responseBody     string
		respDataLocation string
		params           []string
		wantErr          bool
	}{
		{
			name:             "should remove the parameter from the file",
			file:             "http",
			code:             "example1",
			responseBody:     JSONhttpResponseWithout,
			respDataLocation: "response",
			params:           []string{"boolean"},
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			s.Response.ResponseBody = []byte(tt.responseBody)
			if err := s.ValidateResponseBodyJSONFileWithout(
				ctx, tt.code, tt.file, tt.respDataLocation, tt.params); (err != nil) != tt.wantErr {
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
			responseBody: JSON,
			props:        props,
			wantErr:      false,
		},
		{
			name:         "testing validate response body json with incorrect properties",
			responseBody: JSON,
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
			responseBody: JSON,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			s := &Session{}
			s.Response.ResponseBody = []byte(tt.responseBody)
			s.Response.Response = &http.Response{}
			s.Response.Response.ContentLength = 0
			if err := s.ValidateResponseBodyEmpty(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateResponseBodyEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/stretchr/testify/assert"
)

const httpbinURL = "https://httpbin.org"
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
		name      string
		endpoint  string
		method    string
		host      string
		preHeader bool
		wantErr   bool
	}{
		{
			name:      "testing with correct method",
			endpoint:  httpbinURL,
			method:    "POST",
			host:      "",
			preHeader: false,
			wantErr:   false,
		},
		{
			name:      "testing empty endpoint",
			endpoint:  "",
			method:    "POST",
			host:      "",
			preHeader: false,
			wantErr:   true,
		},
		{
			name:      "testing invalid method",
			endpoint:  "httpbinURL",
			method:    "invalid Method",
			host:      "",
			preHeader: false,
			wantErr:   true,
		},
		{
			name:      "testing headers auth",
			endpoint:  "httpbinURL",
			method:    "POST",
			host:      "httpbin.org",
			preHeader: true,
			wantErr:   true,
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

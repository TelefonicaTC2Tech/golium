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
	"fmt"
	"net/http"
	"strings"

	"github.com/TelefonicaTC2Tech/golium/logger"
)

var httpLog *Logger

// Logger logs in a configurable file.
type Logger struct {
	Log *logger.Logger
}

// GetLogger returns the logger for HTTP requests and responses.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	name := "http"
	if httpLog == nil {
		httpLog = &Logger{Log: logger.Factory(name)}
	}
	return httpLog
}

// LogRequest logs an HTTP request in the configured log file.
func (l Logger) LogRequest(req *http.Request, body []byte, corr string) {
	l.Log.Printf("Request [%s]:\n%s\n%s\n%s\n\n",
		corr,
		getRequestFirstLine(req),
		getHeaders(req.Header),
		getBody(body))
}

// LogResponse logs an HTTP response in the configured log file.
func (l Logger) LogResponse(resp *http.Response, body []byte, corr string) {
	l.Log.Printf("Response [%s]:\n%s\n%s\n%s\n\n",
		corr,
		getResponseFirstLine(resp),
		getHeaders(resp.Header),
		getBody(body))
}

// LogTimeout logs an HTTP response with timeout in the configured log file.
func (l Logger) LogTimeout(corr string) {
	l.Log.Print("Response: Timeout\n\n")
}

func getRequestFirstLine(req *http.Request) string {
	return fmt.Sprintf("%s %s", req.Method, req.URL)
}

func getResponseFirstLine(resp *http.Response) string {
	return fmt.Sprintf("%s %s", resp.Proto, resp.Status)
}

func getHeaders(headers map[string][]string) string {
	var fmtHeaders []string
	for key, values := range headers {
		for _, value := range values {
			fmtHeaders = append(fmtHeaders, fmt.Sprintf("%s: %s", key, value))
		}
	}
	return strings.Join(fmtHeaders, "\n")
}

func getBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	return fmt.Sprintf("\n%s", body)
}

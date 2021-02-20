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
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Telefonica/golium"
	"github.com/sirupsen/logrus"
)

var logger *Logger

// GetLogger returns the logger for HTTP requests and responses.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	if logger != nil {
		return logger
	}
	dir := golium.GetConfig().Log.Directory
	path := path.Join(dir, "http.log")
	logger, err := NewLogger(path)
	if err != nil {
		logrus.Fatalf("Error creating HTTP logger with file: '%s'. %s", path, err)
	}
	return logger
}

// Logger logs the HTTP request and response in a configurable file.
type Logger struct {
	log *log.Logger
}

// NewLogger creates an instance of the logger.
// It configures the file path where the HTTP request and response are written.
func NewLogger(path string) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{
		log: log.New(file, "", log.Ldate|log.Lmicroseconds|log.LUTC),
	}, nil
}

// LogRequest logs an HTTP request in the configured log file.
func (l Logger) LogRequest(req *http.Request, body []byte, corr string) {
	l.log.Println(fmt.Sprintf("Request [%s]:\n%s\n%s\n%s\n",
		corr,
		getRequestFirstLine(req),
		getHeaders(req.Header),
		getBody(body)))
}

// LogResponse logs an HTTP response in the configured log file.
func (l Logger) LogResponse(resp *http.Response, body []byte, corr string) {
	l.log.Println(fmt.Sprintf("Response [%s]:\n%s\n%s\n%s\n",
		corr,
		getResponseFirstLine(resp),
		getHeaders(resp.Header),
		getBody(body)))
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

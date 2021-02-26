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

package dns

import (
	"log"
	"os"
	"path"

	"github.com/Telefonica/golium"
	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

var logger *Logger

// GetLogger returns the logger for DNS requests and responses.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	if logger != nil {
		return logger
	}
	dir := golium.GetConfig().Log.Directory
	path := path.Join(dir, "dns.log")
	logger, err := NewLogger(path)
	if err != nil {
		logrus.Fatalf("Error creating DNS logger with file: '%s'. %s", path, err)
	}
	return logger
}

// Logger logs the DNS request and response in a configurable file.
type Logger struct {
	log *log.Logger
}

// NewLogger creates an instance of the logger.
// It configures the file path where the DNS request and response are written.
func NewLogger(path string) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{
		log: log.New(file, "", log.Ldate|log.Lmicroseconds|log.LUTC),
	}, nil
}

// LogRequest logs a DNS request in the configured log file.
func (l Logger) LogRequest(request *dns.Msg, corr string) {
	l.log.Printf("Request [%s]:\n%+v\n\n", corr, request)
}

// LogResponse logs a DNS response in the configured log file.
func (l Logger) LogResponse(response *dns.Msg, corr string) {
	l.log.Printf("Response [%s]:\n%+v\n\n", corr, response)
}

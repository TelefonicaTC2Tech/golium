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
	"github.com/TelefonicaTC2Tech/golium"
	"github.com/miekg/dns"
)

var dnsLog *Logger

// Logger logs the DNS request and response in a configurable file.
type Logger struct {
	Log *golium.Logger
}

// GetLogger returns the logger for DNS requests and responses.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	name := "dns"
	if dnsLog == nil {
		dnsLog = &Logger{Log: golium.LoggerFactory(name)}
	}
	return dnsLog
}

// LogRequest logs a DNS request in the configured log file.
func (l Logger) LogRequest(request *dns.Msg, corr string) {
	l.Log.Printf("Request [%s]:\n%+v\n\n", corr, request)
}

// LogResponse logs a DNS response in the configured log file.
func (l Logger) LogResponse(response *dns.Msg, corr string) {
	l.Log.Printf("Response [%s]:\n%+v\n\n", corr, response)
}

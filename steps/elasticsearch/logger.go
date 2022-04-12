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

package elasticsearch

import (
	"github.com/TelefonicaTC2Tech/golium"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var elasticLog *Logger

// Logger logs the elasticsearch requests and responses in a configurable file.
type Logger struct {
	Log *golium.Logger
}

// GetLogger returns the logger for elasticsearch requests and responses.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	name := "elasticsearch"
	if elasticLog == nil {
		elasticLog = &Logger{Log: golium.LoggerFactory(name)}
	}
	return elasticLog
}

// LogCreateIndex logs a creation in elasticsearch in the configured log file.
func (l Logger) LogCreateIndex(res *esapi.Response, document, index, corr string) {
	l.Log.Printf("Create index '%s' [%s]:\n%s\n\n", index, corr, document)
	l.logResponse(res, corr)
}

// LogSearchIndex logs a search in elasticsearch in the configured log file.
func (l Logger) LogSearchIndex(res *esapi.Response, body, index, corr string) {
	l.Log.Printf("Search index '%s' [%s] with body:\n%s\n\n", index, corr, body)
	l.logResponse(res, corr)
}

func (l Logger) logResponse(res *esapi.Response, corr string) {
	l.Log.Printf("Response [%s]:\n%s\n\n", corr, res.String())
}

// LogError logs a creation in elasticsearch in the configured log file.
func (l Logger) LogError(err error, corr string) {
	l.Log.Printf("Error [%s]:\n%s\n\n", corr, err.Error())
}

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
	"log"
	"os"
	"path"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
)

var logger *Logger

// GetLogger returns the logger for elasticsearch requests and responses.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	if logger != nil {
		return logger
	}
	dir := golium.GetConfig().Log.Directory
	path := path.Join(dir, "elasticsearch.log")
	logger, err := NewLogger(path)
	if err != nil {
		logrus.Fatalf("Error creating elasticsearch logger with file: '%s'. %s", path, err)
	}
	return logger
}

// Logger logs the elasticsearch requests and responses in a configurable file.
type Logger struct {
	log *log.Logger
}

// NewLogger creates an instance of the logger.
// It configures the file path where the elasticsearch requests and responses are written.
func NewLogger(path string) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{
		log: log.New(file, "", log.Ldate|log.Lmicroseconds|log.LUTC),
	}, nil
}

// LogCreateIndex logs a creation in elasticsearch in the configured log file.
func (l Logger) LogCreateIndex(res *esapi.Response, document, index, corr string) {
	l.log.Printf("Create index '%s' [%s]:\n%s\n\n", index, corr, document)
	l.logResponse(res, corr)
}

// LogSearchIndex logs a search in elasticsearch in the configured log file.
func (l Logger) LogSearchIndex(res *esapi.Response, body, index, corr string) {
	l.log.Printf("Search index '%s' [%s] with body:\n%s\n\n", index, corr, body)
	l.logResponse(res, corr)
}

func (l Logger) logResponse(res *esapi.Response, corr string) {
	l.log.Printf("Response [%s]:\n%s\n\n", corr, res.String())
}

// LogError logs a creation in elasticsearch in the configured log file.
func (l Logger) LogError(err error, corr string) {
	l.log.Printf("Error [%s]:\n%s\n\n", corr, err.Error())
}

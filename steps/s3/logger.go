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

package s3steps

import (
	"github.com/TelefonicaTC2Tech/golium/logger"
)

var s3Log *Logger

// Logger logs in a configurable file.
type Logger struct {
	Log *logger.Logger
}

// GetLogger returns the logger for s3 operations.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	name := "s3"
	if s3Log == nil {
		s3Log = &Logger{Log: logger.Factory(name)}
	}
	return s3Log
}

// Log a S3 operation
func (l Logger) LogOperation(operation, bucket, key string) {
	l.Log.Printf("Operation: %s in bucket: %s for key: %s", operation, bucket, key)
}

// Log a S3 message
func (l Logger) LogMessage(message string) {
	l.Log.Printf("%s", message)
}

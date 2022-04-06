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

package golium

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

const (
	SUFFIX = ".log"
)

// Logger logs in a configurable file.
type Logger struct {
	*logrus.Logger
}

// LoggerFactory returns a Logger instance.
func LoggerFactory(name string) *Logger {
	file := configureFile(name)
	return builder(*file)
}

// configureFile configures the file where the logs are written.
func configureFile(name string) *os.File {
	dir := GetConfig().Log.Directory
	logPath := path.Join(dir, fmt.Sprintf("%s%s", name, SUFFIX))

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalf("Error creating '%s' logger with file: '%s'. %s", name, logPath, err)
	}
	os.Chmod(file.Name(), 0766)
	return file
}

// Builder creates an instance of the logger.
func builder(file os.File) *Logger {
	level, err := logrus.ParseLevel(GetConfig().Log.Level)
	if err != nil {
		logrus.Fatalf("Error configuring logging level: '%s'. %s", GetConfig().Log.Level, err)
	}
	loggerFormat := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		DisableQuote:    true,
	}
	return &Logger{
		&logrus.Logger{
			Out:       &file,
			Formatter: loggerFormat,
			Hooks:     make(logrus.LevelHooks),
			Level:     level,
		},
	}
}

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

package redis

import (
	"github.com/TelefonicaTC2Tech/golium/logger"
)

var redisLog *Logger

// Logger logs in a configurable file.
type Logger struct {
	Log *logger.Logger
}

// GetLogger returns the logger for HTTP requests and responses.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	name := "redis-pubsub"
	if redisLog == nil {
		redisLog = &Logger{Log: logger.Factory(name)}
	}
	return redisLog
}

// LogSetKey logs a redis SET command.
func (l Logger) LogSetKey(key, value, corr string) {
	l.Log.Printf("Set key %s to [%s]:\n%s\n\n", key, corr, value)
}

// LogHSetKey logs a redis HSET command.
func (l Logger) LogHSetKey(key string, value interface{}, corr string) {
	l.Log.Printf("HSet key %s to [%s]:\n%+v\n\n", key, corr, value)
}

// LogGetKey logs a redis GET command.
func (l Logger) LogGetKey(key, value string, corr string) {
	l.Log.Printf("Get key %s with value[%s]:\n%s\n\n", key, corr, value)
}

// LogHGetKey logs a redis HGET command.
func (l Logger) LogHGetKey(key string, value interface{}, corr string) {
	l.Log.Printf("HGet key %s with value[%s]:\n%+v\n\n", key, corr, value)
}

// LogExistsKey logs a redis EXISTS command.
func (l Logger) LogExistsKey(key string, exits int, corr string) {
	l.Log.Printf("Exists key %s with value[%s]:\n%d\n\n", key, corr, exits)
}

// LogPublishedMessage logs a redis message published to a topic.
func (l Logger) LogPublishedMessage(msg, topic, corr string) {
	l.Log.Printf("Publish to %s [%s]:\n%s\n\n", topic, corr, msg)
}

// LogReceivedMessage logs a redis message received from a topic.
func (l Logger) LogReceivedMessage(msg, topic, corr string) {
	l.Log.Printf("Received from %s [%s]:\n%s\n\n", topic, corr, msg)
}

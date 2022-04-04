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

package rabbit

import (
	"github.com/TelefonicaTC2Tech/golium/logger"
)

var rabbitLog *Logger

// Logger logs in a configurable file.
type Logger struct {
	Log *logger.Logger
}

// GetLogger returns the logger for rabbit messages in publish/subscribe.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	name := "rabbit-pubsub"
	if rabbitLog == nil {
		rabbitLog = &Logger{Log: logger.Factory(name)}
	}
	return rabbitLog
}

// LogPublishedMessage logs a rabbit message published to a topic.
func (l Logger) LogPublishedMessage(msg, topic, corr string) {
	l.Log.Printf("Publish to %s [%s]:\n%s\n\n", topic, corr, msg)
}

// LogReceivedMessage logs a rabbit message received from a topic.
func (l Logger) LogReceivedMessage(msg, topic, corr string) {
	l.Log.Printf("Received from %s [%s]:\n%s\n\n", topic, corr, msg)
}

// LogSubscribedTopic logs the subscription to a rabbit topic.
func (l Logger) LogSubscribedTopic(topic string) {
	l.Log.Printf("Subscribed to %s \n\n", topic)
}

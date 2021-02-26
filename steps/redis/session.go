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
	"context"
	"fmt"
	"time"

	"github.com/Telefonica/golium"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/sjson"
)

// Session contains the information of a redis session.
type Session struct {
	Client *redis.Client
	// Messages received from the publish/subscribe channel
	Messages []string
	// Correlator is used to correlate the messages for a specific session
	Correlator string
	// Redis PubSub.
	// It is stored after subscription to close the subscription.
	pubsub *redis.PubSub
}

// ConfigureClient creates a redis client based on the configuration in options.
func (s *Session) ConfigureClient(ctx context.Context, options *redis.Options) error {
	s.Client = redis.NewClient(options)
	if err := s.Client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("Error configuring client: %+v. %s", options, err)
	}
	s.Correlator = uuid.New().String()
	return nil
}

// SubscribeTopic subscribes to a redis topic to receive messages via a channel.
func (s *Session) SubscribeTopic(ctx context.Context, topic string) error {
	s.pubsub = s.Client.Subscribe(ctx, topic)
	if _, err := s.pubsub.Receive(ctx); err != nil {
		return fmt.Errorf("Error receiving messages from the topic %s. %s", topic, err)
	}
	channel := s.pubsub.Channel()
	go func() {
		logrus.Debugf("Receiving messages from topic %s...", topic)
		for msg := range channel {
			GetLogger().LogReceivedMessage(msg.Payload, topic, s.Correlator)
			s.Messages = append(s.Messages, msg.Payload)
		}
		logrus.Debugf("Stop receiving messages from topic %s", topic)
	}()
	return nil
}

// UnsubscribeTopic unsubscribes from a redis topic to stop receiving messages via a channel.
// The channel is closed.
// If this method is not invoked, then the goroutine created with SubscribeTopic is never closed
// and will permanently process messages from the topic until the program is finished.
func (s *Session) UnsubscribeTopic(ctx context.Context, topic string) error {
	return s.pubsub.Close()
}

// PublishTextMessage publishes a text message in a redis topic.
func (s *Session) PublishTextMessage(ctx context.Context, topic, message string) error {
	GetLogger().LogPublishedMessage(message, topic, s.Correlator)
	if err := s.Client.Publish(ctx, topic, message).Err(); err != nil {
		return fmt.Errorf("Error publishing the message '%s' to topic '%s'. %s", message, topic, err)
	}
	return nil
}

// PublishJSONMessage publishes a JSON message in a redis topic.
func (s *Session) PublishJSONMessage(ctx context.Context, topic string, props map[string]interface{}) error {
	var json string
	var err error
	for key, value := range props {
		if json, err = sjson.Set(json, key, value); err != nil {
			return fmt.Errorf("Error setting property '%s' with value '%s' in the message. %s", key, value, err)
		}
	}
	return s.PublishTextMessage(ctx, topic, json)
}

// WaitForTextMessage waits up to timeout till the expected message is found in the received messages
// for this session.
func (s *Session) WaitForTextMessage(ctx context.Context, timeout time.Duration, expectedMsg string) error {
	return waitUpTo(timeout, func() error {
		for _, msg := range s.Messages {
			if msg == expectedMsg {
				return nil
			}
		}
		return fmt.Errorf("Not received message: %s", expectedMsg)
	})
}

// WaitForJSONMessageWithProperties waits 1 second and verifies if there is a message received
// in the topic with the requested properties.
func (s *Session) WaitForJSONMessageWithProperties(ctx context.Context, timeout time.Duration, props map[string]interface{}) error {
	return waitUpTo(timeout, func() error {
		for _, msg := range s.Messages {
			logrus.Debugf("Checking message: %s", msg)
			m := golium.NewMapFromJSONBytes([]byte(msg))
			found := func() bool {
				for key, expectedValue := range props {
					value := m.Get(key)
					if value != expectedValue {
						logrus.Debugf("Invalid value: %+v. Expected: %+v", value, expectedValue)
						return false
					}
				}
				return true
			}()
			if found {
				return nil
			}
		}
		return fmt.Errorf("Not received message with JSON properties: %+v", props)
	})
}

func waitUpTo(timeout time.Duration, f func() error) error {
	end := time.Now().Add(timeout)
	var err error
	for time.Now().Before(end) {
		time.Sleep(10 * time.Millisecond)
		if err = f(); err == nil {
			break
		}
	}
	return err
}

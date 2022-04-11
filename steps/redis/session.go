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

	"github.com/TelefonicaTC2Tech/golium"
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
	// Time to live (in milliseconds). If 0, there is no TTL configured for new records
	TTL int
	// Redis PubSub.
	// It is stored after subscription to close the subscription.
	pubsub             *redis.PubSub
	RedisClientService ClientFunctions
}

// ConfigureClient creates a redis client based on the configuration in options.
func (s *Session) ConfigureClient(ctx context.Context, options *redis.Options) error {
	s.Client = redis.NewClient(options)
	if err := s.RedisClientService.Ping(ctx, s.Client); err != nil {
		return fmt.Errorf("failed configuring client '%+v': %w", options, err)
	}
	s.Correlator = uuid.New().String()
	return nil
}

// SelectDatabase select redis database for current connection
func (s *Session) SelectDatabase(ctx context.Context, db int) error {
	if err := s.RedisClientService.Do(ctx, s.Client, "select", db); err != nil {
		return err
	}
	return nil
}

// ConfigureTTL saves TTL (in milliseconds) to apply when setting a value in redis.
func (s *Session) ConfigureTTL(ctx context.Context, ttl int) {
	s.TTL = ttl
}

// SetTextValue sets a redis key with a text value.
// It uses session TTL to establish an expiration time (no expiration if TTL is 0).
func (s *Session) SetTextValue(ctx context.Context, key, value string) error {
	expiration := time.Duration(s.TTL * int(time.Millisecond))
	if err := s.RedisClientService.Set(ctx, s.Client, key, value, expiration); err != nil {
		return err
	}
	GetLogger().LogSetKey(key, value, s.Correlator)
	return nil
}

// SetHashValue sets a redis key with a mapped value.
// It uses session TTL to establish an expiration time (no expiration if TTL is 0).
func (s *Session) SetHashValue(
	ctx context.Context,
	key string, value map[string]interface{},
) error {
	err := s.RedisClientService.HSet(ctx, s.Client, key, value)
	if err != nil {
		return err
	}
	if s.TTL == 0 {
		return nil
	}
	expiration := time.Duration(s.TTL * int(time.Millisecond))
	if err := s.RedisClientService.PExpire(ctx, s.Client, key, expiration); err != nil {
		return err
	}
	GetLogger().LogHSetKey(key, value, s.Correlator)
	return nil
}

// SetJSONValue sets a redis key with a JSON document extracted from a table of properties.
func (s *Session) SetJSONValue(
	ctx context.Context,
	key string, props map[string]interface{},
) error {
	var json string
	var err error
	for key, value := range props {
		if json, err = sjson.Set(json, key, value); err != nil {
			return fmt.Errorf(
				"failed setting property '%s' with value '%s' in the request body: %w",
				key, value, err)
		}
	}
	return s.SetTextValue(ctx, key, json)
}

// ValidateTextValue checks if the text value for a redis key equals the expected value.
// It uses session TTL to establish an expiration time (no expiration if TTL is 0).
func (s *Session) ValidateTextValue(ctx context.Context, key, expectedValue string) error {
	err := s.ValidateNonExistantKey(ctx, key)
	if err == nil {
		return fmt.Errorf("failed validating key: key '%s' does not exist", key)
	}
	value, err := s.RedisClientService.Get(ctx, s.Client, key)
	if err != nil {
		return err
	}
	GetLogger().LogGetKey(key, value, s.Correlator)
	if expectedValue != value {
		return fmt.Errorf(
			"mismatch value for key '%s': expected value '%s', actual value '%s'",
			key, expectedValue, value)
	}
	return nil
}

// ValidateHashValue checks if the mapped value for a redis key equals the expected value.
// It uses session TTL to establish an expiration time (no expiration if TTL is 0).
func (s *Session) ValidateHashValue(
	ctx context.Context,
	key string,
	props map[string]interface{},
) error {
	err := s.ValidateNonExistantKey(ctx, key)
	if err == nil {
		return fmt.Errorf("failed validating key: key '%s' does not exist", key)
	}
	m, err := s.RedisClientService.HGetAll(ctx, s.Client, key)
	if err != nil {
		return err
	}
	GetLogger().LogHGetKey(key, m, s.Correlator)
	for key, expectedValue := range props {
		value, found := m[key]
		if !found {
			return fmt.Errorf("missing property '%s': expected '%s'", key, expectedValue)
		}
		if value != expectedValue {
			return fmt.Errorf(
				"mismatch of json property '%s': expected '%s', actual '%s'",
				key, expectedValue, value)
		}
	}
	return nil
}

// ValidateJSONValue checks if the JSON value for a redis key complies with the table of properties.
func (s *Session) ValidateJSONValue(
	ctx context.Context,
	key string,
	props map[string]interface{},
) error {
	err := s.ValidateNonExistantKey(ctx, key)
	if err == nil {
		return fmt.Errorf("failed validating key: key '%s' does not exist", key)
	}
	value, err := s.RedisClientService.Get(ctx, s.Client, key)
	if err != nil {
		return err
	}
	GetLogger().LogGetKey(key, value, s.Correlator)
	m := golium.NewMapFromJSONBytes([]byte(value))
	for key, expectedValue := range props {
		value := m.Get(key)
		if value != expectedValue {
			return fmt.Errorf(
				"mismatch of json property '%s': expected '%s', actual '%s'",
				key, expectedValue, value)
		}
	}
	return nil
}

// ValidateNonExistantKey checks if the redis key has not value.
func (s *Session) ValidateNonExistantKey(ctx context.Context, key string) error {
	exists, err := s.RedisClientService.Exists(ctx, s.Client, key)
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	GetLogger().LogExistsKey(key, int(exists), s.Correlator)
	if exists == 0 {
		return nil
	}
	return fmt.Errorf("redis key '%s' exists", key)
}

// SubscribeTopic subscribes to a redis topic to receive messages via a channel.
func (s *Session) SubscribeTopic(ctx context.Context, topic string) error {
	s.pubsub = s.RedisClientService.Subscribe(ctx, s.Client, topic)
	if _, err := s.RedisClientService.PubSubReceive(ctx, s.pubsub); err != nil {
		return fmt.Errorf("failed receiving messages from the topic '%s': %w", topic, err)
	}
	channel := s.RedisClientService.PubSubChannel(s.pubsub)
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
	return s.RedisClientService.PubSubClose(s.pubsub)
}

// PublishTextMessage publishes a text message in a redis topic.
func (s *Session) PublishTextMessage(ctx context.Context, topic, message string) error {
	GetLogger().LogPublishedMessage(message, topic, s.Correlator)
	if err := s.RedisClientService.Publish(ctx, s.Client, topic, message); err != nil {
		return fmt.Errorf("failed publishing the message '%s' to topic '%s': %w", message, topic, err)
	}
	return nil
}

// PublishJSONMessage publishes a JSON message in a redis topic.
func (s *Session) PublishJSONMessage(
	ctx context.Context,
	topic string,
	props map[string]interface{},
) error {
	var json string
	var err error
	for key, value := range props {
		if json, err = sjson.Set(json, key, value); err != nil {
			return fmt.Errorf(
				"failed setting property '%s' with value '%s' in the message: %w", key, value, err)
		}
	}
	return s.PublishTextMessage(ctx, topic, json)
}

// WaitForTextMessage waits up to timeout till the expected message
// is found in the received messages for this session.
func (s *Session) WaitForTextMessage(
	ctx context.Context,
	timeout time.Duration,
	expectedMsg string,
) error {
	return waitUpTo(timeout, func() error {
		for _, msg := range s.Messages {
			if msg == expectedMsg {
				return nil
			}
		}
		return fmt.Errorf("not received message '%s'", expectedMsg)
	})
}

// WaitForJSONMessageWithProperties waits 1 second and verifies if there is a message received
// in the topic with the requested properties.
func (s *Session) WaitForJSONMessageWithProperties(
	ctx context.Context,
	timeout time.Duration,
	props map[string]interface{},
) error {
	return waitUpTo(timeout, func() error {
		for _, msg := range s.Messages {
			logrus.Debugf("Checking message: %s", msg)
			if matchMessage(msg, props) {
				return nil
			}
		}
		return fmt.Errorf("not received message with JSON properties '%+v'", props)
	})
}

func matchMessage(msg string, expectedProps map[string]interface{}) bool {
	m := golium.NewMapFromJSONBytes([]byte(msg))
	for key, expectedValue := range expectedProps {
		value := m.Get(key)
		if value != expectedValue {
			logrus.Debugf("Invalid value: %+v. Expected: %+v", value, expectedValue)
			return false
		}
	}
	return true
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

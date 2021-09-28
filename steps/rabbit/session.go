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
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Telefonica/golium"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/tidwall/sjson"
)

// Session contains the information of a rabbit session.
type Session struct {
	Connection *amqp.Connection
	// Messages received from the publish/subscribe channel
	Messages []amqp.Delivery
	// Correlator is used to correlate the messages for a specific session
	Correlator string
	// rabbit channel
	// It is stored after subscription to close the subscription
	channel *amqp.Channel
	// rabbit subscription channel
	subCh <-chan amqp.Delivery
	// rabbit headers to store specific data
	headers amqp.Table
	// rabbit publishing message
	publishing amqp.Publishing
	// rabbit received delivery message
	msg amqp.Delivery
}

// ConfigureConnection creates a rabbit connection based on the URI.
func (s *Session) ConfigureConnection(ctx context.Context, uri string) error {
	var err error
	s.Connection, err = amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("failed configuring connection '%s': %w", uri, err)
	}
	s.Correlator = uuid.New().String()
	s.headers = amqp.Table{}
	s.publishing = amqp.Publishing{}
	return nil
}

// ConfigureHeaders stores a table of rabbit headers in the application context.
func (s *Session) ConfigureHeaders(ctx context.Context, headers map[string]interface{}) error {
	s.headers = headers
	if err := s.headers.Validate(); err != nil {
		return errors.Wrap(err, "failed setting rabbit headers")
	}
	return nil
}

// ConfigureStandardProperties stores a table of rabbit properties in the application context.
func (s *Session) ConfigureStandardProperties(ctx context.Context, props amqp.Publishing) error {
	s.publishing = props
	return nil
}

// SubscribeTopic subscribes to a rabbit topic to receive messages via a channel.
func (s *Session) SubscribeTopic(ctx context.Context, topic string) error {
	GetLogger().LogSubscribedTopic(topic)
	var err error
	s.channel, err = s.Connection.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}
	err = s.channel.ExchangeDeclare(
		topic,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare an exchange")
	}
	q, err := s.channel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare a queue")
	}
	err = s.channel.QueueBind(
		q.Name, // queue name
		"",     // routing key
		topic,  // exchange
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to bind a queue")
	}
	s.subCh, err = s.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	go func() {
		logrus.Debugf("Receiving messages from topic %s...", topic)
		for msg := range s.subCh {
			GetLogger().LogReceivedMessage(string(msg.Body), topic, s.Correlator)
			s.Messages = append(s.Messages, msg)
		}
		logrus.Debugf("Stop receiving messages from topic %s", topic)
	}()
	if err != nil {
		return errors.Wrap(err, "failed to register a consumer")
	}
	return nil
}

// Unsubscribe unsubscribes from rabbit closing the channel asociated.
// If this method is not invoked, then the goroutine created with SubscribeTopic is never closed
// and will permanently processing messages from the topic until the program is finished.
func (s *Session) Unsubscribe(ctx context.Context) error {
	if s.channel == nil {
		return nil
	}
	return s.channel.Close()
}

// PublishTextMessage publishes a text message in a rabbit topic.
func (s *Session) PublishTextMessage(ctx context.Context, topic, message string) error {
	GetLogger().LogPublishedMessage(message, topic, s.Correlator)
	var err error
	s.channel, err = s.Connection.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}
	err = s.channel.ExchangeDeclare(
		topic,    // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange")
	}
	publishing := s.buildPublishingMessage([]byte(message))
	err = s.channel.Publish(
		topic,      // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		publishing, // publishing
	)
	if err != nil {
		return fmt.Errorf("failed publishing the message '%s' to topic '%s': %w", message, topic, err)
	}
	return nil
}

func (s *Session) buildPublishingMessage(body []byte) amqp.Publishing {
	publishing := s.publishing
	if publishing.CorrelationId == "" {
		publishing.CorrelationId = s.Correlator
	}
	if publishing.MessageId == "" {
		publishing.MessageId = uuid.NewString()
	}
	if publishing.ContentType == "" {
		publishing.ContentType = "text/plain"
	}
	publishing.Headers = s.headers
	publishing.Body = []byte(body)
	return publishing
}

// PublishJSONMessage publishes a JSON message in a rabbit topic.
func (s *Session) PublishJSONMessage(ctx context.Context, topic string, props map[string]interface{}) error {
	var json string
	var err error
	for key, value := range props {
		if json, err = sjson.Set(json, key, value); err != nil {
			return fmt.Errorf("failed setting property '%s' with value '%s' in the message: %w", key, value, err)
		}
	}
	s.publishing.ContentType = "application/json"
	return s.PublishTextMessage(ctx, topic, json)
}

// WaitForTextMessage waits up to timeout until the expected message is found in the received messages
// for this session.
func (s *Session) WaitForTextMessage(ctx context.Context, timeout time.Duration, expectedMsg string) error {
	return waitUpTo(timeout, func() error {
		for _, msg := range s.Messages {
			if string(msg.Body) == expectedMsg {
				s.msg = msg
				return nil
			}
		}
		return fmt.Errorf("not received message '%s'", expectedMsg)
	})
}

// WaitForJSONMessageWithProperties waits up to timeout and verifies if there is a message received
// in the topic with the requested properties.
func (s *Session) WaitForJSONMessageWithProperties(ctx context.Context, timeout time.Duration, props map[string]interface{}) error {
	return waitUpTo(timeout, func() error {
		for _, msg := range s.Messages {
			logrus.Debugf("Checking message: %s", msg.Body)
			if matchMessage(string(msg.Body), props) {
				s.msg = msg
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

// WaitForMessagesWithStandardProperties waits for 'count' messages with standard rabbit properties that are equal to the expected values.
func (s *Session) WaitForMessagesWithStandardProperties(ctx context.Context, timeout time.Duration, count int, props amqp.Delivery) error {
	return waitUpTo(timeout, func() error {
		err := fmt.Errorf("no message(s) received match(es) the standard properties")
		if count < 0 {
			return err
		}
		for _, msg := range s.Messages {
			logrus.Debugf("Checking message: %s", msg.Body)
			s.msg = msg
			if err := s.ValidateMessageStandardProperties(ctx, props); err == nil {
				count--
				if count == 0 {
					return nil
				}
			}
		}
		return err
	})
}

// ValidateMessageStandardProperties checks if the message standard rabbit properties are equal the expected values.
func (s *Session) ValidateMessageStandardProperties(ctx context.Context, props amqp.Delivery) error {
	msg := reflect.ValueOf(s.msg)
	expectedMsg := reflect.ValueOf(props)
	t := expectedMsg.Type()
	for i := 0; i < expectedMsg.NumField(); i++ {
		if !expectedMsg.Field(i).IsZero() {
			key := t.Field(i).Name
			value := msg.Field(i).Interface()
			expectedValue := expectedMsg.Field(i).Interface()
			if value != expectedValue {
				return fmt.Errorf("mismatch of standard rabbit property '%s': expected '%s', actual '%s'", key, expectedValue, value)
			}
		}
	}
	return nil
}

// ValidateMessageHeaders checks if the message rabbit headers are equal the expected values.
func (s *Session) ValidateMessageHeaders(ctx context.Context, headers map[string]interface{}) error {
	h := s.msg.Headers
	for key, expectedValue := range headers {
		value, found := h[key]
		if !found {
			return fmt.Errorf("missing rabbit message header '%s'", key)
		}
		if value != expectedValue {
			return fmt.Errorf("mismatch of standard rabbit property '%s': expected '%s', actual '%s'", key, expectedValue, value)
		}
	}
	return nil
}

// ValidateMessageTextBody checks if the message text body is equal to the expected value.
func (s *Session) ValidateMessageTextBody(ctx context.Context, expectedMsg string) error {
	msg := string(s.msg.Body)
	if msg != expectedMsg {
		return fmt.Errorf("mismatch of message text: expected '%s', actual '%s'", expectedMsg, msg)
	}
	return nil
}

// ValidateMessageJSONBody checks if the message json body properties of message in position 'pos' are equal the expected values.
func (s *Session) ValidateMessageJSONBody(ctx context.Context, props map[string]interface{}, pos int) error {
	nMessages := len(s.Messages)
	if pos < 0 || pos >= nMessages {
		return fmt.Errorf("trying to validate message in position: '%d', '%d' messages available", pos, nMessages)
	}
	m := golium.NewMapFromJSONBytes([]byte(s.Messages[pos].Body))
	for key, expectedValue := range props {
		value := m.Get(key)
		if value != expectedValue {
			return fmt.Errorf("mismatch of json property '%s': expected '%s', actual '%s'", key, expectedValue, value)
		}
	}
	return nil
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

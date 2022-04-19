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

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/tidwall/sjson"
)

const (
	convertTableToMapMessage = "failed processing table to a map for the rabbit message: "
	convertTableToMapBody    = "failed processing table to a map for the request body: "
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
	// ampq service
	AMQPServiceClient AMQPServiceFunctions
}

// ConfigureConnection creates a rabbit connection based on the URI.
func (s *Session) ConfigureConnection(ctx context.Context, uri string) error {
	var err error
	s.Connection, err = s.AMQPServiceClient.Dial(uri)
	if err != nil {
		return fmt.Errorf("failed configuring connection '%s': %w", uri, err)
	}
	s.Correlator = uuid.New().String()
	s.headers = amqp.Table{}
	s.publishing = amqp.Publishing{}
	return nil
}

// ConfigureHeaders stores a table of rabbit headers in the application context.
func (s *Session) ConfigureHeaders(ctx context.Context, t *godog.Table) error {
	headers, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return fmt.Errorf(convertTableToMapBody+"%w", err)
	}
	s.headers = headers
	if err := s.headers.Validate(); err != nil {
		return errors.Wrap(err, "failed setting rabbit headers")
	}
	return nil
}

// ConfigureStandardProperties stores a table of rabbit properties in the application context.
func (s *Session) ConfigureStandardProperties(ctx context.Context, t *godog.Table) error {
	var props amqp.Publishing
	if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &props); err != nil {
		return fmt.Errorf("failed configuring rabbit endpoint: %w", err)
	}
	s.publishing = props
	return nil
}

// SubscribeTopic subscribes to a rabbit topic to receive messages via a channel.
func (s *Session) SubscribeTopic(ctx context.Context, topic string) error {
	GetLogger().LogSubscribedTopic(topic)
	var err error
	s.channel, err = s.AMQPServiceClient.ConnectionChannel(s.Connection)
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}
	err = s.AMQPServiceClient.ChannelExchangeDeclare(s.channel, topic)
	if err != nil {
		return errors.Wrap(err, "failed to declare an exchange")
	}
	q, err := s.AMQPServiceClient.ChannelQueueDeclare(s.channel)
	if err != nil {
		return errors.Wrap(err, "failed to declare a queue")
	}
	err = s.AMQPServiceClient.ChannelQueueBind(s.channel, q.Name, topic)
	if err != nil {
		return errors.Wrap(err, "failed to bind a queue")
	}
	s.subCh, err = s.AMQPServiceClient.ChannelConsume(s.channel, q.Name)
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

// Unsubscribe unsubscribes from rabbit closing the channel associated.
// If this method is not invoked, then the goroutine created with SubscribeTopic is never closed
// and will permanently processing messages from the topic until the program is finished.
func (s *Session) Unsubscribe(ctx context.Context) error {
	if s.channel == nil {
		return nil
	}
	return s.AMQPServiceClient.ChannelClose(s.channel)
}

// PublishTextMessage publishes a text message in a rabbit topic.
func (s *Session) PublishTextMessage(ctx context.Context, topic, message string) error {
	GetLogger().LogPublishedMessage(message, topic, s.Correlator)
	var err error
	s.channel, err = s.AMQPServiceClient.ConnectionChannel(s.Connection)
	if err != nil {
		return errors.Wrap(err, "failed to open a channel")
	}
	err = s.AMQPServiceClient.ChannelExchangeDeclare(s.channel, topic)
	if err != nil {
		return fmt.Errorf("failed to declare an exchange")
	}
	publishing := s.buildPublishingMessage([]byte(message))
	err = s.AMQPServiceClient.ChannelPublish(s.channel, topic, publishing)
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
	publishing.Body = body
	return publishing
}

// PublishJSONMessage publishes a JSON message in a rabbit topic.
func (s *Session) PublishJSONMessage(
	ctx context.Context,
	topic string,
	t *godog.Table,
) error {
	props, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return fmt.Errorf(convertTableToMapBody+"%w", err)
	}
	var json string
	for key, value := range props {
		if json, err = sjson.Set(json, key, value); err != nil {
			return fmt.Errorf("failed setting property '%s' with value '%s' in the message: %w",
				key, value, err)
		}
	}
	s.publishing.ContentType = "application/json"
	return s.PublishTextMessage(ctx, topic, json)
}

// WaitForTextMessage waits up to timeout until the expected message is found in
// the received messages for this session.
func (s *Session) WaitForTextMessage(ctx context.Context,
	timeout time.Duration,
	expectedMsg string,
) error {
	return waitUpTo(timeout, func() error {
		for i := range s.Messages {
			if string(s.Messages[i].Body) == expectedMsg {
				s.msg = s.Messages[i]
				return nil
			}
		}
		return fmt.Errorf("not received message '%s'", expectedMsg)
	})
}

// WaitForJSONMessageWithProperties waits up to timeout and verifies if there is a message received
// in the topic with the requested properties.
// When wantError is set to true
func (s *Session) WaitForJSONMessageWithProperties(ctx context.Context,
	timeout time.Duration,
	t *godog.Table,
	wantError bool,
) error {
	props, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return fmt.Errorf(convertTableToMapMessage+"%w", err)
	}
	return waitUpTo(timeout, func() error {
		for i := range s.Messages {
			logrus.Debugf("Checking message: %s", s.Messages[i].Body)
			if matchMessage(string(s.Messages[i].Body), props) {
				s.msg = s.Messages[i]
				if !wantError {
					return nil
				}
				return fmt.Errorf("received message with JSON properties '%+v'", props)
			}
		}
		if !wantError {
			return fmt.Errorf("not received message with JSON properties '%+v'", props)
		}
		return nil
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

// WaitForMessagesWithStandardProperties waits for 'count' messages with standard rabbit properties
// that are equal to the expected values.
func (s *Session) WaitForMessagesWithStandardProperties(
	ctx context.Context,
	timeout time.Duration,
	count int,
	t *godog.Table,
	wantErr bool,
) error {
	return waitUpTo(timeout, func() error {
		err := fmt.Errorf("no message(s) received match(es) the standard properties")
		if count < 0 {
			return err
		}
		for i := range s.Messages {
			logrus.Debugf("Checking message: %s", s.Messages[i].Body)
			s.msg = s.Messages[i]
			if err = s.ValidateMessageStandardProperties(ctx, t, wantErr); err == nil {
				count--
				if count == 0 {
					return nil
				}
			}
		}
		return err
	})
}

// ValidateMessageStandardProperties checks if the message standard rabbit properties are equal
// the expected values.
func (s *Session) ValidateMessageStandardProperties(
	ctx context.Context,
	table *godog.Table,
	wantErr bool,
) error {
	var props amqp.Delivery
	if err := golium.ConvertTableWithoutHeaderToStruct(ctx, table, &props); err != nil {
		return fmt.Errorf("failed configuring rabbit endpoint: %w", err)
	}
	msg := reflect.ValueOf(s.msg)
	expectedMsg := reflect.ValueOf(props)
	t := expectedMsg.Type()
	for i := 0; i < expectedMsg.NumField(); i++ {
		if !expectedMsg.Field(i).IsZero() {
			key := t.Field(i).Name
			value := msg.Field(i).Interface()
			expectedValue := expectedMsg.Field(i).Interface()
			if value == expectedValue {
				if !wantErr {
					return nil
				}
				return fmt.Errorf("received a message with standard rabbit properties '%s'", value)
			} else if !wantErr {
				return fmt.Errorf(
					"mismatch of standard rabbit property '%s': expected '%s', actual '%s'",
					key, expectedValue, value)
			}
		}
	}
	return nil
}

// ValidateMessageHeaders checks if the message rabbit headers are equal the expected values.
func (s *Session) ValidateMessageHeaders(
	ctx context.Context,
	t *godog.Table,
) error {
	headers, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return fmt.Errorf(convertTableToMapMessage+"%w", err)
	}
	h := s.msg.Headers
	for key, expectedValue := range headers {
		value, found := h[key]
		if !found {
			return fmt.Errorf("missing rabbit message header '%s'", key)
		}
		if value != expectedValue {
			return fmt.Errorf(
				"mismatch of standard rabbit property '%s': expected '%s', actual '%s'",
				key, expectedValue, value)
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

// ValidateMessageJSONBody checks if the message json body properties of message in position 'pos'
// are equal the expected values.
// if pos == -1 then it means last message stored, that is the one stored in s.msg
func (s *Session) ValidateMessageJSONBody(ctx context.Context,
	t *godog.Table,
	pos int,
) error {
	props, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return fmt.Errorf(convertTableToMapMessage+"%w", err)
	}
	m := golium.NewMapFromJSONBytes(s.msg.Body)
	if pos != -1 {
		nMessages := len(s.Messages)
		if pos < 0 || pos >= nMessages {
			return fmt.Errorf(
				"trying to validate message in position: '%d', '%d' messages available",
				pos, nMessages)
		}
		m = golium.NewMapFromJSONBytes(s.Messages[pos].Body)
	}
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

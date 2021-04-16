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
	"time"

	"github.com/Telefonica/golium"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/tidwall/sjson"
)

// Session contains the information of a rabbit session.
type Session struct {
	Connection *amqp.Connection
	// Messages received from the publish/subscribe channel
	Messages []string
	// Correlator is used to correlate the messages for a specific session
	Correlator string
	// TransactionID is used to identify each transaction.
	TransactionID string
	// rabbit channel.
	// It is stored after subscription to close the subscription.
	channel *amqp.Channel
}

// ConfigureClient creates a rabbit connection based on the URI.
func (s *Session) ConfigureConnection(ctx context.Context, uri string) error {
	var err error
	s.Connection, err = amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("failed configuring connection '%s': %w", uri, err)
	}
	s.Correlator = uuid.New().String()
	s.TransactionID = uuid.New().String()
	return nil
}

// SubscribeTopic subscribes to a rabbit topic to receive messages via a channel.
func (s *Session) SubscribeTopic(ctx context.Context, topic string) error {
	var err error
	s.channel, err = s.Connection.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel")
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
		return fmt.Errorf("Failed to declare an exchange")
	}
	q, err := s.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("Failed to declare a queue")
	}
	err = s.channel.QueueBind(
		q.Name, // queue name
		"",     // routing key
		topic,  // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to bind a queue")
	}
	msgs, err := s.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("Failed to register a consumer")
	}
	go func() {
		logrus.Debugf("Receiving messages from topic %s...", topic)
		for msg := range msgs {
			GetLogger().LogReceivedMessage(string(msg.Body), topic, s.Correlator)
			s.Messages = append(s.Messages, string(msg.Body))
		}
		logrus.Debugf("Stop receiving messages from topic %s", topic)
	}()
	return nil
}

// UnsubscribeTopic unsubscribes from a rabbit topic to stop receiving messages via a channel.
// The channel is closed.
// If this method is not invoked, then the goroutine created with SubscribeTopic is never closed
// and will permanently process messages from the topic until the program is finished.
func (s *Session) UnsubscribeTopic(ctx context.Context, topic string) error {
	return s.channel.Close()
}

// PublishTextMessage publishes a text message in a rabbit topic.
func (s *Session) PublishTextMessage(ctx context.Context, topic, message string) error {
	GetLogger().LogPublishedMessage(message, topic, s.Correlator)
	var err error
	s.channel, err = s.Connection.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel")
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
		return fmt.Errorf("Failed to declare an exchange")
	}
	err = s.channel.Publish(
		topic, // exchange
		"",    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(message),
			CorrelationId: s.Correlator,
			MessageId:     s.TransactionID,
		})
	if err != nil {
		return fmt.Errorf("failed publishing the message '%s' to topic '%s': %w", message, topic, err)
	}
	return nil
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
		return fmt.Errorf("not received message '%s'", expectedMsg)
	})
}

// WaitForJSONMessageWithProperties waits 1 second and verifies if there is a message received
// in the topic with the requested properties.
func (s *Session) WaitForJSONMessageWithProperties(ctx context.Context, timeout time.Duration, props map[string]interface{}) error {
	return waitUpTo(timeout, func() error {
		for _, msg := range s.Messages {
			logrus.Debugf("Checking message: %s", msg)
			m := golium.NewMapFromJSONBytes([]byte(msg))
			for key, expectedValue := range props {
				value := m.Get(key)
				if value != expectedValue {
					return fmt.Errorf("mismatch of json property '%s': expected '%s', actual '%s'", key, expectedValue, value)
				}
				return nil
			}
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

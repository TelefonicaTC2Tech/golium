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
	"github.com/streadway/amqp"
)

type AMQPServiceFunctions interface {
	Dial(url string) (*amqp.Connection, error)
	ConnectionChannel(c *amqp.Connection) (*amqp.Channel, error)
	ChannelExchangeDeclare(channel *amqp.Channel, name string) error
	ChannelQueueDeclare(channel *amqp.Channel) (amqp.Queue, error)
	ChannelQueueBind(channel *amqp.Channel, name, exchange string) error
	ChannelConsume(channel *amqp.Channel, queue string,
	) (<-chan amqp.Delivery, error)
	ChannelClose(channel *amqp.Channel) error
	ChannelPublish(channel *amqp.Channel,
		exchange string,
		msg amqp.Publishing,
	) error
}

type AMQPService struct{}

func NewAMQPService() *AMQPService {
	return &AMQPService{}
}

func (a AMQPService) Dial(url string) (*amqp.Connection, error) {
	return amqp.DialConfig(url, amqp.Config{Properties: make(amqp.Table)})
}

func (a AMQPService) ConnectionChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	return connection.Channel()
}
func (a AMQPService) ChannelExchangeDeclare(channel *amqp.Channel, name string) error {
	return channel.ExchangeDeclare(name, "fanout", true, false, false, false, nil)
}

func (a AMQPService) ChannelQueueDeclare(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare("", false, true, true, false, nil)
}
func (a AMQPService) ChannelQueueBind(channel *amqp.Channel, name, exchange string) error {
	return channel.QueueBind(name, "", exchange, false, nil)
}

func (a AMQPService) ChannelConsume(channel *amqp.Channel, queue string,
) (<-chan amqp.Delivery, error) {
	return channel.Consume(queue, "", true, false, false, false, nil)
}

func (a AMQPService) ChannelClose(channel *amqp.Channel) error {
	return channel.Close()
}

func (a AMQPService) ChannelPublish(channel *amqp.Channel,
	exchange string,
	msg amqp.Publishing,
) error {
	return channel.Publish(exchange, "", false, false, msg)
}

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

var (
	DialError                   error
	ConnectionChannelError      error
	ChannelExchangeDeclareError error
	ChannelQueueDeclareError    error
	ChannelQueueBindError       error
	ChannelConsumeError         error
	MockSubCh                   <-chan amqp.Delivery
	ChannelPublishError         error
)

type AMQPServiceFuncMock struct{}

func (a AMQPServiceFuncMock) Dial(url string) (*amqp.Connection, error) {
	return nil, DialError
}

func (a AMQPServiceFuncMock) ConnectionChannel(c *amqp.Connection) (*amqp.Channel, error) {
	return nil, ConnectionChannelError
}

func (a AMQPServiceFuncMock) ChannelExchangeDeclare(
	channel *amqp.Channel,
	name, kind string,
	durable, autoDelete, internal, noWait bool,
	args amqp.Table,
) error {
	return ChannelExchangeDeclareError
}

func (a AMQPServiceFuncMock) ChannelQueueDeclare(
	channel *amqp.Channel,
	name string,
	durable, autoDelete, exclusive, noWait bool,
	args amqp.Table,
) (amqp.Queue, error) {
	amqpQueue := amqp.Queue{}
	return amqpQueue, ChannelQueueDeclareError
}

func (a AMQPServiceFuncMock) ChannelQueueBind(channel *amqp.Channel,
	name, key, exchange string,
	noWait bool, args amqp.Table,
) error {
	return ChannelQueueBindError
}

func (a AMQPServiceFuncMock) ChannelConsume(channel *amqp.Channel,
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, error) {
	return MockSubCh, ChannelConsumeError
}
func (a AMQPServiceFuncMock) ChannelClose(channel *amqp.Channel) error {
	return nil
}

func (a AMQPServiceFuncMock) ChannelPublish(channel *amqp.Channel,
	exchange, key string,
	mandatory, immediate bool,
	msg amqp.Publishing,
) error {
	return ChannelPublishError
}

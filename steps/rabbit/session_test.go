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
	"os"
	"testing"
	"time"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
	"github.com/streadway/amqp"
)

const (
	rabbitmq = "amqp://guest:guest@localhost:5672/"
	logsPath = "./logs"
)

func TestConfigureConnection(t *testing.T) {
	tests := []struct {
		name      string
		uri       string
		connError error
		wantErr   bool
	}{
		{
			name:      "Dial error",
			connError: fmt.Errorf("dial error"),
			wantErr:   true,
		},
		{
			name:      "Without connection error",
			uri:       rabbitmq,
			connError: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			ctx := InitializeContext(context.Background())
			s.AMQPServiceClient = AMQPServiceFuncMock{}
			DialError = tt.connError
			if err := s.ConfigureConnection(ctx, tt.uri); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigureConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigureHeaders(t *testing.T) {
	var wrongRabbitHeader = make(map[string]interface{})
	wrongRabbitHeader["wrongParam"] = uint(5)

	var rabbitHeader = make(map[string]interface{})
	rabbitHeader["param"] = "value"
	rabbitHeader["Header1"] = "value1"
	rabbitHeader["Header2"] = "Value2"
	tests := []struct {
		name    string
		headers map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Validate headers error",
			wantErr: true,
			headers: wrongRabbitHeader,
		},
		{
			name:    "No error",
			headers: rabbitHeader,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			if err := s.ConfigureHeaders(context.Background(), tt.headers); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigureHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigureStandardProperties(t *testing.T) {
	rabbitHeaders := amqp.Publishing{}

	tests := []struct {
		name      string
		propTable *godog.Table
	}{
		{
			name:      "Configure",
			propTable: golium.NewTable([][]string{{"ContentType"}, {"application/json"}}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			ctx := InitializeContext(context.Background())
			golium.ConvertTableWithoutHeaderToStruct(ctx, tt.propTable, &rabbitHeaders)
			s.ConfigureStandardProperties(context.Background(), rabbitHeaders)
		})
	}
}

func TestSubscribeTopic(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	tests := []struct {
		name                        string
		topic                       string
		connectionChannelError      error
		channelExchangeDeclareError error
		channelQueueDeclareError    error
		channelQueueBindError       error
		channelConsumeError         error
		subCh                       <-chan amqp.Delivery
		wantErr                     bool
	}{
		{
			name:                   "Connection Channel Error",
			connectionChannelError: fmt.Errorf("connection channel error"),
			wantErr:                true,
		},
		{
			name:                        "Channel Exchange Declare Error",
			connectionChannelError:      nil,
			channelExchangeDeclareError: fmt.Errorf("channel exchange declare error"),
			wantErr:                     true,
		},
		{
			name:                        "Channel Queue Declare Error",
			connectionChannelError:      nil,
			channelExchangeDeclareError: nil,
			channelQueueDeclareError:    fmt.Errorf("channel queue declare error"),
			wantErr:                     true,
		},
		{
			name:                        "Channel Queue Bind Error",
			connectionChannelError:      nil,
			channelExchangeDeclareError: nil,
			channelQueueDeclareError:    nil,
			channelQueueBindError:       fmt.Errorf("channel queue bind error"),
			wantErr:                     true,
		},
		{
			name:                        "Channel Consume Error",
			connectionChannelError:      nil,
			channelExchangeDeclareError: nil,
			channelQueueDeclareError:    nil,
			channelQueueBindError:       nil,
			subCh:                       nil,
			channelConsumeError:         fmt.Errorf("channel queue bind error"),
			wantErr:                     true,
		},
		{
			name:                        "Channel registered without errors",
			connectionChannelError:      nil,
			channelExchangeDeclareError: nil,
			channelQueueDeclareError:    nil,
			channelQueueBindError:       nil,
			subCh:                       nil,
			channelConsumeError:         nil,
			wantErr:                     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goliumCtx := golium.InitializeContext(context.Background())
			ctx := InitializeContext(goliumCtx)
			s := &Session{}
			s.AMQPServiceClient = AMQPServiceFuncMock{}
			ConnectionChannelError = tt.connectionChannelError
			ChannelExchangeDeclareError = tt.channelExchangeDeclareError
			ChannelQueueDeclareError = tt.channelQueueDeclareError
			ChannelQueueBindError = tt.channelQueueBindError
			ChannelConsumeError = tt.channelConsumeError
			MockSubCh = tt.subCh
			if err := s.SubscribeTopic(ctx, tt.topic); (err != nil) != tt.wantErr {
				t.Errorf("Session.SubscribeTopic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnsubscribe(t *testing.T) {
	tests := []struct {
		name    string
		channel *amqp.Channel
		wantErr bool
	}{
		{
			name:    "Nil channel",
			channel: nil,
			wantErr: false,
		},
		{
			name:    "Not nil channel",
			channel: &amqp.Channel{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.AMQPServiceClient = AMQPServiceFuncMock{}
			s.channel = tt.channel
			if err := s.Unsubscribe(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Session.Unsubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPublishTextMessage(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	type args struct {
		topic   string
		message string
	}
	tests := []struct {
		name                        string
		args                        args
		connectionChannelError      error
		channelExchangeDeclareError error
		channelPublishError         error
		wantErr                     bool
	}{
		{
			name:                   "Connection channel error",
			connectionChannelError: fmt.Errorf("connection channel error"),
			wantErr:                true,
		},
		{
			name:                        "Channel exchange declare error",
			connectionChannelError:      nil,
			channelExchangeDeclareError: fmt.Errorf("channel exchange declare error"),
			wantErr:                     true,
		},
		{
			name:                        "Publish error",
			connectionChannelError:      nil,
			channelExchangeDeclareError: nil,
			channelPublishError:         fmt.Errorf("channel publish error"),
			args: args{
				message: "test message",
			},
			wantErr: true,
		},
		{
			name:                        "Publish without error",
			connectionChannelError:      nil,
			channelExchangeDeclareError: nil,
			channelPublishError:         nil,
			args: args{
				message: "test message",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, s := setPublishTestEnv(
				tt.connectionChannelError,
				tt.channelExchangeDeclareError,
				tt.channelPublishError)
			if err := s.PublishTextMessage(
				ctx, tt.args.topic, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Session.PublishTextMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPublishJSONMessage(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	propsOk := make(map[string]interface{})
	propsOk["id"] = "1"
	propsOk["name"] = "test"
	type args struct {
		topic string
		props map[string]interface{}
	}
	tests := []struct {
		name                        string
		connectionChannelError      error
		channelExchangeDeclareError error
		channelPublishError         error
		args                        args
		wantErr                     bool
	}{
		{
			name:                        "Valid props",
			connectionChannelError:      nil,
			channelExchangeDeclareError: nil,
			channelPublishError:         nil,
			args: args{
				topic: "test_topic",
				props: propsOk,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, s := setPublishTestEnv(
				tt.connectionChannelError,
				tt.channelExchangeDeclareError,
				tt.channelPublishError)
			if err := s.PublishJSONMessage(ctx, tt.args.topic, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.PublishJSONMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func setPublishTestEnv(
	conChannelError, channelExchDecError, channelPubError error,
) (context.Context, *Session) {
	goliumCtx := golium.InitializeContext(context.Background())
	ctx := InitializeContext(goliumCtx)
	s := &Session{}
	s.AMQPServiceClient = AMQPServiceFuncMock{}
	ConnectionChannelError = conChannelError
	ChannelExchangeDeclareError = channelExchDecError
	ChannelPublishError = channelPubError
	return ctx, s
}

func TestWaitForTextMessage(t *testing.T) {
	type args struct {
		timeout     time.Duration
		expectedMsg string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Expected message found",
			args: args{
				timeout:     time.Duration(5000),
				expectedMsg: "expected string",
			},
			wantErr: false,
		},
		{
			name: "Expected message not found",
			args: args{
				timeout:     time.Duration(5000),
				expectedMsg: "error expected string",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.Messages = []amqp.Delivery{
				{Body: []byte("expected string")},
			}
			if err := s.WaitForTextMessage(
				context.Background(), tt.args.timeout, tt.args.expectedMsg); (err != nil) != tt.wantErr {
				t.Errorf("Session.WaitForTextMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWaitForJSONMessageWithProperties(t *testing.T) {
	expectedJSON := make(map[string]interface{})
	expectedJSON["id"] = "1"
	wrongExpectedJSON := make(map[string]interface{})
	wrongExpectedJSON["id"] = "5"

	type args struct {
		timeout time.Duration
		props   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Expected json found",
			args: args{
				timeout: time.Duration(5000),
				props:   expectedJSON,
			},
			wantErr: false,
		},
		{
			name: "Expected json not found",
			args: args{
				timeout: time.Duration(5000),
				props:   wrongExpectedJSON,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.Messages = []amqp.Delivery{
				{
					Body: []byte(`{"id": "1"}`),
				},
			}
			if err := s.WaitForJSONMessageWithProperties(
				context.Background(), tt.args.timeout, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.WaitForJSONMessageWithProperties() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWaitForMessagesWithStandardProperties(t *testing.T) {
	type args struct {
		timeout time.Duration
		count   int
		props   amqp.Delivery
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Message count < 0",
			args: args{
				count:   -1,
				timeout: time.Duration(5000),
			},
			wantErr: true,
		},
		{
			name: "Matching properties",
			args: args{
				count:   1,
				timeout: time.Duration(5000),
				props: amqp.Delivery{
					Priority: 5,
				},
			},
			wantErr: false,
		},
		{
			name: "Not matching properties",
			args: args{
				count:   1,
				timeout: time.Duration(5000),
				props: amqp.Delivery{
					Priority: 10,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.Messages = []amqp.Delivery{
				{
					Priority: 5,
				},
			}
			if err := s.WaitForMessagesWithStandardProperties(
				context.Background(),
				tt.args.timeout,
				tt.args.count, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf(
					"Session.WaitForMessagesWithStandardProperties() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessageHeaders(t *testing.T) {
	refHeaders := make(amqp.Table)
	refHeaders["id"] = "1"

	testHeaders := make(map[string]interface{})

	tests := []struct {
		name      string
		testKey   string
		testValue string
		wantErr   bool
	}{
		{
			name:      "Key found and value matches",
			testKey:   "id",
			testValue: "1",
			wantErr:   false,
		},
		{
			name:      "Key not found",
			testKey:   "ids",
			testValue: "1",
			wantErr:   true,
		},
		{
			name:      "Key found wrong value",
			testKey:   "id",
			testValue: "2",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.msg = amqp.Delivery{
				Headers: refHeaders,
			}
			testHeaders[tt.testKey] = tt.testValue
			if err := s.ValidateMessageHeaders(
				context.Background(), testHeaders); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateMessageHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessageTextBody(t *testing.T) {
	tests := []struct {
		name        string
		expectedMsg string
		wantErr     bool
	}{
		{
			name:        "Mismatch of message text",
			expectedMsg: "wrong message",
			wantErr:     true,
		},
		{
			name:        "Right message",
			expectedMsg: "message",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.msg = amqp.Delivery{
				Body: []byte("message"),
			}

			if err := s.ValidateMessageTextBody(
				context.Background(), tt.expectedMsg); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateMessageTextBody() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessageJSONBody(t *testing.T) {
	expectedJSON := make(map[string]interface{})
	expectedJSON["id"] = "1"
	wrongExpectedJSON := make(map[string]interface{})
	wrongExpectedJSON["id"] = "5"
	type args struct {
		props map[string]interface{}
		pos   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "JSON body match with pos -1",
			args: args{
				pos:   -1,
				props: expectedJSON,
			},
			wantErr: false,
		},
		{
			name: "JSON body match with pos != -1",
			args: args{
				pos:   0,
				props: expectedJSON,
			},
			wantErr: false,
		},
		{
			name: "pos != -1 without messages",
			args: args{
				pos:   1,
				props: expectedJSON,
			},
			wantErr: true,
		},
		{
			name: "JSON body mismatch with pos -1",
			args: args{
				pos:   -1,
				props: wrongExpectedJSON,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.Messages = []amqp.Delivery{
				{
					Body: []byte(`{"id": "1"}`),
				},
			}
			s.msg = amqp.Delivery{
				Body: []byte(`{"id": "1"}`),
			}

			if err := s.ValidateMessageJSONBody(
				context.Background(), tt.args.props, tt.args.pos); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateMessageJSONBody() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

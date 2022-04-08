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
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	redisEndpoint = "localhost:6379"
	logsPath      = "./logs"
	key           = "key"
	value         = "value"
	correlator    = "correlator"
)

func TestConfigureClient(t *testing.T) {
	tests := []struct {
		name    string
		options *redis.Options
		pingErr error
		wantErr bool
	}{
		{
			name: "Configuration error",
			options: &redis.Options{
				Addr: redisEndpoint,
				DB:   0,
			},
			pingErr: fmt.Errorf("ping error"),
			wantErr: true,
		},
		{
			name: "Configuration done",
			options: &redis.Options{
				Addr: redisEndpoint,
				DB:   0,
			},
			pingErr: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			PingError = tt.pingErr
			if err := s.ConfigureClient(context.Background(), tt.options); (err != nil) != tt.wantErr {
				t.Errorf("Session.ConfigureClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSelectDatabase(t *testing.T) {
	tests := []struct {
		name    string
		db      int
		doErr   error
		wantErr bool
	}{
		{
			name:    "Select error",
			doErr:   fmt.Errorf("select error"),
			wantErr: true,
		},
		{
			name:    "Select without error",
			doErr:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			DoError = tt.doErr
			if err := s.SelectDatabase(context.Background(), tt.db); (err != nil) != tt.wantErr {
				t.Errorf("Session.SelectDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigureTTL(t *testing.T) {
	tests := []struct {
		name string
		ttl  int
	}{
		{
			name: "Configure",
			ttl:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.ConfigureTTL(context.Background(), tt.ttl)
		})
	}
}

func TestSetTextValue(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		setErr  error
		wantErr bool
	}{
		{
			name:    "Set error",
			setErr:  fmt.Errorf("set error"),
			wantErr: true,
		},
		{
			name: "Set without error",
			args: args{
				key:   key,
				value: value,
			},
			setErr:  nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.TTL = 1000
			s.RedisClientService = ClientServiceFuncMock{}
			SetError = tt.setErr
			if err := s.SetTextValue(
				context.Background(), tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Session.SetTextValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetHashValue(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	testValue := make(map[string]interface{})
	testValue[key] = value
	type args struct {
		key   string
		value map[string]interface{}
	}
	tests := []struct {
		name       string
		hSetErr    error
		pExpireErr error
		ttl        int
		args       args
		wantErr    bool
	}{
		{
			name:    "HSet error",
			hSetErr: fmt.Errorf("HSet error"),
			wantErr: true,
		},
		{
			name:    "TTL is 0",
			hSetErr: nil,
			ttl:     0,
			wantErr: false,
		},
		{
			name:       "PExpire error",
			hSetErr:    nil,
			pExpireErr: fmt.Errorf("PExpire error"),
			ttl:        1000,
			wantErr:    true,
		},
		{
			name: "Set Hash without error",
			args: args{
				key:   key,
				value: testValue,
			},
			hSetErr:    nil,
			pExpireErr: nil,
			ttl:        1000,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.TTL = tt.ttl
			s.RedisClientService = ClientServiceFuncMock{}
			HSetError = tt.hSetErr
			PExpireError = tt.pExpireErr
			if err := s.SetHashValue(
				context.Background(), tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Session.SetHashValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetJSONValue(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	testValue := make(map[string]interface{})
	testValue[key] = value
	type args struct {
		key   string
		props map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set JSON without errors",
			args: args{
				key:   key,
				props: testValue,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.TTL = 1000
			s.RedisClientService = ClientServiceFuncMock{}
			SetError = nil
			if err := s.SetJSONValue(
				context.Background(), tt.args.key, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.SetJSONValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNonExistantKey(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)

	tests := []struct {
		name       string
		existErr   error
		existValue int64
		key        string
		wantErr    bool
	}{
		{
			name:     "Exists common error",
			existErr: fmt.Errorf("exists error"),
			wantErr:  true,
		},
		{
			name:     "Exists redis nil",
			existErr: redis.Nil,
			wantErr:  false,
		},
		{
			name:       "Exists value 0",
			existErr:   nil,
			existValue: 0,
			wantErr:    false,
		},
		{
			name:       "Already exists error",
			existErr:   nil,
			existValue: 1,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			s.Correlator = correlator

			ExistsError = tt.existErr
			ExistsValue = tt.existValue
			if err := s.ValidateNonExistantKey(
				context.Background(), tt.key); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateNonExistantKey() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestValidateTextValue(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	type args struct {
		key           string
		expectedValue string
	}
	tests := []struct {
		name      string
		existsErr error
		getErr    error
		getValue  string
		args      args
		wantErr   bool
	}{
		{
			name:      "Key does not exists",
			existsErr: redis.Nil,
			wantErr:   true,
		},
		{
			name:      "Get error",
			existsErr: fmt.Errorf("key already exists"),
			getErr:    fmt.Errorf("get error"),
			wantErr:   true,
		},
		{
			name:      "Expected value mismatch error",
			existsErr: fmt.Errorf("key already exists"),
			getErr:    nil,
			getValue:  "not expected",
			args: args{
				key:           key,
				expectedValue: "expected",
			},
			wantErr: true,
		},
		{
			name:      "Expected value",
			existsErr: fmt.Errorf("key already exists"),
			getErr:    nil,
			getValue:  "expected",
			args: args{
				key:           key,
				expectedValue: "expected",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			s.Correlator = correlator

			ExistsError = tt.existsErr
			GetError = tt.getErr
			GetValue = tt.getValue
			if err := s.ValidateTextValue(
				context.Background(), tt.args.key, tt.args.expectedValue); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateTextValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateHashValue(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	testProps := make(map[string]interface{})
	testProps[key] = value

	hGetAllMap := make(map[string]string)
	type args struct {
		key   string
		props map[string]interface{}
	}
	tests := []struct {
		name         string
		existsErr    error
		hGetAllErr   error
		hGetMapKey   string
		hGetMapValue string
		args         args
		wantErr      bool
	}{
		{
			name:      "Key does not exists",
			existsErr: redis.Nil,
			wantErr:   true,
		},
		{
			name:       "HGetAll error",
			existsErr:  fmt.Errorf("key exists"),
			hGetAllErr: fmt.Errorf("HGetAll error"),
			wantErr:    true,
		},
		{
			name:         "Key not found error",
			existsErr:    fmt.Errorf("key exists"),
			hGetAllErr:   nil,
			hGetMapKey:   "wrongKey",
			hGetMapValue: value,
			args: args{
				key:   key,
				props: testProps,
			},
			wantErr: true,
		},
		{
			name:         "Mismatch value",
			existsErr:    fmt.Errorf("key exists"),
			hGetAllErr:   nil,
			hGetMapKey:   key,
			hGetMapValue: "wrongValue",
			args: args{
				key:   key,
				props: testProps,
			},
			wantErr: true,
		},
		{
			name:         "Match value",
			existsErr:    fmt.Errorf("key exists"),
			hGetAllErr:   nil,
			hGetMapKey:   key,
			hGetMapValue: value,
			args: args{
				key:   key,
				props: testProps,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			s.Correlator = correlator
			ExistsError = tt.existsErr
			HGetAllError = tt.hGetAllErr
			hGetAllMap[tt.hGetMapKey] = tt.hGetMapValue
			HGetAllMap = hGetAllMap
			if err := s.ValidateHashValue(
				context.Background(), tt.args.key, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateHashValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJSONValue(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	testProps := make(map[string]interface{})
	testProps[key] = value

	type args struct {
		key   string
		props map[string]interface{}
	}
	tests := []struct {
		name      string
		existsErr error
		getErr    error
		getValue  string
		args      args
		wantErr   bool
	}{
		{
			name:      "Key does not exists",
			existsErr: redis.Nil,
			wantErr:   true,
		},
		{
			name:      "HGetAll error",
			existsErr: fmt.Errorf("key exists"),
			getErr:    fmt.Errorf("Get error"),
			wantErr:   true,
		},
		{
			name:      "Mismatch value",
			existsErr: fmt.Errorf("key exists"),
			getErr:    nil,
			getValue:  "wrongValue",
			args: args{
				key:   key,
				props: testProps,
			},
			wantErr: true,
		},
		{
			name:      "Match value",
			existsErr: fmt.Errorf("key exists"),
			getErr:    nil,
			getValue:  `{"key":"value"}`,
			args: args{
				key:   key,
				props: testProps,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			s.Correlator = correlator
			ExistsError = tt.existsErr
			GetError = tt.getErr
			GetValue = tt.getValue
			if err := s.ValidateJSONValue(
				context.Background(), tt.args.key, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.ValidateJSONValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscribeTopic(t *testing.T) {
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	tests := []struct {
		name       string
		receiveErr error
		topic      string
		wantErr    bool
	}{
		{
			name:       "Receive error",
			receiveErr: fmt.Errorf("receive error"),
			wantErr:    true,
		},
		{
			name:       "No error",
			receiveErr: nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			s.Correlator = correlator
			PubSubReceiveError = tt.receiveErr
			if err := s.SubscribeTopic(context.Background(), tt.topic); (err != nil) != tt.wantErr {
				t.Errorf("Session.SubscribeTopic() error = %v, wantErr %v", err, tt.wantErr)
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
		name       string
		publishErr error
		args       args
		wantErr    bool
	}{
		{
			name:       "Publish error",
			publishErr: fmt.Errorf("publish error"),
			wantErr:    true,
		},
		{
			name:       "Publish without error",
			publishErr: nil,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			PublishError = tt.publishErr
			if err := s.PublishTextMessage(
				context.Background(), tt.args.topic, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Session.PublishTextMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPublishJSONMessage(t *testing.T) {
	rightProps := make(map[string]interface{})
	rightProps[key] = value
	os.MkdirAll(logsPath, os.ModePerm)
	defer os.RemoveAll(logsPath)
	type args struct {
		topic string
		props map[string]interface{}
	}
	tests := []struct {
		name       string
		publishErr error
		args       args
		wantErr    bool
	}{
		{
			name:       "Set and publish without error",
			publishErr: nil,
			args: args{
				topic: "topic",
				props: rightProps,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.RedisClientService = ClientServiceFuncMock{}
			PublishError = tt.publishErr
			if err := s.PublishJSONMessage(
				context.Background(), tt.args.topic, tt.args.props); (err != nil) != tt.wantErr {
				t.Errorf("Session.PublishJSONMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
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
			name: "Received message",
			args: args{
				timeout:     time.Duration(5000),
				expectedMsg: value,
			},
			wantErr: false,
		},
		{
			name: "Received message error",
			args: args{
				timeout:     time.Duration(5000),
				expectedMsg: "wrong message",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			s.Messages = append(s.Messages, value)
			if err := s.WaitForTextMessage(
				context.Background(), tt.args.timeout, tt.args.expectedMsg); (err != nil) != tt.wantErr {
				t.Errorf("Session.WaitForTextMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWaitForJSONMessageWithProperties(t *testing.T) {
	expectedProps := make(map[string]interface{})

	type args struct {
		timeout time.Duration
		key     string
		value   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Match msg",
			args: args{
				key:     key,
				value:   value,
				timeout: time.Duration(5000),
			},
			wantErr: false,
		},
		{
			name: "Invalid msg",
			args: args{
				key:     key,
				value:   "wrong value",
				timeout: time.Duration(5000),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{}
			expectedProps[tt.args.key] = tt.args.value
			s.Messages = append(s.Messages, `{"key":"value"}`)
			if err := s.WaitForJSONMessageWithProperties(
				context.Background(), tt.args.timeout, expectedProps); (err != nil) != tt.wantErr {
				t.Errorf("Session.WaitForJSONMessageWithProperties() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

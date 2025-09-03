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
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	PingError          error
	DoError            error
	SetError           error
	HSetError          error
	DelError           error
	PExpireError       error
	ExistsError        error
	ExistsValue        int64
	GetError           error
	GetValue           string
	HGetAllError       error
	HGetAllMap         map[string]string
	SubscribePubSub    *redis.PubSub
	PubSubReceiveError error
	PublishError       error
)

type ClientServiceFuncMock struct{}

func (c ClientServiceFuncMock) Ping(ctx context.Context, client *redis.Client) error {
	return PingError
}

func (c ClientServiceFuncMock) Do(
	ctx context.Context,
	client *redis.Client, op string,
	db int,
) error {
	return DoError
}

func (c ClientServiceFuncMock) Set(
	ctx context.Context,
	client *redis.Client,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	return SetError
}

func (c ClientServiceFuncMock) HSet(
	ctx context.Context,
	client *redis.Client,
	key string,
	values ...interface{},
) error {
	return HSetError
}

func (c ClientServiceFuncMock) Del(
	ctx context.Context,
	client *redis.Client,
	key string,
) error {
	return DelError
}

func (c ClientServiceFuncMock) PExpire(
	ctx context.Context,
	client *redis.Client,
	key string,
	expiration time.Duration,
) error {
	return PExpireError
}

func (c ClientServiceFuncMock) Exists(
	ctx context.Context,
	client *redis.Client,
	keys ...string,
) (int64, error) {
	return ExistsValue, ExistsError
}

func (c ClientServiceFuncMock) Get(
	ctx context.Context,
	client *redis.Client,
	key string,
) (string, error) {
	return GetValue, GetError
}
func (c ClientServiceFuncMock) HGetAll(
	ctx context.Context,
	client *redis.Client,
	key string,
) (map[string]string, error) {
	return HGetAllMap, HGetAllError
}
func (c ClientServiceFuncMock) Subscribe(
	ctx context.Context,
	client *redis.Client,
	channels ...string,
) *redis.PubSub {
	return SubscribePubSub
}

func (c ClientServiceFuncMock) PubSubReceive(
	ctx context.Context,
	pubSub *redis.PubSub,
) (interface{}, error) {
	return nil, PubSubReceiveError
}

func (c ClientServiceFuncMock) PubSubChannel(
	pubSub *redis.PubSub,
) <-chan *redis.Message {
	return nil
}
func (c ClientServiceFuncMock) Publish(ctx context.Context,
	client *redis.Client,
	channel string,
	message interface{},
) error {
	return PublishError
}
func (c ClientServiceFuncMock) PubSubClose(pubSub *redis.PubSub) error {
	return nil
}

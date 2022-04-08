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

	"github.com/go-redis/redis/v8"
)

type ClientFunctions interface {
	Ping(ctx context.Context, client *redis.Client) error
	Do(ctx context.Context, client *redis.Client, op string, db int) error
	Set(
		ctx context.Context,
		client *redis.Client,
		key string,
		value interface{},
		expiration time.Duration,
	) error
	HSet(
		ctx context.Context,
		client *redis.Client,
		key string,
		values ...interface{},
	) error
	PExpire(
		ctx context.Context,
		client *redis.Client,
		key string,
		expiration time.Duration,
	) error
	Exists(
		ctx context.Context,
		client *redis.Client,
		keys ...string,
	) (int64, error)
	Get(
		ctx context.Context,
		client *redis.Client,
		key string,
	) (string, error)
	HGetAll(
		ctx context.Context,
		client *redis.Client,
		key string,
	) (map[string]string, error)
	Subscribe(
		ctx context.Context,
		client *redis.Client,
		channels ...string,
	) *redis.PubSub
	PubSubReceive(
		ctx context.Context,
		pubSub *redis.PubSub,
	) (interface{}, error)
	PubSubChannel(
		pubSub *redis.PubSub,
	) <-chan *redis.Message
	Publish(ctx context.Context,
		client *redis.Client,
		channel string,
		message interface{},
	) error
}
type ClientService struct{}

func NewRedisClientService() *ClientService {
	return &ClientService{}
}

func (c ClientService) Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}

func (c ClientService) Do(ctx context.Context, client *redis.Client, op string, db int) error {
	return client.Do(ctx, op, db).Err()
}

func (c ClientService) Set(
	ctx context.Context,
	client *redis.Client,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	return client.Set(ctx, key, value, expiration).Err()
}

func (c ClientService) HSet(
	ctx context.Context,
	client *redis.Client,
	key string,
	values ...interface{},
) error {
	return client.HSet(ctx, key, values...).Err()
}

func (c ClientService) PExpire(
	ctx context.Context,
	client *redis.Client,
	key string,
	expiration time.Duration,
) error {
	return client.PExpire(ctx, key, expiration).Err()
}

func (c ClientService) Exists(ctx context.Context,
	client *redis.Client,
	keys ...string,
) (int64, error) {
	return client.Exists(ctx, keys...).Result()
}

func (c ClientService) Get(
	ctx context.Context,
	client *redis.Client,
	key string,
) (string, error) {
	return client.Get(ctx, key).Result()
}

func (c ClientService) HGetAll(
	ctx context.Context,
	client *redis.Client,
	key string,
) (map[string]string, error) {
	return client.HGetAll(ctx, key).Result()
}

func (c ClientService) Subscribe(
	ctx context.Context,
	client *redis.Client,
	channels ...string,
) *redis.PubSub {
	return client.Subscribe(ctx, channels...)
}

func (c ClientService) PubSubReceive(
	ctx context.Context,
	pubSub *redis.PubSub,
) (interface{}, error) {
	return pubSub.Receive(ctx)
}
func (c ClientService) PubSubChannel(
	pubSub *redis.PubSub,
) <-chan *redis.Message {
	return pubSub.Channel()
}
func (c ClientService) Publish(ctx context.Context,
	client *redis.Client,
	channel string,
	message interface{},
) error {
	return client.Publish(ctx, channel, message).Err()
}

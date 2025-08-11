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
	"github.com/cucumber/godog"
	"github.com/redis/go-redis/v9"
)

// Steps to initialize common steps.
type Steps struct {
}

// InitializeSteps initializes all the steps.
func (cs Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the HTTP session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)
	// Initialize the steps
	scenCtx.Step(`^the redis endpoint$`, func(t *godog.Table) error {
		var options redis.Options
		if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &options); err != nil {
			return fmt.Errorf("failed configuring redis endpoint: %w", err)
		}
		return session.ConfigureClient(ctx, &options)
	})
	scenCtx.Step(`^I select the redis database "([^"]+)"$`, func(id string) error {
		dbID, err := golium.ValueAsInt(ctx, id)
		if err != nil {
			return err
		}
		return session.SelectDatabase(ctx, dbID)
	})
	scenCtx.Step(`^the redis TTL of "(\d+)" millis`, func(ttl int) {
		session.ConfigureTTL(ctx, ttl)
	})
	scenCtx.Step(`^I set the redis key "([^"]*)" with the text`, func(key string, value *godog.DocString) error {
		return session.SetTextValue(ctx, golium.ValueAsString(ctx, key), golium.ValueAsString(ctx, value.Content))
	})
	scenCtx.Step(`^I set the redis key "([^"]*)" with hash properties`, func(key string, t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the hashed value in redis: %w", err)
		}
		return session.SetHashValue(ctx, golium.ValueAsString(ctx, key), props)
	})
	scenCtx.Step(`^I set the redis key "([^"]*)" with the JSON properties`, func(key string, t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the JSON value in redis: %w", err)
		}
		return session.SetJSONValue(ctx, golium.ValueAsString(ctx, key), props)
	})
	scenCtx.Step(`^I delete the redis key "([^"]*)"`, func(key string) error {
		return session.DeleteKeyValue(ctx, golium.ValueAsString(ctx, key))
	})
	scenCtx.Step(`^the redis key "([^"]*)" must have the text`, func(key string, value *godog.DocString) error {
		return session.ValidateTextValue(ctx, golium.ValueAsString(ctx, key), golium.ValueAsString(ctx, value.Content))
	})
	scenCtx.Step(`^the redis key "([^"]*)" must have hash properties`, func(key string, t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the expected hashed value in redis: %w", err)
		}
		return session.ValidateHashValue(ctx, golium.ValueAsString(ctx, key), props)
	})
	scenCtx.Step(`^the redis key "([^"]*)" must have the JSON properties`, func(key string, t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the expected JSON value in redis: %w", err)
		}
		return session.ValidateJSONValue(ctx, golium.ValueAsString(ctx, key), props)
	})
	scenCtx.Step(`^the redis key "([^"]*)" must not exist`, func(key string) error {
		return session.ValidateNonExistantKey(ctx, golium.ValueAsString(ctx, key))
	})
	scenCtx.Step(`^I subscribe to the redis topic "([^"]*)"$`, func(topic string) error {
		return session.SubscribeTopic(ctx, golium.ValueAsString(ctx, topic))
	})
	scenCtx.Step(`^I unsubscribe from the redis topic "([^"]*)"$`, func(topic string) error {
		return session.UnsubscribeTopic(ctx, golium.ValueAsString(ctx, topic))
	})
	scenCtx.Step(`^I publish a message to the redis topic "([^"]*)" with the text$`, func(topic string, message *godog.DocString) error {
		return session.PublishTextMessage(ctx, golium.ValueAsString(ctx, topic), golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I publish a message to the redis topic "([^"]*)" with the JSON properties$`, func(topic string, t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the request body: %w", err)
		}
		return session.PublishJSONMessage(ctx, golium.ValueAsString(ctx, topic), props)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a redis message with the text$`, func(timeout int, message *godog.DocString) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForTextMessage(ctx, timeoutDuration, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a redis message with the JSON properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the redis message: %w", err)
		}
		return session.WaitForJSONMessageWithProperties(ctx, timeoutDuration, props)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? without a redis message with the JSON properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the redis message: %w", err)
		}
		if err := session.WaitForJSONMessageWithProperties(ctx, timeoutDuration, props); err == nil {
			return fmt.Errorf("received a message with JSON properties '%+v'", props)
		}
		return nil
	})
	return ctx
}

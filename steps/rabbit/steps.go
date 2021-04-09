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
	"github.com/cucumber/godog"
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
	scenCtx.Step(`^the rabbit endpoint "([^"]*)"$`, func(uri string) error {
		return session.ConfigureConnection(ctx, uri)
	})

	scenCtx.Step(`^I subscribe to the rabbit topic "([^"]*)"$`, func(topic string) error {
		return session.SubscribeTopic(ctx, golium.ValueAsString(ctx, topic))
	})
	scenCtx.Step(`^I unsubscribe from the rabbit topic "([^"]*)"$`, func(topic string) error {
		return session.UnsubscribeTopic(ctx, golium.ValueAsString(ctx, topic))
	})
	scenCtx.Step(`^I publish a message to the rabbit topic "([^"]*)" with the text$`, func(topic string, message *godog.DocString) error {
		return session.PublishTextMessage(ctx, golium.ValueAsString(ctx, topic), golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I publish a message to the rabbit topic "([^"]*)" with the JSON properties$`, func(topic string, t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the request body: %w", err)
		}
		return session.PublishJSONMessage(ctx, golium.ValueAsString(ctx, topic), props)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a rabbit message with the text$`, func(timeout int, message *godog.DocString) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForTextMessage(ctx, timeoutDuration, message.Content)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a rabbit message with the JSON properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the rabbit message: %w", err)
		}
		return session.WaitForJSONMessageWithProperties(ctx, timeoutDuration, props)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? without a rabbit message with the JSON properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the rabbit message: %w", err)
		}
		if err := session.WaitForJSONMessageWithProperties(ctx, timeoutDuration, props); err == nil {
			return fmt.Errorf("received a message with JSON properties '%+v'", props)
		}
		return nil
	})
	return ctx
}

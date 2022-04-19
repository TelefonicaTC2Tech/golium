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
	"time"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
)

// Steps to initialize common steps.
type Steps struct {
}

// InitializeSteps initializes all the steps.
func (cs Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the rabbit session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)
	// Initialize the steps
	scenCtx.Step(`^the rabbit endpoint "([^"]*)"$`, func(uri string) error {
		return session.ConfigureConnection(ctx, golium.ValueAsString(ctx, uri))
	})
	scenCtx.Step(`^I subscribe to the rabbit topic "([^"]*)"$`, func(topic string) error {
		return session.SubscribeTopic(ctx, golium.ValueAsString(ctx, topic))
	})
	scenCtx.Step(`^I set rabbit headers$`, func(t *godog.Table) error {
		return session.ConfigureHeaders(ctx, t)
	})
	scenCtx.Step(`^I set standard rabbit properties$`, func(t *godog.Table) error {
		return session.ConfigureStandardProperties(ctx, t)
	})
	scenCtx.Step(`^I publish a message to the rabbit topic "([^"]*)" with the text$`, func(topic string, message *godog.DocString) error {
		return session.PublishTextMessage(ctx, golium.ValueAsString(ctx, topic), golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I publish a message to the rabbit topic "([^"]*)" with the JSON properties$`, func(topic string, t *godog.Table) error {
		return session.PublishJSONMessage(ctx, golium.ValueAsString(ctx, topic), t)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a rabbit message with the text$`, func(timeout int, message *godog.DocString) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForTextMessage(ctx, timeoutDuration, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a rabbit message with the JSON properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForJSONMessageWithProperties(ctx, timeoutDuration, t, false)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? without a rabbit message with the JSON properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForJSONMessageWithProperties(ctx, timeoutDuration, t, true)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a rabbit message with the standard properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForMessagesWithStandardProperties(ctx, timeoutDuration, 1, t, false)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for "(\d+)" rabbit messages with the standard properties$`, func(timeout int, count int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForMessagesWithStandardProperties(ctx, timeoutDuration, count, t, false)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? without a rabbit message with the standard properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForMessagesWithStandardProperties(ctx, timeoutDuration, 1, t, true)
	})
	scenCtx.Step(`^the rabbit message has the rabbit headers$`, func(t *godog.Table) error {
		return session.ValidateMessageHeaders(ctx, t)
	})
	scenCtx.Step(`^the rabbit message has the standard rabbit properties$`, func(t *godog.Table) error {
		return session.ValidateMessageStandardProperties(ctx, t, false)
	})
	scenCtx.Step(`^the rabbit message body has the text$`, func(m *godog.DocString) error {
		message := golium.ValueAsString(ctx, m.Content)
		return session.ValidateMessageTextBody(ctx, message)
	})
	scenCtx.Step(`^the rabbit message body has the JSON properties$`, func(t *godog.Table) error {
		return session.ValidateMessageJSONBody(ctx, t, -1)
	})
	scenCtx.Step(`^the body of the rabbit message in position "(\d+)" has the JSON properties$`, func(pos int, t *godog.Table) error {
		return session.ValidateMessageJSONBody(ctx, t, pos)
	})
	scenCtx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, session.Unsubscribe(ctx)
	})
	return ctx
}

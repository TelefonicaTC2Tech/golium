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

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
	"github.com/streadway/amqp"
)

const (
	convertTableToMapMessage       = "failed processing table to a map for the rabbit message: "
	convertTableToStructProperties = "failed processing table to a map for the standard rabbit properties: "
	convertTableToMapBody          = "failed processing table to a map for the request body: "
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
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf(convertTableToMapBody+"%w", err)
		}
		return session.ConfigureHeaders(ctx, props)
	})
	scenCtx.Step(`^I set standard rabbit properties$`, func(t *godog.Table) error {
		var props amqp.Publishing
		if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &props); err != nil {
			return fmt.Errorf("failed configuring rabbit endpoint: %w", err)
		}
		session.ConfigureStandardProperties(ctx, props)
		return nil
	})

	publishMessageSteps(ctx, session, scenCtx)

	waitRabbitMessageSteps(ctx, session, scenCtx)

	rabbitMessageHasSteps(ctx, session, scenCtx)

	scenCtx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, session.Unsubscribe(ctx)
	})
	return ctx
}
func publishMessageSteps(ctx context.Context, session *Session, scenCtx *godog.ScenarioContext) {
	scenCtx.Step(`^I publish a message to the rabbit topic "([^"]*)" with the text$`, func(topic string, message *godog.DocString) error {
		return session.PublishTextMessage(ctx, golium.ValueAsString(ctx, topic), golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I publish a message to the rabbit topic "([^"]*)" with the JSON properties$`, func(topic string, t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf(convertTableToMapBody+"%w", err)
		}
		return session.PublishJSONMessage(ctx, golium.ValueAsString(ctx, topic), props)
	})
}

func waitRabbitMessageSteps(ctx context.Context, session *Session, scenCtx *godog.ScenarioContext) {
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a rabbit message with the text$`, func(timeout int, message *godog.DocString) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		return session.WaitForTextMessage(ctx, timeoutDuration, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a rabbit message with the JSON properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf(convertTableToMapMessage+"%w", err)
		}
		return session.WaitForJSONMessageWithProperties(ctx, timeoutDuration, props)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? without a rabbit message with the JSON properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf(convertTableToMapMessage+"%w", err)
		}
		if err := session.WaitForJSONMessageWithProperties(ctx, timeoutDuration, props); err == nil {
			return fmt.Errorf("received a message with JSON properties '%+v'", props)
		}
		return nil
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for a rabbit message with the standard properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		var props amqp.Delivery
		if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &props); err != nil {
			return fmt.Errorf(convertTableToStructProperties+"%w", err)
		}
		return session.WaitForMessagesWithStandardProperties(ctx, timeoutDuration, 1, props)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? for "(\d+)" rabbit messages with the standard properties$`, func(timeout int, count int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		var props amqp.Delivery
		if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &props); err != nil {
			return fmt.Errorf(convertTableToStructProperties+"%w", err)
		}
		return session.WaitForMessagesWithStandardProperties(ctx, timeoutDuration, count, props)
	})
	scenCtx.Step(`^I wait up to "(\d+)" seconds? without a rabbit message with the standard properties$`, func(timeout int, t *godog.Table) error {
		timeoutDuration := time.Duration(timeout) * time.Second
		var props amqp.Delivery
		if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &props); err != nil {
			return fmt.Errorf(convertTableToStructProperties+"%w", err)
		}
		if err := session.WaitForMessagesWithStandardProperties(ctx, timeoutDuration, 1, props); err == nil {
			return fmt.Errorf("received a message with standard rabbit properties '%+v'", props)
		}
		return nil
	})
}
func rabbitMessageHasSteps(ctx context.Context, session *Session, scenCtx *godog.ScenarioContext) {
	scenCtx.Step(`^the rabbit message has the standard rabbit properties$`, func(t *godog.Table) error {
		var props amqp.Delivery
		if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &props); err != nil {
			return fmt.Errorf("failed configuring rabbit endpoint: %w", err)
		}
		return session.ValidateMessageStandardProperties(ctx, props)
	})
	scenCtx.Step(`^the rabbit message body has the text$`, func(m *godog.DocString) error {
		message := golium.ValueAsString(ctx, m.Content)
		return session.ValidateMessageTextBody(ctx, message)
	})
	scenCtx.Step(`^the rabbit message body has the JSON properties$`, func(t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf(convertTableToMapMessage+"%w", err)
		}
		return session.ValidateMessageJSONBody(ctx, props, -1)
	})
	scenCtx.Step(`^the rabbit message has the rabbit headers$`, func(t *godog.Table) error {
		headers, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf(convertTableToMapMessage+"%w", err)
		}
		return session.ValidateMessageHeaders(ctx, headers)
	})

	scenCtx.Step(`^the body of the rabbit message in position "(\d+)" has the JSON properties$`, func(pos int, t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf(convertTableToMapMessage+"%w", err)
		}
		return session.ValidateMessageJSONBody(ctx, props, pos)
	})
}

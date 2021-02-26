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

package common

import (
	"context"
	"fmt"
	"time"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
)

// Steps to initialize common steps.
type Steps struct {
}

// InitializeSteps initializes all the steps.
func (cs Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	scenCtx.Step(`^I store "([^"]*)" in context "([^"]*)"$`, func(value, name string) error {
		return StoreValueInContext(ctx, name, value)
	})
	scenCtx.Step(`^I generate a UUID and store it in context "([^"]*)"$`, func(name string) error {
		return GenerateUUIDInContext(ctx, golium.ValueAsString(ctx, name))
	})
	scenCtx.Step(`^I wait for "(\d+)" seconds$`, func(d int) error {
		time.Sleep(time.Duration(d) * time.Second)
		return nil
	})
	scenCtx.Step(`^I wait for "(\d+)" millis$`, func(d int) error {
		time.Sleep(time.Duration(d) * time.Millisecond)
		return nil
	})
	return ctx
}

// StoreValueInContext stores a value in golium.Context using the key name.
func StoreValueInContext(ctx context.Context, name, value string) error {
	golium.GetContext(ctx).Put(name, value)
	return nil
}

// GenerateUUIDInContext generates a UUID and stores it in golium.Context using the key name.
func GenerateUUIDInContext(ctx context.Context, name string) error {
	guid, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("Error generating UUID. %s", err)
	}
	return StoreValueInContext(ctx, name, guid.String())
}

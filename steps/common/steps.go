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
	"net/url"
	"os/exec"
	"strings"
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
		return StoreValueInContext(ctx, golium.ValueAsString(ctx, name), golium.ValueAsString(ctx, value))
	})
	scenCtx.Step(`^I generate a UUID and store it in context "([^"]*)"$`, func(name string) error {
		return GenerateUUIDInContext(ctx, golium.ValueAsString(ctx, name))
	})
	scenCtx.Step(`^I wait for "([^"]*)" seconds$`, func(delay string) error {
		d, err := golium.ValueAsInt(ctx, delay)
		if err != nil {
			return fmt.Errorf("invalid delay '%s': %w", delay, err)
		}
		time.Sleep(time.Duration(d) * time.Second)
		return nil
	})
	scenCtx.Step(`^I wait for "([^"]*)" millis$`, func(delay string) error {
		d, err := golium.ValueAsInt(ctx, delay)
		if err != nil {
			return fmt.Errorf("invalid delay '%s': %w", delay, err)
		}
		time.Sleep(time.Duration(d) * time.Millisecond)
		return nil
	})
	scenCtx.Step(`^I parse the URL "([^"]*)" in context "([^"]*)"$`, func(uri, ctxtPrefix string) error {
		return ParseURL(ctx, golium.ValueAsString(ctx, uri), golium.ValueAsString(ctx, ctxtPrefix))
	})
	scenCtx.Step(`^the value "([^"]*)" must be equal to "([^"]*)"$`, func(value, expectedValue string) error {
		v := golium.Value(ctx, value)
		e := golium.Value(ctx, expectedValue)
		if v == e {
			return nil
		}
		return fmt.Errorf("mismatch of values: expected '%s', actual '%s'", e, v)
	})
	scenCtx.Step(`^I store my local ip in context "([^"]*)"$`, func(key string) error {
		cmd := exec.Command("hostname", "-i")
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error executing `%s` command %v", cmd, string(stdoutStderr))
		}
		ip := strings.Trim(string(stdoutStderr), " \r\n")
		golium.GetContext(ctx).Put(golium.ValueAsString(ctx, key), ip)
		return nil
	})
	scenCtx.Step(`^I store domain "([^"]*)" ip in context "([^"]*)"$`, func(domainParam, key string) error {
		if domainParam == "" {
			return nil
		}
		domain := golium.ValueAsString(ctx, domainParam)
		command := fmt.Sprintf("ping -c 1 %s | head -1 | grep -oe '[0-9]*\\.[0-9]*\\.[0-9]*\\.[0-9]*'", domain)
		cmd := exec.Command("/bin/sh", "-c", command)
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error executing `%s` command %v", cmd, string(stdoutStderr))
		}
		ip := strings.Trim(string(stdoutStderr), " \r\n")
		golium.GetContext(ctx).Put(golium.ValueAsString(ctx, key), ip)
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
	guid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed generating UUID: %w", err)
	}
	return StoreValueInContext(ctx, name, guid.String())
}

// ParseURL parses a URL and stores its values in the context.
// It will store a context value per element parsed from the URL under the context prefix.
func ParseURL(ctx context.Context, uri, ctxtPrefix string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("failed parsing URL '%s': %w", uri, err)
	}
	if err := StoreValueInContext(ctx, fmt.Sprintf("%s.scheme", ctxtPrefix), u.Scheme); err != nil {
		return fmt.Errorf("failed storing scheme of URL '%s': %w", uri, err)
	}
	if err := StoreValueInContext(ctx, fmt.Sprintf("%s.host", ctxtPrefix), u.Host); err != nil {
		return fmt.Errorf("failed storing host of URL '%s': %w", uri, err)
	}
	if err := StoreValueInContext(ctx, fmt.Sprintf("%s.hostname", ctxtPrefix), u.Hostname()); err != nil {
		return fmt.Errorf("failed storing host of URL '%s': %w", uri, err)
	}
	if err := StoreValueInContext(ctx, fmt.Sprintf("%s.path", ctxtPrefix), u.Path); err != nil {
		return fmt.Errorf("failed storing path of URL '%s': %w", uri, err)
	}
	if err := StoreValueInContext(ctx, fmt.Sprintf("%s.rawquery", ctxtPrefix), u.RawQuery); err != nil {
		return fmt.Errorf("failed storing raw query of URL '%s': %w", uri, err)
	}
	return ParseQuery(ctx, u.RawQuery, fmt.Sprintf("%s.query", ctxtPrefix))
}

// ParseQuery parses a query string and stores its values in the context.
// Note that it does not support multiple query params with the same key. It will store only one.
func ParseQuery(ctx context.Context, query, ctxtPrefix string) error {
	v, err := url.ParseQuery(query)
	if err != nil {
		return fmt.Errorf("failed parsing query '%s': %w", query, err)
	}
	for key := range v {
		if err := StoreValueInContext(ctx, fmt.Sprintf("%s.%s", ctxtPrefix, key), v.Get(key)); err != nil {
			return fmt.Errorf("failed storing query param '%s': %w", key, err)
		}
	}
	return nil
}

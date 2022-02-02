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

package http

import (
	"context"
	"fmt"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
)

// Steps type is responsible to initialize the HTTP client steps in godog framework.
type Steps struct {
}

// InitializeSteps adds client HTTP steps to the scenario context.
// It implements StepsInitializer interface.
// It returns a new context (context is immutable) with the HTTP Context.
func (s Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the HTTP session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)
	// Initialize the steps
	scenCtx.Step(`^the HTTP endpoint "([^"]*)"$`, func(endpoint string) error {
		return session.ConfigureEndpoint(ctx, golium.ValueAsString(ctx, endpoint))
	})
	scenCtx.Step(`^an HTTP timeout of "([^"]*)" milliseconds$`, func(timeout string) error {
		to, err := golium.ValueAsInt(ctx, timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout '%s': %w", timeout, err)
		}
		return session.SetHTTPResponseTimeout(ctx, to)
	})
	scenCtx.Step(`^the HTTP path "([^"]*)"$`, func(path string) error {
		return session.ConfigurePath(ctx, golium.ValueAsString(ctx, path))
	})
	scenCtx.Step(`^the HTTP query parameters$`, func(t *godog.Table) error {
		params, err := golium.ConvertTableToMultiMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing query parameters from table: %w", err)
		}
		return session.ConfigureQueryParams(ctx, params)
	})
	scenCtx.Step(`^the HTTP request headers$`, func(t *godog.Table) error {
		headers, err := golium.ConvertTableToMultiMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing HTTP headers from table: %w", err)
		}
		return session.ConfigureHeaders(ctx, headers)
	})
	scenCtx.Step(`^the HTTP request with username "([^"]*)" and password "([^"]*)"$`, func(username, password string) error {
		return session.ConfigureCredentials(ctx, golium.ValueAsString(ctx, username), golium.ValueAsString(ctx, password))
	})
	scenCtx.Step(`^the JSON properties in the HTTP request body$`, func(t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the request body: %w", err)
		}
		return session.ConfigureRequestBodyJSONProperties(ctx, props)
	})
	scenCtx.Step(`^the HTTP request body with the JSON$`, func(message *godog.DocString) error {
		return session.ConfigureRequestBodyJSONText(ctx, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^the HTTP request body with the JSON "([^"]*)" from "([^"]*)" file$`, func(code, file string) error {
		return session.ConfigureRequestBodyJSONFile(ctx, code, file)
	})
	scenCtx.Step(`^the HTTP request body with the JSON "([^"]*)" from "([^"]*)" file without$`, func(code, file string, t *godog.Table) error {
		params, err := golium.ConvertTableColumnToArray(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the request body: %w", err)
		}
		return session.ConfigureRequestBodyJSONFileWithout(ctx, code, file, params)
	})
	scenCtx.Step(`^the HTTP request body with the URL encoded properties$`, func(t *godog.Table) error {
		props, err := golium.ConvertTableToMultiMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the request body: %w", err)
		}
		return session.ConfigureRequestBodyURLEncodedProperties(ctx, props)
	})
	scenCtx.Step(`^the HTTP client does not follow any redirection$`, func() error {
		return session.ConfigureNoRedirection(ctx)
	})
	scenCtx.Step(`^I send a HTTP "([^"]*)" request$`, func(method string) error {
		return session.SendHTTPRequest(ctx, golium.ValueAsString(ctx, method))
	})
	scenCtx.Step(`^the HTTP response timed out$`, func() error {
		return session.ValidateResponseTimedout(ctx)
	})
	scenCtx.Step(`^the HTTP status code must be "(\d+)"$`, func(code int) error {
		return session.ValidateStatusCode(ctx, code)
	})
	scenCtx.Step(`^the HTTP response must contain the headers$`, func(t *godog.Table) error {
		headers, err := golium.ConvertTableToMultiMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing HTTP headers from table: %w", err)
		}
		return session.ValidateResponseHeaders(ctx, headers)
	})
	scenCtx.Step(`^the HTTP response must not contain the headers$`, func(t *godog.Table) error {
		headers, err := golium.ConvertTableColumnToArray(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing HTTP headers from table: %w", err)
		}
		return session.ValidateNotResponseHeaders(ctx, headers)
	})
	scenCtx.Step(`^the HTTP response body must comply with the JSON schema "([^"]*)"$`, func(schema string) error {
		return session.ValidateResponseBodyJSONSchema(ctx, golium.ValueAsString(ctx, schema))
	})
	scenCtx.Step(`^the HTTP response "([^"]*)" must match with the JSON "([^"]*)" from "([^"]*)" file$`, func(respDataLocation, code, file string) error {
		return session.ValidateResponseBodyJSONFile(ctx, code, file, respDataLocation)
	})
	scenCtx.Step(`^the HTTP response "([^"]*)" must match with the JSON "([^"]*)" from "([^"]*)" file without$`, func(respDataLocation, code, file string, t *godog.Table) error {
		params, err := golium.ConvertTableColumnToArray(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the request body: %w", err)
		}
		return session.ValidateResponseBodyJSONFileWithout(ctx, code, file, respDataLocation, params)
	})
	scenCtx.Step(`^the HTTP response body must have the JSON properties$`, func(t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing the table for validating the response body: %w", err)
		}
		return session.ValidateResponseBodyJSONProperties(ctx, props)
	})
	scenCtx.Step(`^the HTTP response body must be empty$`, func() error {
		return session.ValidateResponseBodyEmpty(ctx)
	})
	scenCtx.Step(`^the HTTP response body must be the text$`, func(message *godog.DocString) error {
		return session.ValidateResponseBodyText(ctx, golium.ValueAsString(ctx, message.Content))
	})
	scenCtx.Step(`^I store the element "([^"]*)" from the JSON HTTP response body in context "([^"]*)"$`, func(key string, ctxtKey string) error {
		return session.StoreResponseBodyJSONPropertyInContext(ctx, golium.ValueAsString(ctx, key), golium.ValueAsString(ctx, ctxtKey))
	})
	scenCtx.Step(`^I store the header "([^"]*)" from the HTTP response in context "([^"]*)"$`, func(key string, ctxtKey string) error {
		return session.StoreResponseHeaderInContext(ctx, golium.ValueAsString(ctx, key), golium.ValueAsString(ctx, ctxtKey))
	})
	return ctx
}

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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/tidwall/sjson"
	"github.com/xeipuuv/gojsonschema"
)

// ClientSteps type is responsible to initialize the HTTP client steps in godog framework.
type ClientSteps struct {
}

// InitializeSteps adds client HTTP steps to the scenario context.
// It implements StepsInitializer interface.
// It returns a new context (context is immutable) with the ClientContext.
func (cs ClientSteps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Add the ClientContext to context
	ctx = InitializeClientContext(ctx)
	// Initialize the steps
	scenCtx.Step(`^the HTTP endpoint "([^"]*)"$`, func(endpoint string) error {
		return ConfigureEndpointStep(ctx, golium.ValueAsString(ctx, endpoint))
	})
	scenCtx.Step(`^the HTTP path "([^"]*)"$`, func(path string) error {
		return ConfigurePathStep(ctx, golium.ValueAsString(ctx, path))
	})
	scenCtx.Step(`^the HTTP query parameters$`, func(params *godog.Table) error {
		return ConfigureQueryParamsStep(ctx, params)
	})
	scenCtx.Step(`^the HTTP request headers$`, func(headers *godog.Table) error {
		return ConfigureHeadersStep(ctx, headers)
	})
	scenCtx.Step(`^the JSON properties in the HTTP request body$`, func(properties *godog.Table) error {
		return ConfigureRequestBodyJSONPropertiesStep(ctx, properties)
	})
	scenCtx.Step(`^I send a HTTP "([^"]*)" request$`, func(method string) error {
		return SendHTTPRequestStep(ctx, golium.ValueAsString(ctx, method))
	})
	scenCtx.Step(`^the HTTP status code must be "(\d+)"$`, func(code int) error {
		return ValidateStatusCodeStep(ctx, code)
	})
	scenCtx.Step(`^the HTTP response body must comply with the JSON schema "([^"]*)"$`, func(schema string) error {
		return ValidateResponseBodyJSONSchemaStep(ctx, golium.ValueAsString(ctx, schema))
	})
	scenCtx.Step(`^the HTTP response body must have the JSON properties$`, func(properties *godog.Table) error {
		return ValidateResponseBodyJSONPropertiesStep(ctx, properties)
	})
	return ctx
}

// ConfigureEndpointStep configures the HTTP endpoint.
func ConfigureEndpointStep(ctx context.Context, endpoint string) error {
	clientCtx := GetClientContext(ctx)
	clientCtx.Request.Endpoint = endpoint
	return nil
}

// ConfigurePathStep configures the path of the HTTP endpoint.
// It configures a resource path in the application context.
// The API endpoint and the resource path are composed when invoking the HTTP server.
func ConfigurePathStep(ctx context.Context, path string) error {
	clientCtx := GetClientContext(ctx)
	clientCtx.Request.Path = path
	return nil
}

// ConfigureQueryParamsStep stores a table of query parameters in the application context.
func ConfigureQueryParamsStep(ctx context.Context, t *godog.Table) error {
	clientCtx := GetClientContext(ctx)
	params, err := golium.ConvertTableToMultiMap(t)
	if err != nil {
		return fmt.Errorf("Error processing query parameters from table. %s", err)
	}
	clientCtx.Request.QueryParams = params
	return nil
}

// ConfigureHeadersStep stores a table of HTTP headers in the application context.
func ConfigureHeadersStep(ctx context.Context, t *godog.Table) error {
	clientCtx := GetClientContext(ctx)
	headers, err := golium.ConvertTableToMultiMap(t)
	if err != nil {
		return fmt.Errorf("Error processing HTTP headers from table. %s", err)
	}
	clientCtx.Request.Headers = headers
	return nil
}

// ConfigureRequestBodyJSONPropertiesStep writes the body in the HTTP request as a JSON with properties.
// The content type is forced to application/json.
func ConfigureRequestBodyJSONPropertiesStep(ctx context.Context, t *godog.Table) error {
	clientCtx := GetClientContext(ctx)
	props, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return err
	}
	json := ""
	for key, value := range props {
		if json, err = sjson.Set(json, key, value); err != nil {
			return fmt.Errorf("Error setting property '%s' with value '%s' in the request body. %s", key, value, err)
		}
	}
	clientCtx.Request.RequestBody = []byte(json)
	return nil
}

// SendHTTPRequestStep sends a HTTP request using the configuration in the application context.
func SendHTTPRequestStep(ctx context.Context, method string) error {
	logger := GetLogger()
	clientCtx := GetClientContext(ctx)
	clientCtx.Request.Method = method
	corr := uuid.New().String()
	url, err := clientCtx.URL()
	if err != nil {
		return err
	}
	reqBodyReader := bytes.NewReader(clientCtx.Request.RequestBody)
	req, err := http.NewRequest(method, url.String(), reqBodyReader)
	if err != nil {
		return fmt.Errorf("Error creating the HTTP request with method: '%s' and url: '%s'. %s", method, url, err)
	}
	req.Header = clientCtx.Request.Headers
	logger.LogRequest(req, clientCtx.Request.RequestBody, corr)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error with the HTTP request. %s", err)
	}
	defer resp.Body.Close()
	// This is dangerous for big response bodies, but is read now to make sure that the body reader is closed.
	// TODO: limit the max size of the response body.
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading the response body. %s", err)
	}
	clientCtx.Response.Response = resp
	clientCtx.Response.ResponseBody = respBodyBytes
	logger.LogResponse(resp, respBodyBytes, corr)
	return nil
}

// ValidateStatusCodeStep validates the status code from the HTTP response.
func ValidateStatusCodeStep(ctx context.Context, expectedCode int) error {
	clientCtx := GetClientContext(ctx)
	if expectedCode != clientCtx.Response.Response.StatusCode {
		return fmt.Errorf("Status code mismatch. Expected: %d, actual: %d", expectedCode, clientCtx.Response.Response.StatusCode)
	}
	return nil
}

// ValidateResponseBodyJSONSchemaStep validates the response body against the JSON schema.
func ValidateResponseBodyJSONSchemaStep(ctx context.Context, schema string) error {
	clientCtx := GetClientContext(ctx)
	schemasDir := golium.GetConfig().Dir.Schemas
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s/%s.json", schemasDir, schema))
	documentLoader := gojsonschema.NewStringLoader(string(clientCtx.Response.ResponseBody))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("Error validating response body against schema: %s. %s", schema, err)
	}
	if result.Valid() {
		return nil
	}
	return fmt.Errorf("Invalid response body according to schema: %s. %+v", schema, result.Errors())
}

// ValidateResponseBodyJSONPropertiesStep validates a list of properties in the JSON body of the HTTP response.
func ValidateResponseBodyJSONPropertiesStep(ctx context.Context, t *godog.Table) error {
	clientCtx := GetClientContext(ctx)
	props, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return err
	}
	m := golium.NewMapFromJSONBytes(clientCtx.Response.ResponseBody)
	for key, expectedValue := range props {
		value := m.Get(key)
		if value != expectedValue {
			return fmt.Errorf("Mismatch of json property '%s'. Expected: '%s', actual: '%s'", key, expectedValue, value)
		}
	}
	return nil
}

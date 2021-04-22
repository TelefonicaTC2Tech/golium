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

package elasticsearch

import (
	"context"
	"fmt"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
	"github.com/elastic/go-elasticsearch/v8"
)

// Steps type is responsible to initialize the Elasticsearch client steps in godog framework.
type Steps struct {
}

// InitializeSteps adds client elasticsearch steps to the scenario context.
// It implements StepsInitializer interface.
// It returns a new context (context is immutable) with the elasticsearch Context.
func (s Steps) InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context {
	// Initialize the elasticsearch session in the context
	ctx = InitializeContext(ctx)
	session := GetSession(ctx)
	// Initialize the steps
	scenCtx.Step(`^the elasticsearch server$`, func(t *godog.Table) error {
		var options elasticsearch.Config
		if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &options); err != nil {
			return fmt.Errorf("failed configuring elasticsearch client: %w", err)
		}
		return session.ConfigureClient(ctx, options)
	})
	scenCtx.Step(`^I create the elasticsearch document with index "([^"]*)" and the JSON properties`, func(idx string, t *godog.Table) error {
		index := golium.ValueAsString(ctx, idx)
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the JSON value in elasticsearch: %w", err)
		}
		return session.NewDocument(ctx, index, props)
	})
	scenCtx.Step(`^I search in the elasticsearch index "([^"]*)" with the JSON body$`, func(idx string, b *godog.DocString) error {
		index := golium.ValueAsString(ctx, idx)
		body := golium.ValueAsString(ctx, b.Content)
		return session.SearchDocument(ctx, index, body)
	})
	scenCtx.Step(`^the search result must have the JSON properties`, func(t *godog.Table) error {
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the JSON value in elasticsearch: %w", err)
		}
		return session.ValidateDocumentJSONProperties(ctx, props)
	})
	return ctx
}

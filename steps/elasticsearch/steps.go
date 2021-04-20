package elasticsearch

import (
	"context"
	"fmt"

	"github.com/Telefonica/golium"
	"github.com/cucumber/godog"
	"github.com/elastic/go-elasticsearch/v8"
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
	scenCtx.Step(`^the elasticsearch server$`, func(t *godog.Table) error {
		var options elasticsearch.Config
		if err := golium.ConvertTableWithoutHeaderToStruct(ctx, t, &options); err != nil {
			return fmt.Errorf("failed configuring elasticsearch client: %w", err)
		}
		return session.ConfigureConnection(ctx, options)
	})
	scenCtx.Step(`^I create the elasticsearch document with index "([^"]*)" and the JSON properties`, func(idx string, t *godog.Table) error {
		index := golium.ValueAsString(ctx, idx)
		props, err := golium.ConvertTableToMap(ctx, t)
		if err != nil {
			return fmt.Errorf("failed processing table to a map for the JSON value in elasticsearch: %w", err)
		}
		return session.CreatesDocument(ctx, index, props)
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

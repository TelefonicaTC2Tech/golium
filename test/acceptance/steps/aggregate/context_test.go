package aggregate

import (
	"context"
	"fmt"
	"testing"

	"github.com/Telefonica/golium/test/acceptance/steps/shared"
	"github.com/stretchr/testify/assert"
)

func TestInitializeContext(t *testing.T) {

	var ctx = context.Background()
	var contextKeyValue = "aggregateSession"
	sharedContext := shared.InitializeContext(ctx)
	aggregateContext := InitializeContext(sharedContext)

	aggregatedContextGenerated := aggregateContext.Value(ContextKey(contextKeyValue))

	assert.True(t,
		aggregatedContextGenerated != nil,
		fmt.Sprintf("expected Aggregate context Couldn't be loaded with value \n%s", contextKeyValue),
	)
}

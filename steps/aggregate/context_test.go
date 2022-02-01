package aggregate

import (
	"context"
	"fmt"
	"testing"

	"github.com/Telefonica/golium/steps/http"
	"github.com/stretchr/testify/assert"
)

func TestInitializeContext(t *testing.T) {

	var ctx = context.Background()
	var contextKeyValue = "sharedSession"
	httpContext := http.InitializeContext(ctx)
	aggregateContext := InitializeContext(httpContext)

	httpContextGenerated := aggregateContext.Value(ContextKey(contextKeyValue))

	assert.True(t,
		httpContextGenerated != nil,
		fmt.Sprintf("expected Aggregate context Couldn't be loaded with value \n%s", contextKeyValue),
	)
}

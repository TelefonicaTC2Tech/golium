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

package golium

import (
	"context"

	"github.com/cucumber/godog"
)

// StepsInitializer is an interface to initialize the steps in godog, but extending
// godog initializer with a context.
type StepsInitializer interface {
	// InitializeSteps initializes a set of steps (e.g. http steps) to make them available
	// but using a context.
	InitializeSteps(ctx context.Context, scenCtx *godog.ScenarioContext) context.Context
}

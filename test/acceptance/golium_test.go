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

package main

import (
	"context"
	"testing"

	"github.com/Telefonica/golium"
	"github.com/Telefonica/golium/steps/common"
	"github.com/Telefonica/golium/steps/dns"
	"github.com/Telefonica/golium/steps/http"
	"github.com/Telefonica/golium/steps/jwt"
	"github.com/Telefonica/golium/steps/redis"
	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	launcher := golium.NewLauncher()
	launcher.Launch(InitializeTestSuite, InitializeScenario)
}

func InitializeTestSuite(ctx context.Context, suiteCtx *godog.TestSuiteContext) {
}

func InitializeScenario(ctx context.Context, scenarioCtx *godog.ScenarioContext) {
	stepsInitializers := []golium.StepsInitializer{
		common.Steps{},
		jwt.Steps{},
		http.Steps{},
		dns.Steps{},
		redis.Steps{},
	}
	for _, stepsInitializer := range stepsInitializers {
		ctx = stepsInitializer.InitializeSteps(ctx, scenarioCtx)
	}
}

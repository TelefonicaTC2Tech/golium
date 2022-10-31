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
	"os"
	"testing"

	"github.com/TelefonicaTC2Tech/golium"
	mockhttp "github.com/TelefonicaTC2Tech/golium/mock/http"
	"github.com/TelefonicaTC2Tech/golium/steps/common"
	"github.com/TelefonicaTC2Tech/golium/steps/dns"
	"github.com/TelefonicaTC2Tech/golium/steps/elasticsearch"
	"github.com/TelefonicaTC2Tech/golium/steps/http"
	"github.com/TelefonicaTC2Tech/golium/steps/jwt"
	"github.com/TelefonicaTC2Tech/golium/steps/rabbit"
	"github.com/TelefonicaTC2Tech/golium/steps/redis"

	s3steps "github.com/TelefonicaTC2Tech/golium/steps/aws/s3"
	"github.com/TelefonicaTC2Tech/golium/test/acceptance/steps/aggregated"
	"github.com/TelefonicaTC2Tech/golium/test/acceptance/steps/shared"
	"github.com/cucumber/godog"
)

func TestMain(m *testing.M) {
	launcher := golium.NewLauncher()
	InitializeMocks()
	launcher.Launch(InitializeTestSuite, InitializeScenario)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func InitializeMocks() {
	server := mockhttp.NewServer(9000)
	go server.Start()
}

func InitializeTestSuite(ctx context.Context, suiteCtx *godog.TestSuiteContext) {
}

func InitializeScenario(ctx context.Context, scenarioCtx *godog.ScenarioContext) {
	stepsInitializers := []golium.StepsInitializer{
		common.Steps{},
		jwt.Steps{},
		dns.Steps{},
		redis.Steps{},
		rabbit.Steps{},
		mockhttp.Steps{},
		elasticsearch.Steps{},
		s3steps.Steps{},
		http.Steps{},
		shared.Steps{},
		aggregated.Steps{},
	}
	for _, stepsInitializer := range stepsInitializers {
		ctx = stepsInitializer.InitializeSteps(ctx, scenarioCtx)
	}
	golium.GetContext(ctx).Put("url", http.DefaultTestURL)
}

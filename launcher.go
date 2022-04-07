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
	"os"
	"time"

	"github.com/TelefonicaTC2Tech/golium/cfg"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Launcher is responsible to launch golium (based on godog).
// The default configuration is merged with environment variables.
type Launcher struct {
	log *Logger
}

var goliumLog *Logger

// GetLogger returns the logger for DNS requests and responses.
// If the logger is not created yet, it creates a new instance of Logger.
func GetLogger() *Logger {
	name := "golium"
	if goliumLog == nil {
		goliumLog = LoggerFactory(name)
	}
	return goliumLog
}

// NewLauncher with a default configuration.
func NewLauncher() *Launcher {
	return NewLauncherWithYaml("")
}

// NewLauncherWithYaml with a configuration from a yaml file.
// The yaml file is merged with environment variables.
func NewLauncherWithYaml(path string) *Launcher {
	config := GetConfig()
	if path != "" {
		if err := cfg.LoadYaml(path, config); err != nil {
			logrus.Fatalf("Error configuring golium with yaml file: %s. %s", path, err)
		}
	}
	if err := cfg.LoadEnv(config); err != nil {
		logrus.Fatalf("Error configuring golium with environment variables. %s", err)
	}
	return &Launcher{log: GetLogger()}
}

// Launch golium.
func (l *Launcher) Launch(testSuiteInitializer func(context.Context, *godog.TestSuiteContext),
	scenarioInitializer func(context.Context, *godog.ScenarioContext)) {
	conf := GetConfig()
	godogOpts := godog.Options{
		Output: colors.Colored(os.Stdout),
	}
	godog.BindCommandLineFlags("godog.", &godogOpts)
	pflag.Parse()

	start := time.Now()
	logRecord := l.log.WithField("suite", conf.Suite).WithField("environment", conf.Environment)
	logRecord.Info("Running suite")

	status := godog.TestSuite{
		Name: conf.Suite,
		TestSuiteInitializer: func(suiteContext *godog.TestSuiteContext) {
			ctx := l.initContext()
			testSuiteInitializer(ctx, suiteContext)
		},
		ScenarioInitializer: func(scenarioContext *godog.ScenarioContext) {
			l.configScenarioContext(scenarioContext)
			ctx := l.initContext()
			scenarioInitializer(ctx, scenarioContext)
		},
		Options: &godogOpts,
	}.Run()

	latency := int(time.Since(start).Nanoseconds() / 1000000)
	logRecord = logRecord.WithField("latency", latency).WithField("status", status)
	if status == 0 {
		logRecord.Info("Suite succeeded")
	} else {
		logRecord.Error("Suite failed")
	}
	os.Exit(status)
}

func (l *Launcher) initContext() context.Context {
	ctx := context.Background()
	ctx = InitializeContext(ctx)
	return ctx
}

// configScenarioContext configures the godog.ScenarioContext to include some handlers
// for logging purposes.
// It considers before and after for both steps and scenarios.
func (l *Launcher) configScenarioContext(scenarioContext *godog.ScenarioContext) {
	start := time.Now()
	scenarioContext.StepContext().Before(
		func(ctx context.Context, st *godog.Step) (context.Context, error) {
			l.log.WithField("step", st.Text).Debug("Running step")
			return ctx, nil
		})
	scenarioContext.StepContext().After(
		func(ctx context.Context,
			st *godog.Step,
			status godog.StepResultStatus,
			err error) (context.Context, error) {
			logEntry := l.log.WithField("step", st.Text)
			if err == nil {
				logEntry.Debug("Step succeeded")
			} else {
				logEntry.WithError(err).Error("Step failed")
			}
			return ctx, nil
		})

	scenarioContext.Before(
		func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			l.log.WithField("scenario", sc.Name).Info("Running scenario")
			return ctx, nil
		})
	scenarioContext.After(
		func(ctx context.Context,
			sc *godog.Scenario,
			err error) (context.Context, error) {
			latency := int(time.Since(start).Nanoseconds() / 1000000)
			logEntry := l.log.WithField("latency", latency).WithField("scenario", sc.Name)
			if err == nil {
				logEntry.Info("Scenario succeeded")
			} else {
				logEntry.WithError(err).Error("Scenario failed")
			}
			return ctx, nil
		})
}

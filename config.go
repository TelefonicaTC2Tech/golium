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
	"encoding/json"
	"fmt"
	"path"

	"github.com/TelefonicaTC2Tech/golium/cfg"
	"github.com/sirupsen/logrus"
)

// Global variables for storing the configuration and the environment configuration.
// Global variables simplify the access to this information and this configuration is immutable
// during the test suite execution.

var config = cfg.DefaultConfig
var environment Map

// GetConfig returns the golium configuration.
// This configuration includes relevant information as the environment or the directories
// for some assets or log files.
func GetConfig() *cfg.Config {
	return &config
}

// GetEnvironment returns the environment configuration.
func GetEnvironment() Map {
	if environment == nil {
		environment = initEnvironment()
	}
	return environment
}

// Load the environment configuration from yml files.
// The mandatory environment configuration file is obtained from yml file located at:
//    {config.Dir.Enviroments}/{config.Environment}.yml
// An optional environment configuration file to allow hide or crypt sensitive data could be located at:
//    {config.Dir.Enviroments}/{config.Environment}-private.yml
func initEnvironment() Map {
	pathEnv := path.Join(config.Dir.Environments, fmt.Sprintf("%s.yml", config.Environment))
	logrus.Infof("Loading environment configuration from file: %s", pathEnv)
	env := make(map[string]interface{})
	if err := cfg.LoadYaml(pathEnv, &env); err != nil {
		logrus.Fatalf("Error loading environment configuration from file: %s. %s", pathEnv, err)
	}
	pathEnvPrivate := path.Join(config.Dir.Environments, fmt.Sprintf("%s-private.yml", config.Environment))
	if err := cfg.LoadYaml(pathEnvPrivate, &env); err != nil {
		logrus.Infof("Could not load private environment configuration from file: %s. %s", pathEnvPrivate, err)
	}
	b, err := json.Marshal(env)
	if err != nil {
		logrus.Fatalf("Error converting the yaml to json. %s", err)
	}
	return NewMapFromJSONBytes(b)
}

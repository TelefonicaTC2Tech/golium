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

package cfg

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

// Load reads the configuration from a yml file (at path) and from
// environment variables.
func Load(path string, config interface{}) error {
	if err := LoadYaml(path, config); err != nil {
		return err
	}
	return LoadEnv(config)
}

// LoadYaml reads a yaml file at path in the config struct.
func LoadYaml(path string, config interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	return decoder.Decode(config)
}

// LoadEnv loads the environment variables into config struct.
func LoadEnv(config interface{}) error {
	return envconfig.Process("", config)
}

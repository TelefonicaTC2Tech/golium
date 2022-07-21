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

// Config contains the configuration for golium project.
type Config struct {
	Suite       string    `yaml:"suite" envconfig:"SUITE"`
	Environment string    `yaml:"environment" envconfig:"ENVIRONMENT"`
	Dir         DirConfig `yaml:"dir"`
	Log         LogConfig `yaml:"log"`
}

// DirConfig to configure some configuration directories.
type DirConfig struct {
	Config       string `yaml:"config" envconfig:"DIR_CONFIG"`
	Schemas      string `yaml:"schemas" envconfig:"DIR_SCHEMAS"`
	Environments string `yaml:"environments" envconfig:"DIR_ENVIRONMENTS"`
}

// LogConfig to configure logging.
type LogConfig struct {
	Directory string `yaml:"directory" envconfig:"LOG_DIRECTORY"`
	Level     string `yaml:"level" envconfig:"LOG_LEVEL"`
	Encode    bool   `yaml:"encode" envconfig:"LOG_ENCODE"`
}

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

// DefaultConfig is the default configuration for golium.
// This default configuration can be overriden with environment variables.
var DefaultConfig = Config{
	Suite:       "golium",
	Environment: "local",
	Dir: DirConfig{
		Config:       "./",
		Schemas:      "./schemas",
		Environments: "./environments",
	},
	Log: LogConfig{
		Directory: "./logs",
		Level:     "INFO",
	},
}

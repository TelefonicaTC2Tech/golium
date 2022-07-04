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

package validator

import (
	"errors"
)

var (
	ErrorResponse string
)

type JSONFunctionsMock interface {
	ReplaceMapStringResponse(
		respBody []byte,
		bodyContent interface{},
		replaceValues map[string]interface{},
	) error
	ReplaceStringResponse(respBody []byte,
		bodyContent string,
		replaceValues map[string]interface{},
	) error
}

type JSONServiceMock struct{}

func (m JSONServiceMock) ReplaceMapStringResponse(respBody []byte,
	bodyContent interface{},
	replaceValues map[string]interface{}) error {
	return errors.New(ErrorResponse)
}

func (m JSONServiceMock) ReplaceStringResponse(respBody []byte,
	bodyContent string,
	replaceValues map[string]interface{}) error {
	return errors.New(ErrorResponse)
}

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
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var JSON JSONFunctions = JSONService{}
var unmarshal = json.Unmarshal

type JSONFunctions interface {
	ReplaceMapStringResponse(respBody []byte,
		bodyContent interface{},
		replaceValues map[string]interface{}) error
	ReplaceStringResponse(respBody []byte,
		bodyContent string,
		replaceValues map[string]interface{}) error
}

type JSONService struct{}

// ReplaceMapString
// Validates the response body against the JSON in File replacing values
// when response is a MapString
func (j JSONService) ReplaceMapStringResponse(
	respBody []byte,
	bodyContent interface{},
	replaceValues map[string]interface{},
) error {
	var actual interface{}

	bodyDetails := bodyContent.(map[string]interface{})["details"]
	bodyDetailsMod, _ := bodyDetails.(map[string]interface{})
	newField := fmt.Sprint(replaceValues["field"])
	bodyDetailsMod[newField] = bodyDetailsMod["field_to_replace"]
	delete(bodyDetailsMod, "field_to_replace")
	bodyDetailsMessage := fmt.Sprint(bodyDetailsMod[newField].(map[string]interface{})["message"])

	for key, element := range replaceValues {
		elementStr := fmt.Sprint(element)
		oldStr := fmt.Sprintf("%s_to_replace", key)
		bodyDetailsMessage = strings.Replace(bodyDetailsMessage, oldStr, elementStr, 1)
	}
	bodyDetailsMod[newField].(map[string]interface{})["message"] = bodyDetailsMessage

	if err := unmarshal(respBody, &actual); err != nil {
		return fmt.Errorf("error unmarshalling response body: %w", err)
	}

	if !reflect.DeepEqual(bodyContent, actual) {
		return fmt.Errorf("expected JSON does not match actual, \n%v\n vs \n%s", bodyContent,
			actual)
	}
	return nil
}

// ReplaceString
// Validates the response body against the JSON in File replacing values when response is a String
func (j JSONService) ReplaceStringResponse(
	respBody []byte,
	bodyContent string,
	replaceValues map[string]interface{},
) error {
	for key, element := range replaceValues {
		oldStr := fmt.Sprintf("%s_to_replace", key)
		bodyContent = strings.Replace(bodyContent, oldStr, fmt.Sprint(element), 1)
	}
	if string(respBody) != bodyContent {
		return fmt.Errorf("received body does not match expected, \n%s\n vs \n%s", bodyContent,
			respBody)
	}
	return nil
}

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

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/Telefonica/golium"
)

func GetParamFromJSON(ctx context.Context, file, code, param string) (interface{}, error) {
	data, err := LoadJSONData(file)
	if err != nil {
		return nil, fmt.Errorf("error loading file at %s due to error: %w", file, err)
	}

	dataStruct, err := UnmarshalJSONData(data)
	if err != nil {
		return nil, fmt.Errorf("error unmarsharlling JSON file at %s due to error: %w", file, err)
	}

	paramValue, err := FindValueByCode(dataStruct, code, param)
	if err != nil {
		return nil, fmt.Errorf("Param value: '%s' not found in '%v' due to error: %w", param, dataStruct, err)
	}
	return paramValue, nil
}

func FindValueByCode(dataStruct []map[string]interface{}, code string, param string) (interface{}, error) {
	for _, response := range dataStruct {
		if fmt.Sprint(response["code"]) == code {
			if value, ok := response[param]; ok {
				return value, nil
			}
		}
	}
	return nil, fmt.Errorf("value for param: '%s' with code: '%s' not found", param, code)
}

func LoadJSONData(file string) ([]byte, error) {
	assetsDir := golium.GetConfig().Dir.Schemas
	filePath := fmt.Sprintf("%s%s%s.json", assetsDir, string(os.PathSeparator), file)

	absPath, _ := filepath.Abs(filePath)
	if _, err := os.Stat(absPath); err != nil {
		return nil, fmt.Errorf("file path does not exist: %v", absPath)
	}

	data, readErr := ioutil.ReadFile(absPath)
	if readErr != nil {
		return nil, fmt.Errorf("error reading file at %s due to error: %w", absPath, readErr)
	}
	return data, nil
}

func UnmarshalJSONData(data []byte) ([]map[string]interface{}, error) {
	dataStruct := []map[string]interface{}{}
	err := json.Unmarshal(data, &dataStruct)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data due to error: %w", err)
	}
	return dataStruct, nil
}

// JSONEquals Check if JSON are equal
func JSONEquals(expected interface{}, current interface{}) bool {
	return reflect.DeepEqual(expected, current)
}

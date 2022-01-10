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

	"github.com/Telefonica/golium"
)

func getParamFromJSON(ctx context.Context, file, code, param string) (interface{}, error) {
	assetsDir := golium.GetConfig().Dir.Assets
	filePath := fmt.Sprintf("%s/%s.json", assetsDir, file)
	absPath, _ := filepath.Abs(filePath)
	if _, err := os.Stat(absPath); err != nil {
		return nil, fmt.Errorf("file path does not exist: %v", absPath)
	}

	data, readErr := ioutil.ReadFile(absPath)
	if readErr != nil {
		return nil, fmt.Errorf("error reading file at %s due to error: %w", absPath, readErr)
	}
	dataStruct := []map[string]interface{}{}
	readErr = json.Unmarshal(data, &dataStruct)
	if readErr != nil {
		return nil,
			fmt.Errorf("error unmarshalling JSON at %s due to error: %w", absPath, readErr)
	}

	for _, response := range dataStruct {
		if fmt.Sprint(response["code"]) == code {
			return response[param], nil
		}
	}

	return nil, fmt.Errorf("code not found")
}

package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/cucumber/godog"
)

const (
	parameterError = "error getting parameter from json: %w"
)

// DeleteResponseFields Returns response from schema file without some paramas.
func DeleteResponseFields(
	ctx context.Context,
	code, file string,
	t *godog.Table,
) (interface{}, error) {
	jsonResponseBody, err := GetParam(file, code, "response")
	if err != nil {
		return nil, fmt.Errorf(parameterError, err)
	}
	jsonResponseBodyMap, _ := jsonResponseBody.(map[string]interface{})
	params, err := golium.ConvertTableColumnToArray(ctx, t)
	if err != nil {
		return nil, err
	}
	for _, removeParams := range params {
		delete(jsonResponseBodyMap, removeParams)
	}
	return jsonResponseBody, nil
}

// ModifyResponse Returns modified response from schema file
func ModifyResponse(
	ctx context.Context,
	code, file string,
	t *godog.Table,
) (interface{}, error) {
	jsonResponseBody, err := GetParam(file, code, "response")
	if err != nil {
		return nil, fmt.Errorf("error getting parameter from json: %w", err)
	}
	jsonResponseBodyMap, _ := jsonResponseBody.(map[string]interface{})

	params, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return nil, err
	}
	for key, value := range params {
		if strings.Contains(key, ".") {
			err = processNestedParams(jsonResponseBodyMap, key, value)
			if err != nil {
				return nil, fmt.Errorf("error processing params: %v", err)
			}
		} else {
			_, present := jsonResponseBodyMap[key]
			if !present {
				return nil, fmt.Errorf("error modifying param: param %v does not exists", key)
			}
			jsonResponseBodyMap[key] = value
		}
	}
	return jsonResponseBody, nil
}

func GetBody(
	ctx context.Context,
	code, file string,
) (interface{}, error) {
	message, err := GetParam(file, code, "body")
	if err != nil {
		return nil, fmt.Errorf(parameterError, err)
	}
	return message, nil
}

func ModifyBody(
	ctx context.Context,
	code, file string,
	t *godog.Table,
) (interface{}, error) {
	params, err := golium.ConvertTableToMap(ctx, t)
	if err != nil {
		return nil, err
	}
	message, err := GetParam(file, code, "body")
	if err != nil {
		return nil, fmt.Errorf(parameterError, err)
	}
	messageMap, _ := message.(map[string]interface{})
	for key, value := range params {
		_, present := messageMap[key]
		if !present {
			return nil, fmt.Errorf("error modifying param : param %v does not exists", key)
		}
		messageMap[key] = value
	}
	return message, nil
}

func DeleteBodyFields(
	ctx context.Context,
	code, file string,
	t *godog.Table,
) (interface{}, error) {
	params, err := golium.ConvertTableColumnToArray(ctx, t)
	if err != nil {
		return nil, err
	}
	message, err := GetParam(file, code, "body")
	messageMap, _ := message.(map[string]interface{})
	for _, removeParams := range params {
		delete(messageMap, removeParams)
	}
	if err != nil {
		return nil, fmt.Errorf(parameterError, err)
	}
	return message, nil
}

// processNestedParams Replace nested params split by "." for modified validation.
func processNestedParams(
	jsonResponseBodyMap map[string]interface{},
	key string, value interface{},
) error {
	var aux, lastAux map[string]interface{}
	aux = jsonResponseBodyMap
	keys := strings.Split(key, ".")
	var lastKey string
	for _, key := range keys {
		lastAux = aux
		_, present := aux[key]
		if !present {
			return fmt.Errorf("error modifying nested param: param %v does not exists", key)
		}
		aux, _ = aux[key].(map[string]interface{})
		lastKey = key
	}
	lastAux[lastKey] = value
	return nil
}

// GetParam
// Retrieve values from JSON structure file assets
func GetParam(file, code, param string) (interface{}, error) {
	data, err := LoadData(file)
	if err != nil {
		return nil, fmt.Errorf("error loading file at %s due to error: %w", file, err)
	}

	dataStruct, err := UnmarshalData(data)
	if err != nil {
		return nil, fmt.Errorf("error unmarsharlling JSON file at %s due to error: %w", file, err)
	}

	paramValue, err := FindValueByCode(dataStruct, code, param)
	if err != nil {
		return nil, fmt.Errorf("param value: '%s' not found in '%v' due to error: %w",
			param, dataStruct, err)
	}
	return paramValue, nil
}

// FindValueByCode
// Find value by code and param from dataStruct
func FindValueByCode(dataStruct []map[string]interface{}, code, param string) (interface{}, error) {
	for _, response := range dataStruct {
		if fmt.Sprint(response["code"]) == code {
			if value, ok := response[param]; ok {
				return value, nil
			}
		}
	}
	return nil, fmt.Errorf("value for param: '%s' with code: '%s' not found", param, code)
}

// LoadData
// Load file contents into bytes
func LoadData(file string) ([]byte, error) {
	assetsDir := golium.GetConfig().Dir.Schemas
	filePath := fmt.Sprintf("%s%s%s.json", assetsDir, string(os.PathSeparator), file)

	absPath, _ := filepath.Abs(filePath)
	data, readErr := os.ReadFile(absPath)
	if readErr != nil {
		return nil, fmt.Errorf("error reading file at %s due to error: %w", absPath, readErr)
	}
	return data, nil
}

// UnmarshalData
// Unmarshal bytes to json map struct
func UnmarshalData(data []byte) ([]map[string]interface{}, error) {
	dataStruct := []map[string]interface{}{}
	err := json.Unmarshal(data, &dataStruct)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON data due to error: %w", err)
	}
	return dataStruct, nil
}

// JSONEquals
// Check if unmarshalled JSON maps are equal
func JSONEquals(expected, current interface{}) bool {
	return reflect.DeepEqual(expected, current)
}
